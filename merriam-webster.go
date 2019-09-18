package mydictionary

import (
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const merriamWebsterName = "Merriam Webster"

// MerriamWebsterStruct : Merriam Webster struct
type MerriamWebsterStruct struct {
	cache CacheStruct
}

// GetServiceName : get service name
func (service *MerriamWebsterStruct) GetServiceName() (value string) {
	value = merriamWebsterName
	return
}

// GetCache : get cache
func (service *MerriamWebsterStruct) GetCache() (cache *CacheStruct) {
	cache = &service.cache
	return
}

// Query : query vocabulary
func (service *MerriamWebsterStruct) Query(vocabularyAsk VocabularyAskStruct) (vocabularyAnswer VocabularyAnswerStruct) {
	var (
		err           error
		queryString   string
		item          CacheItemStruct
		document      *goquery.Document
		selection1    *goquery.Selection
		selection2    *goquery.Selection
		selection3    *goquery.Selection
		selection4    *goquery.Selection
	)
	// set
	vocabularyAnswer.SourceName = merriamWebsterName
	vocabularyAnswer.Location.TableType = Online
	vocabularyAnswer.Location.TableIndex = -1
	vocabularyAnswer.Location.ItemIndex = -1
	// query cache
	queryString = url.QueryEscape(vocabularyAsk.Word)
	item, err = service.cache.Query(queryString)
	if err == nil {
		goto SET
	}
	// query online
	// no phrase is contained by Merriam Webster
	if strings.Contains(queryString, "%20") {
		item.Status = "null: MW01"
		goto ADD
	}
	// get page
	document, err = goquery.NewDocument("https://www.merriam-webster.com/dictionary/" + queryString)
	if err != nil {
		item.Status = err.Error()
		goto ADD
	}
	// locate
	selection1 = document.Find(".learners-def")
	if selection1.Nodes == nil {
		item.Status = "null: MW02"
		goto ADD
	}
	selection2 = selection1.Eq(0);
	// get word
	selection3 = selection2.Prev().Find("em")
	if selection3.Nodes == nil {
		item.Status = "null: MW03"
		goto ADD
	}
	item.Word = selection3.Text()
	// get define
	selection4 = selection2.Find(".dtText");
	if selection4.Nodes == nil {
		item.Status = "null: MW04"
		goto ADD
	}
	for i := 0; i < selection4.Size(); i++ {
		item.Definition = append(item.Definition, selection4.Eq(i).Text());
	}
	// mark difference
	if strings.Compare(queryString, url.QueryEscape(item.Word)) != 0 {
		item.Status = "different"
	} else {
		item.Status = Basic
	}
ADD:
	// add to cache
	item.QueryString = queryString
	item.CreationTime = time.Now().Unix()
	service.cache.Add(item)
SET:
	// set
	vocabularyAnswer.Word = item.Word
	vocabularyAnswer.Definition = item.Definition
	vocabularyAnswer.Status = item.Status
	return
}
