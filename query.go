package mydictionary

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zzc-tongji/vocabulary4mydictionary"
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
		vocabularyAsk        vocabulary4mydictionary.VocabularyAskStruct
		vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct
		tm                   time.Time
	)
	vocabularyAsk.Word = "apple"
	vocabularyAsk.Advance = false
	vocabularyAsk.Online = true
	vocabularyAsk.DoNotRecord = true
	vocabularyAnswerList = requestOnline(vocabularyAsk)
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
			if vocabularyAnswerList[i].Status == vocabulary4mydictionary.Basic {
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
func Query(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyResult VocabularyResultStruct, err error) {
	var (
		vocabularyAnswerList        []vocabulary4mydictionary.VocabularyAnswerStruct
		vocabularyAnswerListPrepare []vocabulary4mydictionary.VocabularyAnswerStruct
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
		if vocabularyAnswerListPrepare[i].Status == vocabulary4mydictionary.Basic {
			localNoFound = false
			break
		}
	}
	enableOnline = setting.Online.modeContent.anyway ||
		(setting.Online.modeContent.noFound && localNoFound) ||
		(setting.Online.modeContent.userNeed && vocabularyAsk.Online)
	if enableOnline {
		vocabularyAnswerList = requestOnline(vocabularyAsk)
		vocabularyAnswerListPrepare = append(vocabularyAnswerListPrepare, vocabularyAnswerList...)
	}
	// build result
	for i := 0; i < len(vocabularyAnswerListPrepare); i++ {
		if strings.Compare(vocabularyAnswerListPrepare[i].Status, vocabulary4mydictionary.Advance) == 0 {
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
