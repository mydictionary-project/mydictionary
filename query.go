package mydictionary

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	// name of collection
	collection = "collection"
	// title in .xlsx file
	wd  = "Word"
	def = "Define"
	sn  = "SN"
	qc  = "QC"
	qt  = "QT"
	// BingDictionary : string for "SourceName" in "VocabularyAnswerStruct"
	BingDictionary = "Bing Dictionary"
	// IcibaCollins : string for "SourceName" in "VocabularyAnswerStruct"
	IcibaCollins = "iCIBA Collins"
	// MerriamWebster : string for "SourceName" in "VocabularyAnswerStruct"
	MerriamWebster = "Merriam Webster"
	// Basic : string for "Status" in "VocabularyAnswerStruct"
	Basic = "basic"
	// Advance : string for "Status" in "VocabularyAnswerStruct"
	Advance = "advance"
	// Different : string for "Status" in "VocabularyAnswerStruct"
	Different = "different"
	// Participle : string for "Status" in "VocabularyAnswerStruct"
	Participle = "participle"
)

var (
	initialized    bool
	setting        settingStruct
	dictionaryList dictionaryListSlice
	mutexContent   sync.Mutex
)

func init() {
	initialized = false
}

// Initialize : initialize the library
func Initialize() (information string, err error) {
	var (
		tm      time.Time
		content string
	)
	// return directly if the library has been initialized
	if initialized {
		err = errors.New("the program should be initialized only once")
		return
	}
	// title
	tm = time.Now()
	information = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\nmydictionary\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	// read setting
	content, err = setting.read()
	if err != nil {
		// set flag
		initialized = false
		return
	}
	// information setting
	information += fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n%s\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), content)
	// read collection and dictionary
	err = dictionaryList.read(&setting)
	if err != nil {
		// set flag
		initialized = false
		return
	}
	// set flag
	initialized = true
	return
}

// CheckNetwork : check network
func CheckNetwork() (pass bool, information string) {
	var (
		vocabularyAnswerList []VocabularyAnswerStruct
		tm                   time.Time
	)
	vocabularyAnswerList = queryOnline(VocabularyAskStruct{"apple", false, true, true})
	tm = time.Now()
	information = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	if len(vocabularyAnswerList) == 0 {
		if setting.Online.Mode == 0 {
			information += "offline mode\n\n"
			pass = true
		} else {
			information += "network error\n\n"
			pass = false
		}
	} else {
		pass = true
		for i := 0; i < len(vocabularyAnswerList); i++ {
			if vocabularyAnswerList[i].Status == Basic {
				information += fmt.Sprintf("%s: OK\n\n", vocabularyAnswerList[i].SourceName)
			} else {
				information += fmt.Sprintf("%s: %s\n\n", vocabularyAnswerList[i].SourceName, vocabularyAnswerList[i].Status)
				pass = false
			}
		}
	}
	return
}

// Query : query
func Query(vocabularyAsk VocabularyAskStruct) (vocabularyResult VocabularyResultStruct, err error) {
	var (
		vocabularyAnswerList        []VocabularyAnswerStruct
		vocabularyAnswerListPrepare []VocabularyAnswerStruct
		localNoFound                bool
		enableOnline                bool
	)
	// return directly if the library has not been initialized
	if initialized == false {
		err = errors.New("the program have not been initialized")
		return
	}
	// collection and dictionary: query and update
	for i := 0; i < len(dictionaryList); i++ {
		vocabularyAnswerList = dictionaryList[i].queryAndUpdate(vocabularyAsk)
		vocabularyAnswerListPrepare = append(vocabularyAnswerListPrepare, vocabularyAnswerList...)
	}
	// online: query
	localNoFound = true
	for i := 0; i < len(vocabularyAnswerListPrepare); i++ {
		if vocabularyAnswerListPrepare[i].Status == Basic {
			localNoFound = false
			break
		}
	}
	enableOnline = setting.Online.modeContent.anyway ||
		(setting.Online.modeContent.noFound && localNoFound) ||
		(setting.Online.modeContent.userNeed && vocabularyAsk.Online)
	if enableOnline {
		vocabularyAnswerList = queryOnline(vocabularyAsk)
		vocabularyAnswerListPrepare = append(vocabularyAnswerListPrepare, vocabularyAnswerList...)
	}
	// build result
	for i := 0; i < len(vocabularyAnswerListPrepare); i++ {
		if strings.Compare(vocabularyAnswerListPrepare[i].Status, Advance) == 0 {
			vocabularyResult.Advance = append(vocabularyResult.Advance, vocabularyAnswerListPrepare[i])
		} else {
			vocabularyResult.Basic = append(vocabularyResult.Basic, vocabularyAnswerListPrepare[i])
		}
	}
	if enableOnline {
		// add online to collection
		if vocabularyAsk.DoNotRecord == false {
			dictionaryList[0].add(vocabularyResult.Basic)
		}
	}
	return
}

// Save : save to .xlsx file
func Save() (success bool, information string) {
	var tm time.Time
	// avoid multiple writing (.xlsx file) at the same time
	mutexContent.Lock()
	success, information = dictionaryList.write()
	mutexContent.Unlock()
	if information != "" {
		tm = time.Now()
		information = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second()) + information
	}
	return
}
