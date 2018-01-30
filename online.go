package mydictionary

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zzc-tongji/vocabulary4mydictionary"
)

// quary vocabulary online
func requestOnline(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		vocabularyAnswerChannel chan vocabulary4mydictionary.VocabularyAnswerStruct
		vocabularyAnswer        vocabulary4mydictionary.VocabularyAnswerStruct
	)
	// prepare
	vocabularyAnswerChannel = make(chan vocabulary4mydictionary.VocabularyAnswerStruct, len(onlineList))
	// query
	for i := 0; i < len(onlineList); i++ {
		go func(index int) {
			vocabularyAnswerChannel <- onlineList[index].Query(vocabularyAsk)
		}(i)
	}
	// add to answer list
	for i := 0; i < len(onlineList); i++ {
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
	os.Mkdir(cachePath, 0755)
	for i := 0; i < len(onlineList); i++ {
		err = onlineList[i].GetCache().Read(cachePath+string(filepath.Separator)+onlineList[i].GetServiceName()+".json", Setting.Online.Cache.ShelfLifeDay)
		if err != nil {
			return
		}
	}
	return
}

func saveCache() (success bool, information string) {
	var (
		err  error
		temp string
	)
	success = true
	for i := 0; i < len(onlineList); i++ {
		temp, err = onlineList[i].GetCache().Write()
		if err != nil {
			temp = err.Error() + "\n\n"
			success = false
		}
		information += temp
	}
	return
}
