package grpc

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/solution9th/NSBridge/internal/utils"
)

type RPCServer struct{}

// New new a rpc server
func New() *RPCServer {
	return &RPCServer{}
}

// TypeConvert 类型转换，输入数据库类型，转换成 proto3 类型
// 需要 struct tag json 相同
func TypeConvert(in interface{}, out interface{}) error {

	if reflect.ValueOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("out need ptr")
	}

	body, err := json.Marshal(in)
	if err != nil {
		utils.Error("[type] ma error:", err)
		return err
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		utils.Error("[type] un error:", err)
	}
	return err
}
