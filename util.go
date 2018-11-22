package main

import (
	"strings"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	/*
	//IF WINDOWS
	"runtime"
	"syscall"
	"unsafe"
	//END IF WINDOWS
	*/
)

type GlobalSettingsObject struct {
	webObject WebObject
	globalVariableArray *[]Variable
	globalFunctionArray *[]Function
	globalNativeVarList *[]string
	mySQLSettings map[string]map[string]string
	mySQLResults map[string]map[string][]string //NOTE: TEMPORARY ONLY
	sQLiteSettings map[string]map[string]string
	sQLiteResults map[string]map[string][]string

	/*
	//IF WINDOWS
	consoleInfo CONSOLE_SCREEN_BUFFER_INFO //for windows only
	//END IF WINDOWS
	*/
}

func (globalSettings *GlobalSettingsObject) Init(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string) {
	globalSettings.globalVariableArray = globalVariableArray
	globalSettings.globalFunctionArray = globalFunctionArray
	globalSettings.globalNativeVarList = globalNativeVarList
	globalSettings.webObject = WebObject{}
	globalSettings.webObject.Init(globalSettings)

	globalSettings.mySQLSettings = make(map[string]map[string]string) //TODO: NEED WAY TO CLEAN THIS UP //MAYBE END OF FUNCTION CALLS?
	globalSettings.mySQLResults = make(map[string]map[string][]string) //TODO: NEED WAY TO CLEAN THIS UP //MAYBE END OF FUNCTION CALLS?
	
	globalSettings.sQLiteSettings = make(map[string]map[string]string)
	globalSettings.sQLiteResults = make(map[string]map[string][]string)

	/*
	//IF WINDOWS
	if runtime.GOOS == "windows" {
		//get console handle
		//for windows
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		getConsoleScreenBufferInfoProc := kernel32.NewProc("GetConsoleScreenBufferInfo")
		handle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		_, _, _ = getConsoleScreenBufferInfoProc.Call(uintptr(handle), uintptr(unsafe.Pointer(&globalSettings.consoleInfo)), 0)
	}
	//END IF WINDOWS
	*/
}

func escapeString(str string) string {
	var retStr string

	//escape newline
	retStr = strings.Replace(str, "\\n", "\n", -1)

	//escape carriage return
	retStr = strings.Replace(retStr, "\\r", "\r", -1)

	//escape tab
	retStr = strings.Replace(retStr, "\\t", "\t", -1)

	//escape form feed
	retStr = strings.Replace(retStr, "\\f", "\f", -1)

	//escape bell
	retStr = strings.Replace(retStr, "\\a", "\a", -1)

	//escape backspace
	retStr = strings.Replace(retStr, "\\b", "\b", -1)

	return retStr
}

func unescapeString(str string) string {
	var retStr string

	//escape newline
	retStr = strings.Replace(str, "\n", "\\n", -1)

	//escape carriage return
	retStr = strings.Replace(retStr, "\r", "\\r", -1)

	//escape tab
	retStr = strings.Replace(retStr, "\t", "\\t", -1)

	//escape form feed
	retStr = strings.Replace(retStr, "\f", "\\f", -1)

	//escape bell
	retStr = strings.Replace(retStr, "\a", "\\a", -1)

	//escape backspace
	retStr = strings.Replace(retStr, "\b", "\\b", -1)

	return retStr
}

//Copy code from: https://github.com/otiai10/copy/
func FDCopy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return fdcopy(src, dest, info)
}

func fdcopy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dest, info)
	}
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

func fcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

func dcopy(srcdir, destdir string, info os.FileInfo) error {

	if err := os.MkdirAll(destdir, info.Mode()); err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := fdcopy(cs, cd, content); err != nil {
			return err
		}
	}
	return nil
}

func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}
//end of Copy code