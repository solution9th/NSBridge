package oneapm

import (
	"context"
	"net/http"
	"strings"

	"go-agent/blueware"

	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	OneAPP blueware.Application
)

func AppTransmission(app blueware.Application) {
	OneAPP = app
}

// GinMiddleware Gin 的 APM 中间件
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		path := c.Request.URL.Path

		txn := OneAPP.StartTransaction(path, w, c.Request)
		c.Next()

		txn.End()
	}
}

func NewInterceptor(f func(ctx context.Context) bool) *Interceptor {
	return &Interceptor{
		IgnoreFunc: f,
	}
}

type Interceptor struct {
	IgnoreFunc func(ctx context.Context) bool
}

func (i *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 只有是 web 请求才会被抓取
		// 有一个关键字 is_web = true 才会生效
		r, _ := http.NewRequest("GET", "/pp", nil)

		name := formatName(info.FullMethod)

		txn := OneAPP.StartTransaction(name, nil, r)

		if i.IgnoreFunc != nil && i.IgnoreFunc(ctx) {
			txn.Ignore()
		}
		resp, err := handler(ctx, req)
		if err != nil {
			utils.Error("[grpc] Interceptor error:", err)
		}

		txn.End()
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func (i *Interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		r, _ := http.NewRequest("GET", "/pp", nil)

		name := formatName(info.FullMethod)

		txn := OneAPP.StartTransaction(name, nil, r)

		if i.IgnoreFunc != nil && i.IgnoreFunc(ss.Context()) {
			txn.Ignore()
		}

		err := handler(srv, ss)
		if err != nil {
			utils.Error("[grpc] Interceptor error:", err)
		}
		txn.End()
		return err
	}
}

func formatName(name string) string {

	if name == "" {
		return name
	}

	list := strings.Split(name, "/")
	return "grpc-" + list[len(list)-1]
}
