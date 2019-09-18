package mydictionary

import (
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	bingDictionaryName = "Bing Dictionary"
	different          = "different"
)

// BingDictionaryStruct : Bing Dictionary struct
type BingDictionaryStruct struct {
	cache CacheStruct
}

// GetServiceName : get service name
func (service *BingDictionaryStruct) GetServiceName() (value string) {
	value = bingDictionaryName
	return
}

// GetCache : get cache
func (service *BingDictionaryStruct) GetCache() (cache *CacheStruct) {
	cache = &service.cache
	return
}

// Query : query vocabulary
func (service *BingDictionaryStruct) Query(vocabularyAsk VocabularyAskStruct) (vocabularyAnswer VocabularyAnswerStruct) {
	var (
		err         error
		queryString string
		item        CacheItemStruct
		document    *goquery.Document
		selection1  *goquery.Selection
		selection2  *goquery.Selection
		selection3  *goquery.Selection
		counter     int
	)
	// set
	vocabularyAnswer.SourceName = bingDictionaryName
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
	// get page
	document, err = goquery.NewDocument("https://cn.bing.com/dict/search?q=" + queryString)
	if err != nil {
		item.Status = err.Error()
		goto ADD
	}
	selection1 = document.Find(".qdef")
	if selection1.Nodes == nil {
		item.Status = "null: BD01"
		goto ADD
	}
	// get word
	selection2 = selection1.Find("#headword")
	if selection2.Nodes == nil {
		item.Status = "null: BD02"
		goto ADD
	}
	item.Word = selection2.Text()
	// get define
	selection3 = selection1.Find("ul").Find("li")
	if selection3.Nodes == nil {
		item.Status = "null: BD03"
		goto ADD
	}
	for i := 0; i < selection3.Size(); i++ {
		item.Definition = append(item.Definition, selection3.Eq(i).Find(".pos").Text()+" "+selection3.Eq(i).Find(".def").Text())
	}
	if len(item.Definition) == 1 && strings.Contains(item.Definition[0], "网络 ") {
		item.Status = "null: BD04"
		goto ADD
	}
	counter = 0
	for i := 0; i < len(item.Definition); i++ {
		if strings.Contains(item.Definition[0], "的过去式") ||
			strings.Contains(item.Definition[0], "的过去分词") ||
			strings.Contains(item.Definition[0], "的现在分词") ||
			strings.Contains(item.Definition[0], "的复数") {
			counter++
		}
	}
	if (len(item.Definition) > 1 && counter >= len(item.Definition)-1) ||
		(len(item.Definition) == 1 && counter == 1) {
		item.Status = "participle"
		goto ADD
	}
	// mark difference
	if strings.Compare(document.Find(".in_tip").Text(), "") != 0 {
		item.Status = different
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
