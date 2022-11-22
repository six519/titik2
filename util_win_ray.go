//go:build win && !sdl && ray
// +build win,!sdl,ray

package main

import (
	"database/sql"
	"net"
	"os"

	"github.com/gen2brain/raylib-go/raylib"
	"runtime"
	"syscall"
	"unsafe"
)

type GlobalSettingsObject struct {
	webObject                WebObject
	globalVariableArray      *[]Variable
	globalFunctionArray      *[]Function
	globalNativeVarList      *[]string
	mySQLResults             map[string]map[string][]string //NOTE: TEMPORARY ONLY
	sQLiteSettings           map[string]map[string]string
	sQLiteResults            map[string]map[string][]string
	netConnection            map[string]net.Conn
	netConnectionListener    map[string]net.Listener
	netUDPConnectionListener map[string]*net.UDPConn
	mySQLConnection          map[string]*sql.DB
	fileHandler              map[string]*os.File
	rayImage                 map[string]*rl.Image

	consoleInfo CONSOLE_SCREEN_BUFFER_INFO //for windows only
}

func (globalSettings *GlobalSettingsObject) Init(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string) {
	globalSettings.globalVariableArray = globalVariableArray
	globalSettings.globalFunctionArray = globalFunctionArray
	globalSettings.globalNativeVarList = globalNativeVarList
	globalSettings.webObject = WebObject{}
	globalSettings.webObject.Init(globalSettings)

	globalSettings.mySQLResults = make(map[string]map[string][]string) //TODO: NEED WAY TO CLEAN THIS UP //MAYBE END OF FUNCTION CALLS?

	globalSettings.sQLiteSettings = make(map[string]map[string]string)
	globalSettings.sQLiteResults = make(map[string]map[string][]string)

	globalSettings.netConnection = make(map[string]net.Conn)
	globalSettings.netConnectionListener = make(map[string]net.Listener)
	globalSettings.netUDPConnectionListener = make(map[string]*net.UDPConn)
	globalSettings.mySQLConnection = make(map[string]*sql.DB)
	globalSettings.fileHandler = make(map[string]*os.File)
	globalSettings.rayImage = make(map[string]*rl.Image)

	if runtime.GOOS == "windows" {
		//get console handle
		//for windows
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		getConsoleScreenBufferInfoProc := kernel32.NewProc("GetConsoleScreenBufferInfo")
		handle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		_, _, _ = getConsoleScreenBufferInfoProc.Call(uintptr(handle), uintptr(unsafe.Pointer(&globalSettings.consoleInfo)), 0)
	}
}
