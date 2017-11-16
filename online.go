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

func loadCache() (err error) {
	// Bing Dictionary
	err = bingdictionary4mydictionary.ReadCache(Setting.Online.Cache.Enable, Setting.Online.Cache.ShelfLifeDay)
	if err != nil {
		return
	}
	// iCIBA Collins
	err = icibacollins4mydictionary.ReadCache(Setting.Online.Cache.Enable, Setting.Online.Cache.ShelfLifeDay)
	if err != nil {
		return
	}
	// Merriam Webster
	err = merriamwebster4mydictionary.ReadCache(Setting.Online.Cache.Enable, Setting.Online.Cache.ShelfLifeDay)
	if err != nil {
		return
	}
	// NOTE:
	//
	// 1. Add your loading functions of services above, like the example below.
	// 2. Do not edit this note.
	//
	// Example:
	//
	//    // Example Service
	//    err = example4mydictionary.ReadCache(Setting.Online.Cache.Enable, Setting.Online.Cache.ShelfLifeDay)
	//    if err != nil {
	//      return
	//    }
	//
	return
}

func saveCache() (success bool, information string) {
	var (
		err  error
		temp string
	)
	success = true
	// Bing Dictionary
	temp, err = bingdictionary4mydictionary.WriteCache()
	if err != nil {
		temp = err.Error() + "\n\n"
		success = false
	}
	information += temp
	// iCIBA Collins
	temp, err = icibacollins4mydictionary.WriteCache()
	if err != nil {
		temp = err.Error() + "\n\n"
		success = false
	}
	information += temp
	// Merriam Webster
	temp, err = merriamwebster4mydictionary.WriteCache()
	if err != nil {
		temp = err.Error() + "\n\n"
		success = false
	}
	information += temp
	// NOTE:
	//
	// 1. Add your loading functions of services above, like the example below.
	// 2. Do not edit this note.
	//
	//    // Example Service
	//    temp, err = example4mydictionary.WriteCache()
	//    if err != nil {
	//    	temp = err.Error() + "\n\n"
	//    	success = false
	//    }
	//    information += temp
	//
	return
}
