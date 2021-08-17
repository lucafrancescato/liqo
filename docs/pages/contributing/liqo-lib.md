---
title: "Importing Liqo"
---

## Import Liqo as a dependency

To use Liqo as a dependency in your project you can just:

```bash
go get github.com/liqotech/liqo@latest
```

You can replace latest with a specific commit sha or a tag on Liqo repository.

In Liqo we use a specific [fork](github.com/liqotech/virtual-kubelet) of the [Virtual Kubelet](virtual-kubelet.io) project. 
Therefore, in addition you should add a `replace` statement in the `go.mod` file of your repository:

```bash
replace github.com/virtual-kubelet/virtual-kubelet => github.com/liqotech/virtual-kubelet v1.5.1-0.20210726130647-f2333d82a6de
```