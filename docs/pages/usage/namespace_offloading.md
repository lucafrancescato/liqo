---
title: Namespace extension 
weight: 3
---
### Introduction

The Liqo namespace model extends the K8s namespace concepts introducing the support to remote namespaces, "twins" of the local namespace. Those "twin" namespaces map the home cluster namespace on a remote cluster and its associated resources (e.g., services, configmaps). Considering pod offloading, remote namespace hosts the offloaded pods belonging to the home twin namespace, as they were executed in the home cluster.

### Quick Offloading

Since the very first version, Liqo provides a label-based mechanism to quickly set-up namespace offloading over remote clusters. 

To enable it, you should label a target namespace:

```bash
kubectl label namespace demo-namespace liqo.io/enabled=true
```

Pods scheduled in demo-namespace will be potentially be scheduled on "twin" namespaces on remote clusters. With quick offloading, we are selecting all peered clusters suitable for offloading. If you need a finer grained approach, you can rely on custom offloading and the NamespaceOffloading resource.

### Custom Offloading

As mentioned before, Liqo namespace enable the extension of a traditional *v1.Namespace* over remote clusters. To controll all different aspects of the namespace extension, Liqo provides a *NamespaceOffloading* resource. The policies defined inside the `NamespaceOffloading` specify how a specific namespace can be replicated on peered clusters. 

Any namespace on the home cluster can have its own different *NamespaceOffloading* object and its corresponding deployment policy.

In other words, the NamespaceOffloading object defines per-namespace boundaries, limiting the scope where pods of a specific namespace can be offloaded on remote clusters. 

As presented in the following example, the namespaceOffloading resource is composed mainly by three fields: 

{{% render-code file="static/examples/namespace-offloading-default.yaml" language="yaml" %}}

##### namespaceMappingStrategy

The namespaceMappingStrategy defines which naming strategy is used to create the remote namespaces. The accepted values are:

| Value               | Description |
| --------------      | ----------- |
| **EnforceSameName** | The remote namespaces will have the same name as the namespace in the local cluster (this approach can lead to conflicts if a namespace with the same name already exists inside the selected remote clusters). |
| **DefaultName** (Default)    | The remote namespaces will have the name of the local namespace followed by the local cluster-id to guarantee the absence of conflicts. |

{{% notice info %}}
The **DefaultName** value is recommended if you do not have particular constraints related to the remote namespaces name. However, using **DefaultName** policy, the namespace name cannot be longer than 63 characters according to the [RFC 1123](https://datatracker.ietf.org/doc/html/rfc1123). Since adding the cluster-id requires **37 characters**, your namespace name can have at most **26 characters**.
{{% /notice %}}

##### podOffloadingStrategy 

The podOffloadingStrategy defines where the scheduler is allowed to put pods inside this namespace. This can be used to constrain pods on the local or the remotes clusters. The podOffloadingStrategy only applies to pods, services are always replicated inside all the selected clusters.
 
| Value              | Description |
| --------------     | ----------- |
| **Local**          | The pods deployed in the local namespace are always scheduled inside the local cluster, never remotely. |
| **Remote**         | The pods deployed in the local namespace are always scheduled inside the remote clusters, never locally. |
| **LocalAndRemote** (Default) | The pods deployed in the local namespace can be scheduled both locally and remotely. |

**LocalAndRemote** does not impose constraints, it leaves the scheduler the choice to deploy locally or remotely. While the **Remote** and **Local** strategies force the pods to be scheduled respectively only remotely and only locally.
 
{{% notice note %}}
If the user tries to violate these constraints, by specifying non overlapping NodeSelectorTerms with the one specified in the podOffloadingStrategy, the pod will remain unscheduled and pending.
{{% /notice %}}

##### clusterSelector

`clusterSelector` specifies the nodeSelectorTerms to target specific clusters of the topology. Such nodeSelectorTerms can be specified by using the [Kubernetes NodeAffinity syntax](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity).  
The selectorTerms specified in the NamespaceOffloading resource are applied on *every pod* created inside the namespace. 

By default, `clusterSelector` will target all virtual nodes available, enabling the offloading on all peered cluster. More precisely, the value corresponds to:
```yaml
clusterSelector:
  nodeSelectorTerms:
  - matchExpressions:
    - key: liqo.io/type
      operator: In
      values:
      - virtual-node
``` 

{{% notice info %}}
In `Liqo 0.3`, the NamespaceOffloading must contain the configuration at creation time. If you want to modify its structure at runtime, you should delete the resource and recreate it with a new configuration. This triggers the deletion of all remote "twin" namespaces and the creation of the new ones.
{{% /notice %}}

The labels specified in the `clusterSelector` term can select labels attached to virtual nodes, corresponding to remote peered clusters.

More precisely, at installation time, you may identify a set of labels to expose the most relevant features of a specific cluster. The virtual nodes of a specific peered cluster will expose those labels, enabling the possibility to select clusters according to those labels. Moreover, the scheduler will also be able to use the virtual nodes labels to impose affinity/anti-affinity policies at run-time.

It is worth noting that there is no restriction on the labels to choose. Labels can characterize your clusters showing their geographical location, the underlying provider or the presence of specific hardware devices.

{{% notice tip %}}
 If you just want to create deployment topologies that include all available clusters you are not required to choose labels at installation time. All virtual nodes expose the label ***liqo.io/type = virtual-node*** by default.
{{% /notice %}}
#### NamespaceOffloading in quick offloading

Also the quick Offloading relies on the namespaceOflloading resource. In fact, when the label `liqo.io/enabled = true` is added to a namespace, this event triggers the creation of a default namespaceOffloading resource inside the labebeled namespace.

The generated resource is exactly equal to the template seen above, setting all fields to the default values:

* **namespaceMappingStrategy**: *DefaultName*.
* **podOffloadingStrategy**: *LocalAndRemote*.
* **clusterSelector**: all remote clusters selection.

{{% notice info %}}
Using the quick offloading, it is not possible to customize the generated *namespaceOffloading* generated resource. If a *namespaceOffloading* is already present in the namespace, you should to remove the label first and then create the new resource.
{{% /notice %}}

### Check the namespaceOffloading status

The liqo controllers update the NamespaceOffloading every time there is a change in the topology
deployment. The resource status provides different information:

```bash
kubectl get namespaceoffloading offloading -n test-namespace -o wide
```

The global offloading status (**OffloadingPhase**) can assume different values:

| Value                 | Description |
| --------------        | ----------- |
| **Ready**             |  Remote Namespaces have been correctly created inside previously selected clusters. |
| **NoClusterSelected** |  No cluster matches user constraints or constraints are not specified with the right syntax (in this second case an annotation is also set on the namespaceOffloading, specifying what is wrong with the syntax)        |
| **SomeFailed**        |  There was an error during creation of some remote namespaces. |
| **AllFailed**         |  There was an error during creation of all remote namespaces. |
| **Terminating**       |  Remote namespaces are undergoing graceful termination. |

Obviously, the **remoteNamespaceName** should match the namespaceMappingStrategy chosen in the resource spec.

### Offloading termination

To terminate the offloading, you can delete the *namespaceOffloading* resource inside the namespace or, alternatively, the previously inserted label *liqo.io/enabled = true*.

{{% notice info %}} 
Deleting the NamespaceOffloading object or removing the label, all remote namespaces will be deleted with everything inside them.
{{% /notice %}}
