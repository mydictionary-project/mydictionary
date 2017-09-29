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
	// title in .xlsx file
	wd  = "Word"
	def = "Define"
	sn  = "SN"
	qc  = "QC"
	qt  = "QT"
	nt  = "Note"
)

var (
	initialized    bool
	setting        settingStruct
	collectionList collectionListSlice
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
	information = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\nmydictionary v1.0.1\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	// read setting
	content, err = setting.read()
	if err != nil {
		// set flag
		initialized = false
		return
	}
	// information setting
	information += fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n%s\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), content)
	// read collection
	err = collectionList.read(&setting)
	if err != nil {
		// set flag
		initialized = false
		return
	}
	// read dictionary
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
		timeString           string
		temp                 bool
	)
	vocabularyAsk.Word = "apple"
	vocabularyAsk.Advance = false
	vocabularyAsk.Online = true
	vocabularyAsk.DoNotRecord = true
	tm = time.Now()
	timeString = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	if setting.Online.Mode == 0 {
		// offline mode
		information = timeString + "offline mode\n\n"
		pass = true
	} else {
		// set setting.Online.Debug as false temporarily
		temp = setting.Online.Debug
		setting.Online.Debug = true
		// get result of online query
		vocabularyAnswerList = requestOnline(vocabularyAsk)
		pass = true
		for i := 0; i < len(vocabularyAnswerList); i++ {
			if vocabularyAnswerList[i].Status == vocabulary4mydictionary.Basic {
				information += fmt.Sprintf("%s: OK\n\n", vocabularyAnswerList[i].SourceName)
			} else {
				information += fmt.Sprintf("%s: %s\n\n", vocabularyAnswerList[i].SourceName, vocabularyAnswerList[i].Status)
				pass = false
			}
		}
		if strings.Compare(information, "") != 0 {
			information = timeString + information
		}
		// set setting.Online.Debug as its original value
		setting.Online.Debug = temp
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
	// collection: query and update
	for i := 0; i < len(collectionList); i++ {
		vocabularyAnswerList = collectionList[i].queryAndUpdate(vocabularyAsk)
		vocabularyAnswerListPrepare = append(vocabularyAnswerListPrepare, vocabularyAnswerList...)
	}
	// dictionary: query and update
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
			for i := 0; i < len(collectionList); i++ {
				collectionList[i].add(vocabularyResult.Basic)
			}
		}
	}
	return
}

// Save : save to .xlsx file
func Save() (success bool, information string) {
	var (
		tm                    time.Time
		successCollection     bool
		informationCollection string
		successDictionary     bool
		informationDictionary string
	)
	// avoid multiple writing (.xlsx file) at the same time
	mutexContent.Lock()
	successCollection, informationCollection = collectionList.write()
	successDictionary, informationDictionary = dictionaryList.write()
	mutexContent.Unlock()
	// merge
	success = successCollection && successDictionary
	information = informationCollection + informationDictionary
	// add time
	if information != "" {
		tm = time.Now()
		information = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second()) + information
	}
	return
}
