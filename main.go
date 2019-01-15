package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"

	api "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	discovery "github.com/morvencao/kube-envoy-xds/envoy/service/discovery/v2"
	cache "github.com/morvencao/kube-envoy-xds/pkg/cache"
	resource "github.com/morvencao/kube-envoy-xds/pkg/resource"
	xds "github.com/morvencao/kube-envoy-xds/pkg/server"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 15010, "The gRPC server port.")
var mode = flag.String("mode", "xds", "The mode of envoy management server, can be 'xds' or 'ads'(default value).")

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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	cluster := resource.MakeCluster(*mode, dataserviceClusterName)
	endpoint := resource.MakeEndpoint(dataserviceClusterName, dataserviceEndpointHost, dataserviceEndpointPort)
	listener := resource.MakeHTTPListener(*mode, dataserviceListenerName, dataserviceListenerHost, dataserviceListenerPort, dataserviceTrafficDirection, dataserviceRouteName)
	router := resource.MakeRoute(dataserviceClusterName, dataserviceRouteName, dataservicePrefix)
	
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
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcSvr, xdsSvr)

	// start the gRPC server
	glog.Info("starting the xDS server...")
	grpcSvr.Serve(lis)
}
