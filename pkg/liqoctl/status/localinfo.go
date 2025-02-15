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

package status

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/liqotech/liqo/pkg/liqoctl/install"
	"github.com/liqotech/liqo/pkg/liqoctl/output"
	"github.com/liqotech/liqo/pkg/liqoctl/util"
	"github.com/liqotech/liqo/pkg/utils"
	"github.com/liqotech/liqo/pkg/utils/getters"
)

// LocalInfoChecker implements the Check interface.
// holds the Localinformation about the cluster.
type LocalInfoChecker struct {
	options          *Options
	localInfoSection output.Section
	collectionErrors []collectionError
	getReleaseValues func() (map[string]interface{}, error)
}

const (
	localInfoCheckerName = "Local Cluster Information"
)

// newPodChecker return a new pod checker.
func newLocalInfoChecker(options *Options) *LocalInfoChecker {
	return &LocalInfoChecker{
		localInfoSection: output.NewRootSection(),
		options:          options,
		getReleaseValues: func() (map[string]interface{}, error) {
			return options.HelmClient().GetReleaseValues(install.LiqoReleaseName, true)
		},
	}
}

// newPodChecker return a new pod checker.
func newLocalInfoCheckerTest(options *Options, helmValues map[string]interface{}) *LocalInfoChecker {
	return &LocalInfoChecker{
		localInfoSection: output.NewRootSection(),
		options:          options,
		getReleaseValues: func() (map[string]interface{}, error) {
			return helmValues, nil
		},
	}
}

// Collect implements the collect method of the Checker interface.
// it collects the infos of the local cluster.
func (lic *LocalInfoChecker) Collect(ctx context.Context) error {
	clusterIdentity, err := utils.GetClusterIdentityWithControllerClient(ctx, lic.options.CRClient, lic.options.LiqoNamespace)
	if err != nil {
		lic.addCollectionError("ClusterIdentity", "Getting clusterIdentity", err)
	}
	values, err := lic.getReleaseValues()
	if err != nil {
		lic.addCollectionError("InfoRetrieval", "Getting release values", err)
	}
	clusterLabels, err := util.ExtractValuesFromNestedMaps(values, "discovery", "config", "clusterLabels")
	if err != nil {
		lic.addCollectionError("Cluster Labels", "Getting cluster labels", err)
	}

	clusterIdentitySection := lic.localInfoSection.AddSection("Cluster Identity")
	clusterIdentitySection.AddEntry("Cluster ID", clusterIdentity.ClusterID)
	clusterIdentitySection.AddEntry("Cluster Name", clusterIdentity.ClusterName)
	if clusterLabelsMap, ok := clusterLabels.(map[string]interface{}); ok {
		clusterLabelsSection := clusterIdentitySection.AddSection("Cluster Labels")
		for k, v := range clusterLabelsMap {
			clusterLabelsSection.AddEntry(k, v.(string))
		}
	}

	networkSection := lic.localInfoSection.AddSection("Network")
	ipamStorage, err := getters.GetIPAMStorageByLabel(ctx, lic.options.CRClient, labels.NewSelector())
	if err != nil {
		lic.addCollectionError("Network", "Getting IPAM storage", err)
	} else {
		networkSection.AddEntry("Pod CIDR", ipamStorage.Spec.PodCIDR)
		networkSection.AddEntry("Service CIDR", ipamStorage.Spec.ServiceCIDR)
		networkSection.AddEntry("External CIDR", ipamStorage.Spec.ExternalCIDR)
		if len(ipamStorage.Spec.ReservedSubnets) != 0 {
			networkSection.AddEntry("Reserved Subnets", ipamStorage.Spec.ReservedSubnets...)
		}
	}
	apiServerAddress, err := util.ExtractValuesFromNestedMaps(values, "apiServer", "address")
	if err != nil {
		lic.addCollectionError("Kubernetes API Server", "Getting address", err)
	}
	if apiServerAddressString, ok := apiServerAddress.(string); ok && apiServerAddressString != "" {
		lic.localInfoSection.AddSection("Kubernetes API Server").
			AddEntry("Address", apiServerAddressString)
	}
	return nil
}

// GetTitle implements the getTitle method of the Checker interface.
// it returns the title of the checker.
func (lic *LocalInfoChecker) GetTitle() string {
	return localInfoCheckerName
}

// Format implements the format method of the Checker interface.
// it outputs the infos about the local cluster in a string ready to be
// printed out.
func (lic *LocalInfoChecker) Format() (string, error) {
	text := ""
	var err error
	if len(lic.collectionErrors) == 0 {
		text, err = lic.localInfoSection.SprintForBox(lic.options.Printer)
	} else {
		for _, cerr := range lic.collectionErrors {
			text += lic.options.Printer.Error.Sprintfln(lic.options.Printer.Paragraph.Sprintf("%s\t%s\t%s",
				cerr.appName,
				cerr.appType,
				cerr.err))
		}
		text = strings.TrimRight(text, "\n")
	}
	return text, err
}

// HasSucceeded return true if no errors have been kept.
func (lic *LocalInfoChecker) HasSucceeded() bool {
	return len(lic.collectionErrors) == 0
}

// addCollectionError adds a collection error. A collection error is an error that happens while
// collecting the status of a Liqo component.
func (lic *LocalInfoChecker) addCollectionError(localInfoType, localInfoName string, err error) {
	lic.collectionErrors = append(lic.collectionErrors, newCollectionError(localInfoType, localInfoName, err))
}
