package mydictionary

import (
	"os"
	"path/filepath"
	"strings"
)

// quary vocabulary online
func requestOnline(vocabularyAsk VocabularyAskStruct) (vocabularyAnswerList []VocabularyAnswerStruct) {
	var (
		vocabularyAnswerChannel chan VocabularyAnswerStruct
		vocabularyAnswer        VocabularyAnswerStruct
	)
	// prepare
	vocabularyAnswerChannel = make(chan VocabularyAnswerStruct, len(onlineList))
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
		} else if strings.Compare(vocabularyAnswer.Status, Basic) == 0 {
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
