#!/usr/bin/env bash

KUBE_VERSION=1.23.5
SHELL_COMPLETION='source <(kubectl completion bash)'
SHELL_ALIAS='alias k=kubectl'
QUARANTINE_FOLDER=/tmp/quarantine

set -euo pipefail
if [ "$EUID" -ne 0 ]
  then echo "💔 Please run as root 💔"
  exit
fi

cd ~

echo ". Setup kubectl"
pushd $(dirname $(which kubectl)) &> /dev/null
curl -sLO "https://dl.k8s.io/v${KUBE_VERSION}/bin/linux/amd64/kubectl.sha256"
if ! echo "$(cat kubectl.sha256) kubectl" | sha256sum --check; then
  echo "✝ found a tampered kubectl; moving to quarantine"
  mkdir -p "${QUARANTINE_FOLDER}"
  mv $(which kubectl) "${QUARANTINE_FOLDER}"/kubectl

  popd &> /dev/null
  curl -sLO "https://dl.k8s.io/release/v${KUBE_VERSION}/bin/linux/amd64/kubectl"
  curl -sLO "https://dl.k8s.io/v${KUBE_VERSION}/bin/linux/amd64/kubectl.sha256"
  echo "$(cat kubectl.sha256) kubectl" | sha256sum --check;
  install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
  echo "✓ new kubectl now installed in /usr/local/bin"
else
  echo "✓ kubectl was not tampered with"
fi
echo "✓ kubectl setup done"

exit

if ! grep -q "${SHELL_COMPLETION}" ~/.bashrc; then
  echo "source /usr/share/bash-completion/bash_completion" >>~/.bashrc
  echo "${SHELL_COMPLETION}" >>~/.bashrc
  echo "✓ setup shell completion"
else
  echo "✗ skip shell completion"
fi

if ! grep -q "${SHELL_ALIAS}" ~/.bashrc; then
  echo "${SHELL_ALIAS}" >>~/.bashrc
  echo 'complete -o default -F __start_kubectl k' >>~/.bashrc
  echo "✓ setup shell alias"
else 
  echo "✗ skip shell alias"
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

  echo "✓ setup krew"
else
  echo "✗ skipping krew"
fi

kubectl krew install lineage
echo "✓ setup kubectl plugins"

echo "Please review the PATH and reload shell"
grep PATH ~/.bashrc
