# gob vs json

實驗一下golang的encoding lib對於同一結構先encode再decode後，所產出的差別

## 程式碼

```go
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

type Payload struct {
	Data map[string]map[string]interface{} `json:"data"`
}

func main() {
	t := &Payload{
		Data: map[string]map[string]interface{}{
			"id": {
				"in": []string{
					"abcdefghij",
				},
				"equal": int(1234567890),
			},
		},
	}
	jsonway(t)
	gobway(t)
}

func gobway(src *Payload) {
	fmt.Println("gob way......................")
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		panic(err.Error())
	}
	bs := buf.Bytes()

	dest := &Payload{}
	if err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(dest); err != nil {
		panic(err.Error())
	}

	if value, ok := dest.Data["id"]["in"]; ok {
		getDataType(value)
	}
	if value, ok := dest.Data["id"]["equal"]; ok {
		getDataType(value)
	}
}

func jsonway(src *Payload) {
	fmt.Println("json way......................")
	bs, err := json.Marshal(src)
	if err != nil {
		panic(err.Error())
	}

	dest := &Payload{}
	err = json.Unmarshal(bs, dest)
	if err != nil {
		panic(err.Error())
	}

	if value, ok := dest.Data["id"]["in"]; ok {
		getDataType(value)
	}
	if value, ok := dest.Data["id"]["equal"]; ok {
		getDataType(value)
	}
}

func getDataType(val interface{}) {
	fmt.Printf("value: %v is %T type\n", val, val)
}
```

## 輸出

```txt
json way......................
value: [abcdefghij] is []interface {} type
value: 1.23456789e+09 is float64 type
gob way......................
value: [abcdefghij] is []string type
value: 1234567890 is int type
```

## 行為

在json 的方式當中，[]string的型態被decode成[]interface{}，int則是變成了float64型態，而在gob方式當中，則完整保存了原始數據的型態。
