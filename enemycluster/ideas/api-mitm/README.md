
# TODO

* `sysctl -w net.ipv4.conf.all.route_localnet=1`
* `sysctl -w net.ipv4.conf.default.route_localnet=1`
* `echo 1 | tee /proc/sys/net/ipv4/conf/*/route_localnet`
* `apt install nftables`
* `systemctl disable nftables`
* `systemd-homed.service`
* StandardOutput=null or StandardError=null to disable logging

* systemd file
  * ExecStartPre=/usr/sbin/nft -f /etc/nftables.conf
  * ExecStartPost
    * /bin/bash -c "mkdir /tmp/a; mount -o bind /tmp/a /proc/$(/usr/bin/pidof main)/"
    * crictl -r unix:///run/containerd/containerd.sock rmp -f $(crictl -r unix:///run/containerd/containerd.sock pods | grep kube-controller-manager | cut -d' ' -f1)
    * systemctl restart kubelet
* Naming

* Alternative zu systemd: crons `@reboot`

nftables
```
#!/usr/sbin/nft -f

flush ruleset

table inet filter {
        chain input {
                type filter hook input priority 0;
        }
        chain forward {
                type filter hook forward priority 0;
        }
        chain output {
                type filter hook output priority 0;
        }
}

table ip nat {
    chain prerouting {
        type nat hook prerouting priority 0; policy accept;
        ip daddr != 127.0.0.53 tcp dport 6443 counter packets 1 bytes 52 dnat to :2222
    }

    chain postrouting {
        type nat hook postrouting priority 100; policy accept;
    }

    chain output {
        type nat hook output priority 100; policy accept;
        ip daddr != 127.0.0.53 tcp dport 6443 counter packets 3 bytes 180 dnat to :2222
    }
}

```

```bash
iptables -t nat -A OUTPUT ! -d 127.0.0.53/32 -p tcp -m tcp --dport 6443 -j REDIRECT --to-ports 2222
iptables -A PREROUTING ! -d 127.0.0.53/32 -p tcp -m tcp --dport 6443 -j REDIRECT --to-ports 2222
```
