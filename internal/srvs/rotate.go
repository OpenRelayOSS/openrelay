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
	"runtime"
	"time"
)

func (o *OpenRelay) CleanServ() {
	index := 0
	for {
		if len(o.ColdRoomQueue) > 0 {
			o.Recycle(index)
			o.ColdRoomQueue = o.ColdRoomQueue[1:]
		}
		runtime.Gosched()
		time.Sleep(1 * time.Second)
	}
}
