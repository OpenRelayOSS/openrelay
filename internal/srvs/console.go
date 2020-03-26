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
        "log"
        "net"
        "runtime"
        "time"
)

func (o *OpenRelay) SetBLoopCommand(roomId string) {
	relay, exist := o.RelayQueue[roomId]
	if exist {
		log.Println("start b loop " + roomId + " failed.")
		log.Println("roomId not found.")
		return
	}
	log.Println("start b loop " + roomId)
	relay.EnableBflag = true
}

func (o *OpenRelay) ConsoleServ() {
	listen, err := net.Listen("tcp", ":"+o.AdminPort)
	if err != nil {
		log.Fatal("tcp://" + o.AdminHost + ":" + o.AdminPort + " listen failed.")
	}
	for {
		conn, err := listen.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal("connection accept failed.")
		}
		buf := make([]byte, 1024)
		go func() {
			for {
				n, _ := conn.Read(buf)
				if "" == string(buf[:n]) {
				} else if "setb\r\n" == string(buf[:n]) {
					o.SetBLoopCommand("TODO set room Id here")
					conn.Write([]byte("start b loop\r\n"))
				} else {
					conn.Write([]byte("invalid command >" + string(buf[:n]) + "< "))
				}
				runtime.Gosched()
				time.Sleep(1 * time.Second)
			}
		}()
	}
}
