package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"

	api "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	cache "github.com/morvencao/kube-envoy-xds/pkg/cache"
	resource "github.com/morvencao/kube-envoy-xds/pkg/resource"
	xds "github.com/morvencao/kube-envoy-xds/pkg/server"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 15010, "The gRPC server port")

const (
	version1 = "v1.0"
	version2 = "v2.0"
	dataserviceEnvoyNode = "dataservice-envoy"
	dataserviceClusterName = "dataservice"
	dataserviceEndpointHost = "127.0.0.1"
	dataserviceEndpointPort = 9090
	dataserviceListenerName = "dataservice-http_listener"
	dataserviceListenerHost = "0.0.0.0"
	dataserviceListenerPort = 8300
	dataserviceTrafficDirection = "ingress"
	dataserviceRouteName = "local_route"
	dataservicePrefix = "/data"
)

var (
	cluster = resource.MakeCluster(dataserviceClusterName)
	endpoint = resource.MakeEndpoint(dataserviceClusterName, dataserviceEndpointHost, dataserviceEndpointPort)
	listener = resource.MakeHTTPListener(dataserviceListenerName, dataserviceListenerHost, dataserviceListenerPort, dataserviceTrafficDirection, dataserviceRouteName)
	router = resource.MakeRoute(dataserviceClusterName, dataserviceRouteName, dataservicePrefix)
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	
	snapshotCache := cache.NewSnapshotCache()
	snapshot := resource.NewSnapshot(version1, []resource.Resource{cluster}, []resource.Resource{endpoint}, []resource.Resource{listener}, []resource.Resource{router})
	snapshotCache.SetSnapshot(dataserviceEnvoyNode, snapshot)

	grpcSvr := grpc.NewServer()
	xdsSvr, err := xds.NewXDSServer(snapshotCache)
	if err != nil {
		glog.Fatalf("failed to create GRPC server: %v", err)
	}

	api.RegisterClusterDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterEndpointDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterListenerDiscoveryServiceServer(grpcSvr, xdsSvr)
	api.RegisterRouteDiscoveryServiceServer(grpcSvr, xdsSvr)

	// start the gRPC server
	glog.Info("starting the xDS server...")
	grpcSvr.Serve(lis)
}
