#!/bin/bash

# Usage: ./kind_cilium.sh -k <k8s_version> -n <cluster_name> -s <start_ip> -e <end_ip> -i <image> -p <preload> -m <mount_hostpath> -c <cilium> -t <metalLB>
# Example: ./kind_cilium.sh -k v1.21.1 -n kind -s 200 -e 250 -p true -m true -c false -t false

set -e

# Default values
DEFAULTK8SVERSION="v1.31.0"
DEFAULTNAME="kind"
CILIUM_VERSION="1.15.6"
DEFAULTIMAGE="kindest/node"
DEFAULTSTARTIP=200
DEFAULTENDIP=250
DEFAULTHOSTPATH="${HOME}/temp/kind"
DEFAULTPRELOAD=false
DEFAULTCILIUM=false
DEFAULTMETALLB=false
DEFAULTMOUNTHOSTPTATH=false

# Initialize variables with default values
K8SVERSION=$DEFAULTK8SVERSION
NAME=$DEFAULTNAME
START_IP=$DEFAULTSTARTIP
END_IP=$DEFAULTENDIP
IMAGE=$DEFAULTIMAGE
HOSTPATH=$DEFAULTHOSTPATH
PRELOAD=$DEFAULTPRELOAD
CILIUM=$DEFAULTCILIUM
METALLB=$DEFAULTMETALLB
MOUNTHOSTPATH=$DEFAULTMOUNTHOSTPATH

# Function to display usage
usage() {
    echo "Usage: $0 [-k <k8s_version, v1.30.0>] [-n <cluster_name, kind>] [-s <start_ip for LB service type range, 200>] [-e <end_ip for LB service type range, 250>] [-i <image repo>] [-h <hostpath to be mounted on nodes>] [-p <preload images, bool>] [-m <mount hostpath, bool>] [-c <cilium, bool>] [-t <metalLB, bool>]"
    exit 1
}

# Parse command-line options
while getopts "k:n:s:e:i:h:p:c:t:m:" opt; do
    case ${opt} in
        k ) K8SVERSION=${OPTARG} ;;
        n ) NAME=${OPTARG} ;;
        s ) START_IP=${OPTARG} ;;
        e ) END_IP=${OPTARG} ;;
        i ) IMAGE=${OPTARG} ;;
        h ) HOSTPATH=${OPTARG} ;;
        p ) PRELOAD=${OPTARG} ;;
        c ) CILIUM=${OPTARG} ;;
        m ) MOUNTHOSTPATH=${OPTARG} ;;
        t ) METALLB=${OPTARG} ;;
        \? ) usage ;;
    esac
done

# Create the cluster configuration based on MOUNTHOSTPATH
if [ "${MOUNTHOSTPATH}" = true ]; then
  NODE_CONFIG=$(cat <<EOF
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 8080
    protocol: TCP
  - containerPort: 443
    hostPort: 6443
    protocol: TCP
- role: worker
  extraMounts:
    - hostPath: ${HOSTPATH}
      containerPath: /host
- role: worker
  extraMounts:
    - hostPath: ${HOSTPATH}
      containerPath: /host
- role: worker
  extraMounts:
    - hostPath: ${HOSTPATH}
      containerPath: /host
EOF
  )
else
  NODE_CONFIG=$(cat <<EOF
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 8080
    protocol: TCP
  - containerPort: 443
    hostPort: 6443
    protocol: TCP
- role: worker
- role: worker
- role: worker
EOF
  )
fi

if [ "${CILIUM}" = true ]; then
  NETWORK_CONFIG=$(cat <<EOF
networking:
  disableDefaultCNI: true        # do not install kindnet
  kubeProxyMode: none            # do not run kube-proxy
EOF
  )
else
  NETWORK_CONFIG=$(cat <<EOF
networking:
  disableDefaultCNI: false       # install kindnet
  kubeProxyMode: iptables        # run kube-proxy
EOF
  )
fi


