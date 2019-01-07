package server

import (
	"atomic"
	"context"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"

	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	cache "github.com/morvencao/kube-envoy-xds/pkg/cache"
	resource "github.com/morvencao/kube-envoy-xds/pkg/resource"
)

type GRPCServer interface {
	v2.ClusterDiscoveryServiceServer
	v2.EndpointDiscoveryServiceServer
	v2.ListenerDiscoveryServiceServer
	v2.RouteDiscoveryServiceServer
}

type grpcServer struct {
	cache	*cache.Cache
}

func NewGRPCServer(cache *cache.Cache) (*grpcServer, error) {
	return &grpcServer{cache: cache}, nil
}

type gRPCStream interface {
	grpc.ServerStream
	Recv() (*IncrementalDiscoveryRequest, error)
	Send(*IncrementalDiscoveryResponse) error
}

// watches for all xDS resource types
type watch struct {
	clustersChan	chan *resource.Response
	endpointsChan	chan *resource.Response
	listenersChan	chan *resource.Response
	routersChan		chan *resource.Response

	clustersNonce	string
	endpointsNonce	string
	listenersNonce	string
	routersNonce	string
}

func buildResponse(response resource.Response, typeUrl string) (*v2.DiscoveryResponse, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}
	resources := make([]types.Any, len(resource.Resources))
	for i, resource := range response.Resources {
		data, err := proto.Marshal(resource)
		if err != nil {
			return nil, err
		}
		resources[i] = types.Any{
			TypeUrl: typeUrl
			Value: data
		}
	}
	out := &v2.DiscoveryResponse{
		VersionInfo: response.Version,
		Resources: resources,
		TypeUrl: typeUrl
	}
	return out, nil
}

func (svr *grpcServer) processReq(stream v2.ClusterDiscoveryService_StreamClustersServer, reqCh <-chan *v2.DiscoveryRequest) error {
	var xdsWatch := &watch{
		clustersChan: make(chan *resource.Response),
		endpointsChan: make(chan *resource.Response),
		listenersChan: make(chan *resource.Response),
		routersChan: make(chan *resource.Response),
	}
	streamNonce := int64(0)

	for {
		select {
		case clusters := <-xdsWatch.clustersChan:
			out, err := buildResponse(clusters, resource.ClusterType)
			if err != nil {
				return err
			}
			streamNonce = atomic.addInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(streamNonce, 10)
			xdsWatch.clustersNonce = out.Nonce
			err := stream.Send(out)
			return err
		case endpoints := <-xdsWatch.endpointsChan:
			response, err := buildResponse(endpoints)
			if err != nil {
				return err
			}
			streamNonce = atomic.addInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(streamNonce, 10)
			xdsWatch.endpointsNonce = out.Nonce
			err := stream.Send(out)
			return err
		case listeners := <-xdsWatch.listenersChan:
			response, err := buildResponse(listeners)
			if err != nil {
				return err
			}
			streamNonce = atomic.addInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(streamNonce, 10)
			xdsWatch.listenersNonce = out.Nonce
			err := stream.Send(out)
			return err
		case routers := <-xdsWatch.routersChan:
			response, err := buildResponse(routers)
			if err != nil {
				return err
			}
			streamNonce = atomic.addInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(streamNonce, 10)
			xdsWatch.routersNonce = out.Nonce
			err := stream.Send(out)
			return err
		case req, more := <-reqCh:
			if !more {
				return nil
			}
			if req == nil {
				fmt.Errorf("empty request")
			}
			repNonce := req.GetResponseNonce()
			typeUrl := req.GetTypeUrl()

			switch {
			case typeUrl == resource.ClusterType && (xdsWatch.clustersNonce == "" || xdsWatch.clustersNonce == repNonce):
				s.cache.CreateResponse(req, xdsWatch.clustersChan)
			case typeUrl == resource.EndpointType && (xdsWatch.endpointsNonce == "" || xdsWatch.endpointsNonce == repNonce):
				s.cache.CreateResponse(req, xdsWatch.endpointsChan)
			case typeUrl == resource.EndpointType && (xdsWatch.listenersNonce == "" || xdsWatch.listenersNonce == repNonce):
				s.cache.CreateResponse(req, xdsWatch.listenersChan)
			case typeUrl == resource.EndpointType && (xdsWatch.routersNonce == "" || xdsWatch.routersNonce == repNonce):
				s.cache.CreateResponse(req, xdsWatch.routersChan)
			}
		}
	}
}

func (svr *grpcServer) handler(stream gRPCStream) error {
	reqCh := make(chan *v2.DiscoveryRequest)
	stopFlag := int32(0)

	go func() {
		for {
			req, err := stream.Recv()
			if atomic.LoadInt32(&stopFlag) != 0 {
				return
			}
			if err != nil {
				close(reqCh)
				return
			}
			reqCh <- req
		}
	}()

	err := svr.processReq(stream, reqCh)
	atomic.StoreInt32(&stopFlag, 1)

	return err

func (svr *grpcServer) StreamClusters(stream v2.ClusterDiscoveryService_StreamClustersServer) error {
	return svr.handler(stream)
}

func (svr *grpcServer) IncrementalClusters(stream v2.ClusterDiscoveryService_IncrementalClustersServer) error {

}

func (svr *grpcServer) FetchClusters(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {

}

func (svr *grpcServer) StreamEndpoints(stream v2.EndpointDiscoveryService_StreamEndpointsServer) error {
	return svr.handler(stream)
}

func (svr *grpcServer) FetchEndpoints(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {

}

func (svr *grpcServer) StreamListeners(stream v2.ListenerDiscoveryService_StreamListenersServer) error {
	return svr.handler(stream)
}

func (svr *grpcServer) FetchListeners(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {

}

func (svr *grpcServer) StreamRoutes(stream v2.RouteDiscoveryService_StreamRoutesServer) error {
	return svr.handler(stream)
}
func (svr *grpcServer) IncrementalRoutes(stream v2.RouteDiscoveryService_IncrementalRoutesServer) error {

}
func (svr *grpcServer) FetchRoutes(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {

}
