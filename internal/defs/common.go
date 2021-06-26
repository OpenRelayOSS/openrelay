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
	"fmt"
	"math/rand"
)

type PlayerId uint16
type ObjectId uint16

// user agents.
const UA_NATIVE_CDK = "Native-CDK"
const UA_UNITY_CDK = "Unity-CDK"
const UA_UE4_CDK = "UE4-CDK"

// version.
var Version string
var Shorthash string

// require cdk version.
const REQUIRE_NATIVE_CDK_VERSION = "0.9.8"
const REQUIRE_UNITY_CDK_VERSION = "0.9.8"
const REQUIRE_UE4_CDK_VERSION = "0.9.8"

const FrameVersion = 19
const PropKeyLegacy = "LEGACY"
const PropKeyLegacyLobby = "LEGACY_LOBBY"
const PropKeyGenericPrefix = "OR_SHARE_PROP_"
const PropKeyPlayerPrefix = "OR_PLAYER_PROP_"
const RowSeparator = ";;"
const ColSeparator = "::"

func NewGuid() ([16]byte, error) {
	uuid := [16]byte{}
	_, err := rand.Read(uuid[:])
	if err != nil {
		return uuid, err
	}
	return uuid, nil
}

func GuidFormatString(guid [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", guid[0:4], guid[4:6], guid[6:8], guid[8:10], guid[10:])
}

// TODO implement
//func ValidatePorts(ports string) {
//}

//func ToBytes(Hashtable data) []byte {
//        return nil
//}

func ToExplodeBytes(list []string) []byte {
	return nil
}

func ToStringSlice(bytes []byte) []string {
	var strs []string
	for _, elem := range strings.Split(string(bytes), defs.RowSeparator) {
		strs = append(strs, elem)
	}
	return strs
}
