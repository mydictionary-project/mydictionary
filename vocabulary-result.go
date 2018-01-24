package mydictionary

import "github.com/zzc-tongji/vocabulary4mydictionary"

// VocabularyResultStruct : result
type VocabularyResultStruct struct {
	Basic   []vocabulary4mydictionary.VocabularyAnswerStruct `json:"basic"`
	Advance []vocabulary4mydictionary.VocabularyAnswerStruct `json:"advance"`
}
