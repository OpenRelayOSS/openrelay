/* Copyright (c) 2018 FurtherSystem Co.,Ltd. All rights reserved.

   This program is free software; you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation; version 2 of the License.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program; if not, write to the Free Software
   Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1335  USA */

package srvs

import (
        "openrelay/internal/defs"
)

type OpenRelay struct {
        EntryHost      string
        EntryPort      string
        StfDealHost    string
        StfDealProto   string
        StfDealPorts   string
        StfSubHost     string
        StfSubProto    string
        StfSubPorts    string
        StlDealHost    string
        StlDealProto   string
        StlDealPorts   string
        StlSubHost     string
        StlSubProto    string
        StlSubPorts    string
        AdminHost      string
        AdminPort      string
	ListenIpv4     string
	ListenIpv6     string
	ListenMode     int
	LogLevel       int
	LogDir         string
	RecMode        int
	RepMode        bool
	HeatbeatTimeout int
	JoinTimeout    int
        JoinAllPollingQueue map[string][][]byte
        JoinAllProcessQueue map[string]defs.RoomJoinRequest
        JoinAllTimeoutQueue map[string][]defs.RoomJoinRequest
        JoinProcessTimeStart int64
        RoomQueue map[string]*defs.RoomParameter
        RelayQueue map[string]*defs.RoomInstance
        ReserveRooms map[string][16]byte
        ResolveRoomIds map[string]string
        HotRoomQueue [][16]byte
        ColdRoomQueue [][16]byte
        CleaningRoomQueue [][16]byte
}

func NewOpenRelay(eHost string, ePort string,
                  sfdHost string, sfdProto string, sfdPorts string,
                  sfsHost string, sfsProto string, sfsPorts string,
                  sldHost string, sldProto string, sldPorts string,
                  slsHost string, slsProto string, slsPorts string,
                  aHost string, aPort string,
                  listenIpv4 string, listenIpv6 string,
                  listenMode int, logLevel int, logDir string,
                  recMode int, repMode bool,
                  heatbeatTimeout int, joinTimeout int) *OpenRelay {
        return &OpenRelay{
                   EntryHost: eHost,
                   EntryPort: ePort,
                   StfDealHost: sfdHost,
                   StfDealProto: sfdProto,
                   StfDealPorts: sfdPorts,
                   StfSubHost: sfsHost,
                   StfSubProto: sfsProto,
                   StfSubPorts: sfsPorts,
                   StlDealHost: sldHost,
                   StlDealProto: sldProto,
                   StlDealPorts: sldPorts,
                   StlSubHost: slsHost,
                   StlSubProto: slsProto,
                   StlSubPorts: slsPorts,
                   AdminHost: aHost,
                   AdminPort: aPort,
                   ListenIpv4: listenIpv4,
                   ListenIpv6: listenIpv6,
                   ListenMode: listenMode,
                   LogLevel: logLevel,
                   LogDir: logDir,
                   RecMode: recMode,
                   RepMode: repMode,
                   HeatbeatTimeout: heatbeatTimeout,
                   JoinTimeout: joinTimeout,
                   JoinAllPollingQueue: make(map[string][][]byte, 0),
                   JoinAllProcessQueue: make(map[string]defs.RoomJoinRequest),
                   JoinAllTimeoutQueue: make(map[string][]defs.RoomJoinRequest, 0),
                   JoinProcessTimeStart: 0,
                   RoomQueue: make(map[string]*defs.RoomParameter, 0),
                   RelayQueue: make(map[string]*defs.RoomInstance, 0),
                   ReserveRooms: make(map[string][16]byte, 0),
                   ResolveRoomIds: make(map[string]string, 0),
                   HotRoomQueue: make([][16]byte, 0),
                   ColdRoomQueue: make([][16]byte, 0),
                   CleaningRoomQueue: make([][16]byte, 0),
                }
}
