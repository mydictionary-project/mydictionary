package mydictionary

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/zzc-tongji/rtoa"
)

// setting
type settingStruct struct {
	Collection []struct {
		Name         string `json:"name"`
		FilePath     string `json:"filePath"`
		Readable     bool   `json:"readable"`
		Writable     bool   `json:"writable"`
		OnlineSource string `json:"onlineSource"`
	} `json:"collection"`
	Dictionary []struct {
		Name     string `json:"name"`
		FilePath string `json:"filePath"`
		Readable bool   `json:"readable"`
		Writable bool   `json:"writable"`
	} `json:"dictionary"`
	Online struct {
		Mode        int `json:"mode"`
		modeContent struct {
			userNeed bool
			noFound  bool
			anyway   bool
		}
		Service struct {
			BingDictionary bool `json:"Bing Dictionary"`
			IcibaCollins   bool `json:"iCIBA Collins"`
			MerriamWebster bool `json:"Merriam Webster"`
			// NOTE:
			//
			// 1. Add your services as public members above, like the example below.
			// 2. Add corresponding members in file "mydictionary.setting.json".
			//    For each member, the string of its key should be exactly the same as your service name.
			// 3. Do not edit this note.
			//
			// Example:
			//
			//    ExambleService bool `json:"exambleService"`
			//
		} `json:"service"`
		length int
		Cache  struct {
			Enable       bool  `json:"enable"`
			ShelfLifeDay int64 `json:"shelfLifeDay"`
		} `json:"cache"`
		Debug bool `json:"debug"`
	} `json:"online"`
}

// read setting
func (setting *settingStruct) Read() (content string, err error) {
	var (
		buf  []byte
		path string
	)
	// convert path
	path, err = rtoa.Convert("mydictionary.setting.json", "")
	if err != nil {
		return
	}
	// read
	buf, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, setting)
	if err != nil {
		return
	}
	// set online mode content
	if setting.Online.Mode < 0 {
		setting.Online.Mode = -setting.Online.Mode
	}
	switch setting.Online.Mode % 8 {
	case 0:
		setting.Online.modeContent.userNeed = false
		setting.Online.modeContent.noFound = false
		setting.Online.modeContent.anyway = false
		break
	case 1:
		setting.Online.modeContent.userNeed = true
		setting.Online.modeContent.noFound = false
		setting.Online.modeContent.anyway = false
		break
	case 2:
		setting.Online.modeContent.userNeed = false
		setting.Online.modeContent.noFound = true
		setting.Online.modeContent.anyway = false
		break
	case 3:
		setting.Online.modeContent.userNeed = true
		setting.Online.modeContent.noFound = true
		setting.Online.modeContent.anyway = false
		break
	default:
		setting.Online.modeContent.userNeed = false
		setting.Online.modeContent.noFound = false
		setting.Online.modeContent.anyway = true
		break
	}
	// set online length
	setting.Online.length = 0
	if setting.Online.Service.BingDictionary {
		setting.Online.length++
	}
	if setting.Online.Service.IcibaCollins {
		setting.Online.length++
	}
	if setting.Online.Service.MerriamWebster {
		setting.Online.length++
	}
	// NOTE:
	//
	// 1. Add your services above, like the example below.
	// 2. Do not edit this note.
	//
	// Example:
	//
	//    if setting.Online.Service.ExambleService {
	//    	setting.Online.length++
	//    }
	//
	buf, err = json.MarshalIndent(setting, "", "\t")
	content = string(buf)
	return
}

// Write : write setting
func (setting *settingStruct) Write() (err error) {
	var (
		buf  []byte
		path string
	)
	// convert path
	path, err = rtoa.Convert("mydictionary.setting.json", "")
	if err != nil {
		return
	}
	// write
	buf, err = json.MarshalIndent(setting, "", "\t")
	if err != nil {
		return
	}
	os.Remove(path)
	err = ioutil.WriteFile(path, buf, 0644)
	return
}
