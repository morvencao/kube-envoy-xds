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

func (scache *snapshotCache) CreateResponse(req *v2.DiscoveryRequest, responseChan chan resource.Response) error {
	nodeID := scache.hash.ID(req.Node)
	
	scache.mu.Lock()
	defer scache.mu.Unlock()

	snapshot, exists := scache.snapshots[nodeID]
	if !exists || snapshot.GetResourceVersion(req.TypeUrl) == req.VersionInfo {
		// TODO
	}

	newVersion := snapshot.GetResourceVersion(req.TypeUrl)
	// Create response from cache
	responseChan <- scache.GenerateResponse(req, snapshot.GetResources(req.TypeUrl), newVersion)
}

func (scache *snapshotCache) GenerateResponse(req *v2.DiscoveryRequest, resources map[string]resource.Resource, version string) resource.Response {
	out := make([]resource.Resource, 0, len(resources))
	if len(req.ResourceName) != 0 {
		for i, res := range resources {
			if contains(req.ResourceName, i) {
				out = append(out, res)
			}
		}
	} else {
		for _, res := range resources {
			out = append(out, res)
		}
	}

	return resource.Response{
		Request: req,
		Version: version,
		Resources: out,
	}
}

func contains(s []string, item string) bool {
	for _, n := range s {
		if item == n {
			return true
		}
	}
	return false
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
