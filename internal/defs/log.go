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
	"io"
	"log"
	"os"
	"fmt"
)

const FileSuffix = ".log"
const ServiceLogFilePrefix = "service"
const RelayLogFilePrefix = "relay"
const RelayRecFilePrefix = "record"
const callDepth = 2

type LogLevel byte

const (
	NONE LogLevel = iota
	ERRORONLY
	INFO
	VERBOSE
	VVERBOSE
)

type Logger struct {
	logger    *log.Logger
	logVolume LogLevel
	prefix string
	file *os.File
}

type Recorder struct {
	recorder *log.Logger
	file *os.File
}

func NewLogger(lv LogLevel, dir string, filename string) (*Logger, error) {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	file, err := os.OpenFile(dir+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return nil, err
	}
	logger.SetOutput(io.MultiWriter(file, os.Stdout))
	return &Logger{logger, lv, "", file}, nil
}

func (o *Logger) SetPrefix(p string) {
	o.prefix = p
}

func (o *Logger) Printf(lv LogLevel, format string, v ...interface{}) {
	if lv <= o.logVolume {
		o.logger.Output(callDepth, o.prefix + levelToStr(lv) + fmt.Sprintf(format, v...))
	}
}

func (o *Logger) Println(lv LogLevel, v ...interface{}) {
	if lv <= o.logVolume {
		o.logger.Output(callDepth, o.prefix + levelToStr(lv) + fmt.Sprintln(v...))
	}
}

func levelToStr(lv LogLevel) string {
	lvStr := ""
	switch lv {
		case NONE:
			lvStr = "[NONE] "
		case ERRORONLY:
			lvStr = "[ERRORONLY] "
		case INFO:
			lvStr = "[INFO] "
		case VERBOSE:
			lvStr = "[VERBOSE] "
		case VVERBOSE:
			lvStr = "[VVERBOSE] "
		default:
			lvStr = "[NONE] "
	}
	return lvStr
}

func (o *Logger) Rotate() {
	// TODO Rotate and Truncate logic here.
}

func (o *Logger) Close() {
	o.file.Close()
}

func NewRecorder(dir string, filename string) (*Recorder, error) {
	recorder := log.New(os.Stdout, "", 0)
	file, err := os.OpenFile(dir+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return nil, err
	}
//	defer file.Close()
	recorder.SetOutput(file)
	return &Recorder{recorder, file}, nil
}

func (o *Recorder) Printf(format string, v ...interface{}) {
	o.recorder.Printf(format, v...)
}

func (o *Recorder) Close() {
	o.file.Close()
}

