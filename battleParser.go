package main

// for parsing battle messages from server
func (a *App) parseBattleMessage(roomId string, msg *SplitString) {
	goPrint("incoming battle message")
	msgType := MessageType(msg.Get(0))
	switch msgType {

	case Player:
		//
	default:
		// TODO
	}
}
