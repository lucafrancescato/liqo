// Copyright 2019-2021 The Liqo Authors
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

package workload_test

import (
	"context"
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/liqotech/liqo/pkg/virtualKubelet/forge"
	"github.com/liqotech/liqo/pkg/virtualKubelet/reflection/options"
)

const (
	LocalNamespace  = "local-namespace"
	RemoteNamespace = "remote-namespace"

	LocalClusterID  = "local-cluster"
	RemoteClusterID = "remote-cluster"

	LiqoNodeName = "local-node"
	LiqoNodeIP   = "1.1.1.1"
)

var ctx context.Context

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workload Reflection Suite")
}

var _ = BeforeSuite(func() {
	klog.SetOutput(GinkgoWriter)
	flagset := flag.NewFlagSet("klog", flag.PanicOnError)
	klog.InitFlags(flagset)
	Expect(flagset.Set("v", "4")).To(Succeed())
	klog.LogToStderr(false)

	forge.Init(LocalClusterID, RemoteClusterID, LiqoNodeName, LiqoNodeIP)
})

var _ = BeforeEach(func() { ctx = context.Background() })

var FakeEventHandler = func(options.Keyer) cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(_ interface{}) {},
		UpdateFunc: func(_, obj interface{}) {},
		DeleteFunc: func(_ interface{}) {},
	}
}

var GetPod = func(client kubernetes.Interface, namespace, name string) *corev1.Pod {
	pod, errpod := client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	Expect(errpod).ToNot(HaveOccurred())
	return pod
}

var GetPodError = func(client kubernetes.Interface, namespace, name string) error {
	_, errpod := client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	return errpod
}
