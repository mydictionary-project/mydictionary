package mydictionary

import (
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const icibaCollinsName = "iCIBA Collins"

// IcibaCollinsStruct : iCIBA Collins struct
type IcibaCollinsStruct struct {
	cache CacheStruct
}

// GetServiceName : get service name
func (service *IcibaCollinsStruct) GetServiceName() (value string) {
	value = icibaCollinsName
	return
}

// GetCache : get cache
func (service *IcibaCollinsStruct) GetCache() (cache *CacheStruct) {
	cache = &service.cache
	return
}

// Query : query vocabulary
func (service *IcibaCollinsStruct) Query(vocabularyAsk VocabularyAskStruct) (vocabularyAnswer VocabularyAnswerStruct) {
	var (
		err         error
		queryString string
		item        CacheItemStruct
		document    *goquery.Document
		selection1  *goquery.Selection
		selection2  *goquery.Selection
		selection3  *goquery.Selection
		selection4  *goquery.Selection
		define      string
		selection5  *goquery.Selection
		selection6  *goquery.Selection
		orangeTitle bool
		selection7  *goquery.Selection
	)
	// set
	vocabularyAnswer.SourceName = icibaCollinsName
	vocabularyAnswer.Location.TableType = Online
	vocabularyAnswer.Location.TableIndex = -1
	vocabularyAnswer.Location.ItemIndex = -1
	// query cache
	queryString = url.QueryEscape(vocabularyAsk.Word)
	item, err = service.cache.Query(queryString)
	if err == nil {
		goto SET
	}
	// quert Online
	// no phrase is contained by iCIBA Collins
	if strings.Contains(queryString, "%20") {
		item.Status = "null: IC01"
		goto ADD
	}
	// get page
	document, err = goquery.NewDocument("http://www.iciba.com/" + queryString)
	if err != nil {
		item.Status = err.Error()
		goto ADD
	}
	// get word
	selection1 = document.Find(".keyword")
	if selection1.Nodes == nil {
		item.Status = "null: IC02"
		goto ADD
	}
	item.Word = strings.TrimSpace(selection1.Text())
	// locate
	selection2 = document.Find(".current")
	if selection2.Nodes == nil {
		item.Status = "null: IC03"
		goto ADD
	}
	for i := 0; i < selection2.Size(); i++ {
		if strings.Compare(selection2.Eq(i).Text(), "柯林斯高阶英汉双解学习词典") == 0 {
			selection3 = selection2.Eq(i)
			break
		}
	}
	if selection3 == nil {
		item.Status = "null: IC04"
		goto ADD
	}
	selection4 = selection3.Parent().Parent().Find(".collins-section")
	// get define
	for i := 0; i < selection4.Size(); i++ {
		define = ""
		selection5 = selection4.Eq(i).Find(".section-h")
		define += selection5.Find(".h-order").Text()
		selection6 = selection5.Find(".speech-yellow").Find("span")
		if selection6.Nodes == nil {
			orangeTitle = false
		} else {
			orangeTitle = true
			for j := 0; j < selection6.Size(); j++ {
				define += " "
				define += selection6.Eq(j).Text()
			}
			item.Definition = append(item.Definition, define)
		}
		selection5 = selection4.Eq(i).Find(".section-prep")
		for j := 0; j < selection5.Size(); j++ {
			if orangeTitle {
				define = "  "
			} else {
				define = ""
			}
			selection6 = selection5.Eq(j)
			if j <= 8 {
				define += "0"
			}
			define += selection6.Find(".prep-order-icon").Text()
			selection7 = selection6.Find(".prep-order").Find(".size-chinese").Find("span")
			for k := 0; k < selection7.Size()-1; k++ {
				define += " "
				define += selection7.Eq(k).Text()
			}
			item.Definition = append(item.Definition, define)
			if orangeTitle {
				define = "     "
			} else {
				define = "   "
			}
			define += selection7.Eq(selection7.Size() - 1).Text()
			item.Definition = append(item.Definition, define)
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
