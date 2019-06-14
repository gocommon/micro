# 更新日志

### 修改 api handler 错误返回值

go mod replace github.com/micro/go-micro/api/handler/api  github.com/gocommon/go-micro/api/handler/api

```

// ErrHandler 错误处理
type ErrHandler func(w http.ResponseWriter, r *http.Request, err *errors.Error)


// DefaultErrHandler DefaultErrHandler
var DefaultErrHandler = func(w http.ResponseWriter, r *http.Request, err *errors.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
	return
}
```

### 添加插件

- cors
- gzip
- metrics
- disable_rpc