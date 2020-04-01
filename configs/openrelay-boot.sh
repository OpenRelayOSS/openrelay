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

for OPT in "$@"
do
    case $OPT in
        -dryrun)
            DRYRUN=$2
            shift 2
            ;;
        -perfmode)
            PERFORMANCE_MODE=$2
            shift 2
            ;;
        -standbymode)
            STANDBYMODE=$2
            shift 2
            ;;
        -recmode)
            REC_MODE=$2
            shift 2
            ;;
        -repmode)
            REP_MODE=$2
            shift 2
            ;;
        -log)
            LOG_LEVEL=$2
            shift 2
            ;;
        -logdir)
            LOG_DIRECTORY=$2
            shift 2
            ;;
        -hbtimeout)
            HEATBEAT_TIMEOUT=$2
            shift 2
            ;;
        -jointimeout)
            JOIN_TIMEOUT=$2
            shift 2
            ;;
        -listenmode)
            LISTEN_MODE=$2
            shift 2
            ;;
        -listen_ipv4)
            LISTEN_IPV4=$2
            shift 2
            ;;
        -listen_ipv6)
            LISTEN_IPV6=$2
            shift 2
            ;;
        -ehost)
            ENTRY_LISTEN_ADDR=$2
            shift 2
            ;;
        -eport)
            ENTRY_PORT=$2
            shift 2
            ;;
        -ahost)
            ADMIN_LISTEN_ADDR=$2
            shift 2
            ;;
        -aport)
            ADMIN_PORT=$2
            shift 2
            ;;
        -stf_dproto)
            STATEFULL_DEAL_PROTOCOL=$2
            shift 2
            ;;
        -stf_dhost)
            STATEFULL_DEAL_LISTEN_ADDR=$2
            shift 2
            ;;
        -stf_dports)
            STATEFULL_DEAL_PORTS=$2
            shift 2
            ;;
        -stf_sproto)
            STATEFULL_SUBSCRIBE_PROTOCOL=$2
            shift 2
            ;;
        -stf_shost)
            STATEFULL_SUBSCRIBE_LISTENA_ADDR=$2
            shift 2
            ;;
        -stf_sports)
            STATEFULL_SUBSCRIBE_PORTS=$2
            shift 2
            ;;
        -usestl)
            USE_STATELESS=$2
            shift 2
            ;;
        -)
            shift 1
            break
            ;;
        -*)
            echo "$PROGNAME: illegal option -- '$(echo $1 | sed 's/^-*//')'" 1>&2
            exit 1
            ;;
        *)
            shift 1
            ;;
    esac
done

if [ "x${PERFORMANCE_MODE}" = "x0" ];then
IMAGE_NAME=openrelay
else
IMAGE_NAME=openrelay
fi

export LD_LIBRARY_PATH=/usr/local/openrelay/lib
${DRYRUN} ${IMAGE_PATH}/${IMAGE_NAME} \
-standbymode=${STANDBYMODE} \
-recmode=${REC_MODE} \
-repmode=${REP_MODE} \
-log=${LOG_LEVEL} \
-logdir=${LOG_DIRECTORY} \
-hbtimeout=${HEATBEAT_TIMEOUT} \
-jointimeout=${JOIN_TIMEOUT} \
-listenmode=${LISTEN_MODE} \
-listen_ipv4=${LISTEN_IPV4} \
-listen_ipv6=${LISTEN_IPV6} \
-ehost=${ENTRY_LISTEN_ADDR} \
-eport=${ENTRY_PORT} \
-ahost=${ADMIN_LISTEN_ADDR} \
-aport=${ADMIN_PORT} \
-stf_dproto=${STATEFULL_DEAL_PROTOCOL} \
-stf_dhost="${STATEFULL_DEAL_LISTEN_ADDR}" \
-stf_dports=${STATEFULL_DEAL_PORTS} \
-stf_sproto=${STATEFULL_SUBSCRIBE_PROTOCOL} \
-stf_shost="${STATEFULL_SUBSCRIBE_LISTENA_ADDR}" \
-stf_sports=${STATEFULL_SUBSCRIBE_PORTS} \
-usestl=${USE_STATELESS}

