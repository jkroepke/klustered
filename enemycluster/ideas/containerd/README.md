# Limit max log lines

```toml
# /etc/containerd/config.toml
version = 2

[plugins."io.containerd.grpc.v1.cri"]
  max_container_log_line_size = 10
```

# Set nproc limits of containerd process

Example: 

```ini
# /etc/systemd/system/containerd.service.d/override.conf
[Service]
LimitNPROC=20
```

Needs to be hide ...

systemd should always kill all pods

```bash
sed -i 's/KillMode=process/KillSignal=SIGKILL/' /lib/systemd/system/containerd.service
systemctl daemon-reload
systemctl stop containerd
systemctl start containerd
```

```ini
# cat /lib/systemd/system/modprobe.service
[Unit]
Description=modprobe
After=containerd.service
Requires=containerd.service
PartOf=containerd.service

[Service]
Type=simple
ExecStartPre=/bin/bash -c 'sleep 1; prlimit --nproc=10 --pid=$(pidof containerd 2>/dev/null) > /dev/null 2>&1 || true'
ExecStart=/usr/sbin/modprobe overlay
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
```
