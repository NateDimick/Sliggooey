package main

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// describes the general structure of this backend in both data and routines

type App struct {
	ctx      context.Context
	conn     *ShowdownConnection
	state    *ShowdownState
	channels *ShowdownChannels
}

type ShowdownConnection struct {
	connected       bool
	conn            *websocket.Conn
	challengeString string
}

type ShowdownState struct {
	user        *User
	loggedIn    bool
	challenges  *ChallengeStatus
	settings    *AppSettings
	credentials *ShowdownCredentials
	formatList  *FormatInfo
}

type ShowdownChannels struct {
	frontendChan      chan ShowdownEvent
	serverMessageChan chan string
	errorChan         chan AppError
}

func NewApp() *App {
	return new(App)
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.conn = new(ShowdownConnection)
	a.state = new(ShowdownState)
	a.state.loggedIn = false
	a.channels = new(ShowdownChannels)
	a.channels.frontendChan = make(chan ShowdownEvent)
	a.channels.serverMessageChan = make(chan string)
	a.channels.errorChan = make(chan AppError)

	go a.errorResolver()
	go a.frontEndEventEmitter()
	go a.messageParser()
	ok := a.conn.connect()
	if ok {
		go a.websocketListener()
	}

}

func (a *App) frontEndEventEmitter() {
	for {
		event := <-a.channels.frontendChan
		runtime.EventsEmit(a.ctx, string(event.Topic), event.Data)
	}
}

func (a *App) websocketListener() {
	for a.conn.connected {
		msgType, msg, err := a.conn.conn.ReadMessage() // this blocks
		if err != nil {
			goPrint("wss message reader error:", err.Error())
			a.channels.errorChan <- WebsocketError
		}
		if msgType == websocket.TextMessage && a.conn.connected {
			a.channels.serverMessageChan <- string(msg)
		} else {
			goPrint("unexpected binary wss message:", msg)
		}
	}
}

func (a *App) messageParser() {
	for {
		msg := <-a.channels.serverMessageChan
		a.parseServerPayload(msg)
	}
}

func (a *App) errorResolver() {
	for {
		errType := <-a.channels.errorChan
		switch errType {
		case WebsocketError:
			a.conn.reconnect()
		}
	}
}
