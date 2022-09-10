package main

import (
	"fmt"
	"strings"
)

// exposed functions for the frontend to use
func (a *App) SDLogin(uname, pword string) {
	a.state.credentials = &ShowdownCredentials{uname, pword}
	a.login()
}

func (a *App) SDLogout() {
	a.state = new(ShowdownState)
	a.state.loggedIn = false
	a.logout()
}

func (a *App) GetFormats() FormatInfo {
	return *a.state.formatList
}

func (a *App) SendPM(to, message string) {
	a.conn.SendServerCommand(buildCommand(SendMessage, to, message))
}

func (a *App) SendRoomChat(roomId, message string) {
	a.conn.SendServerMessageToRoom(roomId, message)
}

func (a *App) MakeBattleChoice(choices ...BattleChoice) {
	cmd := FormatBattleChoices(choices...)
	a.conn.SendServerCommand(cmd)

}

func FormatBattleChoices(choices ...BattleChoice) string {
	fmtChoices := make([]string, len(choices))
	for i, c := range choices {
		fmtChoices[i] = c.format()
	}
	return buildCommand(Choose, fmtChoices...)
}

func (a *App) SearchForBattle(formatName, team string) {

}

func (A *App) CancelSearchForBattle() {

}

func (a *App) SendBattleChallengeToUser(user, format, team string) {
	goPrint("sending", format, "challenge to", user, "with team", team)
	a.conn.SendServerCommand(buildCommand(UseTeam, team))
	a.conn.SendServerCommand(buildCommand(Challenge, user, format))
}

func (a *App) CancelBattleChallengeToUser(user string) {
	goPrint("cancelling challenge to", user)
	a.conn.SendServerCommand(buildCommand(CancelChallenge, user))
}

func (a *App) AcceptBattleChallengeFromUser(user, team string) {
	goPrint("accepting challenge from", user, "with team", team)
	a.conn.SendServerCommand(buildCommand(UseTeam, team))
	a.conn.SendServerCommand(buildCommand(AcceptChallenge, user))
}

func (a *App) RejectBattleChallengeFromUser(user string) {
	goPrint("rejecting challenge from", user)
	a.conn.SendServerCommand(buildCommand(RejectChallenge, user))
}

func buildCommand(cmd ChatCommand, parts ...string) string {
	return fmt.Sprintf("%s %s", cmd, strings.Join(parts, ", "))
}
