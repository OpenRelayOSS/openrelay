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

export HOSTIPV4=`ip a | grep inet | grep -v inet6 | grep -v 127.0.0.1 | awk '{print $2}' | head -1 | awk -F'/' '{print $1}'`
export HOSTIPV6=`ip a | grep inet6 | grep -v ::1 | awk '{print $2}' | head -1 | awk -F'/' '{print $1}'`
export OPENRELAY_OPTION="${OPENRELAY_OPTION} -listenmode 0 -endpoint_ipv4 ${HOSTIPV4} -endpoint_ipv6 ${HOSTIPV6}"

${PODMAN_COMPOSE} -f deployments/docker-compose.yml up

cd ${RET_DIR}
