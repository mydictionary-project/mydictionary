package mydictionary

import "github.com/zzc-tongji/rtoa"

// dictionart list
type dictionaryListSlice []dictionaryStruct

// read all collection and dictionary from .xlsx file
func (dictionaryList *dictionaryListSlice) read(setting *settingStruct) (err error) {
	var (
		dictionary dictionaryStruct
		str        string
	)
	// read collection
	dictionary.name = collection
	dictionary.readable = setting.Collection.Readable
	dictionary.writable = setting.Collection.Writable
	dictionary.id = 0
	dictionary.onlineSource = setting.Collection.OnlineSource
	str, err = rtoa.Convert(setting.Collection.FilePath, "")
	if err != nil {
		return
	}
	err = dictionary.read(str)
	if err != nil {
		return
	}
	*dictionaryList = append(*dictionaryList, dictionary)
	// read dictionary
	for i := 0; i < len(setting.Dictionary); i++ {
		dictionary.name = setting.Dictionary[i].Name
		dictionary.readable = setting.Dictionary[i].Readable
		dictionary.writable = setting.Dictionary[i].Writable
		dictionary.id = -(i + 1)
		dictionary.onlineSource = ""
		str, err = rtoa.Convert(setting.Dictionary[i].FilePath, "")
		if err != nil {
			return
		}
		err = dictionary.read(str)
		if err != nil {
			return
		}
		*dictionaryList = append(*dictionaryList, dictionary)
	}
	return
}

// write all collection and dictionary to .xlsx file
func (dictionaryList *dictionaryListSlice) write() (success bool, information string) {
	var (
		err  error
		temp string
	)
	success = true
	for i := 0; i < len(*dictionaryList); i++ {
		temp, err = (*dictionaryList)[i].write()
		if err != nil {
			temp = err.Error() + "\n\n"
			success = false
		}
		information += temp
	}
	return
}
