package server

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/glog"

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
	glog.Infof("creating the xDS server with cache: %v", cache)
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
	clustersChan := make(chan resource.Response, 1)
	endpointsChan := make(chan resource.Response, 1)
	listenersChan := make(chan resource.Response, 1)
	routersChan := make(chan resource.Response, 1)

	var clustersNonce, endpointsNonce, listenersNonce, routersNonce string

	streamNonce := int64(0)

	glog.Info("starting process the xDS request on the stream...")
	for {
		select {
		case clusters, more := <-clustersChan:
			if !more {
				return status.Errorf(codes.Unavailable, "get clusters failed")
			}
			glog.Infof("ready to send the new CDS response with the stream: %v", clusters)
			out, err := buildResponse(&clusters, resource.ClusterType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			clustersNonce = out.Nonce
			if err := stream.Send(out); err!= nil {
				return err
			}
		case endpoints, more := <-endpointsChan:
			if !more {
				return status.Errorf(codes.Unavailable, "get endpoints failed")
			}
			glog.Infof("ready to send the new EDS response with the stream: %v", endpoints)
			out, err := buildResponse(&endpoints, resource.EndpointType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			endpointsNonce = out.Nonce
			if err := stream.Send(out); err!= nil {
				return err
			}
		case listeners, more := <-listenersChan:
			if !more {
				return status.Errorf(codes.Unavailable, "get listeners failed")
			}
			glog.Infof("ready to send the new LDS response with the stream: %v", listeners)
			out, err := buildResponse(&listeners, resource.ListenerType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			listenersNonce = out.Nonce
			if err := stream.Send(out); err!= nil {
				return err
			}
		case routers, more := <-routersChan:
			if !more {
				return status.Errorf(codes.Unavailable, "get routers failed")
			}
			glog.Infof("ready to send the new RDS response with the stream: %v", routers)
			out, err := buildResponse(&routers, resource.RouteType)
			if err != nil {
				return err
			}
			currentNonce := atomic.AddInt64(&streamNonce, 1)
			out.Nonce = strconv.FormatInt(currentNonce, 10)
			routersNonce = out.Nonce
			if err := stream.Send(out); err!= nil {
				return err
			}
		case req, more := <-reqCh:
			if !more {
				glog.Warning("no more xDS request on the stream, exiting...")
				return nil
			}
			if req == nil {
				glog.Error("empty xDS request on the stream, exiting...")
				return fmt.Errorf("empty request")
			}
			repNonce := req.GetResponseNonce()
			typeUrl := req.GetTypeUrl()

			switch {
			case typeUrl == resource.ClusterType && (clustersNonce == "" || clustersNonce == repNonce):
				glog.Infof("last CDS response has been sent, ready to process the new CDS request: %v", req)
				resp, err := svr.cache.CreateResponse(req)
				if err != nil {
					glog.Errorf("failed to create CDS response from cache: %v", err)
					return err
				}
				glog.Infof("CDS response from cache: %v", *resp)
				clustersChan <- *resp
			case typeUrl == resource.EndpointType && (endpointsNonce == "" || endpointsNonce == repNonce):
				glog.Infof("last EDS response has been sent, ready to process the new EDS request: %v", req)
				resp, err := svr.cache.CreateResponse(req)
				if err != nil {
					glog.Errorf("failed to create EDS response from cache: %v", err)
					return err
				}
				endpointsChan <- *resp
			case typeUrl == resource.ListenerType && (listenersNonce == "" || listenersNonce == repNonce):
				glog.Infof("last LDS response has been sent, ready to process the new LDS request: %v", req)
				resp, err := svr.cache.CreateResponse(req)
				if err != nil {
					glog.Errorf("failed to create LDS response from cache: %v", err)
					return err
				}
				listenersChan <- *resp
			case typeUrl == resource.RouteType && (routersNonce == "" || routersNonce == repNonce):
				glog.Infof("last RDS response has been sent, ready to process the new RDS request: %v", req)
				resp, err := svr.cache.CreateResponse(req)
				if err != nil {
					glog.Errorf("failed to create RDS response from cache: %v", err)
					return err
				}
				routersChan <- *resp
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
			glog.Infof("get new xDS request on the stream: %v", req)
			reqCh <- req
		}
	}()

	err := svr.processReq(stream, reqCh)
	atomic.StoreInt32(&stopFlag, 1)

	return err
}

func (svr *xDSServer) StreamClusters(stream v2.ClusterDiscoveryService_StreamClustersServer) error {
	glog.Info("starting handling CDS stream request...")
	return svr.handler(stream)
}

func (svr *xDSServer) IncrementalClusters(stream v2.ClusterDiscoveryService_IncrementalClustersServer) error {
	return fmt.Errorf("not implemented")
}

func (svr *xDSServer) FetchClusters(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	glog.Infof("fetching clusters for single CDS request: %v", req)
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.ClusterType)
}

func (svr *xDSServer) StreamEndpoints(stream v2.EndpointDiscoveryService_StreamEndpointsServer) error {
	glog.Info("starting handling EDS stream request...")
	return svr.handler(stream)
}

func (svr *xDSServer) FetchEndpoints(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	glog.Infof("fetching endpoints for single EDS request: %v", req)
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.EndpointType)
}

func (svr *xDSServer) StreamListeners(stream v2.ListenerDiscoveryService_StreamListenersServer) error {
	glog.Info("starting handling LDS stream request...")
	return svr.handler(stream)
}

func (svr *xDSServer) FetchListeners(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	glog.Infof("fetching listeners for single LDS request: %v", req)
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.ListenerType)
}

func (svr *xDSServer) StreamRoutes(stream v2.RouteDiscoveryService_StreamRoutesServer) error {
	glog.Info("starting handling RDS stream request...")
	return svr.handler(stream)
}
func (svr *xDSServer) IncrementalRoutes(stream v2.RouteDiscoveryService_IncrementalRoutesServer) error {
	return fmt.Errorf("not implemented")
}
func (svr *xDSServer) FetchRoutes(ctx context.Context, req *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	glog.Infof("fetching routes for single RDS request: %v", req)
	response, err := svr.cache.FetchResponse(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(response, resource.RouteType)
}
