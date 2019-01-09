package resource

import (
	proto "github.com/gogo/protobuf/proto"
	core "github.com/morvencao/kube-envoy-xds/envoy/api/v2/core"
	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
)

const (
	typePrefix = "type.googleapis.com/"
	ClusterType = typePrefix + "envoy.api.v2.Cluster"
	EndpointType = typePrefix + "envoy.api.v2.ClusterLoadAssignment"
	ListenerType = typePrefix + "envoy.api.v2.Listener"
	RouteType = typePrefix + "envoy.api.v2.RouteConfiguration"
)

// Resource is the base interface for the xDS payload.
type Resource interface {
	proto.Message
	Equal(interface{}) bool
}

// Resources is a versioned group of resources.
type VersionedResources struct {
	Version string
	Items map[string]Resource
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
	return "anynode"
}
