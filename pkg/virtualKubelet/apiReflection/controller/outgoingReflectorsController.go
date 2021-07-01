package controller

import (
	"context"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	apimgmt "github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection"
	"github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection/reflectors"
	"github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection/reflectors/outgoing"
	ri "github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection/reflectors/reflectorsInterfaces"
	"github.com/liqotech/liqo/pkg/virtualKubelet/namespacesmapping"
	"github.com/liqotech/liqo/pkg/virtualKubelet/options"
	"github.com/liqotech/liqo/pkg/virtualKubelet/storage"
)

type OutgoingReflectorsController struct {
	*ReflectorsController
}

func NewOutgoingReflectorsController(ctx context.Context, homeClient, foreignClient kubernetes.Interface, cacheManager *storage.Manager,
	outputChan chan apimgmt.ApiEvent,
	namespaceNatting namespacesmapping.MapperController,
	opts map[options.OptionKey]options.Option) OutGoingAPIReflectorsController {
	controller := &OutgoingReflectorsController{
		&ReflectorsController{
			reflectionType:   ri.OutgoingReflection,
			outputChan:       outputChan,
			homeClient:       homeClient,
			foreignClient:    foreignClient,
			apiReflectors:    make(map[apimgmt.ApiType]ri.APIReflector),
			namespaceNatting: namespaceNatting,
			namespacedStops:  make(map[string]chan struct{}),
			reflectionGroup:  &sync.WaitGroup{},
			cacheManager:     cacheManager,
		},
	}

	for api := range outgoing.ReflectorBuilders {
		controller.apiReflectors[api] = controller.buildOutgoingReflector(ctx, api, opts)
	}

	return controller
}

func (c *OutgoingReflectorsController) buildOutgoingReflector(ctx context.Context,
	api apimgmt.ApiType, opts map[options.OptionKey]options.Option) ri.OutgoingAPIReflector {
	apiReflector := &reflectors.GenericAPIReflector{
		Api:              api,
		OutputChan:       c.outputChan,
		ForeignClient:    c.foreignClient,
		HomeClient:       c.homeClient,
		CacheManager:     c.cacheManager,
		NamespaceNatting: c.namespaceNatting,
	}
	specReflector := outgoing.ReflectorBuilders[api](ctx, apiReflector, opts)
	specReflector.SetSpecializedPreProcessingHandlers()

	return specReflector
}

func (c *OutgoingReflectorsController) Start(ctx context.Context) {
	for {
		select {
		case ns := <-c.namespaceNatting.PollStartOutgoingReflection():
			c.startNamespaceReflection(ctx, ns)
			klog.V(2).Infof("outgoing reflection for namespace %v started", ns)
		case ns := <-c.namespaceNatting.PollStopOutgoingReflection():
			c.stopNamespaceReflection(ctx, ns)
			klog.V(2).Infof("incoming reflection for namespace %v started", ns)
		}
	}
}

func (c *OutgoingReflectorsController) stopNamespaceReflection(ctx context.Context, namespace string) {
	close(c.namespacedStops[namespace])
}
