package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//RemoveChannelTopicAsset removes topic asset if invoker created topic
func RemoveChannelTopicAsset(w http.ResponseWriter, r *http.Request) {

	// scan parameters
	params := mux.Vars(r)
	topicID := params["topicID"]
	assetID := params["assetID"]

	channelSlot, guid, code, err := getChannelSlot(r, true)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}
	act := &channelSlot.Account

	// load asset
	var asset store.Asset
	if err = store.DB.Preload("Topic.TopicSlot").Where("channel_id = ? AND asset_id = ?", channelSlot.Channel.ID, assetID).First(&asset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusNotFound, err)
		} else {
			ErrResponse(w, http.StatusInternalServerError, err)
		}
		return
	}
	if asset.Topic.TopicSlot.TopicSlotID != topicID {
		ErrResponse(w, http.StatusNotFound, errors.New("invalid topic"))
		return
	}

	// can only update topic if creator
	if asset.Topic.GUID != guid {
		ErrResponse(w, http.StatusUnauthorized, errors.New("topic not created by you"))
		return
	}

	// delete asset record
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Delete(&asset).Error; res != nil {
			return res
		}
		if res := tx.Model(&asset.Topic).Update("detail_revision", act.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&asset.Topic.TopicSlot).Update("revision", act.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&channelSlot.Channel).Update("topic_revision", act.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&channelSlot).Update("revision", act.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(act).Update("channel_revision", act.ChannelRevision+1).Error; res != nil {
			return res
		}
		return nil
	})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	// cleanup files from deleted record
	go garbageCollect(act)

	// determine affected contact list
	cards := make(map[string]store.Card)
	for _, member := range channelSlot.Channel.Members {
		cards[member.Card.GUID] = member.Card
	}
	for _, group := range channelSlot.Channel.Groups {
		for _, card := range group.Cards {
			cards[card.GUID] = card
		}
	}

	// notify
	SetStatus(act)
	for _, card := range cards {
		SetContactChannelNotification(act, &card)
	}

	WriteResponse(w, nil)
}
