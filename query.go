package mydictionary

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/zzc-tongji/service4mydictionary"
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
	workPath     string
	documentPath string
	cachePath    string
	// Setting : mydictionary setting
	Setting        settingStruct
	initialized    bool
	collectionList collectionListSlice
	dictionaryList dictionaryListSlice
	onlineList     []service4mydictionary.ServiceInterface
	mutex          sync.Mutex
)

func init() {
	initialized = false
}

// Initialize : initialize the library
func Initialize(path []string) (success bool, information string) {
	var (
		err     error
		content string
	)
	// lock
	mutex.Lock()
	// path
	if path != nil {
		switch len(path) {
		case 1:
			workPath = filepath.Clean(path[0])
			documentPath = workPath
			cachePath = workPath
			break
		case 2:
			workPath = filepath.Clean(path[0])
			documentPath = filepath.Clean(path[1])
			cachePath = workPath
			break
		case 3:
			workPath = filepath.Clean(path[0])
			documentPath = filepath.Clean(path[1])
			cachePath = filepath.Clean(path[2])
			break
		default:
			information = "the parameter should be a slice with 1-3 item(s)"
			success = false
			// set flag
			initialized = false
			// unlock
			mutex.Unlock()
			return
		}
	} else {
		information = "the parameter should not be nil"
		success = false
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// read Setting
	content, err = Setting.Read()
	if err != nil {
		information = err.Error() + "\n\n"
		success = false
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// information Setting
	information = content + "\n\n"
	// read collection
	err = collectionList.read(&Setting)
	if err != nil {
		information = err.Error() + "\n\n"
		success = false
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// read dictionary
	err = dictionaryList.read(&Setting)
	if err != nil {
		information = err.Error() + "\n\n"
		success = false
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// cache
	err = loadCache()
	if err != nil {
		information = err.Error() + "\n\n"
		success = false
		// set flag
		initialized = false
		// unlock
		mutex.Unlock()
		return
	}
	// success
	success = true
	// set flag
	initialized = true
	// unlock
	mutex.Unlock()
	return
}

// CheckNetwork : check network
func CheckNetwork() (success bool, information string) {
	var err error
	// lock
	mutex.Lock()
	// begin
	if Setting.Online.Mode == 0 {
		// offline mode
		success = true
		information = "network: offline mode\n\n"
	} else {
		_, err = goquery.NewDocument("https://www.baidu.com/")
		if err != nil {
			// network error
			success = false
			information = "network: " + err.Error() + "\n\n"
		} else {
			success = true
			information = "network: OK\n\n"
		}
	}
	// unlock
	mutex.Unlock()
	return
}

// Query : query
func Query(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (success bool, vocabularyResult vocabulary4mydictionary.VocabularyResultStruct) {
	var (
		vocabularyAnswerList        []vocabulary4mydictionary.VocabularyAnswerStruct
		vocabularyAnswerListPrepare []vocabulary4mydictionary.VocabularyAnswerStruct
		localNoFound                bool
		enableOnline                bool
	)
	// lock
	mutex.Lock()
	// return directly if the library has not been initialized
	success = initialized
	if success == false {
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
		successCache          bool
		informationCache      string
	)
	// lock
	mutex.Lock()
	// return directly if the library has not been initialized
	if initialized == false {
		success = false
		information = "MYDICTIONARY has not been initialized.\n\n"
		// unlock
		mutex.Unlock()
		return
	}
	// write
	successCollection, informationCollection = collectionList.write()
	successDictionary, informationDictionary = dictionaryList.write()
	// cache
	successCache, informationCache = saveCache()
	// merge
	success = successCollection && successDictionary && successCache
	information = informationCollection + informationDictionary + informationCache
	// unlock
	mutex.Unlock()
	return
}

// Edit : edit
func Edit(vocabularyEdit vocabulary4mydictionary.VocabularyEditStruct) (success bool, information string) {
	// lock
	mutex.Lock()
	// return directly if the library has not been initialized
	if initialized == false {
		success = false
		information = "MYDICTIONARY has not been initialized.\n\n"
		// unlock
		mutex.Unlock()
		return
	}
	// check
	// table type
	if vocabularyEdit.Location.TableType != vocabulary4mydictionary.Collection &&
		vocabularyEdit.Location.TableType != vocabulary4mydictionary.Dictionary {
		success = false
		information = "invalid variable \"TableType\"\n\n"
		// unlock
		mutex.Unlock()
		return
	}
	// location
	if vocabularyEdit.Location.TableType == vocabulary4mydictionary.Collection {
		if vocabularyEdit.Location.TableIndex < 0 || vocabularyEdit.Location.TableIndex >= len(collectionList) {
			success = false
			information = "invalid variable \"Location.TableIndex\"\n\n"
			// unlock
			mutex.Unlock()
			return
		}
		if vocabularyEdit.Location.ItemIndex < 0 || vocabularyEdit.Location.ItemIndex >= len(collectionList[vocabularyEdit.Location.TableIndex].content) {
			success = false
			information = "invalid variable \"Location.ItemIndex\"\n\n"
			// unlock
			mutex.Unlock()
			return
		}
	} else { // vocabularyEdit.TableType == vocabulary4mydictionary.Dictionary
		if vocabularyEdit.Location.TableIndex < 0 || vocabularyEdit.Location.TableIndex >= len(dictionaryList) {
			success = false
			information = "invalid variable \"Location.TableIndex\"\n\n"
			// unlock
			mutex.Unlock()
			return
		}
		if vocabularyEdit.Location.ItemIndex < 0 || vocabularyEdit.Location.ItemIndex >= len(dictionaryList[vocabularyEdit.Location.TableIndex].content) {
			success = false
			information = "invalid variable \"Location.ItemIndex\"\n\n"
			// unlock
			mutex.Unlock()
			return
		}
	}
	// edit
	if vocabularyEdit.Location.TableType == vocabulary4mydictionary.Collection {
		collectionList[vocabularyEdit.Location.TableIndex].content[vocabularyEdit.Location.ItemIndex].Note = strings.Split(strings.TrimSpace(vocabularyEdit.Note), "\n")
		collectionList[vocabularyEdit.Location.TableIndex].content[vocabularyEdit.Location.ItemIndex].Definition = strings.Split(strings.TrimSpace(vocabularyEdit.Definition), "\n")
	} else { // vocabularyEdit.TableType == vocabulary4mydictionary.Dictionary
		dictionaryList[vocabularyEdit.Location.TableIndex].content[vocabularyEdit.Location.ItemIndex].Note = strings.Split(strings.TrimSpace(vocabularyEdit.Note), "\n")
		dictionaryList[vocabularyEdit.Location.TableIndex].content[vocabularyEdit.Location.ItemIndex].Definition = strings.Split(strings.TrimSpace(vocabularyEdit.Definition), "\n")
	}
	success = true
	information = "edit: OK\n\n"
	// unlock
	mutex.Unlock()
	return
}
