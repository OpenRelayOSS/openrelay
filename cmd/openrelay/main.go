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

package main

import (
	crand "crypto/rand"
	"flag"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"openrelay/internal/defs"
	"openrelay/internal/srvs"
)

var (
	standbyMode  int
	recMode      int
	repMode      bool
	logLevel     int
	logDir       string
	hbTimeout    int
	joinTimeout  int
	listenMode   int
	listenIpv4   string
	listenIpv6   string
	entryHost    string
	entryPort    string
	adminHost    string
	adminPort    string
	stfDealProto string
	stfDealHost  string
	stfDealPorts string
	stfSubProto  string
	stfSubHost   string
	stfSubPorts  string
	useStateless bool
	stlDealProto string
	stlDealHost  string
	stlDealPorts string
	stlSubProto  string
	stlSubHost   string
	stlSubPorts  string
)

func param() {
	flag.IntVar(&standbyMode, "standbymode", -1, "0=allcold, 1<standbymode is pre wake room, -1=allhot")
	flag.IntVar(&recMode, "recmode", 0, "recording mode ... 0=off, 0<recmode is userId ")
	flag.BoolVar(&repMode, "repmode", false, "replay mode ... false=off, true=on ")
	flag.IntVar(&logLevel, "log", 0, "loglevel ... 0=fatalonly, 1=erroronly 2=info, 3=verbose, 4=veryverbose")
	flag.StringVar(&logDir, "logdir", "/var/log/openrelay", "base log directory")
	flag.IntVar(&hbTimeout, "hbtimeout", 30, "heatbeat timeout sec")
	flag.IntVar(&joinTimeout, "jointimeout", 180, "heatbeat timeout sec")
	flag.IntVar(&listenMode, "listenmode", 3, "0=localnetonly, 1=ipv4+ipv6both, 2=ipv6only, 3=ipv4only, 1=ipv4+ipv6bothauto, 2=ipv6onlyauto, 3=ipv4onlyauto")
	flag.StringVar(&listenIpv4, "listen_ipv4", "", "listen global ip addr v4")
	flag.StringVar(&listenIpv6, "listen_ipv6", "", "listen global ip addr v6")
	flag.StringVar(&entryHost, "ehost", "localhost", "entry http service listen host")
	flag.StringVar(&entryPort, "eport", "7000", "entry http service port")
	flag.StringVar(&adminHost, "ahost", "localhost", "admin tcp console listen host")
	flag.StringVar(&adminPort, "aport", "8000", "admin tcp console port")
	flag.StringVar(&stfDealProto, "stf_dproto", "tcp", "statefull dealer protocol tcp or udp")
	flag.StringVar(&stfDealHost, "stf_dhost", "*", "statefull dealer listen host")
	flag.StringVar(&stfDealPorts, "stf_dports", "7001,7003,7005,7007", "statefull dealer port, use separate comma")
	flag.StringVar(&stfSubProto, "stf_sproto", "tcp", "statefull subscribe protocol tcp or udp")
	flag.StringVar(&stfSubHost, "stf_shost", "*", "statefull subscribe listen host")
	flag.StringVar(&stfSubPorts, "stf_sports", "7002,7004,7006,7008", "statefull subscribe port, use separate comma")
	flag.BoolVar(&useStateless, "usestl", false, "enable stateless deal/subscribe services ")
	flag.StringVar(&stlDealProto, "stl_dproto", "tcp", "stateless dealer protocol tcp or udp")
	flag.StringVar(&stlDealHost, "stl_dhost", "*", "stateless dealer listen host")
	flag.StringVar(&stlDealPorts, "stl_dports", "7001,7003,7005,7007", "stateless dealer port, use separate comma")
	flag.StringVar(&stlSubProto, "stl_sproto", "tcp", "stateless subscribe protocol tcp or udp")
	flag.StringVar(&stlSubHost, "stl_shost", "*", "stateless subscribe listen host")
	flag.StringVar(&stlSubPorts, "stl_sports", "7002,7004,7006,7008", "stateless subscribe port, use separate comma")
	flag.Parse()
}

func main() {
	param()
	o := srvs.NewOpenRelay(entryHost, entryPort,
                                 stfDealHost, stfDealProto, stfDealPorts,
                                 stfSubHost, stfSubProto, stfSubPorts,
                                 stlDealHost, stlDealProto, stlDealPorts,
                                 stlSubHost, stlSubProto, stlSubProto,
                                 adminHost, adminPort,
                                 listenIpv4, listenIpv6,
                                 listenMode, logLevel, logDir,
                                 recMode, repMode,
                                 hbTimeout, joinTimeout)
	Initialize(o)
	go o.ConsoleServ()
	go o.EntryServ()
	// TODO setup trap signal
	for {
		time.Sleep(1 * time.Second) // forcely return context
	}
}

func Initialize(o *srvs.OpenRelay) {
        var err error
        seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64)) // TODO mt19937
        rand.Seed(seed.Int64())
        // check stl enable but didn't set
        stfDealPortArray := strings.Split(o.StfDealPorts, ",")
        stfSubPortArray := strings.Split(o.StfSubPorts, ",")
        portCount := len(stfDealPortArray)
        for index := 0; index < portCount; index++ {
                // check port valid
                // check port count
                // check port conflict
                room := defs.RoomParameter{}
                room.ListenMode = byte(listenMode)
                relayInstance := defs.RoomInstance{EnableBflag: false}
                room.Id, err = defs.NewGuid()
                if err != nil {
                        log.Fatal("guid cannot create, initialize faild. ", err)
                }
                roomIdStr := string(room.Id[:])
                var port int
                port, err = strconv.Atoi(stfDealPortArray[index])
                if err != nil {
                        log.Fatal("invalid port, initialize faild. ", err)
                }
                room.StfDealPort = uint16(port)
                port, err = strconv.Atoi(stfSubPortArray[index])
                if err != nil {
                        log.Fatal("invalid port, initialize faild. ", err)
                }
                room.StfSubPort = uint16(port)
                room.UseStateless = false
                o.HotRoomQueue = append(o.HotRoomQueue, room.Id)
                o.RoomQueue[roomIdStr] = &room
                o.RelayQueue[roomIdStr] = &relayInstance
        }
        for _, id := range o.HotRoomQueue {
                idStr := string(id[:])
                o.Clean(o.RelayQueue[idStr], o.RoomQueue[idStr].Id)
                go o.RelayServ(o.RoomQueue[idStr], o.RelayQueue[idStr])
                go o.Heatbeat(o.RelayQueue[idStr], id)
        }
        if o.LogLevel >= defs.INFO {
                log.Printf("available room :%d", len(o.HotRoomQueue))
                log.Printf("initialize ok")
        }
}
