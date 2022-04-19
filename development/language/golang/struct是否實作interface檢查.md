# struct是否實作interface檢查

```go
var _ SomeInterface = new(SomeType)
```

不產生新的global但又可以檢查`SomeType` 是否有實作`SomeInterface`

另一種常見寫法

```go
var _ SomeInterface = (*SomeType)(nil)
```

該檢查會在compile time時進行

可見於

`google.golang.org/grpc@v1.37.0/server.go`

```go
var _ http.Handler = (*Server)(nil)
```

用來檢查Server是否有實作http.Handler
