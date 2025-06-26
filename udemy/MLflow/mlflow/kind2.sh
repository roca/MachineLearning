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
    echo "Usage: $0 [-k <k8s_version, v1.30.0>] [-n <cluster_name, kind>] [-s <start_ip for LB service type range, 200>] [-e <end_ip for LB service type range, 250>] [-i <image repo>] [-h <hostpath to be mounted on nodes>] [-m <mount hostpath, bool>] "
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
        m ) MOUNTHOSTPATH=${OPTARG} ;;
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

NETWORK_CONFIG=$(cat <<EOF
networking:
  disableDefaultCNI: false       # install kindnet
  kubeProxyMode: iptables        # run kube-proxy
EOF
)

# Create the cluster
kind create cluster --name ${NAME} --image kindest/node:${K8SVERSION} --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
${NETWORK_CONFIG}
nodes:
${NODE_CONFIG}
EOF


