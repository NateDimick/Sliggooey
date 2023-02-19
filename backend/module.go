package backend

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func (a *App) JoinRoom(roomId string) {
	a.conn.SendServerCommand(buildCommand(JoinCmd, roomId))
}

func (a *App) LeaveRoom(roomId string) {
	a.conn.SendServerCommand(buildCommand(LeaveCmd, roomId))
}

func (a *App) MakeBattleChoice(roomId string, requestId int, choices []BattleChoice) {
	cmd := FormatBattleChoices(choices...)
	cmd = cmd + fmt.Sprintf("|%d", requestId)
	a.conn.SendServerMessageToRoom(roomId, cmd)

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

func (a *App) JunkyHackyFunctionDoNotUse(debug string) string {
	return fmt.Sprintf("echo! things are okay?")
}

//by reconciling room state in the backend, we can use go structs for in the frontend thanks to wails magic.
//define a struct once, use in both front and back end. it saves work
func (a *App) ReconcileRoomState(updatePayload UpdateRoomStatePayload, presentState RoomState) RoomState {
	//goPrint(fmt.Sprintf("reconciling room state: \n[base] %+v  \n[update] %+v", presentState, updatePayload))
	result := reconcileRoomStateInner(updatePayload, presentState)
	goPrint("room state reconciled")
	return result
}

// get the pokedex json, unparsed, from showdown
func (a *App) GetPokedex() string {
	resp, err := http.Get("https://play.pokemonshowdown.com/data/pokedex.json")
	if err != nil {
		goPrint(err)
		return ""
	}
	defer resp.Body.Close()
	dex, err := io.ReadAll(resp.Body)
	if err != nil {
		goPrint(err)
		return ""
	}
	goPrint("got pokedex JSON")
	return string(dex)
}

func (a *App) GetSpecialSpriteNumbers() map[string]int {
	specialIndex := make(map[string]int)
	resp, err := http.Get("https://play.pokemonshowdown.com/src/battle-dex-data.ts")
	if err != nil {
		goPrint(err)
		return specialIndex
	}
	// okay hear me out - going to parse the source ts file for the BattlePokemonIconIndexes map, and manually parse and extract the values
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	// read up to map definition
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "BattlePokemonIconIndexes") {
			break
		}
	}
	// parse the definitions
	for scanner.Scan() {
		entry := scanner.Text()
		goPrint("reading line ", entry)
		if strings.ContainsRune(entry, '}') {
			// end of the map, end of parsing
			break
		}
		parts := strings.Split(entry, ":")
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		index := 0
		sum := strings.ReplaceAll(parts[1], ",", "")
		for _, n := range strings.Split(sum, "+") {
			num, err := strconv.Atoi(strings.TrimSpace(n))
			if err != nil {
				goPrint("could not parse number", n, "in line", entry)
				continue
			}

			index += num
		}

		specialIndex[key] = index
	}
	goPrint("extracted special icon indices from battle dex code")
	return specialIndex
}
