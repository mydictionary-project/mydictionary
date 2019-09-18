# Service

[简体中文](./service.zh-Hans.md)

### 1. Introduction

cache and *services*

### 2. Cache

file `cache.go`

This module is able to cache query result from *services*. **It can increase the speed of online query considerably.**

##### 2.1. Data Structure

``` go
type CacheItemStruct struct {
	QueryString  string   `json:"queryString"`
	Word         string   `json:"word"`
	Definition   []string `json:"definition"`
	Status       string   `json:"status"`
	CreationTime int64    `json:"creationTime"`
}
```

`CacheItemStruct` is a structure and has got these members:

- `QueryString` indicates the string of query.
- `Word` indicates the word.
- `Definition` indicates definitions.
- `Status` indicates the status.
- `CreationTime` is a unix timestamp which indicates when this item is created.

``` go
type CacheStruct struct {
	path         string
	shelfLifeDay int64
	Content      []CacheItemStruct `json:"content"`
}
```

`CacheStruct` is a structure.

It has got a public member `Content` which stores all cache items.

##### 2.2. Member Function

```go
func (cache *CacheStruct) Read(path string, shelfLifeDay int64) (err error)
```

The function is used for reading cache from a JSON file indicates by `path`, and set the life period in days for all cache items by `shelfLifeDay`.

- If the file indicated by `path` is not existent, create it.
- If `shelfLifeDay` is 0, cache will never expire.

After reading cache, each item will be checked whether it is expired (determined by its `CreationTime`, cache's `shelfLifeDay` and the current time). **Then, all expired items will be removed.**

```go
func (cache *CacheStruct) Query(queryString string) (item ItemStruct, err error)
```

The function is used for searching `queryString` from cache.

```go
func (cache *CacheStruct) Add(item ItemStruct)
```

The function is used for adding `item` to cache.

```go
func (cache *CacheStruct) Write()
```

The function is used for writing cache to the JSON file which it comes from.

### 3. Service

file `service.go`

#### 3.1. Interface

All *services* should obey the following interface.

``` go
type ServiceInterface interface {
	GetServiceName() string
	GetCache() *CacheStruct
	Query(vocabulary4mydictionary.VocabularyAskStruct) vocabulary4mydictionary.VocabularyAnswerStruct
}
```

##### 3.1.1. GetServiceName

``` go
func (service *ServiceInterface) GetServiceName()
```

This function return the name of the service.

##### 3.1.2. GetCache

``` go
func (service *ServiceInterface) GetCache()
```

This function return the pointer of the cache of the *service* .

##### 3.1.3. Query

``` go
func (service *ServiceInterface) Query(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct)
```

This function is used for querying the *vocabulary*.

#### 3.2.  Currently Supported Service

##### 3.2.1. Bing Dictionary

http://cn.bing.com/dict/

``` go
// bing-dictionary.go

type BingDictionaryStruct struct {
	cache CacheStruct
}
```

##### 3.2.2. iCIBA Collins

http://www.iciba.com/

``` go
// iciba-collins.go

type IcibaCollinsStruct struct {
	cache CacheStruct
}
```

##### 3.2.3. Merriam Webster

https://www.merriam-webster.com/

``` go


type MerriamWebsterStruct struct {
	cache CacheStruct
}
```

#### 3.3. Create Service by Yourself

- Create a new ".go" file in this package.
- Copy the code below to the file (it is a good idea to use [goquery](https://github.com/PuerkitoBio/goquery) for grasping content from webpage).

``` go
package service4mydictionary

import (
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const serviceName = "Service"

//
// NOTE:
// 1. Change "Service" above to your service name.
// 2. Rename all "ServiceStruct" below in this file.
// 3. Delete this note in your service.
//

// ServiceStruct : service struct
type ServiceStruct struct {
	cache CacheStruct
}

// GetServiceName : get service name
func (service *ServiceStruct) GetServiceName() (value string) {
	value = serviceName
	return
}

// GetCache : get cache
func (service *ServiceStruct) GetCache() (cache *CacheStruct) {
	cache = &service.cache
	return
}

// Query : query vocabulary
func (service *ServiceStruct) Query(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct) {
	var (
		err         error
		queryString string
		item        CacheItemStruct
	)
	// set
	vocabularyAnswer.SourceName = serviceName
	vocabularyAnswer.Location.TableType = vocabulary4mydictionary.Online
	vocabularyAnswer.Location.TableIndex = -1
	vocabularyAnswer.Location.ItemIndex = -1
	// query cache
	queryString = url.QueryEscape(vocabularyAsk.Word)
	item, err = service.cache.Query(queryString)
	if err == nil {
		goto SET
	}
	// query online
	//
	// NOTE:
	// 1. Add your code here.
	// 2. If success, "item.Status" should be "vocabulary4mydictionary.Basic";
	//    else, set "item.Status" as your debug information and jump to label 'ADD' immidiately.
	// 3. To help you debug this function,
	//    set "online.debug" in file "mydictionary.setting.json" as "true".
	//    At this time, MYDICTIONARY will show off the debug information.
	// 4. Set "online.debug" in file "mydictionary.setting.json" as "false"
	//    to let MYDICTIONARY will hide the debug information
	// 5. Delete this note in your service.
	//
ADD:
	// add to cache
	item.QueryString = queryString
	item.CreationTime = time.Now().Unix()
	service.cache.Add(item)
SET:
	// set
	vocabularyAnswer.Word = item.Word
	vocabularyAnswer.Definition = item.Definition
	vocabularyAnswer.Status = item.Status
	return
}
```

- Modify the file based on annotations begin with `NOTE`.
- Update [setting.go](../setting.go) based on annotations begin with `NOTE`.
- Use pull request to commit your modification (optional).

### 4. Others

- All code files are edited by [Atom](https://atom.io/).
- All ".md" files are edited by [Typora](http://typora.io).
- The style of all ".md" files is [Github Flavored Markdown](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown).
- There is a LF (Linux) at the end of each line.
