package mydictionary

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// quary vocabulary online
func queryOnline(vocabularyAsk VocabularyAskStruct) (vocabularyAnswerList []VocabularyAnswerStruct) {
	var (
		vocabularyAnswerChannel chan VocabularyAnswerStruct
		vocabularyAnswer        VocabularyAnswerStruct
	)
	// prepare
	vocabularyAnswerChannel = make(chan VocabularyAnswerStruct, setting.Online.length)
	// query
	if setting.Online.Service.BingDictionary {
		go func() {
			vocabularyAnswerChannel <- queryBingDictionary(vocabularyAsk)
		}()
	}
	if setting.Online.Service.IcibaCollins {
		go func() {
			vocabularyAnswerChannel <- queryIcibaCollins(vocabularyAsk)
		}()
	}
	if setting.Online.Service.MerriamWebster {
		go func() {
			vocabularyAnswerChannel <- queryMerriamWebster(vocabularyAsk)
		}()
	}
	// add to answer list
	for i := 0; i < setting.Online.length; i++ {
		vocabularyAnswer = <-vocabularyAnswerChannel
		if setting.Online.Debug ||
			strings.Compare(vocabularyAnswer.Status, Basic) == 0 ||
			strings.Compare(vocabularyAnswer.Status, Different) == 0 ||
			strings.Compare(vocabularyAnswer.Status, Participle) == 0 {
			vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
		}
	}
	return
}

// quary vocabulary from Bing Dictionary
func queryBingDictionary(vocabularyAsk VocabularyAskStruct) (vocabularyAnswer VocabularyAnswerStruct) {
	var (
		err        error
		document   *goquery.Document
		selection1 *goquery.Selection
		selection2 *goquery.Selection
		selection3 *goquery.Selection
		counter    int
	)
	vocabularyAnswer.SourceName = BingDictionary
	vocabularyAnswer.sourceID = 1
	// get page
	document, err = goquery.NewDocument("https://cn.bing.com/dict/search?q=" + url.QueryEscape(vocabularyAsk.Word))
	if err != nil {
		vocabularyAnswer.Status = err.Error()
		return
	}
	selection1 = document.Find(".qdef")
	if selection1.Nodes == nil {
		vocabularyAnswer.Status = "null: BD01"
		return
	}
	// get word
	selection2 = selection1.Find("#headword")
	if selection2.Nodes == nil {
		vocabularyAnswer.Status = "null: BD02"
		return
	}
	vocabularyAnswer.Word = selection2.Text()
	// get define
	selection3 = selection1.Find("ul").Find("li")
	if selection3.Nodes == nil {
		vocabularyAnswer.Status = "null: BD03"
		return
	}
	for i := 0; i < selection3.Size(); i++ {
		vocabularyAnswer.Define = append(vocabularyAnswer.Define, selection3.Eq(i).Find(".pos").Text()+" "+selection3.Eq(i).Find(".def").Text())
	}
	if len(vocabularyAnswer.Define) == 1 && strings.Contains(vocabularyAnswer.Define[0], "网络 ") {
		vocabularyAnswer.Status = "null: BD04"
		return
	}
	counter = 0
	for i := 0; i < len(vocabularyAnswer.Define); i++ {
		if strings.Contains(vocabularyAnswer.Define[0], "的过去式") ||
			strings.Contains(vocabularyAnswer.Define[0], "的过去分词") ||
			strings.Contains(vocabularyAnswer.Define[0], "的现在分词") ||
			strings.Contains(vocabularyAnswer.Define[0], "的复数") {
			counter++
		}
	}
	if (len(vocabularyAnswer.Define) > 1 && counter >= len(vocabularyAnswer.Define)-1) ||
		(len(vocabularyAnswer.Define) == 1 && counter == 1) {
		vocabularyAnswer.Status = Participle
		return
	}
	// mark difference
	if strings.Compare(vocabularyAsk.Word, vocabularyAnswer.Word) != 0 || strings.Compare(document.Find(".in_tip").Text(), "") != 0 {
		vocabularyAnswer.Status = Different
	} else {
		vocabularyAnswer.Status = Basic
	}
	return
}

