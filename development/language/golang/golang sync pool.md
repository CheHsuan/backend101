# Golang sync.Pool

## 前言

一個server因業務請求在程式內建立了多個的變數暫時存放資料，或做中繼計算用，到最後總是會迎來GC-自動銷毀沒有在使用的變數及其記憶體。但在一個高併發的環境中，任何的記憶體的宣告和銷毀都會對效能產生影響，因此如何能有效複用已建立的對象至關重要。

## sync.Pool

sync.Pool是一個官方提供的一個工具，用於提供通用物件的取得和釋放，如果該池中已有物件可以直接拿出來使用，用完再回歸於池中，若是需要時物件但池是空的，則會新增一個出來，而若是池中物件空閒太久，則會被GC回收。sync.Pool並沒有具體的大小和上限，並且是treadsafe的，非常適合client-server架構中使用。

**用法**

我們只要定義好`New()`方法，並且使用`Get()`和`Put()`就可以取得或釋出物件。

**注意**

1. 由於GC會清除sync.Pool裡面的物件，所以請勿放置重要的資訊
2. 由於取出的物件可能是被使用過的記憶體空間，因此為避免上次遺留使用後的資訊，在放置回sync.Pool前先行清空內部變數，抑或是取出時更新

## Benchmark

```go
const (
	size = 10000
)

var (
	entryPool = sync.Pool{
		New: func() interface{} {
			return new(Entry)
		},
	}
	bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	payload = make([]byte, size)
)

type Entry struct {
	payload []byte
}

func (e *Entry) Set() {
	if e.payload == nil {
		e.payload = make([]byte, size)
	}
	copy(e.payload[:size], payload)
}

func (e *Entry) Reset() {
	e.payload = e.payload[:0]
}

func BenchmarkSetEntry(b *testing.B) {
	for n := 0; n < b.N; n++ {
		etr := new(Entry)
		etr.Set()
	}
}

func BenchmarkSetEntryWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		etr := entryPool.Get().(*Entry)
		etr.Set()
		etr.Reset()
		entryPool.Put(etr)
	}
}
```

這邊的範例是有一個Entry的struct，裡面的payload是一個byte slice，每經過一次迴圈時，原始的做法每次都會宣告10000個byte，而sync.Pool則是會複用Entry的物件及其宣告好的記憶體空間。

```bash
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: syncpool
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkSetEntry-8           	  907266	      1223 ns/op	   10240 B/op	       1 allocs/op
BenchmarkSetEntryWithPool-8   	13724096	     87.91 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	syncpool	2.428s
```

從memory benchmark結果來看，sync.Pool幾乎每次都複用了池中宣告好的物件，所以數據顯示每次op的平均記憶體使用量為0個byte。

## 結語

在可預期的使用情境下，盡量使用sync.Pool可以大大降低memory的使用量，降低GC的帶來的影響

