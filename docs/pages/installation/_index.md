---
title: Installation 
weight: 2
---

### Pre-Installation

Liqo can be used with different topologies and scenarios. This impacts several installation parameters you will configure (e.g., API Server, Authentication).
Before installing Liqo, you should:
* Provision the clusters you would like to use with Liqo. If you need some advice about how to provision clusters on major providers, we have provided [here](./platforms/) some tips.
* Have a look to the [scenarios page](Advanced/_index.md) presents some common patterns used to expose and interconnect clusters.

### Quick Install

#### Pre-Requirements

To install Liqo, you have to install the liqoctl:

```bash
#(Draft)
OS=linux # possible values: linux,windows,darwin
ARCH=amd64 # possible values: amd64,arm64 
curl -LS https://get.liqo.io/liqoctl-${OS}-${ARCH} -output-file liqoctl && chmod +x liqoctl && sudo cp liqoctl /usr/bin/liqoctl
```

{{% notice note %}}
Liqo only supports Kubernetes >= 1.19.0.
{{% /notice %}}

According to your cluster provider, you may have to perform simple steps before triggering the installation process:

{{< tabs >}}
{{% tab name="K8s (Kubeadm)" %}}

**Optional**: You only have to export the KUBECONFIG environment variable. 
Otherwise, liqoctl will try to use the kubeconfig in kubectl default path (i.e. `${HOME}/.kube/config` )

```bash
export KUBECONFIG=/your/kubeconfig/path
```

{{% /tab %}}

{{% tab name="EKS" %}}

To install Liqo on EKS, you should login using the AWS cli (if you already did that, you can skip this step)
This is widely documented on the [official CLI documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html).

In a nutshell, after having installed the CLI, you just have to run:
```bash
aws configure
```

{{% notice note %}}
To install Liqo on your cluster, the AWS user you configure must have the following permissions: *iam:CreatePolicy*,*iam:AttachUserPolicy*,*iam:CreateUser*,*vpc:DescribeVPC*
{{% /notice %}}

Second, you should export the cluster's kubeconfig if you have not already. You may use the following CLI command:

{{% notice note %}}
To run the following command, you must have permission to use the `eks:DescribeCluster` API action with the cluster you specify.
{{% /notice %}}

```bash
aws eks --region ${EKS_CLUSTER_REGION} update-kubeconfig --name ${EKS_CLUSTER_NAME}
```
{{% /tab %}}

{{% tab name="AKS" %}}
First, you should have the AZ cli installed and your AKS cluster deployed. If you haven't, you can follow the [official guide](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli).

Second, you should log-in:
```bash
az login
```

You also need read-only permissions on AKS cluster and on the Virtual Networks, if your cluster has an Azure CNI.

{{% /tab %}}
{{% tab name="GKE" %}}
To install Liqo on GKE, you should at first create a service account for liqoctl, granting the read rights for the GKE clusters (you may reduce the scope to a specific cluster if you prefer):

```bash
gcloud iam service-accounts create ${SERVICE_ACCOUNT_ID} \
    --description="DESCRIPTION" \
    --display-name="DISPLAY_NAME"
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:SERVICE_ACCOUNT_ID@PROJECT_ID.iam.gserviceaccount.com" \
    --role="ROLE_NAME"
```

Then, you should create and download a service accounts key, as presented [by the official documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys#creating_service_account_keys):
```bash
gcloud iam service-accounts keys create ${SERVICE_ACCOUNT_PATH} \
    --iam-account=s${sa-name}@${project-id}.iam.gserviceaccount.com
```

{{% /tab %}}
{{% tab name="K3s" %}}
**Optional**: You only have to export the KUBECONFIG environment variable.
Otherwise, liqoctl will try to use the kubeconfig in kubectl default path (i.e. `${HOME}/.kube/config` )

```bash
export KUBECONFIG=/your/kubeconfig/path
```
{{% /tab %}}
{{< /tabs >}}

#### Install

{{< tabs >}}
{{% tab name="K8s (Kubeadm)" %}}
```bash
liqoctl install --provider kubeadm
```
{{% /tab %}}
{{% tab name="EKS" %}}

{{% notice note %}}
To run the following command, you must have permission to use the `eks:DescribeCluster` API action with the cluster you specify.
{{% /notice %}}

```bash
liqoctl install --provider eks --eks.region=${EKS_CLUSTER_REGION} --eks.cluster-name=${EKS_CLUSTER_NAME} 
```
{{% /tab %}}
{{% tab name="AKS" %}}
```bash
liqoctl install --provider aks --aks.resource-group-name ${AZURE_RESOURCE_GROUP} --aks.resource-name ${AZURE_RESOURCE_NAME} --aks.subscription-id ${AZURE_SUBSCRIPTION_ID}"
```
{{% /tab %}}
{{% tab name="GKE" %}}
```bash
liqoctl install --provider gke --gke.project-id=${GKE_PROJECT_ID} --gke.cluster-id=${GKE_CLUSTER_ID} --gke.zone=${GKE_CLUSTER_ZONE} --gke.credentials-path=${SERVICE_ACCOUNT_PATH}
```
{{% /tab %}}
{{% tab name="K3s" %}}
```bash
liqoctl install --provider kubeadm
```
{{% /tab %}}
{{< /tabs >}}

#### Next Steps

After you have successfully installed Liqo, you may:

* [Configure](/user/configure): configure the Liqo security, the automatic discovery of new clusters and other system parameters.
* [Use](/user/use) Liqo: orchestrate your applications across multiple clusters.
