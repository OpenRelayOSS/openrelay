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

package defs

import (
	"github.com/zeromq/goczmq"
)

type RelayCode byte

const (
	CONNECT RelayCode = iota
	JOIN
	LEAVE
	RELAY
	TIMEOUT
	REJOIN
	SET_LEGACY_MAP
	GET_LEGACY_MAP
	GET_USERS
	SET_MASTER
	GET_MASTER
	GET_SERVER_TIMESTAMP
	RELAY_LATEST
	GET_LATEST
	SET_LOBBY_MAP
	GET_LOBBY_MAP
	SET_MASK //16
	GET_MASK
	PUSH_STACK
	FETCH_STACK
	REPLAY_JOIN
	RELAY_STREAM
	LOAD_PLAYER
	UPDATE_DIST_MAP
	PICK_DIST_MAP
	NOTIFY_DIST_MAP_LATEST
	// 100 - 199 Platform Dependency RelayCode
	UNITY_CDK_RELAY        = 100
	UNITY_CDK_RELAY_LATEST = 101
	UNITY_CDK_GET_LATEST   = 102
	UE4_CDK_RELAY          = 110
	UE4_CDK_RELAY_LATEST   = 111
	UE4_CDK_GET_LATEST     = 112
	// 200 - 255 User Define RelayCode
)

type ResponseCode uint16

const (
	OPENRELAY_RESPONSE_CODE_OK ResponseCode = iota
	OPENRELAY_RESPONSE_CODE_OK_NO_ROOM
	OPENRELAY_RESPONSE_CODE_OK_ROOM_ASSGIN_AND_CREATED
	OPENRELAY_RESPONSE_CODE_OK_POLLING_CONTINUE
	OPENRELAY_RESPONSE_CODE_NG
	OPENRELAY_RESPONSE_CODE_NG_REQUEST_READ_FAILED
	OPENRELAY_RESPONSE_CODE_NG_RESPONSE_WRITE_FAILED
	OPENRELAY_RESPONSE_CODE_NG_ENTRY_LOGIN_FAILED
	OPENRELAY_RESPONSE_CODE_NG_ENTRY_LOGIN_CLIENT_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_ENTRY_LOGIN_SERVER_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_GET_ROOM_INFO_NOT_FOUND
	OPENRELAY_RESPONSE_CODE_NG_GET_ROOM_INFO_CLIENT_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_GET_ROOM_INFO_SERVER_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_CREATE_ROOM_CAPACITY_OVER
	OPENRELAY_RESPONSE_CODE_NG_CREATE_ROOM_ASSIGN_FAILED
	OPENRELAY_RESPONSE_CODE_NG_CREATE_ROOM_CLIENT_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_CREATE_ROOM_SERVER_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_CREATE_ROOM_ALREADY_EXISTS
	OPENRELAY_RESPONSE_CODE_NG_JOIN_ROOM_NOT_FOUND
	OPENRELAY_RESPONSE_CODE_NG_JOIN_ROOM_CAPACITY_OVER
	OPENRELAY_RESPONSE_CODE_NG_JOIN_ROOM_FAILED
	OPENRELAY_RESPONSE_CODE_NG_JOIN_ROOM_CLIENT_TIMEOUT
	OPENRELAY_RESPONSE_CODE_NG_JOIN_ROOM_SERVER_TIMEOUT
)

const (
	OTHERS = iota
	ALL
	MASTER
	INCLUDE
	EXCLUDE
)

const (
	BLOCK_ROOM_MAX           = iota // Eager join retry
	BLOCK_ROOM_AND_QUEUE_MAX        // Economy join retry
)

type RelayStatus uint8
const (
	LISTEN RelayStatus = iota
	CLOSING
	CLOSED
)

type ABLoop string

const (
	ALoop ABLoop = "A"
	BLoop ABLoop = "B"
)

//16byte/128bit
type Header struct {
	Ver         byte
	RelayCode   RelayCode // = byte
	ContentCode byte
	DestCode    byte // 4byte
	Mask        byte
	SrcUid      PlayerId
	_           byte // 4byte alignment
	SrcOid      ObjectId
	DestLen     uint16 // 4byte
	ContentLen  uint16
	_           [2]byte // 4byte
}
const HeaderBytesLen = 16

type RoomParameter struct {
	Id            [16]byte
	Name          string
	Index         int
	Filter        string
	Capacity      uint16
	QueuingPolicy byte
	Stealth       bool
	ListenMode    byte
	StfDealPort   uint16
	StfSubPort    uint16
	UseStateless  bool
	StlDealPort   uint16
	StlSubPort    uint16
}

type RoomInstance struct {
	Guids         map[string]PlayerId
	Uids          map[PlayerId]string
	Names         map[PlayerId]string
	Hbs           map[PlayerId]int64
	Props         map[string][]byte
	Router        *goczmq.Sock
	Pub           *goczmq.Sock
	LastUid       PlayerId
	MasterUid     PlayerId
	MasterUidNeed bool
	Status        RelayStatus
	Log           *Logger
	Rec           *Recorder
	ABLoop        ABLoop
}

func (r *RoomInstance) ToListen(){
	r.Status = LISTEN
}

func (r *RoomInstance) ToClose(){
	r.Status = CLOSING
}

func (r *RoomInstance) ToClosed() {
	r.Status = CLOSED
}

type RoomResponse struct {
	Id             [16]byte // 16byte
	Capacity       uint16
	UserCount      uint16 // 4byte
	StfDealPort    uint16
	StfSubPort     uint16 // 4byte
	StlDealPort    uint16
	StlSubPort     uint16 // 4byte
	QueuingPolicy  byte
	Flags          byte //stealth | useStateless |x|x|x|x|x|x
	NameLen        byte
	FilterLen      byte      // 4byte
	Name           [256]byte // 256byte
	Filter         [256]byte // 256byte
	ListenMode     byte      // 0 = localnetonly, 1 = ipv4+ipv6both, 2 = ipv6only, 3 = ipv4only
	_              [3]byte   // 4byte alignment
	ListenAddrIpv4 [4]byte
	ListenAddrIpv6 [16]byte
}

type RoomJoinRequest struct {
	Seed      string
	Timestamp int64
}
