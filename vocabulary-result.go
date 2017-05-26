package mydictionary

import "github.com/zzc-tongji/vocabulary4mydictionary"

// VocabularyResultStruct : result
type VocabularyResultStruct struct {
	Basic   []vocabulary4mydictionary.VocabularyAnswerStruct `json:"Basic"`
	Advance []vocabulary4mydictionary.VocabularyAnswerStruct `json:"Advance"`
}
