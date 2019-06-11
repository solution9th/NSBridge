package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/solution9th/NSBridge/api/grpc"
	"github.com/solution9th/NSBridge/internal/config"
	"github.com/solution9th/NSBridge/internal/oneapm"
	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/gin-gonic/gin"
)

func runWeb() error {

	var err error

	port := config.ServerConfig.Port
	if port == 0 {
		port = 8080
	}

	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	mux, err := runGateway()
	if err != nil {
		return err
	}

	r := gin.Default()
	r.Use(oneapm.GinMiddleware())

	r.GET("/api/ping", Ping)
	r.StaticFile("/swagger/swagger.json", "dns_pb/dns.swagger.json")
	r.Any("/api/v1/*any", grpc.AuthMiddleware(), gin.WrapF(mux.ServeHTTP))

	// NewRouter(r, mux)
	NewWebRouter(r)
	fmt.Println("[http] listen:", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen error: %s\n", err)
			utils.Error("http listen error:", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.Error("http shutdown error:", err)
		return err
	}
	log.Println("Server exiting")
	return nil
}
