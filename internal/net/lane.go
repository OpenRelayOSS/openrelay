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
	"net"
	"runtime"
	"time"
)

const SUBSCRIBERS_PEEK = 65535

//const MSG_QUEUE_PEEK = 65535
type SubscriberIndex uint16
type PlayerId uint16

type Client struct {
	id   PlayerId
	conn *net.Conn
}

type Subscriber struct {
	index SubscriberIndex
	prev  *Subscriber
	next  *Subscriber
	cli   *Client
}

type Lane struct {
	Name       string
	Id         byte
	entryPoint *Subscriber
	outPoint   *Subscriber
	newIndex   SubscriberIndex
	chunks     [SUBSCRIBERS_PEEK]Subscriber
	avails     []SubscriberIndex
	adds       []*Client
	marks      []*Client
	msgQueue   []*[]byte
}

func NewLane(laneId byte, laneName string) *Lane {
	initialChunks := [SUBSCRIBERS_PEEK]Subscriber{}
	initialAvails := make([]SubscriberIndex, 0)
	for index := SubscriberIndex(0); index < SUBSCRIBERS_PEEK; index++ {
		initialChunks[index] = Subscriber{index, nil, nil, nil}
		initialAvails = append(initialAvails, index)
	}
	return &Lane{
		Name:       laneName,
		Id:         laneId,
		entryPoint: nil,
		outPoint:   nil,
		newIndex:   0,
		chunks:     initialChunks,
		avails:     nil,
		adds:       nil,
		marks:      nil,
		msgQueue:   nil,
	}
}

//atomic
func (s *Lane) Add(cli *Client) {
	// lock // id check
	s.adds = append(s.adds, cli)
}

//atomic
func (s *Lane) Remove(cli *Client) {
	// lock // id check
	s.marks = append(s.marks, cli)
}

func (s *Lane) Broadcast(msg *[]byte) {
	s.msgQueue = append(s.msgQueue, msg)
}

func (s *Lane) newChunk() *Subscriber {
	index := s.newIndex
	s.newIndex += 1
	s.avails = append(s.avails[:index], s.avails[index+1:]...)
	return &s.chunks[index]
}

func (s *Lane) releaseChunk(index SubscriberIndex) {
	s.chunks[index].next = nil
	s.chunks[index].cli = nil
	// quick recycle
	newAvails := []SubscriberIndex{index}
	newAvails = append(newAvails, s.avails...)
	s.avails = newAvails // or quick recycle?

	// latest recycle
	//s.avails = append(s.avails, index)
}

func (s *Lane) existsCli(cli *Client, index *SubscriberIndex) bool {
	sub := s.entryPoint
	for sub != nil {
		if sub.cli.id == cli.id {
			index = &sub.index
			return true
		}
		sub = sub.next
		runtime.Gosched()
	}
	index = nil
	return false
}

func (s *Lane) maintenanceLoop() {
	for {
		for 0 < len(s.adds) {
			// TODO at first another logic..
			// TODO set entryPoint

			addCli := s.adds[0]
			s.adds = s.adds[1:]
			if s.existsCli(addCli, nil) {
				continue
			}
			sub := s.newChunk()
			sub.cli = addCli
			sub.next = nil
			sub.prev = s.outPoint
			// acvivate
			s.outPoint.next = sub
			// change point
			s.outPoint = sub
		}
		for 0 < len(s.marks) {
			// TODO at first another logic..

			markCli := s.marks[0]
			s.marks = s.marks[1:]
			var index SubscriberIndex
			if !s.existsCli(markCli, &index) {
				continue
			}
			// TODO set entryPoint
			// TODO set outPoint
			// purged
			s.chunks[index].prev.next = s.chunks[index].next
			s.chunks[index].prev = nil
			s.releaseChunk(index)
		}
		time.Sleep(0 * time.Second) // force return context
		// runtime.Gosched()
	}
}

func (s *Lane) broadcastLoop() {
	sub := s.entryPoint
	for {
		for 0 < len(s.msgQueue) {
			msg := s.msgQueue[0]
			s.msgQueue = s.msgQueue[1:]
			for sub != nil {
				fmt.Printf("msg %s -> cli", *msg, sub.cli.id)
				_, err := (*sub.cli.conn).Write(*msg)
				if err != nil {
					fmt.Printf("Failed to write message to %s: %v\n", (*sub.cli.conn).RemoteAddr(), err)
				}
				sub = sub.next
				runtime.Gosched()
			}
			sub = s.entryPoint
			runtime.Gosched()
		}
		time.Sleep(0 * time.Second) // force return context
		// runtime.Gosched()
	}
}
