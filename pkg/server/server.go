package server

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"google.golang.org/grpc"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"

	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	cache "github.com/morvencao/kube-envoy-xds/pkg/cache"
	resource "github.com/morvencao/kube-envoy-xds/pkg/resource"
)

type XDSServer interface {
	v2.ClusterDiscoveryServiceServer
	v2.EndpointDiscoveryServiceServer
	v2.ListenerDiscoveryServiceServer
	v2.RouteDiscoveryServiceServer
}

type xDSServer struct {
	cache	cache.Cache
}

func NewXDSServer(cache cache.Cache) (XDSServer, error) {
	return &xDSServer{cache: cache}, nil
}

type xDSStream interface {
	grpc.ServerStream
	Recv() (*v2.DiscoveryRequest, error)
	Send(*v2.DiscoveryResponse) error
}

func buildResponse(response *resource.Response, typeUrl string) (*v2.DiscoveryResponse, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}
	resources := make([]types.Any, len(response.Resources))
	for i, resource := range response.Resources {
		data, err := proto.Marshal(resource)
		if err != nil {
			return nil, err
		}
		resources[i] = types.Any{
			TypeUrl: typeUrl,
			Value: data,
		}
	}
	out := &v2.DiscoveryResponse{
		VersionInfo: response.Version,
		Resources: resources,
		TypeUrl: typeUrl,
	}
	return out, nil
}

func (svr *xDSServer) processReq(stream v2.ClusterDiscoveryService_StreamClustersServer, reqCh <-chan *v2.DiscoveryRequest) error {
	clustersChan := make(chan resource.Response)
	endpointsChan := make(chan resource.Response)
	listenersChan := make(chan resource.Response)
	routersChan := make(chan resource.Response)

	var clustersNonce, endpointsNonce, listenersNonce, routersNonce string

	streamNonce := int64(0)

	for {
		select {
		case clusters := <-clustersChan:
			out, err := buildResponse(&clusters, resource.ClusterType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			clustersNonce = out.Nonce
			err = stream.Send(out)
			return err
		case endpoints := <-endpointsChan:
			out, err := buildResponse(&endpoints, resource.EndpointType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			endpointsNonce = out.Nonce
			err = stream.Send(out)
			return err
		case listeners := <-listenersChan:
			out, err := buildResponse(&listeners, resource.ListenerType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			listenersNonce = out.Nonce
			err = stream.Send(out)
			return err
		case routers := <-routersChan:
			out, err := buildResponse(&routers, resource.RouteType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			routersNonce = out.Nonce
			err = stream.Send(out)
			return err
		case req, more := <-reqCh:
			if !more {
				return nil
			}
			if req == nil {
				return fmt.Errorf("empty request")
			}
			repNonce := req.GetResponseNonce()
			typeUrl := req.GetTypeUrl()

			switch {
			case typeUrl == resource.ClusterType && (clustersNonce == "" || clustersNonce == repNonce):
				if err := svr.cache.CreateResponse(req, clustersChan); err != nil {
					return err
				}
			case typeUrl == resource.EndpointType && (endpointsNonce == "" || endpointsNonce == repNonce):
				if err := svr.cache.CreateResponse(req, endpointsChan); err != nil {
					return err
				}
			case typeUrl == resource.EndpointType && (listenersNonce == "" || listenersNonce == repNonce):
				if err := svr.cache.CreateResponse(req, listenersChan); err != nil {
					return err
				}
			case typeUrl == resource.EndpointType && (routersNonce == "" || routersNonce == repNonce):
				if err := svr.cache.CreateResponse(req, routersChan); err != nil {
					return err
				}
			}
		}
	}
}

func (svr *xDSServer) handler(stream xDSStream) error {
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
}

func (svr *xDSServer) StreamClusters(stream v2.ClusterDiscoveryService_StreamClustersServer) error {
	return svr.handler(stream)
}

func (svr *xDSServer) IncrementalClusters(stream v2.ClusterDiscoveryService_IncrementalClustersServer) error {
	return fmt.Errorf("not implemented")
}

func (svr *xDSServer) FetchClusters(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.ClusterType)
}

func (svr *xDSServer) StreamEndpoints(stream v2.EndpointDiscoveryService_StreamEndpointsServer) error {
	return svr.handler(stream)
}

func (svr *xDSServer) FetchEndpoints(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.EndpointType)
}

func (svr *xDSServer) StreamListeners(stream v2.ListenerDiscoveryService_StreamListenersServer) error {
	return svr.handler(stream)
}

func (svr *xDSServer) FetchListeners(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.ListenerType)
}

func (svr *xDSServer) StreamRoutes(stream v2.RouteDiscoveryService_StreamRoutesServer) error {
	return svr.handler(stream)
}
func (svr *xDSServer) IncrementalRoutes(stream v2.RouteDiscoveryService_IncrementalRoutesServer) error {
	return fmt.Errorf("not implemented")
}
func (svr *xDSServer) FetchRoutes(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.RouteType)
}
