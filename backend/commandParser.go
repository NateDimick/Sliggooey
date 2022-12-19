package backend

// for parsing command messages
func (a *App) parseServerCommand(from string, to string, cmdMsg string) {
	fromUser := NewUser(from)
	toUser := NewUser(to)
	var withUser *User
	var source PmSource
	if fromUser.UserName == a.state.user.UserName {
		withUser = toUser
		source = Self
	} else {
		withUser = fromUser
		source = Other
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
			goPrint("challenge issued by", fromUser.UserName, "with", withUser.UserName, "to a", challengeArgs.Get(0), "battle!")
			payload := ChallengePayload{withUser.UserName, a.state.user.UserName, challengeArgs.Get(0), source}
			goPrint("sending challenge to frontend", payload)
			a.channels.frontendChan <- ShowdownEvent{ChallengeTopic, payload}
		} else {
			goPrint(toUser.UserName, "cancelled their challenge")
			a.channels.frontendChan <- ShowdownEvent{ChallengeEndTopic, withUser.UserName}
		}
	case Log, CmdError, Text, NoNotify:
		// message is going to the chat but as a special log message
		pm := PrivateMessagePayload{withUser.UserName, a.state.user.UserName, cmd.ReassembleTail(1), System}
		a.channels.frontendChan <- ShowdownEvent{PMTopic, pm}
	default:
		goPrint(cmd.ReassembleTail(0), "is not a special chat command")
		// TODO unsure
	}

}
