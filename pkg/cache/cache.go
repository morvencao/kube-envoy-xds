package cache

import (
	"fmt"
	"sync"

	v2 "github.com/morvencao/kube-envoy-xds/envoy/api/v2"
	resource "github.com/morvencao/kube-envoy-xds/pkg/resource"
)

type Cache interface {
	CreateResponse(*v2.DiscoveryRequest, chan resource.Response) error
	FetchResponse(*v2.DiscoveryRequest) (*resource.Response, error)
}

type SnapshotCache interface {
	Cache
	SetSnapshot(nodeID string, snapshot resource.Snapshot)
	ClearSnapshot(nodeID string)
}

type snapshotCache struct {
	snapshots	map[string]resource.Snapshot
	mu	sync.RWMutex
}

func NewSnapshotCache() SnapshotCache {
	return &snapshotCache{
		snapshots: make(map[string]resource.Snapshot),
	}
}

func (scache *snapshotCache) CreateResponse(req *v2.DiscoveryRequest, responseChan chan resource.Response) error {
	nodeID := resource.GetNodeID(req.Node)

	scache.mu.Lock()
	defer scache.mu.Unlock()

	snapshot, exists := scache.snapshots[nodeID]
	if !exists || snapshot.GetResourceVersion(req.TypeUrl) == req.VersionInfo {
		// TODO: wait until update for resource are ready.
		return fmt.Errorf("no new version found")
	}

	newVersion := snapshot.GetResourceVersion(req.TypeUrl)
	// Create response from cache
	responseChan <- scache.GenerateResponse(req, snapshot.GetResources(req.TypeUrl), newVersion)
	return nil
}

func (scache *snapshotCache) GenerateResponse(req *v2.DiscoveryRequest, resources map[string]resource.Resource, version string) resource.Response {
	out := make([]resource.Resource, 0, len(resources))
	if len(req.ResourceNames) != 0 {
		for i, res := range resources {
			if contains(req.ResourceNames, i) {
				out = append(out, res)
			}
		}
	} else {
		for _, res := range resources {
			out = append(out, res)
		}
	}

	return resource.Response{
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

func (scache *snapshotCache) FetchResponse(req *v2.DiscoveryRequest) (*resource.Response, error) {
	nodeID := resource.GetNodeID(req.Node)

	scache.mu.Lock()
	defer scache.mu.Unlock()

	snapshot, exists := scache.snapshots[nodeID]
	if !exists || snapshot.GetResourceVersion(req.TypeUrl) == req.VersionInfo {
		return nil, fmt.Errorf("no new version found")
	}

	newVersion := snapshot.GetResourceVersion(req.TypeUrl)
	resp := scache.GenerateResponse(req, snapshot.GetResources(req.TypeUrl), newVersion)
	return &resp, nil
}

func (scache *snapshotCache) SetSnapshot(nodeID string, snapshot resource.Snapshot) {
	scache.mu.Lock()
	defer scache.mu.Unlock()

	scache.snapshots[nodeID] = snapshot
}

func (scache *snapshotCache) ClearSnapshot(nodeID string) {
	scache.mu.Lock()
	defer scache.mu.Unlock()

	delete(scache.snapshots, nodeID)
}
