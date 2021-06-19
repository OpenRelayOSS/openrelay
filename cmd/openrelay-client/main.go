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
	"bytes"
	"bufio"
	"encoding/hex"
	"encoding/binary"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"github.com/zeromq/goczmq"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"openrelay/internal/defs"
)

type replay struct {
	Tick  int64
	Frame []byte
}

type client struct {
	Filter string
	Id     int
	PosA   int
	PosB   int
	Deal   *goczmq.Sock
	Sub    *goczmq.Sock
	Err    error
}

var msgSize = 740
var interval = 500
var replaysA []replay
var replaysB []replay

var (
	dropMode       int
	logLevel       int
	destAddr       string
	destDealSchema string
	destDealPort   string
	destSubSchema  string
	destSubPort    string
	errorThreshold int
	startId        int
	wakeCount      int
	wakeIntval     int
	filePath       string
)

func param() {
	flag.IntVar(&logLevel, "log", 0, "loglevel ... 0=lostonly, 1=verbose")
	flag.IntVar(&dropMode, "dropmode", 1, "dropg packet, no check sequence")
	flag.StringVar(&destAddr, "addr", "127.0.0.1", "destination address")
	flag.StringVar(&destDealSchema, "dschm", "tcp://", "destination dealer schema tcp or udp")
	flag.StringVar(&destDealPort, "dport", ":7001", "destination dealer port")
	flag.StringVar(&destSubSchema, "sschm", "tcp://", "destination subscribe schema tcp or udp")
	flag.StringVar(&destSubPort, "sport", ":7002", "destination subscribe port")
	flag.IntVar(&errorThreshold, "errorthreshold", 0, "error Threshold counter")
	flag.IntVar(&startId, "startid", 10, "id start num")
	flag.IntVar(&wakeCount, "wake", 1, "wake client")
	flag.IntVar(&wakeIntval, "wakeint", 30000, "wake client interval (milliseconds)")
	flag.StringVar(&filePath, "filepath", "/var/log/openrelay/replay.log", "replay file fullpath")
	flag.Parse()
}

