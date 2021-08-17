---
title: Peer with a new cluster
weight: 1
---

## Overview

In Liqo, peering establishes an administrative connection between two clusters and enables the resource sharing across them.
It is worth noticing that peering is uni-directional. 
This implies that resources can be shared only from a cluster to another and not the vice-versa. Obviously, it can be optionally be enabled bi-directionally, enabling a two-way resource sharing.

### Peer with a new cluster

#### Generate add command

#### Peer

To peer with a new cluster, you can use *liqoctl*:

```bash
liqoctl add cluster --name=test --token=yzk --url=https://yolo:6433
```

If this command is executed successfully, you have completed a peering. A new node should be available in your cluster:

```bash
kubectl get no
```

#### (Optional): Accept the incoming peering

On the remote cluster, the Liqo control plane will receive a peering request for an incoming foreignCluster.
You can check the status of known foreign clusters with:

```bash 
kubectl get foreignclusters
```

If a incomingPeering is pending, you may decide to accept it by:

```bash
liqoctl accept cluster --name 
```

### Check the peering status

You can in any moment check the healthiness of the peering by running:

```bash
liqoctl check peering-status --name test
```