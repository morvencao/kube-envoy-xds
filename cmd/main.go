package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	
	api "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	cache "github.com/morvencao/kube-envoy-xds/pkg/cache"
	xds "github.com/morvencao/kube-envoy-xds/pkg/server"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 10000, "The gRPC server port")

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	snapshotCache := cache.NewSnapshotCache()
	grpcSvr := grpc.NewServer()
	xdsSvr, err := xds.NewXDSServer(snapshotCache)
	if err != nil {
		log.Fatalf("failed to create GRPC server: %v", err)
	}

	api.RegisterClusterDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterEndpointDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterListenerDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterRouteDiscoveryServiceServer(grpcSvr, xdsSvr)
	
	go func() {
		if err := grpcSvr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
