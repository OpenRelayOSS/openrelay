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

HOME_PATH=${HOME}

cd ${REPO_ROOT_PATH}

ORIGIN_SOURCES_PATH=${REPO_ROOT_PATH}/build/rpms/SOURCES/${IMAGE_FULLNAME}.tar.gz
ORIGIN_SPECS_PATH=${REPO_ROOT_PATH}/build/rpms/SPECS/${IMAGE_FULLNAME}.spec
ORIGIN_RPMS_PATH=${REPO_ROOT_PATH}/build/rpms/RPMS/${IMAGE_ARCH}

RPMBUILD=rpmbuild
RPMBUILD_ROOT_PATH=${HOME}/rpmbuild/
RPMBUILD_SOURCES_PATH=${HOME}/rpmbuild/SOURCES/
RPMBUILD_SPECS_PATH=${HOME}/rpmbuild/SPECS/${IMAGE_FULLNAME}.spec
RPMBUILD_RPMS_PATH=${HOME}/rpmbuild/RPMS/${IMAGE_ARCH}/

cp ${ORIGIN_SOURCES_PATH} ${RPMBUILD_SOURCES_PATH} || die "${ORIGIN_SPECS_PATH} ${RPMBUILD_SOURCES_PATH} copy failed"
cp ${ORIGIN_SPECS_PATH} ${RPMBUILD_SPECS_PATH} || die "${ORIGIN_SPECS_PATH} ${RPMBUILD_SPECS_PATH} copy failed"
${RPMBUILD} -bb --clean ${RPMBUILD_SPECS_PATH} || die "rpmbuild failed"

mkdir -p ${ORIGIN_RPMS_PATH}/ || die "mkdir failed"
cp ${RPMBUILD_RPMS_PATH}/${IMAGE_FULLNAME}.rpm ${ORIGIN_RPMS_PATH}/ || die "copy failed"

cd ${RET_DIR}
