/*
 * Copyright (c) 2014 MongoDB, Inc.
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

package gowin32

import (
	"github.com/winlabs/gowin32/wrappers"

	"syscall"
	"unsafe"
)

type RegRoot syscall.Handle

const (
	RegRootHKCR RegRoot = wrappers.HKEY_CLASSES_ROOT
	RegRootHKCU RegRoot = wrappers.HKEY_CURRENT_USER
	RegRootHKLM RegRoot = wrappers.HKEY_LOCAL_MACHINE
	RegRootHKU  RegRoot = wrappers.HKEY_USERS
	RegRootHKPD RegRoot = wrappers.HKEY_PERFORMANCE_DATA
	RegRootHKCC RegRoot = wrappers.HKEY_CURRENT_CONFIG
	RegRootHKDD RegRoot = wrappers.HKEY_DYN_DATA
)

func DeleteRegValue(root RegRoot, subKey string, valueName string) error {
	var hKey syscall.Handle
	err := wrappers.RegOpenKeyEx(
		syscall.Handle(root),
		syscall.StringToUTF16Ptr(subKey),
		0,
		wrappers.KEY_WRITE,
		&hKey)
	if err != nil {
		return NewWindowsError("RegOpenKeyEx", err)
	}
	defer wrappers.RegCloseKey(hKey)
	if err := wrappers.RegDeleteValue(hKey, syscall.StringToUTF16Ptr(valueName)); err != nil {
		return NewWindowsError("RegDeleteValue", err)
	}
	return nil
}

func GetRegValueDWORD(root RegRoot, subKey string, valueName string) (uint32, error) {
	var hKey syscall.Handle
	err := wrappers.RegOpenKeyEx(
		syscall.Handle(root),
		syscall.StringToUTF16Ptr(subKey),
		0,
		wrappers.KEY_READ,
		&hKey)
	if err != nil {
		return 0, NewWindowsError("RegOpenKeyEx", err)
	}
	defer syscall.RegCloseKey(hKey)
	var valueType uint32
	var size uint32
	err = wrappers.RegQueryValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		&valueType,
		nil,
		&size)
	if err != nil {
		return 0, NewWindowsError("RegQueryValueEx", err)
	}
	if valueType != wrappers.REG_DWORD {
		// use the same error code as RegGetValue, although that function is not used here in order to maintain
		// compatibility with older versions of Windows
		return 0, wrappers.ERROR_UNSUPPORTED_TYPE
	}
	var value uint32
	err = wrappers.RegQueryValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		nil,
		(*byte)(unsafe.Pointer(&value)),
		&size)
	if err != nil {
		return 0, NewWindowsError("RegQueryValueEx", err)
	}
	return value, nil
}

func GetRegValueString(root RegRoot, subKey string, valueName string) (string, error) {
	var hKey syscall.Handle
	err := wrappers.RegOpenKeyEx(
		syscall.Handle(root),
		syscall.StringToUTF16Ptr(subKey),
		0,
		wrappers.KEY_READ,
		&hKey)
	if err != nil {
		return "", NewWindowsError("RegOpenKeyEx", err)
	}
	defer wrappers.RegCloseKey(hKey)
	var valueType uint32
	var size uint32
	err = wrappers.RegQueryValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		&valueType,
		nil,
		&size)
	if err != nil {
		return "", NewWindowsError("RegQueryValueEx", err)
	}
	if valueType != wrappers.REG_SZ {
		// use the same error code as RegGetValue, although that function is not used here in order to maintain
		// compatibility with older versions of Windows
		return "", wrappers.ERROR_UNSUPPORTED_TYPE
	}
	buf := make([]uint16, size/2)
	err = wrappers.RegQueryValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		nil,
		(*byte)(unsafe.Pointer(&buf[0])),
		&size)
	if err != nil {
		return "", NewWindowsError("RegQueryValueEx", err)
	}
	return syscall.UTF16ToString(buf), nil
}

func SetRegValueDWORD(root RegRoot, subKey string, valueName string, data uint32) error {
	var hKey syscall.Handle
	err := wrappers.RegCreateKeyEx(
		syscall.Handle(root),
		syscall.StringToUTF16Ptr(subKey),
		0,
		nil,
		0,
		wrappers.KEY_WRITE,
		nil,
		&hKey,
		nil)
	if err != nil {
		return NewWindowsError("RegCreateKeyEx", err)
	}
	defer wrappers.RegCloseKey(hKey)
	err = wrappers.RegSetValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		0,
		wrappers.REG_DWORD,
		(*byte)(unsafe.Pointer(&data)),
		uint32(unsafe.Sizeof(data)))
	if err != nil {
		return NewWindowsError("RegSetValueEx", err)
	}
	return nil
}

func SetRegValueString(root RegRoot, subKey string, valueName string, data string) error {
	var hKey syscall.Handle
	err := wrappers.RegCreateKeyEx(
		syscall.Handle(root),
		syscall.StringToUTF16Ptr(subKey),
		0,
		nil,
		0,
		wrappers.KEY_WRITE,
		nil,
		&hKey,
		nil)
	if err != nil {
		return NewWindowsError("RegCreateKeyEx", err)
	}
	defer wrappers.RegCloseKey(hKey)
	err = wrappers.RegSetValueEx(
		hKey,
		syscall.StringToUTF16Ptr(valueName),
		0,
		wrappers.REG_SZ,
		(*byte)(unsafe.Pointer(syscall.StringToUTF16Ptr(data))),
		uint32(2*(len(data)+1)))
	if err != nil {
		return NewWindowsError("RegSetValueEx", err)
	}
	return nil
}
