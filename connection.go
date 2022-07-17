package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// "low level" interface with pokemon showdown

var smogonUrl string = "wss://sim3.psim.us/showdown/websocket"

func (c *ShowdownConnection) connect() bool {
	wsDialer := new(websocket.Dialer)
	wsDialer.Jar = http.DefaultClient.Jar
	wsConn, httpResp, err := wsDialer.Dial(smogonUrl, nil)
	if err != nil {
		goPrint("websocket connect error", err)
		c.connected = false
	} else {
		goPrint("websocket connect http response", httpResp)

		c.conn = wsConn
		c.connected = true
	}
	return c.connected
}

func (c *ShowdownConnection) reconnect() {
	// TODO
	goPrint("reconnect error received")
}

func (c *ShowdownConnection) disconnect() {
	c.conn.Close()
	c.connected = false
	c.challengeString = ""
}

func (c *ShowdownConnection) SendServerMessage(msg string) {
	goPrint("sending message:", msg)
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *ShowdownConnection) SendServerCommand(cmd string) {
	c.SendServerMessage("|" + cmd)
}

func (c *ShowdownConnection) SendServerMessageToRoom(roomId string, msg string) {
	formattedMsg := roomId + "|" + msg
	c.SendServerMessage(formattedMsg)
}

func (a *App) login() {
	// TODO re-login with cookies
	url := "https://play.pokemonshowdown.com/~~showdown/action.php"
	postData := fmt.Sprintf("act=login&name=%s&pass=%s&challstr=%s", a.state.credentials.UserName, a.state.credentials.Password, a.conn.challengeString)
	goPrint(url)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		goPrint("error posting login info")
		a.channels.frontendChan <- ShowdownEvent{LoginFail, LoginFailurePayload{err.Error()}}
		return
	} else if resp.StatusCode != 200 {
		var errMsg []byte
		resp.Body.Read(errMsg)
		a.channels.frontendChan <- ShowdownEvent{LoginFail, LoginFailurePayload{string(errMsg)}}
		return
	}
	loginInfo := new(LoginResponse)
	body, err := ioutil.ReadAll(resp.Body)
	goPrint(string(body[1:]))
	if err != nil {
		goPrint("error reading login response")
		a.channels.frontendChan <- ShowdownEvent{LoginFail, LoginFailurePayload{err.Error()}}
		return
	}
	err = json.Unmarshal(body[1:], loginInfo)
	if err != nil {
		goPrint("error parsing login response")
		a.channels.frontendChan <- ShowdownEvent{LoginFail, LoginFailurePayload{err.Error()}}
		return
	}
	confirmationMsg := fmt.Sprintf("/trn %s,0,%s", a.state.credentials.UserName, loginInfo.Assertion)
	a.conn.SendServerCommand(confirmationMsg)
}

func (a *App) logout() {

}
