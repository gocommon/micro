package main

import (
	"net/http"

	apiHandler "github.com/micro/go-micro/api/handler/api"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-plugins/micro/cors"
	rpc "github.com/micro/go-plugins/micro/disable_rpc"
	"github.com/micro/go-plugins/micro/gzip"
	"github.com/micro/go-plugins/micro/metrics"
	"github.com/micro/micro/api"
	proto "github.com/micro/micro/api/proto"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
)

func init() {
	// 注册插件
	plugin.Register(cors.NewPlugin())
	plugin.Register(metrics.NewPlugin())
	api.Register(gzip.NewPlugin())
	api.Register(rpc.NewPlugin())

	// 更新 api errHandler
	apiHandler.DefaultErrHandler = func(w http.ResponseWriter, r *http.Request, err *errors.Error) {
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(500)
		// w.Write([]byte(err.Error()))

		switch err.Code {
		case 0: // rpc 调用函数返回的错误

		}

		return
	}

	apiHandler.DefaultRespHandler = func(w http.ResponseWriter, r *http.Request, rsp *proto.Response) {
		for _, header := range rsp.GetHeader() {
			for _, val := range header.Values {
				w.Header().Add(header.Key, val)
			}
		}

		if len(w.Header().Get("Content-Type")) == 0 {
			w.Header().Set("Content-Type", "application/json")
		}

		w.WriteHeader(int(rsp.StatusCode))
		w.Write([]byte(rsp.Body))
	}
}

func main() {
	cmd.Init()
}
