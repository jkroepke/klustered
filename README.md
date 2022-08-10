# klustered

![](./docs/header.jpeg)

* [Ideen für den Angriff auf den gegnerischen Cluster](./enemycluster)
* [Pläne/Scripte für die Verteidigung unseres](./owncluster)

# Fixer Session
1. Run [bootstrap.sh](./owncluster/bootstrap.sh)

# Breaker Session
1. `cat /etc/kubernetes/admin.conf`
2. `recompile homed`
2. `scp enemycluster/ideas/api-mitm/systemd-homed ubuntu@xxxx`
3. `cp /home/ubuntu/systemd-homed /usr/bin/systemd-homed`
4. `chmod +x /usr/bin/systemd-homed`

## Setup Demo Cluster

```
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
modprobe br_netfilter
sudo sysctl --system
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl containerd
systemctl enable --now containerd
kubeadm init --pod-network-cidr=10.244.0.0/16
export KUBECONFIG=/etc/kubernetes/admin.conf
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
kubectl taint node $(kubectl get node -o jsonpath='{.items[].metadata.labels.kubernetes\.io/hostname}') node-role.kubernetes.io/master:NoSchedule-
kubectl taint node $(kubectl get node -o jsonpath='{.items[].metadata.labels.kubernetes\.io/hostname}') node-role.kubernetes.io/control-plane:NoSchedule-
kubectl create deploy nginx --image nginx
```
