package app

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	ggrpc "github.com/solution9th/NSBridge/api/grpc"
	pb "github.com/solution9th/NSBridge/dns_pb"

	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func runGateway() (*runtime.ServeMux, error) {

	gRPCPort := config.GRpcConfig.Port

	ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	pemPath, ok := utils.FindFile("keys/grpc/grpc.pem", "/etc/ns_bridge/grpc.pem")
	if !ok {
		panic("not found pem")
	}

	dcreds, err := credentials.NewClientTLSFromFile(pemPath, "dev")
	if err != nil {
		utils.Error("[grpc gateway] TLS server error:", err)
		return nil, err
	}

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(ggrpc.GatewayHaderMatcher),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.HTTPBodyMarshaler{
				Marshaler: &runtime.JSONPb{OrigName: true, EmitDefaults: true},
			},
		),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	err = pb.RegisterDNSServerHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", gRPCPort), opts)
	if err != nil {
		utils.Error("[grpc gateway] error:", err)
		return nil, err
	}

	// port := 8080
	// return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)

	return mux, nil
}
