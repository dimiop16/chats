package databag

import (
	"databag/internal/store"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

//GetChannels retrieves channels under account or shared with contact
func GetChannels(w http.ResponseWriter, r *http.Request) {
	var channelRevisionSet bool
	var channelRevision int64
	var viewRevisionSet bool
	var viewRevision int64
	var typesSet bool
	var types []string

	channel := r.FormValue("channelRevision")
	if channel != "" {
		var err error
		channelRevisionSet = true
		if channelRevision, err = strconv.ParseInt(channel, 10, 64); err != nil {
			ErrResponse(w, http.StatusBadRequest, err)
			return
		}
	}

	view := r.FormValue("viewRevision")
	if view != "" {
		var err error
		viewRevisionSet = true
		if viewRevision, err = strconv.ParseInt(view, 10, 64); err != nil {
			ErrResponse(w, http.StatusBadRequest, err)
			return
		}
	}

	schemas := r.FormValue("types")
	if schemas != "" {
		var err error
		typesSet = true
		dec := json.NewDecoder(strings.NewReader(schemas))
		if dec.Decode(&types) != nil {
			ErrResponse(w, http.StatusBadRequest, err)
			return
		}
	}

	response := []*Channel{}
	tokenType := ParamTokenType(r)
	if tokenType == APPTokenAgent {

		account, code, err := ParamAgentToken(r, false)
		if err != nil {
			ErrResponse(w, code, err)
			return
		}

		var slots []store.ChannelSlot
		if channelRevisionSet {
			if err := store.DB.Preload("Channel").Where("account_id = ? AND revision > ?", account.ID, channelRevision).Find(&slots).Error; err != nil {
				ErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			if err := store.DB.Preload("Channel.Topics", func(db *gorm.DB) *gorm.DB {
				return store.DB.Order("topics.id DESC")
			}).Preload("Channel.Members.Card.CardSlot").Preload("Channel.Groups.GroupSlot").Where("account_id = ? AND channel_id != 0", account.ID).Find(&slots).Error; err != nil {
				ErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		}

		for _, slot := range slots {
			if !typesSet || hasChannelType(types, slot.Channel) || slot.Channel == nil {
				if channelRevisionSet {
					response = append(response, getChannelRevisionModel(&slot, true))
				} else if slot.Channel != nil {
          video := getBoolConfigValue(CNFEnableVideo, true);
          audio := getBoolConfigValue(CNFEnableAudio, true);
          image := getBoolConfigValue(CNFEnableImage, true);
					response = append(response, getChannelModel(&slot, true, true, image, audio, video))
				}
			}
		}

		w.Header().Set("Channel-Revision", strconv.FormatInt(account.ChannelRevision, 10))

	} else if tokenType == APPTokenContact {

		card, code, err := ParamContactToken(r, true)
		if err != nil {
			ErrResponse(w, code, err)
			return
		}

		if viewRevisionSet || channelRevisionSet {
			if viewRevision != card.ViewRevision {
				ErrResponse(w, http.StatusGone, errors.New("channel view has changed"))
				return
			}
		}

		account := &card.Account
		var slots []store.ChannelSlot
		if channelRevisionSet {
			if err := store.DB.Preload("Channel.Members.Card").Preload("Channel.Groups.Cards").Where("account_id = ? AND revision > ?", account.ID, channelRevision).Find(&slots).Error; err != nil {
				ErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			if err := store.DB.Preload("Channel.Topics", func(db *gorm.DB) *gorm.DB {
				return store.DB.Order("topics.id DESC")
			}).Preload("Channel.Members.Card").Preload("Channel.Groups.Cards").Where("account_id = ? AND channel_id != 0", account.ID).Find(&slots).Error; err != nil {
				ErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		}

		for _, slot := range slots {
			if !typesSet || hasChannelType(types, slot.Channel) || slot.Channel == nil {
				shared := isChannelShared(card.GUID, slot.Channel)
				if channelRevisionSet {
					response = append(response, getChannelRevisionModel(&slot, shared))
				} else if shared {
          video := getBoolConfigValue(CNFEnableVideo, true);
          audio := getBoolConfigValue(CNFEnableAudio, true);
          image := getBoolConfigValue(CNFEnableImage, true);
					response = append(response, getChannelModel(&slot, true, false, image, audio, video))
				}
			}
		}

		w.Header().Set("Channel-Revision", strconv.FormatInt(account.ChannelRevision, 10))
		w.Header().Set("View-Revision", strconv.FormatInt(card.ViewRevision, 10))

	} else {
		ErrResponse(w, http.StatusBadRequest, errors.New("invalid token type"))
		return
	}

	WriteResponse(w, response)
}

func isChannelShared(guid string, channel *store.Channel) bool {
	if channel == nil {
		return false
	}
	for _, member := range channel.Members {
		if guid == member.Card.GUID {
			return true
		}
	}
	for _, group := range channel.Groups {
		for _, card := range group.Cards {
			if guid == card.GUID {
				return true
			}
		}
	}
	return false
}

func hasChannelType(types []string, channel *store.Channel) bool {
	if channel == nil {
		return false
	}
	for _, schema := range types {
		if schema == channel.DataType {
			return true
		}
	}
	return false
}
