package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//GetCardProfile retrieves contact's public profile
func GetCardProfile(w http.ResponseWriter, r *http.Request) {

	account, code, err := ParamAgentToken(r, false)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}
	cardID := mux.Vars(r)["cardID"]

	var slot store.CardSlot
	if err := store.DB.Preload("Card.Groups").Where("account_id = ? AND card_slot_id = ?", account.ID, cardID).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusNotFound, err)
		} else {
			ErrResponse(w, http.StatusInternalServerError, err)
		}
		return
	}
	if slot.Card == nil {
		ErrResponse(w, http.StatusNotFound, errors.New("referenced empty card slot"))
		return
	}

	WriteResponse(w, getCardProfileModel(&slot))
}
