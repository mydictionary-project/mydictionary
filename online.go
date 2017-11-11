package mydictionary

import (
	"strings"

	"github.com/zzc-tongji/bingdictionary4mydictionary"
	"github.com/zzc-tongji/icibacollins4mydictionary"
	"github.com/zzc-tongji/merriamwebster4mydictionary"
	"github.com/zzc-tongji/vocabulary4mydictionary"
	// NOTE:
	//
	// 1. Add your packages of services above, like the example below.
	// 2. Do not edit this note.
	//
	// Example:
	//
	//    "github.com/zzc-tongji/example4mydictionary"
	//
)

// quary vocabulary online
func requestOnline(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		vocabularyAnswerChannel chan vocabulary4mydictionary.VocabularyAnswerStruct
		vocabularyAnswer        vocabulary4mydictionary.VocabularyAnswerStruct
	)
	// prepare
	vocabularyAnswerChannel = make(chan vocabulary4mydictionary.VocabularyAnswerStruct, Setting.Online.length)
	// query
	if Setting.Online.Service.BingDictionary {
		go func() {
			vocabularyAnswerChannel <- bingdictionary4mydictionary.Request(vocabularyAsk)
		}()
	}
	if Setting.Online.Service.IcibaCollins {
		go func() {
			vocabularyAnswerChannel <- icibacollins4mydictionary.Request(vocabularyAsk)
		}()
	}
	if Setting.Online.Service.MerriamWebster {
		go func() {
			vocabularyAnswerChannel <- merriamwebster4mydictionary.Request(vocabularyAsk)
		}()
	}
	// NOTE:
	//
	// 1. Add your functions of services above, like the example below.
	// 2. Do not edit this note.
	//
	// Example:
	//
	//    if Setting.Online.Service.ExambleService {
	//    	go func() {
	//    		vocabularyAnswerChannel <- example4mydictionary.Request(vocabularyAsk)
	//    		}()
	//    }()
	//
	// add to answer list
	for i := 0; i < Setting.Online.length; i++ {
		vocabularyAnswer = <-vocabularyAnswerChannel
		if Setting.Online.Debug {
			vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
		} else if strings.Compare(vocabularyAnswer.Status, vocabulary4mydictionary.Basic) == 0 {
			vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
		}
	}
	return
}
