package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	
	api "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	xds "github.com/morvencao/kube-envoy-xds/pkg/server"

	"google.golang.org/grpc"
)

var (
	port := flag.Int("port", 10000, "The gRPC server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcSvr = grpc.NewServer()
	xdsSvr := xds.NewServer()

	api.RegisterClusterDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterEndpointDiscoveryServiceServer(grpcSvr, xdsSvr)

	go func() {
		if err := grpcSvr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
