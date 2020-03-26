#!/bin/sh
# Copyright (c) 2018 FurtherSystem Co.,Ltd. All rights reserved.
#
#   This program is free software; you can redistribute it and/or modify
#   it under the terms of the GNU General Public License as published by
#   the Free Software Foundation; version 2 of the License.
#
#   This program is distributed in the hope that it will be useful,
#   but WITHOUT ANY WARRANTY; without even the implied warranty of
#   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#   GNU General Public License for more details.
#
#   You should have received a copy of the GNU General Public License
#   along with this program; if not, write to the Free Software
#   Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1335  USA

source `dirname $0`/common.env

cd ${REPO_ROOT_PATH}

sudo ${DNF} -y install epel-release
sudo ${DNF} -y install tar wget make gcc gcc-c++ libtool automake autoconf git pkgconfig libunwind libunwind-devel rpm-build 

mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

if [ "x${DNF}" = "xdnf" ];then
sudo ${DNF} config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/containerd.io-1.2.2-3.3.el7.x86_64.rpm
sudo ${DNF} -y install --nobest docker-ce docker-ce-cli
sudo ${DNF} -y install containerd.io-1.2.2-3.3.el7.x86_64.rpm
rm -rf containerd.io-1.2.2-3.3.el7.x86_64.rpm
sudo systemctl enable docker
sudo systemctl start docker
sudo curl -L "https://github.com/docker/compose/releases/download/1.25.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo firewall-cmd --add-masquerade --permanent
sudo firewall-cmd --reload
else
sudo ${DNF} -y install docker docker-compose
fi

curl https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz | tar zx -C ~/
echo export GO111MDULE=on >> ~/.bash_profile
echo export GOPATH=~/go >> ~/.bash_profile
echo 'export PATH=${PATH}:${GOPATH}/bin' >> ~/.bash_profile
echo export PATH >> ~/.bash_profile
source ~/.bash_profile

cd ${RET_DIR}
