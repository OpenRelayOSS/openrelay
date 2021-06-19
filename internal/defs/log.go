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
	"io"
	"log"
	"os"
	"runtime"
)

const FileSuffix = ".log"
const ServiceLogFilePrefix = "service"
const RelayLogFilePrefix = "relay"
const RelayRecFilePrefix = "record"
const stackDepth = 2

type LogLevel byte

const (
	NONE LogLevel = iota
	INFO
	NOTICE
	VERBOSE
	VVERBOSE
)

const (
	CALLIN =  "< callin  > "
	CALLOUT = "< callout > "
	WATCH =   "< watch   > "
)

type Logger struct {
	logger    *log.Logger
	logVolume LogLevel
	prefix    string
	file      *os.File
}

type Recorder struct {
	recorder *log.Logger
	file     *os.File
}

func NewLogger(lv LogLevel, dir string, filename string, needStdout bool) (*Logger, error) {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	file, err := os.OpenFile(dir+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return nil, err
	}
	if needStdout {
		logger.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		logger.SetOutput(file)
	}
	return &Logger{logger, lv, "", file}, nil
}

func (l *Logger) PrintRaw(lv LogLevel, rawString string) {
	if lv <= l.logVolume {
		l.logger.SetFlags(0)
		l.logger.Output(stackDepth, rawString)
		l.logger.SetFlags(log.LstdFlags|log.Lmicroseconds)
	}
}

func (l *Logger) Printf(lv LogLevel, format string, v ...interface{}) {
	if lv <= l.logVolume {
		l.logger.Output(stackDepth, levelToStr(lv) + l.prefix+ " | " + fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Println(lv LogLevel, v ...interface{}) {
	if lv <= l.logVolume {
		l.logger.Output(stackDepth, levelToStr(lv) + l.prefix + " | " + fmt.Sprintln(v...))
	}
}

func (l *Logger) Error(v ...interface{}) {
	l.logger.Output(stackDepth, l.prefix+"| ERROR   | "+fmt.Sprintln(v...))
	l.printStacktrace(stackDepth + 1)
}

func (l *Logger) Panic(v ...interface{}) {
	l.logger.Output(stackDepth, l.prefix+"| PANIC   | "+fmt.Sprintln(v...))
	l.printStacktrace(stackDepth + 1)
	panic(l.prefix + " CALLED PANIC.")
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logger.Output(stackDepth, l.prefix+"| FATAL   | "+fmt.Sprintln(v...))
	l.printStacktrace(stackDepth + 1)
	os.Exit(1)
}

func (l *Logger) SetPrefix(p string) {
	l.prefix = p
}

func (l *Logger) MuteStdout() {
	l.logger.SetOutput(l.file)
}

func (l *Logger) UnmuteStdout() {
	l.logger.SetOutput(io.MultiWriter(l.file, os.Stdout))
}

func (l *Logger) printStacktrace(stackDepth int) {
	stackMax := 20
	for stack := 0; stack < stackMax; stack++ {
		if stack < stackDepth {
			continue
		}
		point, file, line, ok := runtime.Caller(stack)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(point).Name()
		l.logger.Printf(l.prefix+"[STACKTRACE] file=%s, line=%d, func=%v\n", file, line, funcName)
	}
}

func levelToStr(lv LogLevel) string {
	lvStr := ""
	switch lv {
	case INFO:
		lvStr = "| INFO     "
	case NOTICE:
		lvStr = "| NOTICE   "
	case VERBOSE:
		lvStr = "| VERBOSE  "
	case VVERBOSE:
		lvStr = "| VVERBOSE "
	default:
		lvStr = "| NONE     "
	}
	return lvStr
}

func (l *Logger) Rotate() {
	// TODO Rotate and Truncate logic here.
}

func (l *Logger) Close() {
	l.file.Close()
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

func (r *Recorder) Printf(format string, v ...interface{}) {
	r.recorder.Printf(format, v...)
}

func (r *Recorder) Close() {
	r.file.Close()
}
