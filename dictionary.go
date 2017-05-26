package mydictionary

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/zzc-tongji/vocabulary4mydictionary"
)

// dictionary
type dictionaryStruct struct {
	name         string
	readable     bool
	writable     bool
	id           int // distinction: id < 0 --> dictionary, id == 0 --> collection, else --> online (impossible)
	onlineSource string
	xlsx         *excelize.File
	columnIndex  map[string]int
	content      []vocabulary4mydictionary.VocabularyAnswerStruct
}

// open and check .xlsx file
func (dictionary *dictionaryStruct) check(filePath string) (err error) {
	var (
		contentTemp  [][]string
		columnNumber int
	)
	// file -> ram image
	dictionary.xlsx, err = excelize.OpenFile(filePath)
	if err != nil {
		return
	}
	contentTemp = dictionary.xlsx.GetRows("sheet1")
	if contentTemp == nil {
		err = fmt.Errorf("incorrect format of file \"%s\": the 1st sheet is empty", dictionary.xlsx.Path)
		return
	}
	columnNumber = len(contentTemp[0])
	// check existence of sheet header (column) in row 1
	dictionary.columnIndex = map[string]int{wd: -1, def: -1, sn: -1, qc: -1, qt: -1}
	for i := 0; i < columnNumber; i++ {
		switch contentTemp[0][i] {
		case wd:
			dictionary.columnIndex[wd] = i
			break
		case def:
			dictionary.columnIndex[def] = i
			break
		case sn:
			dictionary.columnIndex[sn] = i
			break
		case qc:
			dictionary.columnIndex[qc] = i
			break
		case qt:
			dictionary.columnIndex[qt] = i
			break
		default:
			break
		}
	}
	if dictionary.columnIndex[wd] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", dictionary.xlsx.Path, wd)
		return
	}
	if dictionary.columnIndex[def] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", dictionary.xlsx.Path, def)
		return
	}
	if dictionary.columnIndex[sn] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", dictionary.xlsx.Path, sn)
		return
	}
	if dictionary.columnIndex[qc] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", dictionary.xlsx.Path, qc)
		return
	}
	if dictionary.columnIndex[qt] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", dictionary.xlsx.Path, qt)
		return
	}
	return
}

// read data from .xlsx file and put to collection and dictionary
func (dictionary *dictionaryStruct) read(filePath string) (err error) {
	var (
		str              string
		vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct
	)
	if dictionary.readable {
		// check
		err = dictionary.check(filePath)
		if err != nil {
			return
		}
		// get space of content
		dictionary.content = make([]vocabulary4mydictionary.VocabularyAnswerStruct, 0)
		// ram image -> content
		for i := 2; ; i++ {
			// `xlsx:wd` -> .Word
			str = dictionary.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[wd]), i))
			if strings.Compare(str, "") == 0 {
				break
			}
			vocabularyAnswer.Word = str
			// `xlsx:def` -> .Define
			str = dictionary.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[def]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Define = strings.Split(str, "\n")
			// `xlsx:sn` -> .SerialNumber
			vocabularyAnswer.SerialNumber, err = strconv.Atoi(dictionary.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[sn]), i)))
			if err != nil {
				vocabularyAnswer.SerialNumber = i
			}
			// `xlsx:qc` -> .QueryCounter
			vocabularyAnswer.QueryCounter, err = strconv.Atoi(dictionary.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qc]), i)))
			if err != nil {
				vocabularyAnswer.QueryCounter = 0
			}
			// `xlsx:qt` -> .QueryTime
			vocabularyAnswer.QueryTime = dictionary.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qt]), i))
			/**/
			// others
			vocabularyAnswer.SourceName = dictionary.name
			vocabularyAnswer.Type = vocabulary4mydictionary.Dictionary
			vocabularyAnswer.Status = ""
			// add to dictionary
			dictionary.content = append(dictionary.content, vocabularyAnswer)
		}
	}
	err = nil
	return
}

