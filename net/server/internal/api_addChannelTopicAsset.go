package databag

import (
	"databag/internal/store"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"hash/crc32"
	"io"
	"net/http"
	"os"
	"strings"
)

//AddChannelTopicAsset adds an asset to a topic and queues it for appropriate transform
func AddChannelTopicAsset(w http.ResponseWriter, r *http.Request) {

	// scan parameters
	params := mux.Vars(r)
	topicID := params["topicID"]
	var transforms []string
	if r.FormValue("transforms") != "" {
		if err := json.Unmarshal([]byte(r.FormValue("transforms")), &transforms); err != nil {
			ErrResponse(w, http.StatusBadRequest, errors.New("invalid asset transform"))
			return
		}
	}

	channelSlot, guid, code, err := getChannelSlot(r, true)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}
	act := &channelSlot.Account

	// check storage
	if full, err := isStorageFull(act); err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	} else if full {
		ErrResponse(w, http.StatusNotAcceptable, errors.New("storage limit reached"))
		return
	}

	// load topic
	var topicSlot store.TopicSlot
	if err = store.DB.Preload("Topic").Where("channel_id = ? AND topic_slot_id = ?", channelSlot.Channel.ID, topicID).First(&topicSlot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusNotFound, err)
		} else {
			ErrResponse(w, http.StatusInternalServerError, err)
		}
		return
	}
	if topicSlot.Topic == nil {
		ErrResponse(w, http.StatusNotFound, errors.New("referenced empty topic"))
		return
	}

	// can only update topic if creator
	if topicSlot.Topic.GUID != guid {
		ErrResponse(w, http.StatusUnauthorized, errors.New("topic not created by you"))
		return
	}

	// avoid async cleanup of file before record is created
	garbageSync.Lock()
	defer garbageSync.Unlock()

	// save new file
	id := uuid.New().String()
	path := getStrConfigValue(CNFAssetPath, APPDefaultPath) + "/" + channelSlot.Account.GUID + "/" + id
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	file, _, err := r.FormFile("asset")
	if err != nil {
		ErrResponse(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	crc, size, err := saveAsset(file, path)
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	assets := []Asset{}
	asset := &store.Asset{}
	asset.AssetID = id
	asset.AccountID = channelSlot.Account.ID
	asset.ChannelID = channelSlot.Channel.ID
	asset.TopicID = topicSlot.Topic.ID
	asset.Status = APPAssetReady
	asset.Size = size
	asset.Crc = crc
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Save(asset).Error; res != nil {
			return res
		}
		assets = append(assets, Asset{AssetID: id, Status: APPAssetReady})
		for _, transform := range transforms {
			asset := &store.Asset{}
			asset.AssetID = uuid.New().String()
			asset.AccountID = channelSlot.Account.ID
			asset.ChannelID = channelSlot.Channel.ID
			asset.TopicID = topicSlot.Topic.ID
			asset.Status = APPAssetWaiting
			asset.TransformID = id
			t := strings.Split(transform, ";")
			if len(t) > 0 {
				asset.Transform = t[0]
			}
			if len(t) > 1 {
				asset.TransformQueue = t[1]
			}
			if len(t) > 2 {
				asset.TransformParams = t[2]
			}
			if res := tx.Save(asset).Error; res != nil {
				return res
			}
			assets = append(assets, Asset{AssetID: asset.AssetID, Transform: transform, Status: APPAssetWaiting})
		}
		if res := tx.Model(&topicSlot.Topic).Update("detail_revision", act.ChannelRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&topicSlot).Update("revision", act.ChannelRevision+1).Error; res != nil {
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

	// invoke transcoder
	transcode()

	WriteResponse(w, &assets)
}

func isStorageFull(act *store.Account) (full bool, err error) {

	storage := getNumConfigValue(CNFStorage, 0)
	if storage == 0 {
		return
	}

	var assets []store.Asset
	if err = store.DB.Where("account_id = ?", act.ID).Find(&assets).Error; err != nil {
		return
	}

	var size int64
	for _, asset := range assets {
		size += asset.Size
	}
	if size >= storage {
		full = true
	}

	return
}

func saveAsset(src io.Reader, path string) (crc uint32, size int64, err error) {

	output, res := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if res != nil {
		err = res
		return
	}
	defer output.Close()

	// prepare hash
	table := crc32.MakeTable(crc32.IEEE)

	// compute has as data is saved
	data := make([]byte, 4096)
	for {
		n, res := src.Read(data)
    if n > 0 {
      crc = crc32.Update(crc, table, data[:n])
      output.Write(data[:n])
    }
		if res != nil {
			if res == io.EOF {
				break
			}
			err = res
			return
		}
	}

	// read size
	info, ret := output.Stat()
	if ret != nil {
		err = ret
		return
	}
	size = info.Size()
	return
}
