package databag

import (
	"databag/internal/store"
	"gorm.io/gorm"
	"net/http"
	"os"
)

//RemoveProfile removes account
func RemoveProfile(w http.ResponseWriter, r *http.Request) {

  account, code, err := ParamAgentToken(r, true)
  if err != nil {
    PrintMsg(r)
    ErrResponse(w, code, err)
    return
  }

	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Tag{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.TagSlot{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Asset{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Topic{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.TopicSlot{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.ChannelSlot{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Channel{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Article{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.ArticleSlot{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.CardSlot{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Card{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Group{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.GroupData{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.Group{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.App{}).Error; res != nil {
			return res
		}
		if res := tx.Where("account_id = ?", account.ID).Delete(&store.AccountToken{}).Error; res != nil {
			return res
		}
		if res := tx.Delete(&store.AccountDetail{}, account.AccountDetailID).Error; res != nil {
			return res
		}
		if res := tx.Delete(account).Error; res != nil {
			return res
		}
		return nil
	})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	// delete asset files
	path := getStrConfigValue(CNFAssetPath, APPDefaultPath) + "/" + account.GUID
	if err = os.RemoveAll(path); err != nil {
		ErrMsg(err)
	}

	WriteResponse(w, nil)
}
