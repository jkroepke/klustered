#!/usr/bin/env bash

KUBE_VERSION=1.24.3
SHELL_COMPLETION='source <(kubectl completion bash)'
SHELL_ALIAS='alias k=kubectl'
QUARANTINE_FOLDER=/tmp/quarantine

set -euo pipefail
if [ "$EUID" -ne 0 ]
  then echo "💔 Please run as root 💔"
  exit
fi

mkdir -p "${QUARANTINE_FOLDER}"
cd ~

function log.SUCCESS () {
  echo -e "\e[32m✓ ${1} \e[0m"
}

function log.ERROR () {
  echo -e "\e[31m✝ ${1} \e[0m"
}

function log.SKIP () {
  echo -e "\e[33m✗ ${1} \e[0m"
}

if ! [ -f ~/.vimrc ]; then
  mkdir -p ~/.vim/colors
  curl -s https://raw.githubusercontent.com/morhetz/gruvbox/master/colors/gruvbox.vim > ~/.vim/colors/gruvbox.vim
  cat << EOF > ~/.vimrc
set bg=dark
set nu
set ai et cuc cul sw=2 ts=2
colo gruvbox
filetype plugin indent on
set list
EOF
  log.SUCCESS "vimrc setup done"
else
  log.SKIP "skip vimrc"
fi

pushd /usr/bin &> /dev/null
curl -sLO "https://dl.k8s.io/v${KUBE_VERSION}/bin/linux/amd64/kubectl.sha256"
if ! echo "$(cat kubectl.sha256) kubectl" | sha256sum --check; then
  log.ERROR "found a tampered kubectl; moving to quarantine"
  mv kubectl "${QUARANTINE_FOLDER}"/kubectl

  popd &> /dev/null
  curl -sLO "https://dl.k8s.io/release/v${KUBE_VERSION}/bin/linux/amd64/kubectl"
  curl -sLO "https://dl.k8s.io/v${KUBE_VERSION}/bin/linux/amd64/kubectl.sha256"
  echo "$(cat kubectl.sha256) kubectl" | sha256sum --check;
  install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
  log.SUCCESS "new kubectl now installed in /usr/local/bin"
else
  log.SUCCESS "kubectl was not tampered with"
fi
log.SUCCESS "kubectl setup done"

if ! grep -q "${SHELL_COMPLETION}" ~/.bashrc; then
  echo "source /usr/share/bash-completion/bash_completion" >>~/.bashrc
  echo "${SHELL_COMPLETION}" >>~/.bashrc
  log.SUCCESS "setup shell completion"
else
  log.SKIP "skip shell completion"
fi

if ! grep -q "${SHELL_ALIAS}" ~/.bashrc; then
  echo "${SHELL_ALIAS}" >>~/.bashrc
  echo 'complete -o default -F __start_kubectl k' >>~/.bashrc
  echo 'export EDITOR=vim' >> ~/.bashrc
  echo 'export KUBE_EDITOR=vim' >> ~/.bashrc
  echo 'export KUBECONFIG=/etc/kubernetes/admin.conf' >> ~/.bashrc
  log.SUCCESS "setup shell alias"
else
  log.SKIP "skip shell alias"
fi

if ! command -v ~/.krew/bin/kubectl-krew &> /dev/null; then
  cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
  echo "export PATH=\"${KREW_ROOT:-$HOME/.krew}/bin:$PATH\"" >> ~/.bashrc
  export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"

  kubectl krew install lineage

  log.SUCCESS "setup krew"
else
  log.SKIP "skipping krew"
fi

pushd $(dirname $(which kubelet)) &> /dev/null
curl -sLO "https://dl.k8s.io/v${KUBE_VERSION}/bin/linux/amd64/kubelet.sha256"
if ! echo "$(cat kubelet.sha256) kubelet" | sha256sum --check; then
  log.ERROR "found a tampered kubelet; handle with care; Download new one:"
  cp kubelet "${QUARANTINE_FOLDER}/kubelet"
  echo "  curl -sLO \"https://dl.k8s.io/release/v${KUBE_VERSION}/bin/linux/amd64/kubectl\""
  echo '  echo "$(cat kubelet.sha256) kubelet" | sha256sum --check;'
fi

log.SUCCESS "Please review the PATH and reload shell"
grep PATH ~/.bashrc
