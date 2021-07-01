package outgoing

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"k8s.io/klog"

	liqoconst "github.com/liqotech/liqo/pkg/consts"
	"github.com/liqotech/liqo/pkg/liqonet"
	apimgmt "github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection"
	ri "github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection/reflectors/reflectorsInterfaces"
	"github.com/liqotech/liqo/pkg/virtualKubelet/options"
	"github.com/liqotech/liqo/pkg/virtualKubelet/options/types"
)

// ReflectorBuilders is a map for outgoing reflectors collection. It indexes reflectorBuilders by resource apimgmt resource types.
var ReflectorBuilders = map[apimgmt.ApiType]func(ctx context.Context, reflector ri.APIReflector,
	opts map[options.OptionKey]options.Option) ri.OutgoingAPIReflector{
	apimgmt.Configmaps:     configmapsReflectorBuilder,
	apimgmt.EndpointSlices: endpointslicesReflectorBuilder,
	apimgmt.Secrets:        secretsReflectorBuilder,
	apimgmt.Services:       servicesReflectorBuilder,
}

func configmapsReflectorBuilder(ctx context.Context, reflector ri.APIReflector, _ map[options.OptionKey]options.Option) ri.OutgoingAPIReflector {
	return &ConfigmapsReflector{APIReflector: reflector}
}

func endpointslicesReflectorBuilder(ctx context.Context, reflector ri.APIReflector,
	opts map[options.OptionKey]options.Option) ri.OutgoingAPIReflector {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", opts[options.OptionKey(types.LiqoIpamServer)].Value(), liqoconst.NetworkManagerIpamPort),
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		klog.Error(err)
	}
	ipamClient := liqonet.NewIpamClient(conn)

	return &EndpointSlicesReflector{
		APIReflector:    reflector,
		VirtualNodeName: opts[types.VirtualNodeName],
		IpamClient:      ipamClient,
	}
}

func secretsReflectorBuilder(ctx context.Context, reflector ri.APIReflector, _ map[options.OptionKey]options.Option) ri.OutgoingAPIReflector {
	return &SecretsReflector{APIReflector: reflector}
}

func servicesReflectorBuilder(ctx context.Context, reflector ri.APIReflector, _ map[options.OptionKey]options.Option) ri.OutgoingAPIReflector {
	return &ServicesReflector{APIReflector: reflector}
}