// quary vocabulary from iCIBA Collins
func queryIcibaCollins(vocabularyAsk VocabularyAskStruct) (vocabularyAnswer VocabularyAnswerStruct) {
	var (
		err         error
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
	vocabularyAnswer.SourceName = IcibaCollins
	vocabularyAnswer.sourceID = 2
	// no phrase is contained by iCIBA Collins
	if strings.Contains(vocabularyAsk.Word, " ") {
		vocabularyAnswer.Status = "null: IC01"
		return
	}
	// get page
	document, err = goquery.NewDocument("http://www.iciba.com/" + url.QueryEscape(vocabularyAsk.Word))
	if err != nil {
		vocabularyAnswer.Status = err.Error()
		return
	}
	// get word
	selection1 = document.Find(".keyword")
	if selection1.Nodes == nil {
		vocabularyAnswer.Status = "null: IC02"
		return
	}
	vocabularyAnswer.Word = strings.TrimSpace(selection1.Text())
	// locate
	selection2 = document.Find(".current")
	if selection2.Nodes == nil {
		vocabularyAnswer.Status = "null: IC03"
		return
	}
	for i := 0; i < selection2.Size(); i++ {
		if strings.Compare(selection2.Eq(i).Text(), "柯林斯高阶英汉双解学习词典") == 0 {
			selection3 = selection2.Eq(i)
			break
		}
	}
	if selection3 == nil {
		vocabularyAnswer.Status = "null: IC04"
		return
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
			vocabularyAnswer.Define = append(vocabularyAnswer.Define, define)
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
			vocabularyAnswer.Define = append(vocabularyAnswer.Define, define)
			if orangeTitle {
				define = "     "
			} else {
				define = "   "
			}
			define += selection7.Eq(selection7.Size() - 1).Text()
			vocabularyAnswer.Define = append(vocabularyAnswer.Define, define)
		}
	}
	// mark difference
	if strings.Compare(vocabularyAsk.Word, vocabularyAnswer.Word) != 0 {
		vocabularyAnswer.Status = Different
	} else {
		vocabularyAnswer.Status = Basic
	}
	return
}

// quary vocabulary from Merriam Webster
func queryMerriamWebster(vocabularyAsk VocabularyAskStruct) (vocabularyAnswer VocabularyAnswerStruct) {
	var (
		err           error
		document      *goquery.Document
		selection1    *goquery.Selection
		indexList     []int
		selectionList []*goquery.Selection
		selection2    *goquery.Selection
		selection3    *goquery.Selection
	)
	vocabularyAnswer.SourceName = MerriamWebster
	vocabularyAnswer.sourceID = 3
	// no phrase is contained by Merriam Webster
	if strings.Contains(vocabularyAsk.Word, " ") {
		vocabularyAnswer.Status = "null: MW01"
		return
	}
	// get page
	document, err = goquery.NewDocument("https://www.merriam-webster.com/dictionary/" + url.QueryEscape(vocabularyAsk.Word))
	if err != nil {
		vocabularyAnswer.Status = err.Error()
		return
	}
	selection1 = document.Find(".inner-box-wrapper").Find(".card-box-title")
	if selection1.Nodes == nil {
		vocabularyAnswer.Status = "null: MW02"
		return
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
		vocabularyAnswer.Status = "null: MW03"
		return
	}
	// get word
	selection2 = selectionList[0].Find(".word-and-pronunciation").Find("h2")
	if selection2.Nodes == nil {
		vocabularyAnswer.Status = "null: MW04"
		return
	}
	vocabularyAnswer.Word = strings.TrimSpace(selection2.Text())
	// get define
	for i := 0; i < len(selectionList); i++ {
		vocabularyAnswer.Define = append(vocabularyAnswer.Define, selectionList[i].Find(".main-attr").Text())
		selection3 = selectionList[i].Find(".definition-inner-item")
		for j := 0; j < selection3.Size(); j++ {
			vocabularyAnswer.Define = append(vocabularyAnswer.Define, strings.TrimSpace(selection3.Eq(j).Text()))
		}
	}
	// mark difference
	if strings.Compare(vocabularyAsk.Word, vocabularyAnswer.Word) != 0 {
		vocabularyAnswer.Status = Different
	} else {
		vocabularyAnswer.Status = Basic
	}
	return
}
