---
title: Peer to a foreign cluster
weight: 3
---

Once Liqo is installed in your cluster, you can start establishing new *peerings*. This tutorial relies on [LAN Discovery](/configuration/discovery#lan-discovery) since our Kind clusters are in the same L2 broadcast domain.

## LAN Discovery

Liqo can automatically discover any available clusters or make any clusters discoverable on the same LAN.

Using kubectl, you can also manually obtain the list of discovered foreign clusters:

```bash
kubectl get foreignclusters.discovery.liqo.io
NAME                                   OUTGOING PEERING PHASE   INCOMING PEERING PHASE   NETWORKING STATUS   AUTHENTICATION STATUS   AGE
64c6d802-a552-446e-a597-a6f4df9c937b   Established              Established              Established         Established             104s
```

The `foreigncluster` object is used by Liqo to model the discovered remote clusters, here we can check
the status of the interconnection betwwen the two clusters:

* The `Outgoing Peering Phase` shows if the outgoing peering is enabled. If th outgoing peering is enabled, the home cluster can offload pods to the remote one.
* The `Incoming Peering Phase` shows if the outgoing peering is enabled. If the incoming peering is enabled, the home cluster can receive offloaded pods from the foreign cluster
* The `Networking Status` shows if the network interconnection between the two clusters was established.
* The `Authentication Status` shows if the home and the foreign cluster have completed the mutual authentication.

## Peering checking

### Presence of the virtual-node

If the peering has succeded, you should see a virtual node (named `liqo-*`) in addition to your physical nodes:

```
kubectl get nodes

NAME                                      STATUS   ROLES
master-node                               Ready    master
worker-node-1                             Ready    <none>
worker-node-2                             Ready    <none>
liqo-64c6d802-a552-446e-a597-a6f4df9c937b READY    agent    <-- This is the virtual node
```

## Verify that the resulting infrastructure works correctly

You are now ready to verify that the resulting infrastructure works correctly, as presented in the [next step](../test).
