package mydictionary

import (
	"fmt"
	"time"
)

// VocabularyAnswerStruct : define and information
type VocabularyAnswerStruct struct {
	Word         string   `json:"word"`         // `xlsx:wd`
	Define       []string `json:"define"`       // `xlsx:def`
	SerialNumber int      `json:"serialNumber"` // `xlsx:sn`
	QueryCounter int      `json:"queryCounter"` // `xlsx:qc`
	QueryTime    string   `json:"queryTime"`    // `xlsx:qt`
	SourceName   string   `json:"sourceName"`
	sourceID     int
	Status       string `json:"status"`
}

func (VocabularyAnswer *VocabularyAnswerStruct) clean() {
	VocabularyAnswer.Word = ""
	VocabularyAnswer.Define = nil
	VocabularyAnswer.SerialNumber = 0
	VocabularyAnswer.QueryCounter = 0
	VocabularyAnswer.QueryTime = ""
	VocabularyAnswer.SourceName = ""
	VocabularyAnswer.sourceID = 255
	VocabularyAnswer.Status = ""
}

// update query counter and query time of a vocabulary which is found in collection or dictionary
func (VocabularyAnswer *VocabularyAnswerStruct) update() {
	var tm time.Time = time.Now()
	VocabularyAnswer.QueryCounter++
	VocabularyAnswer.QueryTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
}

// do something before adding a new vocabulary to collection
func (VocabularyAnswer *VocabularyAnswerStruct) prepare(dictionary *dictionaryStruct) {
	var tm time.Time = time.Now()
	VocabularyAnswer.SerialNumber = len(dictionary.content) + 1
	VocabularyAnswer.QueryCounter = 1
	VocabularyAnswer.QueryTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	VocabularyAnswer.SourceName = collection
	VocabularyAnswer.sourceID = 0
	VocabularyAnswer.Status = ""
}
