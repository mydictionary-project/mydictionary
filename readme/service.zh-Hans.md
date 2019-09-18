# 在线服务

[English](./service.md)

### 1. 简介

缓存和*在线服务*

### 2. 缓存

文件`cache.go`

该模块能够缓存从*在线服务*获得的查询结果。**这会显著增加联网查询的速度。**

##### 2.1. 数据结构

```go
type ItemStruct struct {
	QueryString  string   `json:"queryString"`
	Word         string   `json:"word"`
	Definition   []string `json:"definition"`
	Status       string   `json:"status"`
	CreationTime int64    `json:"creationTime"`
}
```

`ItemStruct`是一个结构体，包含下列成员：

- `QueryString`指示了查询字符串。
- `Word`指示了词汇。
- `Definition`指示了释义。
- `Status`指示了状态。
- `CreationTime`是一个UNIX时间戳，它指示了本缓存项的创建时间。

```go
type CacheStruct struct {
	path         string
	shelfLifeDay int64
	Content      []ItemStruct `json:"content"`
}
```

`CacheStruct`是一个结构体。

它具有一个公开的成员`Content`用于存放所有的缓存项。

##### 2.2. 成员函数

```go
func (cache *CacheStruct) Read(path string, shelfLifeDay int64) (err error)
```

这个函数用于从被`path`指定的JSON文件中读取缓存，同时将`shelfLifeDay`设定为所有缓存项的保存期限（天）。

- 如果`path`指定的文件不存在，那么创建之。
- 如果`shelfLifeDay`为0，那么缓存永不过期。

在读取缓存之后，每个缓存项会被检查以确定是否过期（由缓存项的`CreationTime`、缓存的`shelfLifeDay`和当前时间共同决定）。**随后，所有过期的缓存项将被移除。**

```go
func (cache *CacheStruct) Query(queryString string) (item ItemStruct, err error)
```

这个函数用于在缓存中搜索`queryString`。

```go
func (cache *CacheStruct) Add(item ItemStruct)
```

这个函数用于添加`item`到缓存。

```go
func (cache *CacheStruct) Write()
```

这个函数用于将缓存写回JSON文件。

### 3. 在线服务

文件 `service.go`

#### 3.1. 接口

所有*在线服务*都应当满足以下接口。

``` go
type ServiceInterface interface {
	GetServiceName() string
	GetCache() *CacheStruct
	Query(vocabulary4mydictionary.VocabularyAskStruct) vocabulary4mydictionary.VocabularyAnswerStruct
}
```

##### 2.2.1. GetServiceName

```go
func (service *ServiceInterface) GetServiceName()
```

该函数返回*在线服务*的名称。

##### 3.1.2. GetCache

```go
func (service *ServiceInterface) GetCache()
```

该函数返回*在线服务*缓存的指针。

##### 3.1.3. Query

```go
func (service *ServiceInterface) Query(vocabularyAsk vocabulary4mydictionary.VocabularyAskStruct) (vocabularyAnswer vocabulary4mydictionary.VocabularyAnswerStruct)
```

该函数用来查询*词条*。

#### 3.2. 目前支持的在线服务

##### 3.2.1. 必应词典

http://cn.bing.com/dict/

```go
// bing-dictionary.go

type BingDictionaryStruct struct {
	cache CacheStruct
}
```

##### 3.2.2. 金山词霸-柯林斯词典

http://www.iciba.com/

```go
// iciba-collins.go

type IcibaCollinsStruct struct {
	cache CacheStruct
}
```

##### 3.2.3. 韦氏词典

https://www.merriam-webster.com/

```go
// merriam-webster.go

type MerriamWebsterStruct struct {
	cache CacheStruct
}
```

#### 3.3. 创建自定义在线服务

- 在本包内新建一个".go"文件。
- 复制下列代码到这个新文件（推荐使用[goquery](https://github.com/PuerkitoBio/goquery)来获取网页内容）。

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

- 根据`NOTE`开头的注释，修改上述文件。
- 根据`NOTE`开头的注释，更新文件[setting.go](../setting.go)文件。
- 使用 pull request 请求来提交你的修改（可选）。

### 4. 其他

- 所以代码文件是用[Atom](https://atom.io/)编写的。
- 所有".md"文件是用[Typora](http://typora.io)编写的。
- 所有".md"文件的风格是[Github Flavored Markdown](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown)。
- 各行以LF（Linux）结尾。

