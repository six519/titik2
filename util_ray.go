//go:build !win && !sdl && ray
// +build !win,!sdl,ray

package main

import (
	"database/sql"
	"github.com/gen2brain/raylib-go/raylib"
	"net"
	"os"
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
	rayTexture               map[string]rl.Texture2D
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
	globalSettings.rayTexture = make(map[string]rl.Texture2D)

}
