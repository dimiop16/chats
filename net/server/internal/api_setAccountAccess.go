package databag

import (
	"databag/internal/store"
	"encoding/hex"
	"errors"
  "time"
	"github.com/theckman/go-securerandom"
	"gorm.io/gorm"
	"net/http"
)

//SetAccountAccess creates token to gain access to account
func SetAccountAccess(w http.ResponseWriter, r *http.Request) {

	token, _, res := AccessToken(r)
	if res != nil || token.TokenType != APPTokenReset {
    time.Sleep(APPUsernameWait * time.Millisecond);
		ErrResponse(w, http.StatusUnauthorized, res)
		return
	}
	if token.Account == nil {
		ErrResponse(w, http.StatusUnauthorized, errors.New("invalid reset token"))
		return
	}
	account := token.Account

  // parse authentication token
  appName := r.FormValue("appName")
  appVersion := r.FormValue("appVersion")
  platform := r.FormValue("platform")
  deviceToken := r.FormValue("deviceToken")

  // parse requested notifications
  var notifications []Notification
  if err := ParseRequest(r, w, &notifications); err != nil {
    ErrMsg(err);
  }

	// gernate app token
	data, err := securerandom.Bytes(APPTokenSize)
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	access := hex.EncodeToString(data)

	// create app entry
	session := store.Session{
		AccountID:   account.GUID,
		Token:       access,
    AppName:     appName,
    AppVersion:  appVersion,
    Platform:    platform,
    PushToken:   deviceToken,
    PushEnabled: true,
	}

	// save app and delete token
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&session).Error; res != nil {
			return res
		}
    for _, notification := range notifications {
      pushEvent := &store.PushEvent{}
      pushEvent.SessionID = session.ID
      pushEvent.Event = notification.Event
      pushEvent.MessageTitle = notification.MessageTitle
      pushEvent.MessageBody = notification.MessageBody
      if res := tx.Save(pushEvent).Error; res != nil {
        return res
      }
    }
		if res := tx.Delete(token).Error; res != nil {
			return res
		}
		return nil
	})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

  login := LoginAccess{
    GUID: account.GUID,
    AppToken: account.GUID + "." + access,
    Created:  session.Created,
    PushSupported: getBoolConfigValue(CNFPushSupported, true),
  }

	WriteResponse(w, login)
}
