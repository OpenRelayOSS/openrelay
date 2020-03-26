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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/zeromq/goczmq"
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"openrelay/internal/defs"
//	"github.com/pion/dtls/examples/util"
)

func (o *OpenRelay) RelayServ(room *defs.RoomParameter, relay *defs.RoomInstance) {
	var err error

	startTime := time.Now()
	relay.Guids = make(map[string]defs.PlayerId)
	relay.Uids = make(map[defs.PlayerId]string)
	relay.Names = make(map[defs.PlayerId]string)
	relay.Hbs = make(map[defs.PlayerId]int64)
	relay.Props = make(map[string][]byte)
	relay.LastUid = 0
	relay.MasterUid = 0
	relay.MasterUidNeed = true
	relay.EnableBflag = false

	roomIdStr := string(room.Id[:])
	joinPollingQueue := o.JoinAllPollingQueue[roomIdStr]
	joinPollingQueue = make([][]byte, 0)
	o.JoinAllPollingQueue[roomIdStr] = joinPollingQueue

	recPrefix := "[RECMODE]"
	logfile, err := os.OpenFile(o.LogDir+"/relay-trace-"+defs.GuidFormatString(room.Id)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Panicf("cannot open "+o.LogDir+"/relay-trace-"+defs.GuidFormatString(room.Id)+".log:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))


	//addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: int(room.StlDealPort)}
	//config := &dtls.Config{
	//	PSK: func(hint []byte) ([]byte, error) {
	//		fmt.Printf("Client's hint: %s \n", hint)
	//		return []byte{0x7C, 0xCD, 0xE1, 0x4A, 0x5C, 0xF3, 0xB7, 0x1C, 0x0C, 0x08, 0xC8, 0xB7, 0xF9, 0xE5}, nil
	//	},
	//	PSKIdentityHint:      []byte("oFIrQFrW8EWcZ5u7eGfrkw"),
	//	CipherSuites:         []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
	//	ExtendedMasterSecret: dtls.DisableExtendedMasterSecret, //dtls.RequireExtendedMasterSecret,
	//	ConnectTimeout:       dtls.ConnectTimeoutOption(30 * time.Second),
	//}
	//listener, err := dtls.Listen("udp", addr, config)
	//util.Check(err)
	//defer func() {
	//	util.Check(listener.Close(5 * time.Second))
	//}()
	//fmt.Println("Listening")
	// Simulate a chat session
	//hub := util.NewHub()
	//hub.Chat()

	//go func() {
	//	for {
	//		// Wait for a connection.
	//		conn, err := listener.Accept()
	//		util.Check(err)
	//		// defer conn.Close() // TODO: graceful shutdown
	//
	//		// Register the connection with the chat hub
	//		hub.Register(conn)
	//	}
	//}()

	relay.Router, err = goczmq.NewRouter(o.StfDealProto + "://" + o.StfDealHost + ":" + strconv.Itoa(int(room.StfDealPort)))
	if err != nil {
		log.Printf("fail roomId " + defs.GuidFormatString(room.Id))
		//log.Fatal("relay.Router create failed. "+o.StfDealProto+"://"+o.StfDealHost+":"+strconv.Itoa(int(room.StfDealPort)), err)
		log.Printf("relay.Router create failed. " + o.StfDealProto + "://" + o.StfDealHost + ":" + strconv.Itoa(int(room.StfDealPort)) + err.Error())
		return
	}
	defer relay.Router.Destroy()

	relay.Pub, err = goczmq.NewPub(o.StfSubProto + "://" + o.StfSubHost + ":" + strconv.Itoa(int(room.StfSubPort)))
	if err != nil {
		//log.Fatal("relay.Pub create failed. "+o.StfSubProto+"://"+o.StfSubHost+":"+strconv.Itoa(int(room.StfSubPort)), err)
		log.Printf("relay.Pub create failed. " + o.StfSubProto + "://" + o.StfSubHost + ":" + strconv.Itoa(int(room.StfSubPort)) + err.Error())
		return
	}
	defer relay.Pub.Destroy()

	header := defs.Header{}
	for {
		request, err := relay.Router.RecvMessage()
		if err != nil {
			if o.LogLevel >= defs.ERRORONLY {
				log.Println("relay.Router recv failed. ", err)
			}
			continue
		}
		if o.LogLevel >= defs.VVERBOSE {
			log.Printf("relay.Router received '%s' from '%v'", hex.EncodeToString(request[1]), request[0])
		}
		if request == nil || len(request) < 2 {
			if o.LogLevel >= defs.ERRORONLY {
				log.Println("invalid request.. ")
			}
			continue
		}

		readBuf := bytes.NewReader(request[1])
		header = defs.Header{}
		err = binary.Read(readBuf, binary.LittleEndian, &header)
		if err != nil {
			if o.LogLevel >= defs.ERRORONLY {
				log.Println("binary read failed. ", err)
			}
			continue
		}

		if header.Ver != defs.FrameVersion {
			if o.LogLevel >= defs.ERRORONLY {
				log.Printf("invalid FrameVersion %d != %d", defs.FrameVersion, header.Ver)
			}
			continue
		}

		if o.LogLevel >= defs.VVERBOSE {
			log.Printf("received header.Ver: '%d' ", header.Ver)
			log.Printf("received header.RelayCode: '%d' ", header.RelayCode)
			log.Printf("received header.ContentCode: '%d' ", header.ContentCode)
			log.Printf("received header.DestCode: '%d' ", header.DestCode)
			log.Printf("received header.Mask: '%d' ", header.Mask)
			log.Printf("received header.SrcUid: '%d' ", header.SrcUid)
			log.Printf("received header.SrcOid: '%d' ", header.SrcOid)
			log.Printf("received header.DestLen: '%d' ", header.DestLen)
			log.Printf("received header.ContentLen: '%d' ", header.ContentLen)
		}

		switch header.RelayCode {
		case defs.RELAY, defs.RELAY_STREAM, defs.UNITY_CDK_RELAY, defs.UE4_CDK_RELAY:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			//destUids := make([]byte, header.DestLen)
			//content := make([]byte, header.ContentLen)
			//err = binary.Read(readBuf, binary.LittleEndian, &destUids)
			//err = binary.Read(readBuf, binary.LittleEndian, &content)

			err = relay.Pub.SendFrame(request[1], goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay %d '%s' ", header.RelayCode, hex.EncodeToString(request[1]))
			}
			if o.RecMode == int(header.SrcUid) && relay.EnableBflag {
				log.Printf(recPrefix+"%d;%s;%s", time.Now().UnixNano(), "B", request[1])
			} else if o.RecMode == int(header.SrcUid) && !relay.EnableBflag {
				log.Printf("relay.LastUid: %d", relay.LastUid)
				log.Printf(recPrefix+"%d;%s;%s", time.Now().UnixNano(), "A", request[1])
			}

		case defs.JOIN:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			alignmentLen := uint16(0)
			alignment := []byte{}

			var seedLen uint16
			err = binary.Read(readBuf, binary.LittleEndian, &seedLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("received join seedLen: '%d' ", seedLen)
			}

			var nameLen uint16
			err = binary.Read(readBuf, binary.LittleEndian, &nameLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("received join nameLen: '%d' ", nameLen)
			}

			joinSeed := make([]byte, seedLen)
			err = binary.Read(readBuf, binary.LittleEndian, &joinSeed)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}

			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("received join seed: '%s' ", hex.EncodeToString(joinSeed))
			}

			//read adjust alignment at seedLen
			alignmentLen = seedLen % 4
			if alignmentLen != 0 {
				alignment = make([]byte, alignmentLen)
				err = binary.Read(readBuf, binary.LittleEndian, &alignment)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("binary read failed. ", err)
					}
					continue
				}
			}

			name := make([]byte, nameLen)
			err = binary.Read(readBuf, binary.LittleEndian, &name)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}

			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("received join name: '%s' ", string(name))
			}

			assginUid := relay.Guids[string(joinSeed)]
			relay.Names[relay.LastUid] = string(name)
			header.SrcUid = relay.LastUid
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, assginUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, relay.MasterUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, seedLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, nameLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = binary.Write(writeBuf, binary.LittleEndian, joinSeed)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			//write adjust alignment at seedLen.
			alignmentLen = seedLen % 4
			if alignmentLen != 0 {
				alignment = make([]byte, alignmentLen)
				err = binary.Write(writeBuf, binary.LittleEndian, alignment)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("binary write failed. ", err)
					}
					continue
				}
			}

			err = binary.Write(writeBuf, binary.LittleEndian, name)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(writeBuf.Bytes()))
			}

			if o.RecMode == int(relay.LastUid) {
				log.Printf("relay.LastUid: %d", relay.LastUid)
				log.Printf(recPrefix+"%d;%s;%s", time.Now().UnixNano(), "A", request[1])
			}

		case defs.LEAVE:
			joinSeed := make([]byte, header.ContentLen)
			err = binary.Read(readBuf, binary.LittleEndian, &joinSeed)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			srcUid := relay.Guids[string(joinSeed)]
			if srcUid != header.SrcUid {
				if o.LogLevel >= defs.ERRORONLY {
					log.Printf("invalid srcUid %l !=  %l", srcUid, header.SrcUid)
				}
				continue
			}
			delete(relay.Guids, string(joinSeed))
			delete(relay.Uids, srcUid)
			delete(relay.Names, srcUid)
			delete(relay.Hbs, srcUid)

			if len(relay.Guids) == 0 {
				o.Clean(relay, room.Id)
			} else if relay.MasterUid == srcUid {
				for i, _ := range relay.Uids {
					relay.MasterUid = i
					break
				}
			}

			header.ContentLen = 0
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, relay.MasterUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.INFO {
				log.Println("-> leave ", srcUid)
			}

		case defs.TIMEOUT:
		case defs.REJOIN:
		case defs.SET_LEGACY_MAP:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("invalid srcUid: ", header.SrcUid)
				}
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			var keysLen uint16
			err = binary.Read(readBuf, binary.LittleEndian, &keysLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("received join keysLen: '%d' ", keysLen)
			}

			var propsLen uint16
			err = binary.Read(readBuf, binary.LittleEndian, &propsLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("received join propsLen: '%d' ", propsLen)
			}

			keysBytes := make([]byte, keysLen)
			err = binary.Read(readBuf, binary.LittleEndian, &keysBytes)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}

			//read adjust alignment at keysLen
			var alignmentLen = keysLen % 4
			if alignmentLen != 0 {
				var alignment = make([]byte, alignmentLen)
				err = binary.Read(readBuf, binary.LittleEndian, &alignment)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("binary read failed. ", err)
					}
					continue
				}
			}

			properties := make([]byte, propsLen)
			err = binary.Read(readBuf, binary.LittleEndian, &properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			relay.Props[defs.PropKeyLegacy] = properties
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = binary.Write(writeBuf, binary.LittleEndian, keysLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, propsLen)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = binary.Write(writeBuf, binary.LittleEndian, keysBytes)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			//write adjust alignment at keysLen.
			alignmentLen = keysLen % 4
			if alignmentLen != 0 {
				var alignment = make([]byte, alignmentLen)
				err = binary.Write(writeBuf, binary.LittleEndian, alignment)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("binary write failed. ", err)
					}
					continue
				}
			}

			err = binary.Write(writeBuf, binary.LittleEndian, properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("set legacy map %s \n", relay.Props[defs.PropKeyLegacy])
			}

		case defs.GET_LEGACY_MAP:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("invalid srcUid: ", header.SrcUid)
				}
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			header.ContentLen = uint16(len(relay.Props[defs.PropKeyLegacy]))
			properties := relay.Props[defs.PropKeyLegacy]
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("get legacy map %s \n", relay.Props[defs.PropKeyLegacy])
			}

		case defs.GET_USERS:
		case defs.SET_MASTER:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			header.ContentLen = 0
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, relay.MasterUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(request[1]))
			}

		case defs.GET_MASTER:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			header.ContentLen = 2
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, relay.MasterUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(request[1]))
			}

		case defs.GET_SERVER_TIMESTAMP:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			timestamp := uint16(time.Since(startTime) / time.Second)
			header.ContentLen = 2
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, timestamp)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(request[1]))
			}

		case defs.RELAY_LATEST, defs.UNITY_CDK_RELAY_LATEST, defs.UE4_CDK_RELAY_LATEST:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			properties := make([]byte, header.ContentLen)
			err = binary.Read(readBuf, binary.LittleEndian, &properties)
			relay.Props[defs.PropKeyPlayerPrefix+strconv.Itoa(int(header.SrcUid))] = properties
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(request[1]))
			}

		case defs.GET_LATEST, defs.UNITY_CDK_GET_LATEST, defs.UE4_CDK_GET_LATEST:
			if _, ok := relay.Hbs[header.SrcUid]; !ok {
				log.Println("source uid is invalid ", header.SrcUid)
				continue
			}
			relay.Hbs[header.SrcUid] = time.Now().Unix()

			var targetUid uint16
			err = binary.Read(readBuf, binary.LittleEndian, &targetUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("get latest uid:%d latest stack", targetUid)
			}

			properties := relay.Props[defs.PropKeyPlayerPrefix+strconv.Itoa(int(header.SrcUid))]
			header.ContentLen = uint16(len(properties))
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(request[1]))
			}

		case defs.SET_LOBBY_MAP:
			//if _, ok := relay.Hbs[header.SrcUid]; !ok {
			//	if o.LogLevel >= defs.ERRORONLY {log.Println("invalid srcUid: ", header.SrcUid) }
			//	continue
			//}
			//relay.Hbs[header.SrcUid] = time.Now().Unix()

			properties := make([]byte, header.ContentLen)
			err = binary.Read(readBuf, binary.LittleEndian, &properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary read failed. ", err)
				}
				continue
			}
			relay.Props[defs.PropKeyLegacyLobby] = properties
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("set lobby map %s \n", relay.Props[defs.PropKeyLegacyLobby])
			}

		case defs.GET_LOBBY_MAP:
			//if _, ok := relay.Hbs[header.SrcUid]; !ok {
			//	if o.LogLevel >= defs.ERRORONLY {log.Println("invalid srcUid: ", header.SrcUid) }
			//	continue
			//}
			//relay.Hbs[header.SrcUid] = time.Now().Unix()

			header.ContentLen = uint16(len(relay.Props[defs.PropKeyLegacyLobby]))
			properties := relay.Props[defs.PropKeyLegacyLobby]
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, properties)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("get legacy map %s \n", relay.Props[defs.PropKeyLegacy])
			}

		case defs.REPLAY_JOIN:
			relay.LastUid += 1
			if relay.MasterUidNeed {
				relay.MasterUidNeed = false
				relay.MasterUid = relay.LastUid
			}
			joinSeed := make([]byte, header.ContentLen)
			err = binary.Read(readBuf, binary.LittleEndian, &joinSeed)
			if err != nil {
				log.Println("read joinseed failed. ", err)
				continue
			}
			assginUid := relay.LastUid
			//relay.MasterUid := relay.MasterUid
			joinedUids := []defs.PlayerId{}
			for k, _ := range relay.Uids {
				joinedUids = append(joinedUids, k)
			}
			relay.Guids[string(joinSeed)] = relay.LastUid
			relay.Uids[relay.LastUid] = string(joinSeed)
			writeBuf := new(bytes.Buffer)
			err = binary.Write(writeBuf, binary.LittleEndian, header)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, assginUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, relay.MasterUid)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			err = binary.Write(writeBuf, binary.LittleEndian, joinedUids)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("binary write failed. ", err)
				}
				continue
			}
			relay.Hbs[relay.LastUid] = time.Now().Unix()

			err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
			if err != nil {
				if o.LogLevel >= defs.ERRORONLY {
					log.Println("frame send failed. ", err)
				}
				continue
			}
			if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> relay '%s' ", hex.EncodeToString(request[1]))
			}
			if o.RecMode == int(relay.LastUid) {
				log.Printf("relay.LastUid: %d", relay.LastUid)
				log.Printf(recPrefix+"%d;%s;%s", time.Now().UnixNano(), "A", request[1])
			}

		case defs.PUSH_STACK:
			log.Printf("message code defs.PUSH_STACK ... %d\n", header.RelayCode)
		case defs.FETCH_STACK:
			log.Printf("message code defs.FETCH_STACK ... %d\n", header.RelayCode)
		case defs.CONNECT:
		default:
			log.Printf("invalid message code ... %d\n", header.RelayCode)
		}

		time.Sleep(0 * time.Second) // return context
	}
}

