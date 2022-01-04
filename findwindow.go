//go:build windows

/*
Package w32find provides a set of interface to win32 APIs that can be used to find windows and their controls.
*/
package w32find

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
	procFindWindowW    = user32.MustFindProc("FindWindowW")
	procFindWindowExW  = user32.MustFindProc("FindWindowExW")
)

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindWindowW(className, windowName *uint16) syscall.Handle {
	ret, _, _ := syscall.Syscall(procFindWindowW.Addr(), 2,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		0)
	return syscall.Handle(ret)
}

func FindWindowExW(hwndParent, hwndChildAfter syscall.Handle, className, windowName *uint16) syscall.Handle {
	ret, _, _ := procFindWindowExW.Call(
		uintptr(hwndParent),
		uintptr(hwndChildAfter),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)))
	return syscall.Handle(ret)
}

// FindWindowFromEnum finds a window from user32!EnumWindows().
//
// 	h1, err1 := FindWindowFromEnum("Your parent window title")
//
func FindWindowFromEnum(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	EnumWindows(syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	}), 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

// FindWindow finds hwnd by name.
//
// 	hwnd := FindWindow(nil, "Your window title")
//
func FindWindow(className, windowName string) (syscall.Handle, error) {
	var (
		ptrStrClassName  *uint16 = nil
		ptrStrWindowName *uint16 = nil
	)
	if className != "" {
		var err error
		ptrStrClassName, err = syscall.UTF16PtrFromString(className)
		if err != nil {
			return 0, err
		}
	}
	if windowName != "" {
		var err error
		ptrStrWindowName, err = syscall.UTF16PtrFromString(windowName)
		if err != nil {
			return 0, err
		}
	}
	h := FindWindowW(ptrStrClassName, ptrStrWindowName)
	if h == 0 {
		return 0, errors.New("Window not found")
	}
	return h, nil
}

// FindWindowEx finds a child window.
//
// 	h2, err2 := FindWindowEx(h1, 0, "", "Your child window title")
//
func FindWindowEx(hwndParent, hwndChildAfter syscall.Handle, className, windowName string) (syscall.Handle, error) {
	var (
		ptrStrClassName  *uint16 = nil
		ptrStrWindowName *uint16 = nil
	)
	if className != "" {
		var err error
		ptrStrClassName, err = syscall.UTF16PtrFromString(className)
		if err != nil {
			return 0, err
		}
	}
	if windowName != "" {
		var err error
		ptrStrWindowName, err = syscall.UTF16PtrFromString(windowName)
		if err != nil {
			return 0, err
		}
	}
	h := FindWindowExW(hwndParent, hwndChildAfter, ptrStrClassName, ptrStrWindowName)
	if h == 0 {
		return 0, errors.New("Window not found")
	}
	return h, nil
}
