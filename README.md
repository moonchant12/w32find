# w32find
[![Go Reference](https://pkg.go.dev/badge/github.com/moonchant12/w32find)](https://pkg.go.dev/github.com/moonchant12/w32find)

Package w32find provides a set of interface to win32 APIs that can be used to find windows and their controls.

## Install
```cmd
go get -v github.com/moonchant12/w32find
```

## Import
```Go
import "github.com/moonchant12/w32find"
```

## Usage
```Go
hwnd1, err1 := FindWindowFromEnum("Your parent window title")
hwnd2, err2 := FindWindowEx(hwnd1, 0, "", "Your child window title")
```
