package resource

// Snapshot is the state of resources set at a particular point in time
type Snapshot struct {
	Clusters	VersionedResources
	Endpoints	VersionedResources
	Listeners	VersionedResources
	Routers		VersionedResources
}

func NewSnapshot(version string, 
	clusters []Resource, 
	endpoints []Resource, 
	listeners []Resource, 
	routers []Resource) Snapshot {
		return Snapshot {
			Clusters: GetVersionedResources(version, clusters),
			Endpoints: GetVersionedResources(version, endpoints),
			Listeners: GetVersionedResources(version, listeners),
			Routers: GetVersionedResources(version, routers),
		}
}

func (snapshot *Snapshot) GetResourceVersion(typeUrl string) string {
	if snapshot == nil {
		return ""
	}
	switch typeUrl {
	case ClusterType:
		return snapshot.Clusters.Version
	case EndpointType:
		return snapshot.Endpoints.Version
	case ListenerType:
		return snapshot.Listeners.Version
	case RouteType:
		return snapshot.Routers.Version
	}
	return ""
}

func (snapshot *Snapshot) GetResources(typeUrl string) map[string]Resource {
	if snapshot == nil {
		return nil
	}
	switch typeUrl {
	case ClusterType:
		return snapshot.Clusters.Items
	case EndpointType:
		return snapshot.Endpoints.Items
	case ListenerType:
		return snapshot.Listeners.Items
	case RouteType:
		return snapshot.Routers.Items
	}
	return nil
}
