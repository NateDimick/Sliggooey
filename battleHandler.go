package main

import "encoding/json"

// details functions to aid battleParser

// unmarshals battle request and forwards to front end
func (a *App) handleBattleRequest(roomId, request string) {
	req := new(BattleRequest)
	err := json.Unmarshal([]byte(request), req)
	if err != nil {
		a.channels.errorChan <- JsonUnmarshalError
		return
	}
	a.channels.frontendChan <- ShowdownEvent{BattleRequestTopic, req}
}
