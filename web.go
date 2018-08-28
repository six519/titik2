package main

type WebObject struct {
	IsProcessing bool
}

func (webObject *WebObject) Init() {
	webObject.IsProcessing = false
}