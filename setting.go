package mydictionary

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// setting
type settingStruct struct {
	path       string
	Collection []struct {
		Name         string `json:"name"`
		FileName     string `json:"fileName"`
		Readable     bool   `json:"readable"`
		Writable     bool   `json:"writable"`
		OnlineSource string `json:"onlineSource"`
	} `json:"collection"`
	Dictionary []struct {
		Name     string `json:"name"`
		FileName string `json:"fileName"`
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
			//
			// NOTE:
			//
			// 1. Add your services as public members above, like the example below.
			// 2. Add corresponding members in file "mydictionary.setting.json".
			//    For each member, the string of its key should be exactly the same as your service name.
			// 3. Do not edit this note.
			//
			// Example:
			//
			//    Service bool `json:"Service"`
			//
		} `json:"service"`
		Cache struct {
			Enable       bool  `json:"enable"`
			ShelfLifeDay int64 `json:"shelfLifeDay"`
		} `json:"cache"`
		Debug bool `json:"debug"`
	} `json:"online"`
}

// read setting
func (setting *settingStruct) Read() (content string, err error) {
	var buf []byte
	// read
	setting.path = workPath + string(filepath.Separator) + "mydictionary.setting.json"
	buf, err = ioutil.ReadFile(setting.path)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, setting)
	if err != nil {
		return
	}
	// set online
	if setting.Online.Service.BingDictionary {
		onlineList = append(onlineList, new(BingDictionaryStruct))
	}
	if setting.Online.Service.IcibaCollins {
		onlineList = append(onlineList, new(IcibaCollinsStruct))
	}
	if setting.Online.Service.MerriamWebster {
		onlineList = append(onlineList, new(MerriamWebsterStruct))
	}
	//
	// NOTE:
	//
	// 1. Add your functions of services above, like the example below.
	// 2. Do not edit this note.
	//
	// Example:
	//
	//    if setting.Online.Service.Service {
	//    	onlineList = append(onlineList, new(Service))
	//    }
	//
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
	buf, err = json.MarshalIndent(setting, "", "\t")
	content = string(buf)
	return
}

// Write : write setting
func (setting *settingStruct) Write() (err error) {
	var buf []byte
	// write
	buf, err = json.MarshalIndent(setting, "", "\t")
	if err != nil {
		return
	}
	os.Remove(setting.path)
	err = ioutil.WriteFile(setting.path, buf, 0644)
	return
}