func (o *OpenRelay) Clean(relay *defs.RoomInstance, roomId [16]byte) {
	roomIdStr := string(roomId[:])
	roomName := o.ResolveRoomIds[roomIdStr]
	delete(o.ReserveRooms, roomName)
	delete(o.ResolveRoomIds, roomIdStr)

	relay.MasterUidNeed = true
	relay.Guids = make(map[string]defs.PlayerId)
	relay.Uids = make(map[defs.PlayerId]string)
	relay.Names = make(map[defs.PlayerId]string)
	relay.Hbs = make(map[defs.PlayerId]int64)
	relay.Props = make(map[string][]byte)
	relay.LastUid = 0
	relay.MasterUid = 0
	relay.MasterUidNeed = true
	relay.EnableBflag = false
	o.JoinAllProcessQueue[roomIdStr] = defs.RoomJoinRequest{Seed:"", Timestamp:0}

	joinPollingQueue := make([][]byte, 0)
	o.JoinAllPollingQueue[roomIdStr] = joinPollingQueue

	// restart here. relay, hbckeck

//	o.HotRoomQueue = append(o.HotRoomQueue, roomId)
	if o.LogLevel >= defs.INFO {
		log.Printf("cleaning room ok, id:%s", defs.GuidFormatString(roomId))
	}
}

