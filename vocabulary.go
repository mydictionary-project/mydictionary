package mydictionary

const (
	// Basic : string for "Status" in "VocabularyAnswerStruct"
	Basic = "basic"
	// Advance : string for "Status" in "VocabularyAnswerStruct"
	Advance = "advance"
	// Collection : int for "TableType" in "LocationStruct"
	Collection = 1
	// Dictionary : int for "TableType" in "LocationStruct"
	Dictionary = 2
	// Online : int for "TableType" in "LocationStruct"
	Online = 3
)

// VocabularyAskStruct : content and option for query
type VocabularyAskStruct struct {
	Word        string `json:"word"`
	Advance     bool   `json:"advance"`
	Online      bool   `json:"online"`
	DoNotRecord bool   `json:"doNotRecord"`
}

// VocabularyAnswerStruct : vocabulary, include word, definition, note and other information
type VocabularyAnswerStruct struct {
	Word         string         `json:"word"`         // `xlsx:wd`
	Definition   []string       `json:"definition"`   // `xlsx:def`
	SerialNumber int            `json:"serialNumber"` // `xlsx:sn`
	QueryCounter int            `json:"queryCounter"` // `xlsx:qc`
	QueryTime    string         `json:"queryTime"`    // `xlsx:qt`
	Note         []string       `json:"note"`         // `xlsx:nt`
	SourceName   string         `json:"sourceName"`
	Status       string         `json:"status"`
	Location     LocationStruct `json:"location"`
}

// VocabularyResultStruct : set of query result
type VocabularyResultStruct struct {
	Basic   []VocabularyAnswerStruct `json:"basic"`
	Advance []VocabularyAnswerStruct `json:"advance"`
}

// VocabularyEditStruct : editor for definition and note of vocabulary
type VocabularyEditStruct struct {
	Location   LocationStruct `json:"location"`
	Definition string         `json:"definition"`
	Note       string         `json:"note"`
}

// LocationStruct : location of vocabulary
type LocationStruct struct {
	TableType  int `json:"tableType"`
	TableIndex int `json:"tableIndex"`
	ItemIndex  int `json:"itemIndex"`
}
