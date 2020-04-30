#!/usr/bin/env bash

GOZIP="go1.14.1.linux-amd64.tar.gz"
DLDIR="/vagrant/install/"

GOURL="https://dl.google.com/go/${GOZIP}"

cd $DLDIR

if [ ! -f "${GOZIP}" ]; then
  echo "[+] Downloading ${GOZIP} from ${GOURL} to ${DLDIR}"
  wget -O "${GOZIP}" $GOURL
fi

if [ ! -f "/usr/local/go/bin/go" ]; then
  echo "[+] Extracting go to /usr/local"
  tar -C /usr/local -xzf "${GOZIP}"
fi

grep -q "/usr/local/go/bin" /etc/profile || echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
grep -q "/usr/local/go/bin" /root/.profile || echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.profile
