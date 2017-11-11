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
	// version
	version = "v2.1.0"
	// title in .xlsx file
	wd  = "Word"
	def = "Define"
	sn  = "SN"
	qc  = "QC"
	qt  = "QT"
	nt  = "Note"
)

var (
	// Setting : mydictionary setting
	Setting        settingStruct
	initialized    bool
	tm             time.Time
	timeString     string
	collectionList collectionListSlice
	dictionaryList dictionaryListSlice
	mutex          sync.Mutex
)

func init() {
	initialized = false
}

// Initialize : initialize the library
func Initialize() (information string, err error) {
	var content string
	// lock
	mutex.Lock()
	// return directly if the library has been initialized
	if initialized {
		err = errors.New("the program should be initialized only once")
		// unlock
		mutex.Unlock()
		return
	}
	// get time
	tm = time.Now()
	timeString = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	// title
	information = timeString + "mydictionary " + version + "\n\n"
	// read Setting
	content, err = Setting.Read()
	if err != nil {
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// information Setting
	information += timeString + content + "\n\n"
	// read collection
	err = collectionList.read(&Setting)
	if err != nil {
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// read dictionary
	err = dictionaryList.read(&Setting)
	if err != nil {
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// set flag
	initialized = true
	// unlock
	mutex.Unlock()
	return
}

// CheckNetwork : check network
func CheckNetwork() (pass bool, information string) {
	var (
		vocabularyAsk        vocabulary4mydictionary.VocabularyAskStruct
		vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct
		temp                 bool
	)
	// lock
	mutex.Lock()
	// begin
	if Setting.Online.Mode == 0 {
		// offline mode
		information = "offline mode\n\n"
		pass = true
	} else {
		// set word for online query
		vocabularyAsk.Word = "apple"
		vocabularyAsk.Advance = false
		vocabularyAsk.Online = true
		vocabularyAsk.DoNotRecord = true
		// set Setting.Online.Debug as false temporarily
		temp = Setting.Online.Debug
		Setting.Online.Debug = true
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
		if strings.Compare(information, "") == 0 {
			information = "online mode, but no service is enabled\n\n"
		}
		// set Setting.Online.Debug as its original value
		Setting.Online.Debug = temp
	}
	// get time
	tm = time.Now()
	timeString = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	information = timeString + information
	// unlock
	mutex.Unlock()
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
	// lock
	mutex.Lock()
	// return directly if the library has not been initialized
	if initialized == false {
		err = errors.New("the program have not been initialized")
		// unlock
		mutex.Unlock()
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
	enableOnline = Setting.Online.modeContent.anyway ||
		(Setting.Online.modeContent.noFound && localNoFound) ||
		(Setting.Online.modeContent.userNeed && vocabularyAsk.Online)
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
	// unlock
	mutex.Unlock()
	return
}

// Save : save to .xlsx file
func Save() (success bool, information string) {
	var (
		successCollection     bool
		informationCollection string
		successDictionary     bool
		informationDictionary string
	)
	// lock
	mutex.Lock()
	// write
	successCollection, informationCollection = collectionList.write()
	successDictionary, informationDictionary = dictionaryList.write()
	// merge
	success = successCollection && successDictionary
	information = informationCollection + informationDictionary
	// get
	if strings.Compare(information, "") != 0 {
		tm = time.Now()
		timeString = fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]\n\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
		information = timeString + information
	}
	// unlock
	mutex.Unlock()
	return
}
