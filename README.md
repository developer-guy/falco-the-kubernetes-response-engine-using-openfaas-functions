![openfaas_introduction](.res/civo.png)

# Kubernetes Response Engine powered by OpenFaaS

Although Falco can be used to detect any _intrusion_ attempts and sends alerts to channels according to the given rules afterwards, it does not have any _remediation_ system.
This is why we need something called Kubernetes Response Engine. It simply aims to catch alerts and take actions on it. These actions can be designed as _fine-grained_ serverless functions.

FIXME: A pod that has pwned can be indicated as an example.

Table of Contents
=================
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [Overall Architecture](#overall-architecture)
- [Prerequisites](#prerequisites)
- [Target Audience](#target-audience)
- [What is ...?](#what-is-)
  - [CIVO](#civo)
  - [K3S](#k3s)
  - [Falco](#falco)
  - [OpenFaaS](#openfaas)
- [Hands-on Demonstration](#hands-on-demonstration)
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Overall Architecture

```bash
                +-----------+
                |   Falco   +
                +-----^-----+
                      |
              +-------v-------+
              >   OpenFaaS    +
              +-------v-------+
+-----------+         |          +-----------+
| notify-fn <---------+----------> delete-fn |
+-----v-----+ notice     warning +-----+-----+
      |                                |
      | send alert          delete pod |
      |                                |
+-----v-----+                    +-----v-----+       
|   Slack   |                    | Pwned Pod |
+-----------+                    +-----------+
```

# Prerequisites

* <img src="https://www.civo.com/brand-assets/logo/full-colour/civo-logo-fullcolour.svg" height="16" width="16"/> civo-cli v0.7.6
* <img src="https://cncf-branding.netlify.app/img/projects/helm/horizontal/color/helm-horizontal-color.svg" height="16" width="16" /> Helm v3.5.1
* <img src="https://github.com/openfaas/openfaas.github.io/blob/master/images/openfaas/logo.svg" height="16" width="16"/> faas-cli
* <img src="https://raw.githubusercontent.com/cncf/artwork/master/other/illustrations/ashley-mcnamara/kubectl/kubectl.svg" height="16" width="16"/> kubectl v1.20.2

> We are going to do this demo on macOS Catalina 1.15.7, you can find the prerequisites on [brew](https://brew.sh).

# Target Audience

If you want to:

* Set up a K3S cluster on CIVO
* Set up the Falco
* Create custom serverless functions using OpenFaaS
* Subscribe Falco events from serverless functions

# What is ...?

## [CIVO](https://www.civo.com/)

_Built for speed and simplicity, with K3s under the hood_

**Join the [#Kube100](https://www.civo.com/kube100) beta: [Apply to join today](https://www.civo.com/signup) _(Get free credit to test-drive the world’s first K3s-powered, managed Kubernetes service)_**

* Fast, but powerful
> Spin up Kubernetes in under 2 minutes, without the bloat, using the lightweight K3s distribution

* Management CLI
> Manage your clusters with our custom-built CLI. Streamline deploys with our clean REST API.

* Application marketplace
  
> Launch clusters with preinstalled applications, or install on the fly with our Kubernetes marketplace.

## [K3S](https://k3s.io/)

_The certified Kubernetes distribution built for IoT & Edge computing_

* Perfect for Edge
> K3s is a highly available, certified Kubernetes distribution designed for production workloads in unattended, resource-constrained, remote locations or inside IoT appliances.
Simplified & Secure

* Simplified & Secure
> K3s is packaged as a single <40MB binary that reduces the dependencies and steps needed to install, run and auto-update a production Kubernetes cluster.

* Optimized for ARM
> Both ARM64 and ARMv7 are supported with binaries and multiarch images available for both. K3s works great from something as small as a Raspberry Pi to an AWS a1.4xlarge 32GiB server.

## [Falco](https://falco.org/)

_Cloud-Native runtime security, de facto Kubernetes threat detection engine_

* Strengthen container security
> The flexible rules engine allows you to describe any type of host or container behavior or activity.

* Reduce risk via immediate alerts
> You can immediately respond to policy violation alerts and integrate Falco within your response workflows.

* Leverage most current detection rules
> Falco out-of-the box rules alert on malicious activity and CVE exploits.

## [OpenFaaS](https://www.openfaas.com/)

_OpenFaaS® makes it simple to deploy both functions and existing code to Kubernetes_

* Anywhere
> Avoid lock-in through the use of Docker. Run on any public or private cloud.

* Any code
> Build both microservices & functions in any language. Legacy code and binaries.

* Any scale
> Auto-scale for demand or to zero when idle.

# Hands-on Demonstration

## Create CIVO Playground

* Save an API Key on [Security Dashboard](https://www.civo.com/account/security)

![civo_dashboard_settings](.res/civo-settings.png)

```bash
$ civo apikey save my-awesome-key $KEY
Saved the API Key $KEY as my-awesome-key
```

```bash
$ civo apikey list
+----------------+---------+
| Name           | Default |
+----------------+---------+
| my-awesome-key | $KEY    |
+----------------+---------+
```

```bash
$ civo apikey current my-awesome-key
Set the default API Key to be my-awesome-key
```

* Create a cluster

```bash
$ civo kubernetes create playground --size=g3.k3s.medium --nodes=3 --region NYC1 --wait
Created Kubernetes cluster playground
```

* Show the playground on [Kubernetes Dashboard](https://www.civo.com/account/kubernetes)

```bash
$ civo kubernetes show playground
          ID : 79435efe-2dac-403d-bfd2-f6644988830a
        Name : playground
       Nodes : 3
        Size : g3.k3s.medium
      Status : ACTIVE
     Version : 1.20.0-k3s2
API Endpoint : https://212.2.243.151:6443
   Master IP : 212.2.243.151
DNS A record : 79435efe-2dac-403d-bfd2-f6644988830a.k8s.civo.com
Nodes:
+-------------+---------------+--------+---------------+-----------+------+----------+
| Name        | IP            | Status | Size          | Cpu Cores | Ram  | SSD disk |
+-------------+---------------+--------+---------------+-----------+------+----------+
| master-7c8a | 212.2.243.151 | ACTIVE | g3.k3s.medium |         2 | 4096 |       25 |
| node-04ed   |               | ACTIVE | g3.k3s.medium |         2 | 4096 |       25 |
| node-5258   |               | ACTIVE | g3.k3s.medium |         2 | 4096 |       25 |
+-------------+---------------+--------+---------------+-----------+------+----------+
```

![civo_dashboard_settings](.res/civo-clusters.png)


* Configure the playground

```bash
$ civo kubernetes config playground --save --local-path ./kubeconfig
Access your cluster with:
KUBECONFIG=./kubeconfig kubectl get node
```

* Ensure all is OK

```bash
$ KUBECONFIG=./kubeconfig kubectl get node
NAME                                  STATUS   ROLES                  AGE   VERSION
k3s-playground-66b18d51-node-04ed     Ready    <none>                 40h   v1.20.2+k3s1
```

You can find more details about **civo-cli** [here](https://github.com/civo/cli).