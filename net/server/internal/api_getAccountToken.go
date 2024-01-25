package databag

import (
	"net/http"
)

//GetAccountToken retrieve type of account token
func GetAccountToken(w http.ResponseWriter, r *http.Request) {

	accountToken, err := BearerAccountToken(r)
	if err != nil {
		LogMsg("token not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	WriteResponse(w, accountToken.TokenType)
}
