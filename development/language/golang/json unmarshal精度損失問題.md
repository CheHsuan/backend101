# json Unmarshal精度損失問題

golang json.Unmarshal()時對數字預設的型別為float64，對於大數字的處理時，會丟失精度，導致後續程式行為錯誤。

```golang
package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	Id int64
}

type B struct {
	Id json.Number 
}

type C struct {
	Id float64
}

func main() {
	bs := []byte(`{"id": 148480899330089071}`)
	
	a := A{}
	json.Unmarshal(bs, &a)
	fmt.Println(a)
	
	b := B{}
	json.Unmarshal(bs, &b)
	fmt.Println(b)
	
	c := C{}
	json.Unmarshal(bs, &c)
	fmt.Println(c)
}

```

執行結果為

```bash
{148480899330089071}
{148480899330089071}
{1.4848089933008906e+17}
```

因此在unmarshal沒有強制使用int64或json.Number去接的話，而是interface {}的話，在處理這個問題可以使用json.Decoder

```golang
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type D struct {
	Id interface{}
}

func main() {
	bs := []byte(`{"id": 148480899330089071}`)

	d := D{}
	decoder := json.NewDecoder(bytes.NewReader(bs))
	decoder.UseNumber()
	decoder.Decode(&d)

	fmt.Println(d)
  
	v, _ := d.Id.(json.Number).Int64()
	fmt.Println(v)
}
```

輸出為

```
{148480899330089071}
148480899330089071
```

看一下source code

```golang
// A Number represents a JSON number literal.
type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
    return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
    return strconv.ParseInt(string(n), 10, 64)
}
```

使用decode.UseNumber的好處是可以將數字轉換為json.Number，而後在根據需求轉換成float64或int64


