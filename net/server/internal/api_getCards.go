package databag

import (
	"databag/internal/store"
	"net/http"
	"strconv"
)

//GetCards retrieves all accounts contacts
func GetCards(w http.ResponseWriter, r *http.Request) {
	var cardRevisionSet bool
	var cardRevision int64

	account, code, err := ParamAgentToken(r, false)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}

	card := r.FormValue("revision")
	if card != "" {
		cardRevisionSet = true
		if cardRevision, err = strconv.ParseInt(card, 10, 64); err != nil {
			ErrResponse(w, http.StatusBadRequest, err)
			return
		}
	}

	response := []*Card{}
	if cardRevisionSet {
		var slots []store.CardSlot
		if err := store.DB.Preload("Card").Where("account_id = ? AND revision > ?", account.ID, cardRevision).Find(&slots).Error; err != nil {
			ErrResponse(w, http.StatusInternalServerError, err)
			return
		}
		for _, slot := range slots {
			response = append(response, getCardRevisionModel(&slot))
		}
	} else {
		var slots []store.CardSlot
		if err := store.DB.Preload("Card.Groups.GroupSlot").Where("account_id = ? AND card_id != 0", account.ID).Find(&slots).Error; err != nil {
			ErrResponse(w, http.StatusInternalServerError, err)
			return
		}
		for _, slot := range slots {
			response = append(response, getCardModel(&slot))
		}
	}

	w.Header().Set("Card-Revision", strconv.FormatInt(account.CardRevision, 10))
	WriteResponse(w, response)
}
