* fake API server


* corrupt CNI binaries
* corrupt CRI implemation
  * for example, decrease cgroup default limits (java gives some funny error messags here)

* corrupt kube-root-ca
* corrupt sa.key/sa.pub
* change default scheduler (hard variant: through re-compile)
* manulle lock leases mit hohem expire (wenn das geht)
* probes werden vom kubelet ausgeführt, ggf. network probes verhindern. (aber easy zu entdecken)
* kubelet certificates expired
* Webhooks mit hohem delay und timeouts
* using wrong arch (sx390) image for one etcd? pod (hard variant: local retag to have consitent image digest)
  * in combination with logs are dumped to /dev/null

* running iptables inside pods (different netnamespace) (restart of pod would solve this)
* change location of static pods, but the original location is preseved.