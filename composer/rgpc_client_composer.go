package composer

import (
	sctx "github.com/viettranx/service-context"
	"golang-ai-management/common"
	"golang-ai-management/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ComposeUserAuthRPCClient(serviceCtx sctx.ServiceContext) pb.UserAuthServiceClient {
	configComp := serviceCtx.MustGet(common.KeyCompConf).(common.Config)

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial(configComp.GetGRPCUserAddress(), opts)

	if err != nil {
		log.Fatal(err)
	}

	return pb.NewUserAuthServiceClient(clientConn)
}

func ComposeProfileRPCClient(serviceCtx sctx.ServiceContext) pb.UserServiceClient {
	configComp := serviceCtx.MustGet(common.KeyCompConf).(common.Config)

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial(configComp.GetGRPCProfileAddress(), opts)

	if err != nil {
		log.Fatal(err)
	}

	return pb.NewUserServiceClient(clientConn)
}
