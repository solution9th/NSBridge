package app

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	ggrpc "github.com/solution9th/NSBridge/api/grpc"
	pb "github.com/solution9th/NSBridge/dns_pb"
	"github.com/solution9th/NSBridge/oneapm"

	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func runGRpc() error {

	port := config.GRpcConfig.Port

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		utils.Error("[grpc] failed to listen:", err)
		return err
	}

	pemPath, ok := utils.FindFile("keys/grpc/grpc.pem", "/etc/ns_bridge/grpc.pem")
	if !ok {
		panic("not found pem")
	}

	keyPath, ok := utils.FindFile("keys/grpc/grpc.key", "/etc/ns_bridge/grpc.key")
	if !ok {
		panic("not found key")
	}

	creds, err := credentials.NewServerTLSFromFile(pemPath, keyPath)
	if err != nil {
		utils.Error("[grpc] TLS server error:", err)
		return err
	}

	oneInterceptor := oneapm.NewInterceptor(ggrpc.IsHTTPRequest)

	s := grpc.NewServer(grpc.Creds(creds),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
				oneInterceptor.StreamServerInterceptor(),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				oneInterceptor.UnaryServerInterceptor(),
			),
		),
	)
	// s := grpc.NewServer()

	pb.RegisterDNSServerServer(s, ggrpc.New())

	fmt.Println("[grpc] run:", port)

	if err = s.Serve(lis); err != nil {
		utils.Error("[grpc] failed to serve:", err)
		return err
	}

	return nil
}
