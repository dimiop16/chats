package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//GetChannelDetail retrieve channel top level attributes
func GetChannelDetail(w http.ResponseWriter, r *http.Request) {

	// scan parameters
	params := mux.Vars(r)
	channelID := params["channelID"]

	var guid string
	var act *store.Account
	tokenType := ParamTokenType(r)
	if tokenType == APPTokenAgent {
		account, code, err := ParamAgentToken(r, false)
		if err != nil {
			ErrResponse(w, code, err)
			return
		}
		act = account
	} else if tokenType == APPTokenContact {
		card, code, err := ParamContactToken(r, true)
		if err != nil {
			ErrResponse(w, code, err)
			return
		}
		act = &card.Account
		guid = card.GUID
	} else {
		ErrResponse(w, http.StatusBadRequest, errors.New("unknown token type"))
		return
	}

	// load channel
	var slot store.ChannelSlot
	if err := store.DB.Preload("Channel.Members.Card.CardSlot").Preload("Channel.Groups.Cards").Preload("Channel.Groups.GroupSlot").Where("account_id = ? AND channel_slot_id = ?", act.ID, channelID).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusNotFound, err)
		} else {
			ErrResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

  video := getBoolConfigValue(CNFEnableVideo, true);
  audio := getBoolConfigValue(CNFEnableAudio, true);
  image := getBoolConfigValue(CNFEnableImage, true);

	// return model data
	if guid != "" {
		if isChannelShared(guid, slot.Channel) {
			WriteResponse(w, getChannelDetailModel(&slot, false, image, audio, video))
		} else {
			ErrResponse(w, http.StatusNotFound, errors.New("channel not shared with requestor"))
			return
		}
	} else {
		WriteResponse(w, getChannelDetailModel(&slot, true, image, audio, video))
	}
}
