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

package relay

import (
	"fmt"
)

const bufSize = 8192

// Receiver is a helper to handle one to many chat
type Receiver struct {
	clients map[PlayerId]Client
}

// NewReceiver builds a new hub
func NewReceiver() *Receiver {
	return &Receiver{clients: make(map[PlayerId]Client)}
}

// Register adds a new conn to the Receiver
func (r *Receiver) Add(cli Client) {
	fmt.Printf("Connected to %s\n", (*cli.conn).RemoteAddr())
	fmt.Printf("Connected to %v\n", cli.id)
	r.clients[cli.id] = cli

	go r.readLoop(cli)
}

func (r *Receiver) Remove(cli Client) {
	delete(r.clients, cli.id)
	err := (*cli.conn).Close()
	if err != nil {
		fmt.Println("Failed to disconnect %s", (*cli.conn).RemoteAddr())
		fmt.Println("Failed to disconnect %v %s", cli.id, err)
	} else {
		fmt.Println("Disconnected ", (*cli.conn).RemoteAddr())
	}
}

func (r *Receiver) readLoop(cli Client) {
	b := make([]byte, bufSize)
	for {
		n, err := (*cli.conn).Read(b)
		if err != nil {
			r.Remove(cli)
			return
		}
		fmt.Printf("Got message: %s\n", string(b[:n]))
	}
}

// TODO need queing message logic.
// TODO need Client dispose logic.
