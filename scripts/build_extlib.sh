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

EXTLIB_PATH=${REPO_ROOT_PATH}/extlib/
export GO111MDULE=on
export GOPATH=/go
export PATH=${PATH}:${GOPATH}/bin
export PATH

mkdir -p ${EXTLIB_PATH}
cd ${EXTLIB_PATH}
export LD_LIBRARY_PATH=${EXTLIB_PATH}

git clone git://github.com/jedisct1/libsodium.git -b 1.0.18-RELEASE --depth 1 2>/dev/null || die "git clone failed."
git clone git://github.com/zeromq/libzmq.git -b v4.3.2 --depth 1 2>/dev/null || die "git clone failed."
git clone git://github.com/zeromq/czmq.git -b v4.2.0 --depth 1 2>/dev/null || die "git clone failed."
cd libsodium
./autogen.sh && ./configure && make || die "libsodium make failed."
make check || :
sudo make install || die "libsodium make failed."
cp src/libsodium/.libs/libsodium.so.*.* ../
cd -
ln -s libsodium.so.*.* libsodium.so
cd libzmq
./autogen.sh && ./configure && make  || die "libzmq make failed."
make check || :
sudo make install || die "libzmq make failed."
cp src/.libs/libzmq.so.*.* ../
cd -
ln -s libzmq.so.*.* libzmq.so
cd czmq
./autogen.sh && ./configure && make  || die "libczmq make failed."
make check || :
sudo make install || die "libczmq make failed."
cp src/.libs/libczmq.so.*.* ../
cd -
ln -s libczmq.so.*.* libczmq.so

cd ${RET_DIR}
