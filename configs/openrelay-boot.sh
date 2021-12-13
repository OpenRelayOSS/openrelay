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
            shift 1
            ;;
        -perfmode)
            PERFORMANCE_MODE=$2
            shift 1
            ;;
        -standbymode)
            STANDBYMODE=$2
            shift 1
            ;;
        -recmode)
            REC_MODE=$2
            shift 1
            ;;
        -repmode)
            REP_MODE=$2
            shift 1
            ;;
        -log)
            LOG_LEVEL=$2
            shift 1
            ;;
        -logdir)
            LOG_DIRECTORY=$2
            shift 1
            ;;
        -hbtimeout)
            HEATBEAT_TIMEOUT=$2
            shift 1
            ;;
        -jointimeout)
            JOIN_TIMEOUT=$2
            shift 1
            ;;
        -listenmode)
            LISTEN_MODE=$2
            shift 1
            ;;
        -endpoint_ipv4)
            ENDPOINT_IPV4=$2
            shift 1
            ;;
        -endpoint_ipv6)
            ENDPOINT_IPV6=$2
            shift 1
            ;;
        -ehost)
            ENTRY_LISTEN_ADDR=$2
            shift 1
            ;;
        -eport)
            ENTRY_PORT=$2
            shift 1
            ;;
        -ahost)
            ADMIN_LISTEN_ADDR=$2
            shift 1
            ;;
        -aport)
            ADMIN_PORT=$2
            shift 1
            ;;
        -stf_dproto)
            STATEFULL_DEAL_PROTOCOL=$2
            shift 1
            ;;
        -stf_dhost)
            STATEFULL_DEAL_LISTEN_ADDR=$2
            shift 1
            ;;
        -stf_dports)
            STATEFULL_DEAL_PORTS=$2
            shift 1
            ;;
        -stf_sproto)
            STATEFULL_SUBSCRIBE_PROTOCOL=$2
            shift 1
            ;;
        -stf_shost)
            STATEFULL_SUBSCRIBE_LISTEN_ADDR=$2
            shift 1
            ;;
        -stf_sports)
            STATEFULL_SUBSCRIBE_PORTS=$2
            shift 1
            ;;
        -usestl)
            USE_STATELESS=$2
            shift 1
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

#
# Endpoint auto detect logics for local
#
if [ "x${LISTEN_MODE}" = "x0" ];then
  if [ "x${ENDPOINT_IPV4}" = "x" ];then
    RESULT=`${IMAGE_PATH}/lipcheck -4 -allowprivate`
    ENDPOINT_IPV4=`echo -n ${RESULT} | awk '{print $3}' | head -1`
    if [ "${ENDPOINT_IPV4}" != "x" ];then
      echo -- Endpoint_ipv4 detected: ${ENDPOINT_IPV4}
      echo -- Please set this Endpoint_ipv4 for /etc/sysconfig/openrelay.env
    else
      echo -- Endpoint_ipv4 detect failed, please check environment.
      DETECT_ENDPOINT_IPV4_FAILED=1
    fi
  else
      echo -- Endpoint_ipv4 fixed: ${ENDPOINT_IPV4} ok.
  fi
  if [ "x${ENDPOINT_IPV6}" = "x" ];then
    RESULT=`${IMAGE_PATH}/lipcheck -6 -allowlinklocal`
    ENDPOINT_IPV6=`echo -n ${RESULT} | awk '{print $3}' | head -1`
    if [ "${ENDPOINT_IPV6}" != "x" ];then
      echo -- Endpoint_ipv6 detected: ${ENDPOINT_IPV6}
      echo -- Please set this Endpoint_ipv6 for /etc/sysconfig/openrelay.env
    else
      echo -- Endpoint_ipv6 detect failed, please check environment.
      DETECT_ENDPOINT_IPV6_FAILED=1
    fi
  else
      echo -- Endpoint_ipv6 fixed: ${ENDPOINT_IPV6} ok.
  fi
fi

#
# Endpoint auto detect logics for global
#
if [ "x${LISTEN_MODE}" = "x1" ];then
  if [ "x${ENDPOINT_IPV4}" = "x" ];then
    RESULT=`${IMAGE_PATH}/gipcheck -4 -url https://ifconfig.io/ip`
    STATUS=`echo -n ${RESULT} | awk '{print $1}'`
    if [ "x${STATUS}" = "x200" ];then
      ENDPOINT_IPV4=`echo -n ${RESULT} | awk '{print $2}'`
      echo -- Endpoint_ipv4 detected: ${ENDPOINT_IPV4}
      echo -- Please set this Endpoint_ipv4 for /etc/sysconfig/openrelay.env
    else
      echo -- Endpoint_ipv4 detect failed, please check environment.
      DETECT_ENDPOINT_IPV4_FAILED=1
    fi
  else
      echo -- Endpoint_ipv4 fixed ${ENDPOINT_IPV4} ok.
  fi
  if [ "x${ENDPOINT_IPV6}" = "x" ];then
    RESULT=`${IMAGE_PATH}/gipcheck -6 -url https://ifconfig.io/ip`
    STATUS=`echo -n ${RESULT} | awk '{print $1}'`
    if [ "x${STATUS}" = "x200" ];then
      ENDPOINT_IPV6=`echo -n ${RESULT} | awk '{print $2}'`
      echo -- Endpoint_ipv6 detected: ${ENDPOINT_IPV6}
      echo -- Please set this Endpoint_ipv6 for /etc/sysconfig/openrelay.env
    else
      echo -- Endpoint_ipv6 detect failed, please check environment.
      DETECT_ENDPOINT_IPV6_FAILED=1
    fi
  else
      echo -- Endpoint_ipv6 fixed: ${ENDPOINT_IPV6} ok.
  fi
fi

if [  "x${DETECT_ENDPOINT_IPV4_FAILED}" = "x1" -a "x${DETECT_ENDPOINT_IPV6_FAILED}" = "x1" ];then
   echo -- ipv4/ipv6 both endpoint detect failed. cannot wakeup.
   exit 1
fi

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
-endpoint_ipv4=${ENDPOINT_IPV4} \
-endpoint_ipv6=${ENDPOINT_IPV6} \
-ehost=${ENTRY_LISTEN_ADDR} \
-eport=${ENTRY_PORT} \
-ahost=${ADMIN_LISTEN_ADDR} \
-aport=${ADMIN_PORT} \
-stf_dproto=${STATEFULL_DEAL_PROTOCOL} \
-stf_dhost="${STATEFULL_DEAL_LISTEN_ADDR}" \
-stf_dports=${STATEFULL_DEAL_PORTS} \
-stf_sproto=${STATEFULL_SUBSCRIBE_PROTOCOL} \
-stf_shost="${STATEFULL_SUBSCRIBE_LISTEN_ADDR}" \
-stf_sports=${STATEFULL_SUBSCRIBE_PORTS} \
-usestl=${USE_STATELESS}

