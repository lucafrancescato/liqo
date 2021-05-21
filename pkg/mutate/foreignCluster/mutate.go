package foreigncluster

import (
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
)

type patchType struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func MutateForeignCluster(body []byte) ([]byte, error) {
	var err error

	// unmarshal request into AdmissionReview struct
	admReview := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	var foreignCluster *discoveryv1alpha1.ForeignCluster

	responseBody := []byte{}
	admissionReviewRequest := admReview.Request
	resp := admissionv1.AdmissionResponse{}

	if admissionReviewRequest == nil {
		return responseBody, fmt.Errorf("received admissionReview with empty request")
	}

	// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
	if err := json.Unmarshal(admissionReviewRequest.Object.Raw, &foreignCluster); err != nil {
		return nil, fmt.Errorf("unable unmarshal ForeignCluster json object %v", err)
	}

	// set response options
	resp.Allowed = true
	resp.UID = admissionReviewRequest.UID
	patchStrategy := admissionv1.PatchTypeJSONPatch
	resp.PatchType = &patchStrategy

	patches := []patchType{}

	if needsClusterIdentityDefaulting(foreignCluster) {
		defaultingPatches, err := clusterIdentityDefaulting(foreignCluster)
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		patches = append(patches, defaultingPatches...)
	}

	if resp.Patch, err = json.Marshal(patches); err != nil {
		return nil, err
	}

	resp.Result = &metav1.Status{
		Status: "Success",
	}

	admReview.Response = &resp
	responseBody, err = json.Marshal(admReview)
	if err != nil {
		return nil, err
	}

	klog.Infof("foreignCluster %s patched", foreignCluster.Name)
	klog.V(8).Infof("response: %s", string(responseBody))

	return responseBody, nil
}
