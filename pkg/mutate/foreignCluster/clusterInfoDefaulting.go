package foreigncluster

import (
	"k8s.io/klog/v2"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	"github.com/liqotech/liqo/internal/discovery/utils"
)

// needsClusterIdentityDefaulting checks if the ForeignCluster CR does not have a value in one of the required fields (i.e. ClusterID)
// and needs a value defaulting.
func needsClusterIdentityDefaulting(foreignCluster *discoveryv1alpha1.ForeignCluster) bool {
	return foreignCluster.Spec.ClusterIdentity.ClusterID == ""
}

// clusterIdentityDefaulting loads the default values for that ForeignCluster basing on the AuthUrl value, an HTTP request is sent and
// the retrieved values are applied for the following fields (if they are empty): ClusterIdentity.ClusterID,
// ClusterIdentity.Namespace and the TrustMode
// if it returns no error, the ForeignCluster CR has been updated.
func clusterIdentityDefaulting(foreignCluster *discoveryv1alpha1.ForeignCluster) ([]patchType, error) {
	klog.V(4).Infof("Defaulting ClusterIdentity values for ForeignCluster %v", foreignCluster.Name)
	ids, trustMode, err := utils.GetClusterInfo(foreignCluster.Spec.AuthURL)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	patches := []patchType{}

	clusterIdentity := foreignCluster.Spec.ClusterIdentity
	if clusterIdentity.ClusterID == "" {
		patches = append(patches, patchType{
			Op:    "add",
			Path:  "spec/clusterIdentity/clusterID",
			Value: ids.ClusterID,
		})
	}
	if clusterIdentity.ClusterName == "" {
		patches = append(patches, patchType{
			Op:    "add",
			Path:  "spec/clusterIdentity/clusterName",
			Value: ids.ClusterName,
		})
	}

	foreignCluster.Spec.TrustMode = trustMode
	if foreignCluster.Spec.TrustMode != trustMode {
		patches = append(patches, patchType{
			Op:    "replace",
			Path:  "spec/trustMode",
			Value: trustMode,
		})
	}

	klog.V(4).Infof("New values:\n\tClusterId:\t%v\n\tClusterName:\t%v\n\tTrustMode:\t%v",
		ids.ClusterID, ids.ClusterName, trustMode)

	return patches, nil
}
