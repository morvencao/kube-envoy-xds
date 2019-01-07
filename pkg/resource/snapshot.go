package resource

// Snapshot is the state of resources set at a particular point in time
type Snapshot struct {
	Clusters	VersionedResources
	Endpoints	VersionedResources
	Listeners	VersionedResources
	Routers		VersionedResources
}

func (snapshot *Snapshot) GetResourceVersion(string typeUrl) string {
	switch typeUrl {
	case ClusterType:
		return snapshot.Clusters.Version
	case EndpointType:
		return snapshot.Endpoints.Version
	case ListenersType:
		return snapshot.Listeners.Version
	case RoutesType:
		return snapshot.Routers.Version
	}
}

func (snapshot *Snapshot) GetResources(string typeUrl) {
	
}
