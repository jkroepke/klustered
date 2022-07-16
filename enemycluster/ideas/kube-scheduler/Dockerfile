FROM busybox AS busybox

ADD config.yaml /etc/kubernetes/kube-scheduler.conf
ADD kube-scheduler /usr/local/sbin/

FROM k8s.gcr.io/kube-scheduler:v1.24.3

COPY --from=busybox / /
