# IDEAS


| Idea                                                                     | Resolve                                                                                | Ulli | Joe |
|--------------------------------------------------------------------------|----------------------------------------------------------------------------------------|------|-----|
| containerd: Limit log lines                                              | `rm /etc/containerd/config.toml && systemctl restart containerd`                       | [ ]  | [ ] |
| containerd: Limit nprocs                                                 | `systemctl disable --now modprobe && systemctl restart containerd`                     | [x]  | [ ] |
| kube-apiserver: MITM dry-run proxy                                       | `iptables & nft tables flushen` @Ulli: nft binary umbenennen? web3?                    | [x]  | [ ] |
| impersonation                                                            | `kubectl --as-group=system:masters --as=kubernetes-admins` <br> kubeadm geht nicht OOB | [x]  | [ ] |
| kube-scheduler: Change name of default scheduler                         | `ctr --namespace=k8s.io images pull k8s.gcr.io/kube-scheduler:v1.24.3`                 | [x]  | [ ] |
| kube-controller: Disable replicaset or deployment controller though flag | Remove flag from `/etc/kubernetes/manifests/kube-controller-manager.yaml`.             | [x]  | [ ] |
