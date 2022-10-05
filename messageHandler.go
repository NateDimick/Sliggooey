package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (a *App) formatList(list string) {
	sectionPattern := "^,\\d$"
	var currentSection int
	var currentTitle string
	nextSectionIsTitle := false
	formats := make(FormatInfo, 1)
	for _, section := range strings.Split(list, "|") {
		fmt.Print(section)
		if nextSectionIsTitle {
			//goPrint(" is section title")
			nextSectionIsTitle = false
			currentTitle = section
			formats[currentSection][currentTitle] = make([]Format, 1)
			continue
		}
		isSectionHeader, _ := regexp.MatchString(sectionPattern, section)
		if isSectionHeader {
			//goPrint(" is section header")
			nextSectionIsTitle = true
			currentSection, _ = strconv.Atoi(string(section[1:]))
			formats = append(formats, make(FormatSection))
		} else {
			//goPrint(" is format name")
			chunkedFormat := strings.Split(section, ",")
			formats[currentSection][currentTitle] = append(formats[currentSection][currentTitle], Format{chunkedFormat[0], chunkedFormat[1]})
		}
	}
	//goPrint("parsed formats", formats)
	a.state.formatList = &formats
}

func (a *App) updateUserResultEvent(uname string, isGuest int, avatarId int, settings string) {
	if isGuest == 1 && a.state.loggedIn == false {
		goPrint("User logged in! notifying UI", uname, avatarId)
		a.channels.frontendChan <- ShowdownEvent{LoginSuccess, LoginSuccessPayload{strings.TrimSpace(uname)}}
	} else if isGuest == 0 {
		return
	}
	user := NewUser(uname)
	user.IsUser = true
	a.state.user = user
	// TODO send user settings json to the frontend to store as state
}

func (a *App) popupMessageEvent(message string) {
	formattedMsg := strings.ReplaceAll(message, "||", "\n")
	goPrint("Todo: emit popup event with this message", formattedMsg)
	// might require application context
}

func (a *App) userCountEvent(users int) {
	goPrint("todo wow! so many users!", users)
}

func (a *App) nameTakenEvent(takenName string, errReason string) {
	goPrint(takenName, "is taken. todo: try again.", errReason)
}

func (a *App) updateChallengeEvent(challJson string) {
	var status *ChallengeStatus
	err := json.Unmarshal([]byte(challJson), status)
	if err != nil {
		goPrint("could not parse challenge status", challJson)
		return
	}
	// TODO check for existing challenge status
	// if a new challenge appears, notify UI
	// if challenge is accepted, notify UI
	a.state.challenges = status
}

func (a *App) privateMessageEvent(from string, to string, message string) {
	var pm PrivateMessagePayload
	fromUser := NewUser(from)
	toUser := NewUser(to)
	goPrint("pm from", from, "to", to, message, "and I am", a.state.user.UserName)
	if fromUser.UserName == a.state.user.UserName {
		pm = PrivateMessagePayload{toUser.UserName, fromUser.UserName, message, Self}
	} else {
		pm = PrivateMessagePayload{fromUser.UserName, toUser.UserName, message, Other}
	}
	a.channels.frontendChan <- ShowdownEvent{PMTopic, pm}
}
