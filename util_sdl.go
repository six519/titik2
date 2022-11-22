//go:build !win && sdl
// +build !win,sdl

package main

import (
	"database/sql"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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
	sdlWindow                map[string]*sdl.Window
	sdlSurface               map[string]*sdl.Surface
	sdlRect                  map[string]sdl.Rect
	sdlEvent                 map[string]sdl.Event
	sdlFont                  map[string]*ttf.Font
	sdlMusic                 map[string]*mix.Music
	sdlChunk                 map[string]*mix.Chunk
	sdlRenderer              map[string]*sdl.Renderer
	sdlTexture               map[string]*sdl.Texture
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
	globalSettings.sdlWindow = make(map[string]*sdl.Window)
	globalSettings.sdlSurface = make(map[string]*sdl.Surface)
	globalSettings.sdlRect = make(map[string]sdl.Rect)
	globalSettings.sdlEvent = make(map[string]sdl.Event)
	globalSettings.sdlFont = make(map[string]*ttf.Font)
	globalSettings.sdlMusic = make(map[string]*mix.Music)
	globalSettings.sdlChunk = make(map[string]*mix.Chunk)
	globalSettings.sdlRenderer = make(map[string]*sdl.Renderer)
	globalSettings.sdlTexture = make(map[string]*sdl.Texture)

}
