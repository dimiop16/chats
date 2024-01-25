package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//SetGroupSubject sets subject of group under account
func SetGroupSubject(w http.ResponseWriter, r *http.Request) {

	account, code, err := ParamAgentToken(r, true)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}

	// scan parameters
	params := mux.Vars(r)
	groupID := params["groupID"]

	var subject Subject
	if err := ParseRequest(r, w, &subject); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	// load referenced group
	var slot store.GroupSlot
	if err := store.DB.Preload("Group.GroupData").Where("account_id = ? AND group_slot_id = ?", account.ID, groupID).First(&slot).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusInternalServerError, err)
		} else {
			ErrResponse(w, http.StatusNotFound, err)
		}
		return
	}
	if slot.Group == nil {
		ErrResponse(w, http.StatusNotFound, errors.New("group has been deleted"))
		return
	}

	// save and update contact revision
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&slot.Group.GroupData).Update("data", subject.Data).Error; res != nil {
			return res
		}
		if res := tx.Model(&slot.Group).Update("data_type", subject.DataType).Error; res != nil {
			return res
		}
		if res := tx.Model(&slot).Update("revision", account.GroupRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&account).Update("group_revision", account.GroupRevision+1).Error; res != nil {
			return res
		}
		return nil
	})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	SetStatus(account)
	WriteResponse(w, getGroupModel(&slot))
}
