package main

// for parsing command messages
func (a *App) parseServerCommand(from string, to string, cmdMsg string) {
	fromUser := NewUser(from)
	toUser := NewUser(to)
	if toUser.UserName != a.state.user.User.UserName {
		fromUser = toUser
		toUser = a.state.user.User
	}
	cmd := NewSplitString(cmdMsg, " ")
	commandType := ChatCommand(cmd.Get(0))
	switch commandType {
	// TODO decide if there are worthy command parsing cases
	case Challenge:
		// /challenge format name|format name|?|?|?
		//                 0     |     1     |2|3|4
		// user has been challenged - or the challenge was cancelled. prompt them to accept or decline
		challengeArgs := NewSplitServerMessage(cmd.ReassembleTail(1))
		if challengeArgs.Get(0) != "" {
			goPrint("challenged by", fromUser.UserName, "to a", challengeArgs.Get(0), "battle!")
		} else {
			goPrint(toUser.UserName, "cancelled their challenge")
		}
	case Log, CmdError, Text, NoNotify:
		// message is going to the chat but as a special log message
		pm := PrivateMessagePayload{fromUser.UserName, toUser.UserName, cmd.ReassembleTail(1), System}
		a.channels.frontendChan <- ShowdownEvent{PMTopic, pm}
	default:
		goPrint(cmd.ReassembleTail(0), "is not a special chat command")
		// TODO unsure
	}

}
