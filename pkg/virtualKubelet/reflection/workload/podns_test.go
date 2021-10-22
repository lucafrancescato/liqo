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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"

	fakeipam "github.com/liqotech/liqo/pkg/liqonet/ipam/fake"
	"github.com/liqotech/liqo/pkg/virtualKubelet/reflection/manager"
	"github.com/liqotech/liqo/pkg/virtualKubelet/reflection/options"
	"github.com/liqotech/liqo/pkg/virtualKubelet/reflection/workload"
)

var _ = Describe("Namespaced Pod Reflection Tests", func() {

	Describe("pod handling", func() {
		const PodName = "name"

		var (
			reflector manager.NamespacedReflector
			client    kubernetes.Interface

			ipam   *fakeipam.IPAMClient
			cancel context.CancelFunc
		)

		BeforeEach(func() { ctx, cancel = context.WithCancel(ctx) })
		AfterEach(func() { cancel() })

		JustBeforeEach(func() {
			ipam = fakeipam.NewIPAMClient("192.168.200.0/24", "192.168.201.0/24", true)

			client = fake.NewSimpleClientset()
			factory := informers.NewSharedInformerFactory(client, 10*time.Hour)

			metricsFactory := func(string) metricsv1beta1.PodMetricsInterface { return nil }
			rfl := workload.NewPodReflector(nil, metricsFactory, ipam, 0)
			rfl.Start(ctx, options.New(client, factory.Core().V1().Pods()))
			reflector = rfl.NewNamespaced(options.NewNamespaced().
				WithLocal(LocalNamespace, client, factory).
				WithRemote(RemoteNamespace, client, factory).
				WithHandlerFactory(FakeEventHandler))

			factory.Start(ctx.Done())
			factory.WaitForCacheSync(ctx.Done())
		})

		Context("address translation", func() {
			var (
				input, output string
				podinfo       workload.PodInfo
				err           error
			)

			BeforeEach(func() { input = "192.168.0.25" })

			When("translating a remote to a local address", func() {
				JustBeforeEach(func() {
					output, err = reflector.(*workload.NamespacedPodReflector).MapPodIP(ctx, &podinfo, input)
				})

				It("should succeed", func() { Expect(err).ToNot(HaveOccurred()) })
				It("should return the correct translations", func() { Expect(output).To(BeIdenticalTo("192.168.201.25")) })

				When("translating again the same set of IP addresses", func() {
					JustBeforeEach(func() {
						output, err = reflector.(*workload.NamespacedPodReflector).MapPodIP(ctx, &podinfo, input)
					})

					// The IPAMClient is configured to return an error if the same translation is requested twice.
					It("should succeed (i.e., use the cached values)", func() { Expect(err).ToNot(HaveOccurred()) })
					It("should return the same translations", func() { Expect(output).To(BeIdenticalTo("192.168.201.25")) })
				})
			})
		})

		Context("pod restarts inference", func() {
			Status := func(name string, restarts int32) *corev1.PodStatus {
				return &corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: name, RestartCount: restarts}}}
			}

			DescribeTable("the InferAdditionalRestarts function",
				func(local, remote *corev1.PodStatus, expected int) {
					Expect(reflector.(*workload.NamespacedPodReflector).InferAdditionalRestarts(local, remote)).
						To(BeNumerically("==", expected))
				},
				Entry("when the local status is not yet configured", &corev1.PodStatus{}, &corev1.PodStatus{}, 0),
				Entry("when the local restarts are higher than the remote ones", Status("foo", 5), Status("foo", 3), 2),
				Entry("when the local restarts are equal to the remote ones", Status("foo", 5), Status("foo", 5), 0),
				Entry("when the local restarts are lower than the remote ones", Status("foo", 3), Status("foo", 5), 0),
				Entry("when the container names do not match", Status("foo", 5), Status("bar", 1), 0),
			)
		})
	})
})
