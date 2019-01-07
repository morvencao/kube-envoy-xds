package cache

import (
	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	resource "github.com/morvencao/kube-envoy-xds/pkg/resource"
)

type Cache interface {
	CreateResponse(*v2.DiscoveryRequest, chan *resource.Response) error
}

type SnapshotCache interface {
	Cache
	SetSnapshot(nodeID string, snapshot *resource.Snapshot) error
	ClearSnapshot(nodeID string) error
}

type snapshotCache struct {
	hash resource.NodeHash
	snapshots	map[string]*resource.Snapshot
	mu	sync.RWMutex
}

func NewSnapshotCache() *SnapshotCache {
	return &SnapshotCache{
		snapshots: make(map[string]*resource.Snapshot),
	}
}

func (scache *snapshotCache) CreateResponse(req *v2.DiscoveryRequest, responseChan chan *resource.Response) error {
	nodeID := scache.hash.ID(req.Node)
	
	scache.mu.Lock()
	defer scache.mu.Unlock()
	
	responseChan := make(chan *v2.DiscoveryRequest, 1)

	snapshot, exists := scache.snapshots[nodeID]
	if !exists || snapshot.GetResourceVersion(req.TypeUrl) == req.VersionInfo {
		return responseChan
	}

	//TODO
	responseChan := scache.GenerateResponse(req, responseChan, snapshot.GetResources(req.TypeUrl))
}
 
func (scache *snapshotCache) SetSnapshot(nodeID string, snapshot *resource.Snapshot) error {
	scache.mu.Lock()
	defer scache.mu.Unlock()

	scache[nodeID] = snapshot
}

func (scache *snapshotCache) ClearSnapshot(nodeID string) error {
	scache.mu.Lock()
	defer scache.mu.Unlock()
	
	delete(scache.snapshots, nodeID)
}
