package mydictionary

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/zzc-tongji/vocabulary4mydictionary"
)

// collection
type collectionStruct struct {
	name         string
	readable     bool
	writable     bool
	onlineSource string
	xlsx         *excelize.File
	columnIndex  map[string]int
	content      []vocabulary4mydictionary.VocabularyAnswerStruct
}

// open and check .xlsx file
func (collection *collectionStruct) check(filePath string) (err error) {
	var (
		contentTemp  [][]string
		columnNumber int
	)
	// file -> ram image
	collection.xlsx, err = excelize.OpenFile(filePath)
	if err != nil {
		return
	}
	contentTemp = collection.xlsx.GetRows("sheet1")
	if contentTemp == nil {
		err = fmt.Errorf("incorrect format of file \"%s\": the 1st sheet is empty", collection.xlsx.Path)
		return
	}
	columnNumber = len(contentTemp[0])
	// check existence of sheet header (column) in row 1
	collection.columnIndex = map[string]int{wd: -1, def: -1, sn: -1, qc: -1, qt: -1, nt: -1}
	for i := 0; i < columnNumber; i++ {
		switch contentTemp[0][i] {
		case wd:
			collection.columnIndex[wd] = i
			break
		case def:
			collection.columnIndex[def] = i
			break
		case sn:
			collection.columnIndex[sn] = i
			break
		case qc:
			collection.columnIndex[qc] = i
			break
		case qt:
			collection.columnIndex[qt] = i
			break
		case nt:
			collection.columnIndex[nt] = i
		default:
			break
		}
	}
	if collection.columnIndex[wd] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", collection.xlsx.Path, wd)
		return
	}
	if collection.columnIndex[def] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", collection.xlsx.Path, def)
		return
	}
	if collection.columnIndex[sn] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", collection.xlsx.Path, sn)
		return
	}
	if collection.columnIndex[qc] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", collection.xlsx.Path, qc)
		return
	}
	if collection.columnIndex[qt] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", collection.xlsx.Path, qt)
		return
	}
	if collection.columnIndex[nt] == -1 {
		err = fmt.Errorf("incorrect format of file \"%s\": missing cell \"%s\" in row 1", collection.xlsx.Path, nt)
		return
	}
	return
}

// read data from .xlsx file and put to collection and collection
func (collection *collectionStruct) read(filePath string) (err error) {
	var (
		str              string
		vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct
	)
	if collection.readable {
		// check
		err = collection.check(filePath)
		if err != nil {
			return
		}
		// get space of content
		collection.content = make([]vocabulary4mydictionary.VocabularyAnswerStruct, 0)
		// ram image -> content
		for i := 2; ; i++ {
			// `xlsx:wd` -> .Word
			str = collection.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[wd]), i))
			if strings.Compare(str, "") == 0 {
				break
			}
			vocabularyAnswer.Word = str
			// `xlsx:def` -> .Define
			str = collection.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[def]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Define = strings.Split(str, "\n")
			if len(vocabularyAnswer.Define) == 1 &&
				strings.Compare(vocabularyAnswer.Define[0], "") == 0 {
				vocabularyAnswer.Define = nil
			}
			// `xlsx:sn` -> .SerialNumber
			vocabularyAnswer.SerialNumber, err = strconv.Atoi(collection.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[sn]), i)))
			if err != nil {
				vocabularyAnswer.SerialNumber = i
			}
			// `xlsx:qc` -> .QueryCounter
			vocabularyAnswer.QueryCounter, err = strconv.Atoi(collection.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qc]), i)))
			if err != nil {
				vocabularyAnswer.QueryCounter = 0
			}
			// `xlsx:qt` -> .QueryTime
			vocabularyAnswer.QueryTime = collection.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qt]), i))
			// `xlsx:nt` -> .Note
			str = collection.xlsx.GetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[nt]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Note = strings.Split(str, "\n")
			if len(vocabularyAnswer.Note) == 1 &&
				strings.Compare(vocabularyAnswer.Note[0], "") == 0 {
				vocabularyAnswer.Note = nil
			}
			// others
			vocabularyAnswer.SourceName = collection.name
			vocabularyAnswer.Type = vocabulary4mydictionary.Collection
			vocabularyAnswer.Status = ""
			// add to collection
			collection.content = append(collection.content, vocabularyAnswer)
		}
	}
	err = nil
	return
}

