package resource

import (
	"time"

	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	core "github.com/morvencao/kube-envoy-xds/envoy/api/v2/core"
	endpoint "github.com/morvencao/kube-envoy-xds/envoy/api/v2/endpoint"
	listener "github.com/morvencao/kube-envoy-xds/envoy/api/v2/listener"
	route "github.com/morvencao/kube-envoy-xds/envoy/api/v2/route"
	al "github.com/morvencao/kube-envoy-xds/envoy/config/accesslog/v2"
	fal "github.com/morvencao/kube-envoy-xds/envoy/config/filter/accesslog/v2"
	hcm "github.com/morvencao/kube-envoy-xds/envoy/config/filter/network/http_connection_manager/v2"
	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	util "github.com/morvencao/kube-envoy-xds/pkg/util"
)

const (
	typePrefix = "type.googleapis.com/"
	ClusterType = typePrefix + "envoy.api.v2.Cluster"
	EndpointType = typePrefix + "envoy.api.v2.ClusterLoadAssignment"
	ListenerType = typePrefix + "envoy.api.v2.Listener"
	RouteType = typePrefix + "envoy.api.v2.RouteConfiguration"
	AnyType = ""  // AnyType is used in ADS
)

const (
	AnyNode = "anynode"
	XdsCluster = "xds_cluster"
	XDS = "xds"
	ADS = "ads"
	HTTPConnectionManager = "envoy.http_connection_manager"
	HTTPGRPCAccessLog = "envoy.http_grpc_access_log"
	Router = "envoy.router"
)

// Resource is the base interface for the xDS payload.
type Resource interface {
	proto.Message
	Equal(interface{}) bool
}

func GetResourceName(res Resource) string {
	switch resType := res.(type) {
	case *v2.Cluster:
		return resType.GetName()
	case *v2.ClusterLoadAssignment:
		return resType.GetClusterName()
	case *v2.Listener:
		return resType.GetName()
	case *v2.RouteConfiguration:
		return resType.GetName()
	default:
		return ""
	}
}

// Resources is a versioned group of resources.
type VersionedResources struct {
	Version string
	Items map[string]Resource
}

func GetVersionedResources(version string, resources []Resource) VersionedResources {
	return VersionedResources{
		Version: version,
		Items: GetResourceMapfromSlice(resources),
	}
}

func GetResourceMapfromSlice(resources []Resource) map[string]Resource {
	out := make(map[string]Resource, len(resources))
	for _, resource := range resources {
		out[GetResourceName(resource)] = resource
	}
	return out
} 

// Request is an alias for the discovery request type.
type Request v2.DiscoveryRequest

// Response is a pre-serialized xDS response.
type Response struct {
	Version string
	Resources []Resource
}

func GetNodeID(node *core.Node) string {
	if node != nil {
		return node.Id
	}
	return AnyNode
}

// MakeCluster creates cluster of EDS with given cluster name.
func MakeCluster(mode, clusterName string) *v2.Cluster {
	var edsSource *core.ConfigSource
	switch mode {
	case XDS:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType: core.ApiConfigSource_GRPC,
					GrpcServices: []*core.GrpcService{{
						TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
								ClusterName: XdsCluster,
							},
						},
					}},
				},
			},
		}
	case ADS:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_Ads{
				Ads: &core.AggregatedConfigSource{},
			},
		}
	}
	return &v2.Cluster{
		Name: clusterName,
		ConnectTimeout: 5*time.Second,
		Type: v2.Cluster_EDS,
		EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
			EdsConfig: edsSource,
		},
	}
}

// MakeEndpoint creates  endpoint on given host and port.
func MakeEndpoint(clusterName, host string, port uint32) *v2.ClusterLoadAssignment {
	return &v2.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: []endpoint.LbEndpoint{{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Protocol: core.TCP,
								Address: host,
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: port,
								},
							},
						},
					},
				},
			}},
		}},
	}
}

// MakeHTTPListener creates listener with given listener name, host, port and route
func MakeHTTPListener(mode, listenerName, host string, port uint32, trafficDirection, routeName string) *v2.Listener {
	// RDS configuration source
	var rdsSource core.ConfigSource
	switch mode {
	case XDS:
		rdsSource = core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType: core.ApiConfigSource_GRPC,
					GrpcServices: []*core.GrpcService{{
						TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
						},
					}},
				},
			},
		}
	case ADS:
		rdsSource = core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_Ads{
				Ads: &core.AggregatedConfigSource{},
			},
		}
	}
	
	// access log configuration
	accessLogConfig := &al.HttpGrpcAccessLogConfig{
		CommonConfig: &al.CommonGrpcAccessLogConfig{
			LogName: "echo",
			GrpcService: &core.GrpcService{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
						ClusterName: XdsCluster,
					},
				},
			},
		},
	}

	accessLogConfigStruct, err := util.MessageToStruct(accessLogConfig)
	if err != nil {
		panic(err)
	}

	var tracingOperationName hcm.HttpConnectionManager_Tracing_OperationName
	switch {
	case trafficDirection == "ingress":
		tracingOperationName = hcm.INGRESS
	case trafficDirection == "egress":
		tracingOperationName = hcm.EGRESS
	default:
		tracingOperationName = hcm.INGRESS
	}
	
	// HTTP filter configuration
	hcManager := &hcm.HttpConnectionManager{
		GenerateRequestId: &types.BoolValue{
			Value: true,
		},
		Tracing: &hcm.HttpConnectionManager_Tracing{
			OperationName: tracingOperationName,
		},
		CodecType: hcm.AUTO,
		StatPrefix: trafficDirection + "_http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: routeName,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: Router,
		}},
		AccessLog: []*fal.AccessLog{{
			Name: HTTPGRPCAccessLog,
			ConfigType: &fal.AccessLog_Config{
				Config: accessLogConfigStruct,
			},
		}},
	}
	hcManagerStruct, err := util.MessageToStruct(hcManager)
	if err != nil {
		panic(err)
	}

	return &v2.Listener{
		Name: listenerName,
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address: host,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name: HTTPConnectionManager,
				ConfigType: &listener.Filter_Config{
					Config: hcManagerStruct,
				},
			}},
		}},
	}
}

// MakeRoute create an HTTP route to given cluster
func MakeRoute(clusterName, routeName, prefix string) *v2.RouteConfiguration {
	return &v2.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []route.VirtualHost{{
			Name:    routeName,
			Domains: []string{"*"},
			Routes: []route.Route{{
				Match: route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: prefix,
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
					},
				},
			}},
		}},
	}
}