// get data from collection and dictionary and write to .xlsx file
func (dictionary *dictionaryStruct) write() (information string, err error) {
	if dictionary.readable && dictionary.writable {
		// content -> ram image
		for i := 0; i < len(dictionary.content); i++ {
			if dictionary.id == 0 {
				// set row height
				dictionary.xlsx.SetRowHeight("sheet1", i+1, dictionary.xlsx.GetRowHeight("sheet1", 0))
				// .Word -> `xlsx:wd`
				dictionary.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[wd]), i+2), dictionary.content[i].Word)
				// .Define -> `xlsx:def`
				dictionary.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[def]), i+2), strings.Join(dictionary.content[i].Define, "\n"))
				// .SerialNumber -> `xlsx:sn`
				dictionary.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[sn]), i+2), dictionary.content[i].SerialNumber)
			}
			// .QueryCounter -> `xlsx:qc`
			dictionary.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qc]), i+2), dictionary.content[i].QueryCounter)
			// .QueryTime -> `xlsx:qt`
			dictionary.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qt]), i+2), dictionary.content[i].QueryTime)
		}
		// ram image -> file
		err = dictionary.xlsx.Save()
		if err != nil {
			return
		}
		// output
		if dictionary.id == 0 {
			information = fmt.Sprintf("Collection \"%s\" has been updated.\n\n", dictionary.xlsx.Path)
		} else {
			information = fmt.Sprintf("Dictionary \"%s\" has been updated.\n\n", dictionary.xlsx.Path)
		}
	}
	return
}

func (dictionary *dictionaryStruct) queryAndUpdate(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct
		tm               time.Time
	)
	if dictionary.readable {
		for i := 0; i < len(dictionary.content); i++ {
			// basic
			if strings.Compare(dictionary.content[i].Word, vocabularyAsk.Word) == 0 {
				if dictionary.writable {
					// update dictionary or collection
					if vocabularyAsk.DoNotRecord == false {
						// uodate
						tm = time.Now()
						dictionary.content[i].QueryCounter++
						dictionary.content[i].QueryTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())

					}
				}
				vocabularyAnswer = dictionary.content[i]
				vocabularyAnswer.Status = vocabulary4mydictionary.Basic
				vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
				if vocabularyAsk.Advance {
					continue
				} else {
					break
				}
			}
			if vocabularyAsk.Advance {
				// advance
				if strings.Contains(dictionary.content[i].Word, vocabularyAsk.Word) {
					vocabularyAnswer = dictionary.content[i]
					vocabularyAnswer.Status = vocabulary4mydictionary.Advance
					vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
					goto ADVANCE_END
				}
				for j := 0; j < len(dictionary.content[i].Define); j++ {
					if strings.Contains(dictionary.content[i].Define[j], vocabularyAsk.Word) {
						vocabularyAnswer = dictionary.content[i]
						vocabularyAnswer.Status = vocabulary4mydictionary.Advance
						vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
						goto ADVANCE_END
					}
				}
			ADVANCE_END:
			}
		}
	}
	return
}

// add vocabulary which is not existent in dictionary to collection
func (dictionary *dictionaryStruct) add(vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		existent         bool
		index            int
		vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct
		tm               time.Time
	)
	if dictionary.id == 0 && dictionary.readable && dictionary.writable {
		// only available for collection, which is readable and writable
		existent = false
		index = -1
		for i := 0; i < len(vocabularyAnswerList); i++ {
			if strings.Compare(vocabularyAnswerList[i].Status, vocabulary4mydictionary.Basic) == 0 {
				// only for vocabulary with define from basic query
				if vocabularyAnswerList[i].Type != vocabulary4mydictionary.Online {
					// from collection or dictionary: check existence
					existent = true
				} else {
					// from online: check whether online source index match or not
					if strings.Compare(vocabularyAnswerList[i].SourceName, dictionary.onlineSource) == 0 {
						index = i
					}
				}
			}
		}
		if existent == false && index != -1 {
			// add to collection
			vocabularyAnswer = vocabularyAnswerList[index]
			// prepare
			tm = time.Now()
			vocabularyAnswer.SerialNumber = len(dictionary.content) + 1
			vocabularyAnswer.QueryCounter = 1
			vocabularyAnswer.QueryTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
			vocabularyAnswer.SourceName = collection
			vocabularyAnswer.Type = vocabulary4mydictionary.Collection
			vocabularyAnswer.Status = ""
			// add
			dictionary.content = append(dictionary.content, vocabularyAnswer)
		}
	}
}
