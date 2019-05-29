package app

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/grpc-ecosystem/grpc-gateway/runtime"
// 	"github.com/solution9th/NSBridge/api/grpc"
// )

// // NewRouter new router 【舍弃】
// // 只要是 /api/v1 前缀请求都会被 gRPC gateway 接管，并且会运行 Gin 中间件
// func NewRouter(r *gin.Engine, mux *runtime.ServeMux) {

// 	r.Any("/api/v1/*any", grpc.AuthMiddleware(), gin.WrapF(mux.ServeHTTP))

// 	return

// 	// apiRouter := r.Group("/api/v1")
// 	// apiRouter.Use(v1.AuthMiddleware())
// 	// {
// 	// 	apiRouter.GET("/ping", v1.Ping)

// 	// 	// 其他接口
// 	// 	apiRouter.GET("/types", v1.GetTypes)                       // 获取可以添加的记录类型
// 	// 	apiRouter.GET("/lines/all", v1.GetLindeIDs)                // 获得所有线路列表
// 	// 	apiRouter.GET("/lines/continental", v1.GetLineContinental) // 获取洲线路列表
// 	// 	apiRouter.GET("/lines/isp", v1.GetLineISP)                 // 获得 ISP 线路列表
// 	// 	apiRouter.GET("/lines/country", v1.GetlineCountry)         // 获得国家线路列表
// 	// 	apiRouter.GET("/lines/province", v1.GetLineProvince)       // 获得 省 列表
// 	// 	apiRouter.GET("/lines/outcity", v1.GetLineOutCity)         // 获得国外城市列表
// 	// 	apiRouter.GET("/lines/cityisp", v1.GetLineChinaCityISP)    // 获得国内城市与运营商组合列表

// 	// 	// 域名类接口
// 	// 	apiRouter.GET("/domains", v1.GetDomainLists)     // 域名列表
// 	// 	apiRouter.POST("/domains", v1.CreateDomain)      // 添加域名
// 	// 	apiRouter.DELETE("/domain/:id", v1.DeleteDomain) // 删除域名
// 	// 	apiRouter.GET("/status", v1.IsTakeOver)          // 检查域名是否托管

// 	// 	// 解析记录
// 	// 	apiRouter.GET("/records", v1.GetRecordList)            // 获取解析记录列表
// 	// 	apiRouter.POST("/records", v1.CreateRecord)            // 新增解析记录
// 	// 	apiRouter.GET("/record/:recordid", v1.GetRecord)       // 获得解析记录详情
// 	// 	apiRouter.PUT("/record/:recordid", v1.UpdateRecord)    // 更新解析记录
// 	// 	apiRouter.DELETE("/record/:recordid", v1.DeleteRecord) // 删除解析记录
// 	// 	apiRouter.PATCH("/record/:recordid", v1.DisableRecord) // 启动暂停记录
// 	// }
// }
