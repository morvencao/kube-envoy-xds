package server

import (
	"context"

	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"

)

type grpcServer struct {

}

func NewServer() (*grpcServer, error) {
	return &grpcServer{}, nil
}

func (svr *grpcServer) StreamClusters(stream v2.ClusterDiscoveryService_StreamClustersServer) error {
	
}

func (svr *grpcServer) IncrementalClusters(stream v2.ClusterDiscoveryService_IncrementalClustersServer) error {

}

func (svr *grpcServer) FetchClusters(ctx context.Context, dsReq *DiscoveryRequest) (*DiscoveryResponse, error) {

}
