package incoming

import (
	"context"

	apimgmt "github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection"
	ri "github.com/liqotech/liqo/pkg/virtualKubelet/apiReflection/reflectors/reflectorsInterfaces"
	"github.com/liqotech/liqo/pkg/virtualKubelet/options"
)

// ReflectorBuilder is a map for incoming reflectors collection. It indexes reflectorBuilders by resource apimgmt resource types.
var ReflectorBuilder = map[apimgmt.ApiType]func(ctx context.Context, reflector ri.APIReflector,
	opts map[options.OptionKey]options.Option) ri.IncomingAPIReflector{
	apimgmt.Pods:        podsReflectorBuilder,
	apimgmt.ReplicaSets: replicaSetsReflectorBuilder,
}

func podsReflectorBuilder(ctx context.Context, reflector ri.APIReflector, opts map[options.OptionKey]options.Option) ri.IncomingAPIReflector {
	return &PodsIncomingReflector{
		APIReflector: reflector}
}

func replicaSetsReflectorBuilder(ctx context.Context, reflector ri.APIReflector, _ map[options.OptionKey]options.Option) ri.IncomingAPIReflector {
	return &ReplicaSetsIncomingReflector{
		APIReflector: reflector,
	}
}
