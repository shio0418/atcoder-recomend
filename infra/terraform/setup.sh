#!/bin/bash
set -e

export DEBIAN_FRONTEND=noninteractive

apt update
apt -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" upgrade -y

# Docker導入
apt install -y ca-certificates curl gnupg
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null
apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# ファイアウォール
apt install -y ufw
ufw allow 22
ufw allow 80
ufw allow 443
ufw --force enable

# deployユーザー作成
if ! id deploy &>/dev/null; then
  adduser --disabled-password --gecos "" deploy
  usermod -aG sudo deploy
  usermod -aG docker deploy
  echo "deploy ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/deploy
  chmod 440 /etc/sudoers.d/deploy

  mkdir -p /home/deploy/.ssh
  cp /root/.ssh/authorized_keys /home/deploy/.ssh/
  chown -R deploy:deploy /home/deploy/.ssh
  chmod 700 /home/deploy/.ssh
  chmod 600 /home/deploy/.ssh/authorized_keys
fi

# rootのSSHログイン禁止
sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config
systemctl restart ssh
