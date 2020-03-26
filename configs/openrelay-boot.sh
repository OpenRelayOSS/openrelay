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

ENV_FILE=/etc/sysconfig/openrelay.env
source $ENV_FILE

if [ "x${PERFORMANCE_MODE}" = "x0" ];then
IMAGE_NAME=openrelay
else
IMAGE_NAME=openrelay
fi

export LD_LIBRARY_PATH=/usr/local/openrelay/lib
# -recmode ${REC_MODE} \
# -repmode ${REP_MODE} \
${DRYRUN} ${IMAGE_PATH}/${IMAGE_NAME} \
-standbymode ${STANDBYMODE} \
-log ${LOG_LEVEL} \
-logdir ${LOG_DIRECTORY} \
-hbtimeout ${HEATBEAT_TIMEOUT} \
-jointimeout ${JOIN_TIMEOUT} \
-listenmode ${LISTEN_MODE} \
-listen_ipv4 ${LISTEN_IPV4} \
-listen_ipv6 ${LISTEN_IPV6} \
-ehost ${ENTRY_LISTEN_ADDR} \
-eport ${ENTRY_PORT} \
-ahost ${ADMIN_LISTEN_ADDR} \
-aport ${ADMIN_PORT} \
-stf_dproto ${STATEFULL_DEAL_PROTOCOL} \
-stf_dhost "${STATEFULL_DEAL_LISTEN_ADDR}" \
-stf_dports ${STATEFULL_DEAL_PORTS} \
-stf_sproto ${STATEFULL_SUBSCRIBE_PROTOCOL} \
-stf_shost "${STATEFULL_SUBSCRIBE_LISTENA_ADDR}" \
-stf_sports ${STATEFULL_SUBSCRIBE_PORTS} \
-usestl ${USE_STATELESS}

