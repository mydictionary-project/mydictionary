package mydictionary

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// collection
type collectionStruct struct {
	name         string
	readable     bool
	writable     bool
	onlineSource string
	index        int
	xlsx         *excelize.File
	sheetName    string
	columnIndex  map[string]int
	content      []VocabularyAnswerStruct
}

// open and check .xlsx file
func (collection *collectionStruct) check(filePath string) (err error) {
	var (
		recheckCounter    int
		recheckUpperLimit int
		contentTemp       [][]string
		columnNumber      int
	)
	// file -> ram image
	collection.xlsx, err = excelize.OpenFile(filePath)
	if err != nil {
		return
	}
	// recheck
	recheckCounter = 0
	recheckUpperLimit = 10
RECHECK:
	collection.sheetName = collection.xlsx.GetSheetMap()[1]
	if strings.Compare(collection.sheetName, "") == 0 {
		// no worksheet in workbook: create a worksheet
		if strings.Compare(collection.name, "") == 0 {
			collection.xlsx.NewSheet("collection")
		} else {
			collection.xlsx.NewSheet(collection.name)
		}
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	contentTemp = collection.xlsx.GetRows(collection.sheetName)
	if len(contentTemp) == 0 {
		// empty worksheet: create row 1
		collection.xlsx.SetCellValue(collection.sheetName, "A1", sn)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
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
		collection.xlsx.SetCellValue(collection.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", wd)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if collection.columnIndex[def] == -1 {
		collection.xlsx.SetCellValue(collection.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", def)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if collection.columnIndex[sn] == -1 {
		collection.xlsx.SetCellValue(collection.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", sn)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if collection.columnIndex[qc] == -1 {
		collection.xlsx.SetCellValue(collection.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", qc)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if collection.columnIndex[qt] == -1 {
		collection.xlsx.SetCellValue(collection.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", qt)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if collection.columnIndex[nt] == -1 {
		collection.xlsx.SetCellValue(collection.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", nt)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	return
}

// read data from .xlsx file and put to collection and collection
func (collection *collectionStruct) read(filePath string) (err error) {
	var (
		str              string
		vocabularyAnswer VocabularyAnswerStruct
	)
	if collection.readable {
		// check
		err = collection.check(filePath)
		if err != nil {
			return
		}
		// get space of content
		collection.content = make([]VocabularyAnswerStruct, 0)
		// ram image -> content
		for i := 2; ; i++ {
			// `xlsx:wd` -> .Word
			str = collection.xlsx.GetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[wd]), i))
			if strings.Compare(str, "") == 0 {
				break
			}
			vocabularyAnswer.Word = str
			// `xlsx:def` -> .Define
			str = collection.xlsx.GetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[def]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Definition = strings.Split(str, "\n")
			if len(vocabularyAnswer.Definition) == 1 &&
				strings.Compare(vocabularyAnswer.Definition[0], "") == 0 {
				vocabularyAnswer.Definition = nil
			}
			// `xlsx:sn` -> .SerialNumber
			vocabularyAnswer.SerialNumber, err = strconv.Atoi(collection.xlsx.GetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[sn]), i)))
			if err != nil {
				vocabularyAnswer.SerialNumber = i - 1
			}
			// `xlsx:qc` -> .QueryCounter
			vocabularyAnswer.QueryCounter, err = strconv.Atoi(collection.xlsx.GetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qc]), i)))
			if err != nil {
				vocabularyAnswer.QueryCounter = 0
			}
			// reset err
			err = nil
			// `xlsx:qt` -> .QueryTime
			vocabularyAnswer.QueryTime = collection.xlsx.GetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qt]), i))
			// `xlsx:nt` -> .Note
			str = collection.xlsx.GetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[nt]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Note = strings.Split(str, "\n")
			if len(vocabularyAnswer.Note) == 1 &&
				strings.Compare(vocabularyAnswer.Note[0], "") == 0 {
				vocabularyAnswer.Note = nil
			}
			// others
			vocabularyAnswer.SourceName = collection.name
			vocabularyAnswer.Location.TableType = Collection
			vocabularyAnswer.Status = ""
			// add to collection
			collection.content = append(collection.content, vocabularyAnswer)
			// set location
			collection.content[len(collection.content)-1].Location.TableIndex = collection.index
			collection.content[len(collection.content)-1].Location.ItemIndex = len(collection.content) - 1
		}
	}
	return
}

// get data from collection and collection and write to .xlsx file
func (collection *collectionStruct) write() (information string, err error) {
	if collection.readable && collection.writable {
		// content -> ram image
		for i := 0; i < len(collection.content); i++ {
			// set row height
			collection.xlsx.SetRowHeight(collection.sheetName, i+1, collection.xlsx.GetRowHeight(collection.sheetName, 0))
			// .Word -> `xlsx:wd`
			collection.xlsx.SetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[wd]), i+2), collection.content[i].Word)
			// .Define -> `xlsx:def`
			collection.xlsx.SetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[def]), i+2), strings.Join(collection.content[i].Definition, "\n"))
			// .SerialNumber -> `xlsx:sn`
			collection.xlsx.SetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[sn]), i+2), collection.content[i].SerialNumber)
			// .QueryCounter -> `xlsx:qc`
			collection.xlsx.SetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qc]), i+2), collection.content[i].QueryCounter)
			// .QueryTime -> `xlsx:qt`
			collection.xlsx.SetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[qt]), i+2), collection.content[i].QueryTime)
			// .Note -> `xlsx:nt`
			collection.xlsx.SetCellValue(collection.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(collection.columnIndex[nt]), i+2), strings.Join(collection.content[i].Note, "\n"))
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
func (collection *collectionStruct) queryAndUpdate(vocabularyAsk VocabularyAskStruct) (vocabularyAnswerList []VocabularyAnswerStruct) {
	var (
		vocabularyAnswer VocabularyAnswerStruct
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
				vocabularyAnswer.Status = Basic
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
					vocabularyAnswer.Status = Advance
					vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
					goto ADVANCE_END
				}
				for j := 0; j < len(collection.content[i].Definition); j++ {
					if strings.Contains(collection.content[i].Definition[j], vocabularyAsk.Word) {
						vocabularyAnswer = collection.content[i]
						vocabularyAnswer.Status = Advance
						vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
						goto ADVANCE_END
					}
				}
				for j := 0; j < len(collection.content[i].Note); j++ {
					if strings.Contains(collection.content[i].Note[j], vocabularyAsk.Word) {
						vocabularyAnswer = collection.content[i]
						vocabularyAnswer.Status = Advance
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
func (collection *collectionStruct) add(vocabularyAnswerList []VocabularyAnswerStruct) {
	var (
		existent         bool
		index            int
		vocabularyAnswer VocabularyAnswerStruct
		tm               time.Time
	)
	if collection.readable && collection.writable {
		// only available for collection which is readable and writable
		existent = false
		index = -1
		for i := 0; i < len(vocabularyAnswerList); i++ {
			if strings.Compare(vocabularyAnswerList[i].Status, Basic) == 0 {
				// only for vocabulary with define from basic query
				if vocabularyAnswerList[i].Location.TableType == Online {
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
			vocabularyAnswer.Location.TableType = Collection
			vocabularyAnswer.Status = ""
			// add
			collection.content = append(collection.content, vocabularyAnswer)
		}
	}
}
