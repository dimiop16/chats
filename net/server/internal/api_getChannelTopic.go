package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//GetChannelTopic retrieves channel topic
func GetChannelTopic(w http.ResponseWriter, r *http.Request) {

	// scan parameters
	params := mux.Vars(r)
	topicID := params["topicID"]

	channelSlot, _, code, err := getChannelSlot(r, false)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}

	// load topic
	var topicSlot store.TopicSlot
	if err = store.DB.Preload("Topic.Assets").Where("channel_id = ? AND topic_slot_id = ?", channelSlot.Channel.ID, topicID).First(&topicSlot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
		ErrResponse(w, code, err)
		return
	}

	WriteResponse(w, getTopicModel(&topicSlot))
}

func isMember(guid string, members []store.Member) bool {
	for _, member := range members {
		if guid == member.Card.GUID {
			return true
		}
	}
	return false
}

func isViewer(guid string, groups []store.Group) bool {
	for _, group := range groups {
		for _, card := range group.Cards {
			if guid == card.GUID {
				return true
			}
		}
	}
	return false
}

func getChannelSlot(r *http.Request, member bool) (slot store.ChannelSlot, guid string, code int, err error) {

	// scan parameters
	params := mux.Vars(r)
	channelID := params["channelID"]

	// validate contact access
	var account *store.Account
	tokenType := ParamTokenType(r)
	if tokenType == APPTokenAgent {
		account, code, err = ParamAgentToken(r, false)
		if err != nil {
			return
		}
		guid = account.GUID
	} else if tokenType == APPTokenContact {
		var card *store.Card
		card, code, err = ParamContactToken(r, true)
		if err != nil {
			return
		}
		account = &card.Account
		guid = card.GUID
	} else {
		err = errors.New("unknown token type")
		code = http.StatusBadRequest
		return
	}

	// load channel
	if err = store.DB.Preload("Account").Preload("Channel.Members.Card").Preload("Channel.Groups.Cards").Where("account_id = ? AND channel_slot_id = ?", account.ID, channelID).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
		return
	}
	if slot.Channel == nil {
		err = errors.New("referenced empty channel")
		code = http.StatusNotFound
		return
	}

	// validate access to channel
	if tokenType == APPTokenContact {
		if member && !isMember(guid, slot.Channel.Members) {
			err = errors.New("contact is not a channel member")
			code = http.StatusUnauthorized
			return
		} else if !isViewer(guid, slot.Channel.Groups) && !isMember(guid, slot.Channel.Members) {
			err = errors.New("contact is not a channel viewer")
			code = http.StatusUnauthorized
			return
		}
	}

	return
}
