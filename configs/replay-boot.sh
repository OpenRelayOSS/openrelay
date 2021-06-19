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
        -log)
            LOG_LEVEL=$2
            shift 1
            ;;
        -dropmode)
            DROP_MODE=$2
            shift 1
            ;;
        -filepath)
            FILEPATH=$2
            shift 1
            ;;
        -addr)
            DEST_ADDR=$2
            shift 1
            ;;
        -dschm)
            DEST_DEAL_SCHEMA=$2
            shift 1
            ;;
        -dport)
            DEST_PORT=$2
            shift 1
            ;;
        -sschm)
            DEST_SUB_SCHEMA=$2
            shift 1
            ;;
        -sport)
            DEST_SUB_PORT=$2
            shift 1
            ;;
        -errorthreshold)
            ERROR_THRESHOLD=$2
            shift 1
            ;;
        -startid)
            START_ID=$2
            shift 1
            ;;
        -wake)
            WAKE=$2
            shift 1
            ;;
        -wakeint)
            WAKE_INT=$2
            shift 1
            ;;
        -logdir)
            LOG_DIRECTORY=$2
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

if [ "x${PERFORMANCE_MODE}" = "x0" ];then
IMAGE_NAME=replay
else
IMAGE_NAME=replay
fi

export LD_LIBRARY_PATH=/usr/local/openrelay/lib
${DRYRUN} ${IMAGE_PATH}/${IMAGE_NAME} -filepath ${FILEPATH} -addr ${DEST_ADDR} -startid ${START_ID} -wake ${WAKE} -wakeint ${WAKE_INT} -log ${LOG_LEVEL}
