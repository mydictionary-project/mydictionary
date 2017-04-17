package mydictionary

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/zzc-tongji/rtoa"
)

// setting
type settingStruct struct {
	Collection struct {
		FilePath     string `json:"filePath"`
		Readable     bool   `json:"readable"`
		Writable     bool   `json:"writable"`
		OnlineSource int    `json:"onlineSource"`
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
			BingDictionary bool `json:"bingDictionary"`
			IcibaCollins   bool `json:"icibaCollins"`
			MerriamWebster bool `json:"merriamWebster"`
		} `json:"service"`
		length int
		Debug  bool `json:"debug"`
	} `json:"online"`
}

// read setting
func (setting *settingStruct) read() (content string, err error) {
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
	// check collection online source
	if setting.Collection.OnlineSource <= 0 || setting.Collection.OnlineSource >= 4 {
		setting.Collection.OnlineSource = 1
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
	content = strings.TrimRight(string(buf), "\n")
	return
}
