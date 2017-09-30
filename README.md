# MYDICTIONARY

### 1. Introduction

MYDICTIONARY is a library designed by golang. It provides the API for developers to build applications of excel-based and online dictionaries.

### 2. Basic Information

#### 2.1. Vocabulary

The item contains word, definition and other information called ***vocabulary***.

#### 2.2. Online

MYDICTIONARY can grasp pages from the website and extract information to build *vocabularies*. This process called ***service***. By doing this, MYDICTIONARY enables us to get *vocabularies* which are not included in local dictionaries.

Here are *services* the library provided now:

- [Bing Dictionary](http://cn.bing.com/dict/): repository [bingdictionary4mydictionary](https://github.com/zzc-tongji/bingdictionary4mydictionary)
- [iCIBA Collins](http://www.iciba.com/): repository [icibacollins4mydictionary](https://github.com/zzc-tongji/icibacollins4mydictionary)
- [Merriam Webster](https://www.merriam-webster.com/): repository [merriamwebster4mydictionary](https://github.com/zzc-tongji/merriamwebster4mydictionary)

**Declaration:**

- **Copyrights of all information in aforementioned websites belong to corresponding companies.**
- **All of these information are prohibited for any form of commercial use.**
- **I am NOT accountable for misappropriation of these information.**

Developers can designed their own *service*. Get further information from repository [example4mydictionary](https://github.com/zzc-tongji/example4mydictionary).

#### 2.3. Collection & Dictionary

All dictionary applications need files to store words and defines. As far as I know, most of these files are specially designed, which means users could not view or edit without using corresponding applications. It puzzles me sometimes. As a matter of fact, I choose to store data in ".xlsx" files when designing MYDICTIONARY, which means that users can easily read and write by using Microsoft Excel.

All ".xlsx" files which march the following conditions could become ***collection files*** or ***dictionary files*** of this library.

- Contain at least one worksheet.
- In the first worksheet, cells of the first line should contain: "SN", "Word", "Define", "Note", "QC" (query counter) and "QT" (last query time).

Here are examples:

- *collection file* "bing-dictionary.xlsx"

![bing-dictionary](./README.picture/bing-dictionary.png)

- *dictionary file* "animal.xlsx"

![animal](./README.picture/animal.png)

***Collections*** and ***Dictionaries*** are mainly composed of *vocabularies*. They are RAM images of *collection files* and *dictionary files*. Each line of *collection files* and *dictionary files* is converted to a *vocabulary* in *collections* and *dictionaries*. The data structure of the *vocabulary* is [here](https://github.com/zzc-tongji/vocabulary4mydictionary#4-answer).

*Collections* and *dictionary* have similar structures, but there are still some differences:

- *Dictionaries* could provided *vocabularies* which are queried. We can let *dictionaries* update QC and QT (add 1 to QC and set the current time to QT), but others are all read-only for this library.
- *Collections* could also provided *vocabularies* which are queried, too. But the main function of *collections* is recording *vocabularies* from the *service* which are not existed in all *dictionaries*. So, the library can not only update QC and QT of existent *vocabularies* in *collections*, but also add new *vocabularies*.

The action of updating and/or adding is called ***record***.

#### 2.4. Configuration

The ***configuration*** file `mydictionary.setting.json` should be placed at the location of the application with MYDICTIONARY.

Here is an example:

```json
{
	"collection":
	[
		{
			"name":"bing-dictionary",
			"filePath":"data/bing-dictionary.xlsx",
			"readable":true,
			"writable":true,
			"onlineSource":"Bing Dictionary"
		},
		{
			"name":"iciba-collins",
			"filePath":"data/iciba-collins.xlsx",
			"readable":true,
			"writable":true,
			"onlineSource":"iCIBA Collins"
		},
		{
			"name":"merriam-webster",
			"filePath":"data/merriam-webster.xlsx",
			"readable":true,
			"writable":true,
			"onlineSource":"Merriam Webster"
		}
	],
	"dictionary":
	[
		{
			"name":"animal",
			"filePath":"data/animal.xlsx",
			"readable":true,
			"writable":true
		},
		{
			"name":"fruit",
			"filePath":"data/fruit.xlsx",
			"readable":true,
			"writable":true
		}
	],
	"online":
	{
		"mode":3,
		"service":
		{
			"bingDictionary":true,
			"icibaCollins":true,
			"merriamWebster":true
		},
		"debug":true
	}
}
```

There are 3 structures in the *configuration*: `"collection"`, `"dictionary"` and `"online"`.

##### 2.4.1. collection

`"collection"` is an array and each item in this array has got such members:

- String `"name"`: it is the name of the *collection*.
- String `"filePath"`: it is the file path of the *collection*. It can be a relative path base on the location of the application.
- Boolean `"readable"`: If it is `false`, the *collection* will be ignored by the library. By setting this, we can disable the *collection* without removing the whole item.
- Boolean `"writable"`: If it is `true`, the library will be allowed to *record* *vocabularies* of the *collection*.
- String `"onlineSource"`: each *collection* is only able to record *vocabularies* from one *service*, but the library can get *vocabularies* from several different *services*. So we need this member to indicate the corresponding relation between the *collection* and the *service*.

##### 2.4.2. dictionary

`"dictionary"` is an array and each item in this array has got such members:

- String `"name"`: it is the name of the *dictionary*.
- String `"filePath"`: it is the file path of the *dictionary*. It can be a relative path base on the location of the application.
- Boolean `"readable"`: If it is `false`, the *dictionary* will be ignored by the library. By setting this, we can disable the *dictionary* without removing the whole item.
- Boolean `"writable"`: If it is `true`, the library will be allowed to *record* *vocabularies* of the *dictionary*.

##### 2.4.3. online

`"online"` is a structure and has got these members:

###### 2.4.3.1. mode

`"mode"` is an integer. It determines on what condition the library should provide *vocabularies* from *services* (query online).

Here is the possible value:

- `0`: the library will never query online.
- `1`: the library will query online, if users need.
- `2`: the library will query online, if the *vocabulary* is not existent in all *collections* and *dictionaries*.
- `3`: the library will query online, if users need OR the *vocabulary* is not existent in all *collections* and *dictionaries*.
- `4`: the library will always query online.

**If it's hard to understand, set `"mode"` as `3` by default.**

###### 2.3.4.2. service

`"service"` is an array. The key-value (string-boolean) pair of each item determines whether the *service* is enabled.

###### 2.3.4.3. debug

`"debug"` indicates whether the library is in debug mode. The default value is `false`. Do not modify it if you are not developer. Get further information from repository [example4mydictionary](https://github.com/zzc-tongji/example4mydictionary).

### 3. API

#### 3.1. Vocabulary

##### 3.1.1. Ask & Answer

See repository [vocabulary4mydictionary](https://github.com/zzc-tongji/vocabulary4mydictionary).

##### 3.1.2. Result

```go
type VocabularyResultStruct struct {
	Basic   []vocabulary4mydictionary.VocabularyAnswerStruct `json:"Basic"`
	Advance []vocabulary4mydictionary.VocabularyAnswerStruct `json:"Advance"`
}
```

`Basic` are made up by *vocabularies* come from *basic query*.

`Advance` are made up by *vocabularies* come from *advanced query*.

#### 3.2. Function

##### 3.2.1. Initialize

```go
func Initialize() (information string, err error)
```

The function is used to initialize the library. Here is the procedure:

- Read and parse the *configuration*.
- Read *collection files* and *dictionaries files* and build their RAM images (*collections* and *dictionaries*).
- Check network.

**Note:**

- **The function should be called before any other functions.**
- **The function should be called only once.**

Return values:

- If success, the content of the *configuration* will be returned by `information` with the current time, and `err` will be `nil`.
- If failure, the content of `err` will be returned by `information` with the current time.

##### 3.2.2. CheckNetwork

```go
func CheckNetwork() (pass bool, information string)
```

The function requests all enabled *services* for word "apple" to check network.

Return values:

- If all *services* return normal responses, `pass` will be `true`.
- If not, `pass` will be `false`.
- `information` will provide further information.

##### 3.2.3. Query

```go
func Query(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyResult VocabularyResultStruct, err error)
```

It is the core function of the API. Here is the procedure:

1. Query the word in `vocabularyAsk` in all *collections* and *dictionaries* and record offline results.
2. Determine whether the function should query the word online based on `online.mode`, option in `vocabularyAsk` and offline results.
3. If "2" is true, query the word online and record online results. **The process is concurrent** by using [goroutines](https://golang.org/doc/faq#goroutines), which means that the query time depends on the longest request time of all enabled *services*.
4. Return offline and online results by `vocabularyResult`.

Return values:

- If success, the result will be returned by `vocabularyResult`, and `err` will be `nil`.
- If failure, error will be returned by `err`.

##### 3.2.4. Save

```go
func Save() (success bool, information string)
```

The function is used to write RAM images (*collections* and *dictionaries*) back to their corresponding *collection files* and *dictionary files*.

Return values:

- If all files are written successfully, `success` will be `true`.
- If not, `success` will be `false`.
- `information` will provide further information.

### 4. Thanks

- Thanks to [xuri](https://github.com/xuri) for providing [excelize](https://github.com/360EntSecGroup-Skylar/excelize), a golang library for reading and writing Microsoft Excel™ (XLSX) files.

### 5. Communication

- [issues](https://github.com/zzc-tongji/mydictionary-local-cli/issues) of this repository
- QQ group: 657218106

![657218106](./README.picture/657218106.png)

### 6. Others

- All ".md" files are edited by [Typora](http://typora.io).
- The style of all ".md" files is [Github Flavored Markdown](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown).
- There is a LF (Linux) at the end of each line.
