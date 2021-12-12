#!/bin/sh -e
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

${REPO_ROOT_PATH}/scripts/build_extlib.sh || die "build extlib failed."

${REPO_ROOT_PATH}/scripts/build_img.sh || die "build image failed."

${PODMAN_COMPOSE} ${REPO_ROOT_PATH}/deployments/docker-compose.yml build || die "docker-compose failed."

cd -
