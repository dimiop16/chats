package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//SetChannelSubject updates channel subject for account holder
func SetChannelSubject(w http.ResponseWriter, r *http.Request) {

	account, code, err := ParamAgentToken(r, false)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}

	// scan parameters
	params := mux.Vars(r)
	channelID := params["channelID"]

	var subject Subject
	if err := ParseRequest(r, w, &subject); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	// load referenced channel
	var slot store.ChannelSlot
	if err := store.DB.Preload("Channel.Members.Card").Preload("Channel.Groups.Cards").Where("account_id = ? AND channel_slot_id = ?", account.ID, channelID).First(&slot).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusInternalServerError, err)
		} else {
			ErrResponse(w, http.StatusNotFound, err)
		}
		return
	}
	if slot.Channel == nil {
		ErrResponse(w, http.StatusNotFound, errors.New("channel has been deleted"))
		return
	}

	// determine affected contact list
	cards := make(map[string]store.Card)
	for _, member := range slot.Channel.Members {
		cards[member.Card.GUID] = member.Card
	}
	for _, group := range slot.Channel.Groups {
		for _, card := range group.Cards {
			cards[card.GUID] = card
		}
	}

	// save and update contact revision
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&slot.Channel).Update("data", subject.Data).Error; res != nil {
			return res
		}
		if res := tx.Model(&slot.Channel).Update("data_type", subject.DataType).Error; res != nil {
			return res
		}
		if res := tx.Model(&slot.Channel).Update("detail_revision", account.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&slot).Update("revision", account.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&account).Update("channel_revision", account.ChannelRevision+1).Error; res != nil {
			return res
		}
		return nil
	})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	// notify contacts of content change
	SetStatus(account)
	for _, card := range cards {
		SetContactChannelNotification(account, &card)
	}

  video := getBoolConfigValue(CNFEnableVideo, true);
  audio := getBoolConfigValue(CNFEnableAudio, true);
  image := getBoolConfigValue(CNFEnableImage, true);
	WriteResponse(w, getChannelModel(&slot, true, true, image, audio, video))
}
