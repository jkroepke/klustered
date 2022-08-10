# Exploits

## Chmod kubectl
kubectl is missing it's execution bit, as well as chmod

### Fix
`/lib64/ld-linux-x86-64.so.2 /usr/bin/chmod +x /usr/bin/chmod`

## Containerd Logs
-> [Idea](./containerd/README.md) <-

### Fix
1. Remove /etc/containerd/config.toml
2. Restart containerd

## Api Mitm
-> [Idea](./api-mitm/README.md) <-

### Fix
1. Remove iptables (`iptables -t nat -D PREROUTING 2`)
2. Remove nftables (`nft flush ruleset`)

## Impersonation
-> [Idea](./impersonation/REAMDE.md) <-

### Fix
`kubeadm kubeconfig user --client-name kubernetes-admin --org system:masters --config <(kubeadm config print init-defaults) > .kube/config`

or

`k --as kubernetes-admins --as-group system:masters`

## Deploymentcontroller

### Fix
1. Remove missing deployment controller from kube-controller-manager

## Scheduler

### Fix
1. Rename scheduler to `nondefault`
or schedule with `nodeName`

