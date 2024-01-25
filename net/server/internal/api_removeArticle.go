package databag

import (
	"databag/internal/store"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//RemoveArticle removes article from account
func RemoveArticle(w http.ResponseWriter, r *http.Request) {

	account, code, err := ParamAgentToken(r, false)
	if err != nil {
		ErrResponse(w, code, err)
		return
	}

	// scan parameters
	params := mux.Vars(r)
	articleID := params["articleID"]

	// load referenced article
	var slot store.ArticleSlot
	if err := store.DB.Preload("Article.Groups.Cards").Where("account_id = ? AND article_slot_id = ?", account.ID, articleID).First(&slot).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ErrResponse(w, http.StatusInternalServerError, err)
		} else {
			ErrResponse(w, http.StatusNotFound, err)
		}
		return
	}
	if slot.Article == nil {
		ErrResponse(w, http.StatusNotFound, errors.New("article has been deleted"))
		return
	}

	// determine affected contact list
	cards := make(map[string]*store.Card)
	for _, group := range slot.Article.Groups {
		for _, card := range group.Cards {
			cards[card.GUID] = &card
		}
	}

	// save and update contact revision
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&slot.Article).Association("Groups").Clear(); res != nil {
			return res
		}
		if res := tx.Delete(&slot.Article).Error; res != nil {
			return res
		}
		slot.Article = nil
		if res := tx.Model(&slot).Update("article_id", 0).Error; res != nil {
			return res
		}
		if res := tx.Model(&slot).Update("revision", account.ArticleRevision+1).Error; res != nil {
			return res
		}
		if res := tx.Model(&account).Update("article_revision", account.ArticleRevision+1).Error; res != nil {
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
		SetContactArticleNotification(account, card)
	}

	WriteResponse(w, getArticleModel(&slot, true, true))
}
