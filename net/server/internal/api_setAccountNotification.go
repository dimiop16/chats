package databag

import (
	"databag/internal/store"
	"gorm.io/gorm"
	"net/http"
)

//SetAccountNotification sets whether notifications should be received
func SetAccountNotification(w http.ResponseWriter, r *http.Request) {

	session, code, err := GetSession(r)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}

	var flag bool
	if err := ParseRequest(r, w, &flag); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(session).Update("push_enabled", flag).Error; res != nil {
			return res
		}
    session.Account.AccountRevision += 1;
		if res := tx.Model(session.Account).Update("account_revision", session.Account.AccountRevision).Error; res != nil {
			return res
		}
		return nil
	})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	SetStatus(&session.Account)
	WriteResponse(w, nil)
}
