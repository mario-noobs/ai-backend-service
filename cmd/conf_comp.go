package cmd

import (
	"flag"
	sctx "github.com/viettranx/service-context"
	"golang-ai-management/common"
)

type config struct {
	grpcPort           int    // for server port listening
	grpcServerAddress  string // for client make grpc client connection
	grpcUserAddress    string
	grpcProfileAddress string
}

func NewConfig() *config {
	return &config{}
}

func (c *config) ID() string {
	return common.KeyCompConf
}

func (c *config) InitFlags() {
	flag.IntVar(
		&c.grpcPort,
		"grpc-port",
		3100,
		"gRPC Port. Default: 3100",
	)

	flag.StringVar(
		&c.grpcServerAddress,
		"grpc-server-address",
		"localhost:3101",
		"gRPC server address. Default: localhost:3101",
	)

	flag.StringVar(
		&c.grpcUserAddress,
		"grpc-auth-user-address",
		"localhost:3101",
		"gRPC user address. Default: localhost:3201",
	)

	flag.StringVar(
		&c.grpcProfileAddress,
		"grpc-profile-address",
		"localhost:3201",
		"gRPC user address. Default: localhost:3201",
	)
}

func (c *config) Activate(_ sctx.ServiceContext) error {
	return nil
}

func (c *config) Stop() error {
	return nil
}

func (c *config) GetGRPCPort() int {
	return c.grpcPort
}

func (c *config) GetGRPCServerAddress() string {
	return c.grpcServerAddress
}

func (c *config) GetGRPCUserAddress() string {
	return c.grpcUserAddress
}
func (c *config) GetGRPCProfileAddress() string {
	return c.grpcProfileAddress
}
