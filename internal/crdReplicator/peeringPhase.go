// Copyright 2019-2022 The Liqo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crdreplicator

import (
	"k8s.io/klog/v2"

	"github.com/liqotech/liqo/internal/crdReplicator/resources"
	"github.com/liqotech/liqo/pkg/consts"
)

// getPeeringPhase returns the peering phase for a cluster given its clusterID.
func (c *Controller) getPeeringPhase(clusterID string) consts.PeeringPhase {
	c.peeringPhasesMutex.RLock()
	defer c.peeringPhasesMutex.RUnlock()
	if c.peeringPhases == nil {
		return consts.PeeringPhaseNone
	}
	if phase, ok := c.peeringPhases[clusterID]; ok {
		return phase
	}
	return consts.PeeringPhaseNone
}

// setPeeringPhase sets the peering phase for a given clusterID.
func (c *Controller) setPeeringPhase(clusterID string, phase consts.PeeringPhase) {
	c.peeringPhasesMutex.Lock()
	defer c.peeringPhasesMutex.Unlock()
	if c.peeringPhases == nil {
		c.peeringPhases = map[string]consts.PeeringPhase{}
	}
	c.peeringPhases[clusterID] = phase
}

// isReplicationAllowed indicates if the given peering phase matches the one required by the given resource
func isReplicationAllowed(peeringPhase consts.PeeringPhase, resource *resources.Resource) bool {
	switch resource.PeeringPhase {
	case consts.PeeringPhaseNone:
		return false
	case consts.PeeringPhaseAuthenticated:
		return peeringPhase != consts.PeeringPhaseNone
	case consts.PeeringPhaseBidirectional:
		return peeringPhase == consts.PeeringPhaseBidirectional
	case consts.PeeringPhaseIncoming:
		return peeringPhase == consts.PeeringPhaseBidirectional || peeringPhase == consts.PeeringPhaseIncoming
	case consts.PeeringPhaseOutgoing:
		return peeringPhase == consts.PeeringPhaseBidirectional || peeringPhase == consts.PeeringPhaseOutgoing
	case consts.PeeringPhaseEstablished:
		bidirectional := peeringPhase == consts.PeeringPhaseBidirectional
		incoming := peeringPhase == consts.PeeringPhaseIncoming
		outgoing := peeringPhase == consts.PeeringPhaseOutgoing
		induced := peeringPhase == consts.PeeringPhaseInduced
		return bidirectional || incoming || outgoing || induced
	default:
		klog.Warning("Unknown peering phase %v", resource.PeeringPhase)
		return false
	}
}