func  (o *OpenRelay) Heatbeat(relay *defs.RoomInstance, roomId [16]byte) {
	var err error

	interval := time.Duration(500)
	timeout := int64(o.HeatbeatTimeout)
	for {
		for k, v := range relay.Hbs {
			if v+timeout < time.Now().Unix() {
				g := relay.Uids[k]
				delete(relay.Guids, g)
				delete(relay.Uids, k)
				delete(relay.Names, k)
				delete(relay.Hbs, k)

				if len(relay.Guids) > 0 && relay.MasterUid == k {
					for i, _ := range relay.Uids {
						relay.MasterUid = i
						break
					}
				}
				header := defs.Header{}
				header.Ver = 0
				header.RelayCode = defs.LEAVE
				header.ContentCode = 0
				header.DestCode = defs.ALL
				header.Mask = 0
				header.SrcUid = k
				header.DestLen = 0
				header.ContentLen = 0
				writeBuf := new(bytes.Buffer)
				err = binary.Write(writeBuf, binary.LittleEndian, header)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("binary write failed. ", err)
					}
					break
				}
				err = binary.Write(writeBuf, binary.LittleEndian, relay.MasterUid)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("binary write failed. ", err)
					}
					continue
				}

				err = relay.Pub.SendFrame(writeBuf.Bytes(), goczmq.FlagNone)
				if err != nil {
					if o.LogLevel >= defs.ERRORONLY {
						log.Println("frame send failed. ", err)
					}
				}
				if o.LogLevel >= defs.INFO {
					log.Printf("-> timeout force logout %s %d", hex.EncodeToString([]byte(g)), k)
				}

				if len(relay.Guids) == 0 {
					o.Clean(relay, roomId)
				}

			} else if o.LogLevel >= defs.VVERBOSE {
				log.Printf("-> heatbeat check ok uid: %d time: %d < %d \n", k, v+timeout, time.Now().Unix())
			}
		}
		time.Sleep(interval * time.Millisecond) // return context
	}
}
