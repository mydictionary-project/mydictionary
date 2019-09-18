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
		indexList     []int
		selectionList []*goquery.Selection
		selection2    *goquery.Selection
		selection3    *goquery.Selection
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
	selection1 = document.Find(".inner-box-wrapper").Find(".card-box-title")
	if selection1.Nodes == nil {
		item.Status = "null: MW02"
		goto ADD
	}
	// locate
	for i := 0; i < selection1.Size(); i++ {
		if strings.Contains(selection1.Eq(i).Text(), "English Language Learners") {
			indexList = append(indexList, i)
		}
	}
	for i := 0; i < len(indexList); i++ {
		selectionList = append(selectionList, selection1.Eq(indexList[i]).Parent().Parent().Parent())
	}
	if selectionList == nil {
		item.Status = "null: MW03"
		goto ADD
	}
	// get word
	selection2 = selectionList[0].Find(".word-and-pronunciation").Find("h2")
	if selection2.Nodes == nil {
		item.Status = "null: MW04"
		goto ADD
	}
	item.Word = strings.TrimSpace(selection2.Text())
	// get define
	for i := 0; i < len(selectionList); i++ {
		item.Definition = append(item.Definition, selectionList[i].Find(".main-attr").Text())
		selection3 = selectionList[i].Find(".definition-inner-item")
		for j := 0; j < selection3.Size(); j++ {
			item.Definition = append(item.Definition, strings.TrimSpace(selection3.Eq(j).Text()))
		}
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
