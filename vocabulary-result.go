package mydictionary

// VocabularyResultStruct : result
type VocabularyResultStruct struct {
	Basic   []VocabularyAnswerStruct `json:"Basic"`
	Advance []VocabularyAnswerStruct `json:"Advance"`
}
