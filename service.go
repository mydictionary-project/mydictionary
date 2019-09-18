package mydictionary

// ServiceInterface : service interface
type ServiceInterface interface {
	GetServiceName() string
	GetCache() *CacheStruct
	Query(VocabularyAskStruct) VocabularyAnswerStruct
}