// get data from collection and collection and write to .xlsx file
func (collection *collectionStruct) write() (information string, err error) {
	if collection.readable && collection.writable {
		// content -> ram image
		for i := 0; i < len(collection.content); i++ {
			// set row height
			collection.xlsx.SetRowHeight("sheet1", i+1, collection.xlsx.GetRowHeight("sheet1", 0))
			// .Word -> `xlsx:wd`
			collection.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[wd]), i+2), collection.content[i].Word)
			// .Define -> `xlsx:def`
			collection.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[def]), i+2), strings.Join(collection.content[i].Define, "\n"))
			// .SerialNumber -> `xlsx:sn`
			collection.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[sn]), i+2), collection.content[i].SerialNumber)
			// .QueryCounter -> `xlsx:qc`
			collection.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qc]), i+2), collection.content[i].QueryCounter)
			// .QueryTime -> `xlsx:qt`
			collection.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qt]), i+2), collection.content[i].QueryTime)
			// .Note -> `xlsx:nt`
			collection.xlsx.SetCellValue("sheet1", fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[nt]), i+2), strings.Join(collection.content[i].Note, "\n"))
		}
		// ram image -> file
		err = collection.xlsx.Save()
		if err != nil {
			return
		}
		// output
		information = fmt.Sprintf("Collection \"%s\" has been updated.\n\n", collection.xlsx.Path)
	}
	return
}

// query and update
func (collection *collectionStruct) queryAndUpdate(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct
		tm               time.Time
	)
	if collection.readable {
		for i := 0; i < len(collection.content); i++ {
			// basic
			if strings.Compare(collection.content[i].Word, vocabularyAsk.Word) == 0 {
				if collection.writable {
					// update collection or collection
					if vocabularyAsk.DoNotRecord == false {
						// update
						tm = time.Now()
						collection.content[i].QueryCounter++
						collection.content[i].QueryTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())

					}
				}
				vocabularyAnswer = collection.content[i]
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
				if strings.Contains(collection.content[i].Word, vocabularyAsk.Word) {
					vocabularyAnswer = collection.content[i]
					vocabularyAnswer.Status = vocabulary4mydictionary.Advance
					vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
					goto ADVANCE_END
				}
				for j := 0; j < len(collection.content[i].Define); j++ {
					if strings.Contains(collection.content[i].Define[j], vocabularyAsk.Word) {
						vocabularyAnswer = collection.content[i]
						vocabularyAnswer.Status = vocabulary4mydictionary.Advance
						vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
						goto ADVANCE_END
					}
				}
				for j := 0; j < len(collection.content[i].Note); j++ {
					if strings.Contains(collection.content[i].Note[j], vocabularyAsk.Word) {
						vocabularyAnswer = collection.content[i]
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

// add vocabulary to collection
func (collection *collectionStruct) add(vocabularyAnswerList []vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		existent         bool
		index            int
		vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct
		tm               time.Time
	)
	if collection.readable && collection.writable {
		// only available for collection which is readable and writable
		existent = false
		index = -1
		for i := 0; i < len(vocabularyAnswerList); i++ {
			if strings.Compare(vocabularyAnswerList[i].Status, vocabulary4mydictionary.Basic) == 0 {
				// only for vocabulary with define from basic query
				if vocabularyAnswerList[i].Type == vocabulary4mydictionary.Online {
					// from online: check whether online source index match or not
					if strings.Compare(vocabularyAnswerList[i].SourceName, collection.onlineSource) == 0 {
						index = i
					}
				} else {
					// from dictionary: check existence
					existent = true
				}
			}
		}
		if existent == false && index != -1 {
			// add to collection
			vocabularyAnswer = vocabularyAnswerList[index]
			// prepare
			tm = time.Now()
			vocabularyAnswer.SerialNumber = len(collection.content) + 1
			vocabularyAnswer.QueryCounter = 1
			vocabularyAnswer.QueryTime = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
			vocabularyAnswer.SourceName = collection.name
			vocabularyAnswer.Type = vocabulary4mydictionary.Collection
			vocabularyAnswer.Status = ""
			// add
			collection.content = append(collection.content, vocabularyAnswer)
		}
	}
}