# Create the cluster
kind create cluster --name ${NAME} --image kindest/node:${K8SVERSION} --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
${NETWORK_CONFIG}
nodes:
${NODE_CONFIG}
EOF

 # kubectl config set clusters.kind-kind.insecure-skip-tls-verify true
 # Kubernetes cluster unreachable: specifying a root certificates file with the insecure flag is not allowed

# Preload images if PRELOAD is true
if [ "${PRELOAD}" = true ]; then
    echo "Pre-loading images..."
    export DOCKER_CLI_HINTS=false
    for image in \
    docker.io/kindest/local-path-provisioner:v20240202-8f1494ea \
    registry.k8s.io/metrics-server/metrics-server:v0.7.1 \
    nicolaka/netshoot \
    quay.io/cilium/cilium:v${CILIUM_VERSION} \
    quay.io/cilium/hubble-relay:v${CILIUM_VERSION} \
    quay.io/cilium/operator-generic:v${CILIUM_VERSION} \
    quay.io/frrouting/frr:9.0.2 \
    quay.io/metallb/speaker:v0.14.5 ; do docker pull $image; kind load docker-image $image -n $NAME ; done >/dev/null 2>&1
fi

if [ "${CILIUM}" = true ]; then
helm upgrade --install --namespace kube-system --version ${CILIUM_VERSION} --repo https://helm.cilium.io cilium cilium --values - <<EOF
kubeProxyReplacement: strict
k8sServiceHost: ${NAME}-control-plane # api server lb, kind-external-load-balancer if multimaster
k8sServicePort: 6443                        # api server port
hostServices:
  enabled: true
externalIPs:
  enabled: true
nodePort:
  enabled: true
hostPort:
  enabled: true
image:
  pullPolicy: IfNotPresent
ipam:
  mode: kubernetes
hubble:
  enabled: trune
  relay:
    enabled: true
EOF

kubectl wait po -n kube-system --timeout=600s -l k8s-app=cilium -l app.kubernetes.io/name=cilium-agent --for condition=Ready

fi


helm repo add --kube-insecure-skip-tls-verify bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install metrics-server bitnami/metrics-server
helm upgrade --insecure-skip-tls-verify  --install --set args={--kubelet-insecure-tls} metrics-server --repo https://kubernetes-sigs.github.io/metrics-server/ metrics-server --namespace kube-system

if [ "${METALLB}" = true ]; then
# MetalLB, modified for podman. Needs https://github.com/jasonmadigan/podman-mac-net-connect
helm upgrade --insecure-skip-tls-verify --install --namespace metallb-system --create-namespace --repo https://metallb.github.io/metallb metallb metallb

kubectl wait po -n metallb-system --timeout=600s -l app.kubernetes.io/name=metallb -l app.kubernetes.io/component=controller --for condition=Ready
sleep 5
#Mac
#KIND_NET_CIDR=$(docker network inspect kind | jq -r '.[0].subnets.[1].subnet'| cut -d'/' -f1)
#Linux, check if Config.[0] or Config.[1] corresponds to the IPv4 Network
KIND_NET_CIDR=$(docker network inspect kind | jq -r '.[0].IPAM.Config.[0].Subnet'|cut -d '/' -f1)

echo $KIND_NET_CIDR
METALLB_IP_START=$(echo ${KIND_NET_CIDR} | sed "s/\(.*\.\)0$/\1$START_IP/")
echo $METALLB_IP_START
METALLB_IP_END=$(echo ${KIND_NET_CIDR} | sed "s/\(.*\.\)0$/\1$END_IP/")
echo $METALLB_IP_END
METALLB_IP_RANGE="${METALLB_IP_START}-${METALLB_IP_END}"
echo $METALLB_IP_RANGE
kubectl apply -f - <<EOF
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: lbpool
  namespace: metallb-system
spec:
  addresses:
  - ${METALLB_IP_RANGE}
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: l2adv
  namespace: metallb-system
EOF
fi
