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
	name        string
	readable    bool
	writable    bool
	index       int
	xlsx        *excelize.File
	sheetName   string
	columnIndex map[string]int
	content     []vocabulary4mydictionary.VocabularyAnswerStruct
}

// open and check .xlsx file
func (dictionary *dictionaryStruct) check(filePath string) (err error) {
	var (
		recheckCounter    int
		recheckUpperLimit int
		contentTemp       [][]string
		columnNumber      int
	)
	// file -> ram image
	dictionary.xlsx, err = excelize.OpenFile(filePath)
	if err != nil {
		return
	}
	// recheck
	recheckCounter = 0
	recheckUpperLimit = 10
RECHECK:
	dictionary.sheetName = dictionary.xlsx.GetSheetMap()[1]
	if strings.Compare(dictionary.sheetName, "") == 0 {
		// no worksheet in workbook: create a worksheet
		if strings.Compare(dictionary.name, "") == 0 {
			dictionary.xlsx.NewSheet("dictionary")
		} else {
			dictionary.xlsx.NewSheet(dictionary.name)
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
	contentTemp = dictionary.xlsx.GetRows(dictionary.sheetName)
	if len(contentTemp) == 0 {
		// empty worksheet: create row 1
		dictionary.xlsx.SetCellValue(dictionary.sheetName, "A1", sn)
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
	dictionary.columnIndex = map[string]int{wd: -1, def: -1, sn: -1, qc: -1, qt: -1, nt: -1}
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
		case nt:
			dictionary.columnIndex[nt] = i
		default:
			break
		}
	}
	if dictionary.columnIndex[wd] == -1 {
		dictionary.xlsx.SetCellValue(dictionary.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", wd)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if dictionary.columnIndex[def] == -1 {
		dictionary.xlsx.SetCellValue(dictionary.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", def)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if dictionary.columnIndex[sn] == -1 {
		dictionary.xlsx.SetCellValue(dictionary.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", sn)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if dictionary.columnIndex[qc] == -1 {
		dictionary.xlsx.SetCellValue(dictionary.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", qc)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if dictionary.columnIndex[qt] == -1 {
		dictionary.xlsx.SetCellValue(dictionary.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", qt)
		// recheck
		if recheckCounter < recheckUpperLimit {
			recheckCounter++
			goto RECHECK
		} else {
			err = fmt.Errorf("there are format errors in file \"%s\"", filePath)
			return
		}
	}
	if dictionary.columnIndex[nt] == -1 {
		dictionary.xlsx.SetCellValue(dictionary.sheetName, excelize.ToAlphaString(len(contentTemp[0]))+"1", nt)
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

// read data from .xlsx file and put to dictionary
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
			str = dictionary.xlsx.GetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[wd]), i))
			if strings.Compare(str, "") == 0 {
				break
			}
			vocabularyAnswer.Word = str
			// `xlsx:def` -> .Define
			str = dictionary.xlsx.GetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[def]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Definition = strings.Split(str, "\n")
			if len(vocabularyAnswer.Definition) == 1 &&
				strings.Compare(vocabularyAnswer.Definition[0], "") == 0 {
				vocabularyAnswer.Definition = nil
			}
			// `xlsx:sn` -> .SerialNumber
			vocabularyAnswer.SerialNumber, err = strconv.Atoi(dictionary.xlsx.GetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[sn]), i)))
			if err != nil {
				vocabularyAnswer.SerialNumber = i - 1
			}
			// `xlsx:qc` -> .QueryCounter
			vocabularyAnswer.QueryCounter, err = strconv.Atoi(dictionary.xlsx.GetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qc]), i)))
			if err != nil {
				vocabularyAnswer.QueryCounter = 0
			}
			// reset err
			err = nil
			// `xlsx:qt` -> .QueryTime
			vocabularyAnswer.QueryTime = dictionary.xlsx.GetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qt]), i))
			// `xlsx:nt` -> .Note
			str = dictionary.xlsx.GetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[nt]), i))
			str = strings.TrimSpace(str)
			vocabularyAnswer.Note = strings.Split(str, "\n")
			if len(vocabularyAnswer.Note) == 1 &&
				strings.Compare(vocabularyAnswer.Note[0], "") == 0 {
				vocabularyAnswer.Note = nil
			}
			// others
			vocabularyAnswer.SourceName = dictionary.name
			vocabularyAnswer.Location.TableType = vocabulary4mydictionary.Dictionary
			vocabularyAnswer.Status = ""
			// add to dictionary
			dictionary.content = append(dictionary.content, vocabularyAnswer)
			// set location
			dictionary.content[len(dictionary.content)-1].Location.TableIndex = dictionary.index
			dictionary.content[len(dictionary.content)-1].Location.ItemIndex = len(dictionary.content) - 1
		}
	}
	return
}

// get data dictionary and write to .xlsx file
func (dictionary *dictionaryStruct) write() (information string, err error) {
	if dictionary.readable && dictionary.writable {
		// content -> ram image
		for i := 0; i < len(dictionary.content); i++ {
			// .QueryCounter -> `xlsx:qc`
			dictionary.xlsx.SetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qc]), i+2), dictionary.content[i].QueryCounter)
			// .QueryTime -> `xlsx:qt`
			dictionary.xlsx.SetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[qt]), i+2), dictionary.content[i].QueryTime)
			// .Note -> `xlsx:nt`
			dictionary.xlsx.SetCellValue(dictionary.sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(dictionary.columnIndex[nt]), i+2), strings.Join(dictionary.content[i].Note, "\n"))
		}
		// ram image -> file
		err = dictionary.xlsx.Save()
		if err != nil {
			return
		}
		// output
		information = fmt.Sprintf("Dictionary \"%s\" has been updated.\n\n", dictionary.xlsx.Path)
	}
	return
}

// query and update
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
					// update dictionary
					if vocabularyAsk.DoNotRecord == false {
						// update
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
				for j := 0; j < len(dictionary.content[i].Definition); j++ {
					if strings.Contains(dictionary.content[i].Definition[j], vocabularyAsk.Word) {
						vocabularyAnswer = dictionary.content[i]
						vocabularyAnswer.Status = vocabulary4mydictionary.Advance
						vocabularyAnswerList = append(vocabularyAnswerList, vocabularyAnswer)
						goto ADVANCE_END
					}
				}
				for j := 0; j < len(dictionary.content[i].Note); j++ {
					if strings.Contains(dictionary.content[i].Note[j], vocabularyAsk.Word) {
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
