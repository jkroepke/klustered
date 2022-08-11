FROM busybox AS busybox

ADD config.yaml /etc/kubernetes/kube-scheduler.conf
ADD kube-scheduler /usr/local/bin/kube-scheduler

COPY --from=k8s.gcr.io/kube-scheduler:v1.23.3 /usr/local/bin/kube-scheduler /usr/bin/kube-scheduler

FROM k8s.gcr.io/kube-scheduler:v1.23.3 AS kube-scheduler

COPY --from=busybox  / /

FROM scratch

COPY --from=kube-scheduler  / /
