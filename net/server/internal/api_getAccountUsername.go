package databag

import (
	"databag/internal/store"
	"errors"
  "time"
	"net/http"
	"strings"
)

type accountUsername struct {
	Username string
}

//GetAccountUsername determines if username is assignable
func GetAccountUsername(w http.ResponseWriter, r *http.Request) {

  if r.FormValue("agent") != "" {
    _, code, res := ParamAgentToken(r, false)
    if res != nil {
      ErrResponse(w, code, res)
      return
    }
  } else if r.FormValue("token") != "" {
		token, _, res := AccessToken(r)
		if res != nil || token.TokenType != APPTokenCreate {
      time.Sleep(APPUsernameWait * time.Second);
			ErrResponse(w, http.StatusUnauthorized, res)
			return
		}
	} else {
		if available, err := getAvailableAccounts(); err != nil {
			ErrResponse(w, http.StatusInternalServerError, err)
			return
		} else if available == 0 {
			ErrResponse(w, http.StatusForbidden, errors.New("no open accounts available"))
			return
		}
	}

	username := r.URL.Query().Get("name")
	if username == "" {
    WriteResponse(w, true);
		return
	}

	if strings.Contains(username, " ") || strings.Contains(username, "\t") {
		WriteResponse(w, false)
		return
	}

	var accounts []accountUsername
	if err := store.DB.Model(&store.Account{}).Where("username = ?", username).Find(&accounts).Error; err != nil {
		LogMsg("failed to query accounts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(accounts) != 0 {
		WriteResponse(w, false)
		return
	}

	handle := strings.ToLower(username)
	if err := store.DB.Model(&store.Account{}).Where("handle = ?", handle).Find(&accounts).Error; err != nil {
		LogMsg("failed to query accounts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(accounts) != 0 {
		WriteResponse(w, false)
		return
	}

	WriteResponse(w, true)
}
