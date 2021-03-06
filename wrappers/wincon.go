/*
 * Copyright (c) 2014-2019 MongoDB, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the license is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wrappers

import (
	"syscall"
	"unsafe"
)

const (
	FOREGROUND_BLUE      = 0x0001
	FOREGROUND_GREEN     = 0x0002
	FOREGROUND_RED       = 0x0004
	FOREGROUND_INTENSITY = 0x0008
	BACKGROUND_BLUE      = 0x0010
	BACKGROUND_GREEN     = 0x0020
	BACKGROUND_RED       = 0x0040
	BACKGROUND_INTENSITY = 0x0080
)

const (
	CTRL_C_EVENT        = 0
	CTRL_BREAK_EVENT    = 1
	CTRL_CLOSE_EVENT    = 2
	CTRL_LOGOFF_EVENT   = 5
	CTRL_SHUTDOWN_EVENT = 6
)

const (
	ENABLE_PROCESSED_INPUT        = 0x0001
	ENABLE_LINE_INPUT             = 0x0002
	ENABLE_ECHO_INPUT             = 0x0004
	ENABLE_WINDOW_INPUT           = 0x0008
	ENABLE_MOUSE_INPUT            = 0x0010
	ENABLE_INSERT_MODE            = 0x0020
	ENABLE_QUICK_EDIT_MODE        = 0x0040
	ENABLE_EXTENDED_FLAGS         = 0x0080
	ENABLE_VIRTUAL_TERMINAL_INPUT = 0x0200
)

const (
	ENABLE_PROCESSED_OUTPUT            = 0x0001
	ENABLE_WRAP_AT_EOL_OUTPUT          = 0x0002
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	DISABLE_NEWLINE_AUTO_RETURN        = 0x0008
	ENABLE_LVB_GRID_WORLDWIDE          = 0x0010
)

var (
	procGenerateConsoleCtrlEvent = modkernel32.NewProc("GenerateConsoleCtrlEvent")
	procGetConsoleMode = modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = modkernel32.NewProc("SetConsoleMode")
)

func GenerateConsoleCtrlEvent(ctrlEvent uint32, processGroupId uint32) error {
	r1, _, e1 := syscall.Syscall(
		procGenerateConsoleCtrlEvent.Addr(),
		2,
		uintptr(ctrlEvent),
		uintptr(processGroupId),
		0)
	if r1 == 0 {
		if e1 != ERROR_SUCCESS {
			return e1
		} else {
			return syscall.EINVAL
		}
	}
	return nil
}

func GetConsoleMode(consoleHandle syscall.Handle, mode *uint32) error {
	r1, _, e1 := syscall.Syscall(
		procGetConsoleMode.Addr(),
		2,
		uintptr(consoleHandle),
		uintptr(unsafe.Pointer(mode)),
		0)
	if r1 == 0 {
		if e1 != ERROR_SUCCESS {
			return e1
		} else {
			return syscall.EINVAL
		}
	}
	return nil
}

func SetConsoleMode(consoleHandle syscall.Handle, mode uint32) error {
	r1, _, e1 := syscall.Syscall(
		procSetConsoleMode.Addr(),
		2,
		uintptr(consoleHandle),
		uintptr(mode),
		0)
	if r1 == 0 {
		if e1 != ERROR_SUCCESS {
			return e1
		} else {
			return syscall.EINVAL
		}
	}
	return nil
}