func main() {
	param()
	replaysA = make([]replay, 0, 3000)
	replaysB = make([]replay, 0, 3000)
	replayFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error "+err.Error())
		os.Exit(1)
	}
	defer replayFile.Close()
	scanner := bufio.NewScanner(replayFile)
	for scanner.Scan() {
		rep := replay{}
		line := strings.Split(scanner.Text(), "\t")
		rep.Tick = 0
		msg, err := hex.DecodeString(line[3])
		if err != nil {
			fmt.Println("error "+err.Error())
			os.Exit(1)
		}
		rep.Frame = []byte(msg)
		if line[1] == "A" {
			replaysA = append(replaysA, rep)
		} else if line[1] == "B" {
			replaysB = append(replaysB, rep)
		}
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
	}
	log.Printf("A frame count %d\n", len(replaysA))
	log.Printf("B frame count %d\n", len(replaysB))
	i := startId
	for {
		if i < startId+wakeCount {
			cli := client{}
			cli.Id = i
			cli.PosA = 0
			cli.PosB = 0
			go Send(&cli)
			go Recv(&cli)
			time.Sleep(time.Duration(wakeIntval) * time.Millisecond)
			i++
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	os.Exit(0)
}

func GetNextAFrame(index int, uid int) ([]byte, error) {
	if len(replaysA)-1 < index {
		return []byte{}, errors.New("out of index")
	}
	rewritedMsg, err := rewriteSrcUid(replaysA[index].Frame, defs.PlayerId(uid))
	if err != nil {
		return nil, err
	}
	return rewritedMsg, nil
}

func GetNextBFrame(index int, uid int) ([]byte, error) {
	if len(replaysB)-1 < index {
		return []byte{}, errors.New("out of index")
	}
	rewritedMsg, err := rewriteSrcUid(replaysB[index].Frame, defs.PlayerId(uid))
	if err != nil {
		return nil, err
	}
	return rewritedMsg, nil
}

func rewriteSrcUid(record []byte, uid defs.PlayerId) ([]byte, error) {
	var err error
	readBuf := bytes.NewReader([]byte(record))
	header := defs.Header{}
	err = binary.Read(readBuf, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
        }
	header.SrcUid = uid
	writeBuf := new(bytes.Buffer)
	err = binary.Write(writeBuf, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	err = binary.Write(writeBuf, binary.LittleEndian, record[defs.HeaderBytesLen:])
	if err != nil {
		return nil, err
	}
	return writeBuf.Bytes(), nil
}

func Send(cli *client) {
	var err error
	cli.Deal, err = goczmq.NewDealer(destDealSchema + destAddr + destDealPort)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Deal.Destroy()
	var frame []byte
	var joinFrame []byte
	var isABloop = false

	joinFrame, err = createReplayJoinMessage(cli.Id)
	if err != nil {
		log.Fatal(err)
	}

        readBuf := bytes.NewReader(joinFrame)
        header := defs.Header{}
        err = binary.Read(readBuf, binary.LittleEndian, &header)
        if err != nil {
		log.Fatal(err)
        }

        if header.Ver != defs.FrameVersion {
        }

        log.Printf("header.Ver: '%d' ", header.Ver)
        log.Printf("header.RelayCode: '%d' ", header.RelayCode)
        log.Printf("header.ContentCode: '%d' ", header.ContentCode)
        log.Printf("header.DestCode: '%d' ", header.DestCode)
        log.Printf("header.Mask: '%d' ", header.Mask)
        log.Printf("header.SrcUid: '%d' ", header.SrcUid)
        log.Printf("header.SrcOid: '%d' ", header.SrcOid)
        log.Printf("header.DestLen: '%d' ", header.DestLen)
        log.Printf("header.ContentLen: '%d' ", header.ContentLen)

	cli.Deal.SendFrame(joinFrame, goczmq.FlagNone)
	time.Sleep(time.Duration(3) * time.Second)

	for {
		if isABloop {
			frame, err = GetNextBFrame(cli.PosB, cli.Id)
			if err != nil {
				fmt.Println("err:", err)
				cli.PosB = 0
				frame, err = GetNextBFrame(cli.PosB, cli.Id)
			}
		} else {
			frame, err = GetNextAFrame(cli.PosA, cli.Id)
			if err != nil {
				fmt.Println("err:", err)
				isABloop = true
			}
		}
		err = cli.Deal.SendFrame(frame, goczmq.FlagNone)
		if err != nil {
			log.Fatal(err)
		}
		if logLevel > 0 {
			log.Printf("<- %s\n", hex.EncodeToString(frame))
			log.Printf("# id:%d A:%d B:%d\n", cli.Id, cli.PosA, cli.PosB)
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
		runtime.Gosched()
		if isABloop {
			cli.PosB++
		} else {
			cli.PosA++
		}
	}
}

func Recv(cli *client) {
	var err error
	cli.Sub, err = goczmq.NewSub(destSubSchema+destAddr+destSubPort, cli.Filter)
	//cli.Sub.SetSubscribe(cli.Filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Sub.Destroy()
	for {
		if dropMode == 1 {
			cli.Sub.RecvMessage()
		} else {
			reply, err := cli.Sub.RecvMessage()
			if err != nil {
				log.Fatal(err)
			}
			if logLevel > 0 {
				log.Printf("-> %s", string(reply[0]))
			}
		}
		runtime.Gosched()
	}
}

func createReplayJoinMessage(id int) ([]byte, error) {
	var err error
	joinSeed, err := newUUID()
	if err != nil {
		return nil, err
	}
	header := defs.Header{}
	header.Ver = defs.FrameVersion
	header.RelayCode = defs.REPLAY_JOIN
	header.ContentCode = 0
	header.DestCode = defs.ALL
	header.Mask = 0
	header.DestLen = 0
	header.SrcUid = defs.PlayerId(id)
	header.ContentLen = 16 // size of uint16
	writeBuf := new(bytes.Buffer)
	err = binary.Write(writeBuf, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	err = binary.Write(writeBuf, binary.LittleEndian, joinSeed)
	if err != nil {
		return nil, err
	}

	return writeBuf.Bytes(), nil
}

func createUUID() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return err.Error()
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func newUUID() ([]byte, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return uuid, err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return uuid, nil
}
