# Pokemon Sliggooey

The goal of this project is to build a Pokemon showdown client that offers a better user experience than the default client. By building in Go with Wails, it also will be available for Windows, Mac or Linux, whereas the default client is only available for Windows and Mac.

Improvements that this project seeks to add are:

* Better use of screen real estate
* Chat command hot keys
* Automatic good sportmanship (GGs and GLHFs)
* Custom avatars (?)
* Damage calculator integration @smogon/calc
* Pikalytics integration (?)
* Quick EV/IV presets in team builder (?)
* Looking good

## Timeline

Milestones will be as follows:

* Milestone 1: can play random battles
* Milestone 2: Damage Calculator integration
* Milestone 3: can play any format by loading a team (and maybe teambuilder too)
* Milestone 4: Multiple battles at once

## About

Gooey is a play on words of "GUI" (If you're not in the know, many people pronounce GUI as "Gooey" rather than "gee yoo eye") and also allows us to use the Goodra line as our mascot.

## Development Details

This project is built in [Go 1.18](go.dev) with a [Svelte](svelte.dev)-Typescript front end with [Wails version 2](https://wails.io/)

Contributions are welcome and appreciated. Make an issue, a PR or hit me up on Twitter if you have an idea for a feature to add.

See the Pokemon Showdown websocket protocol and api reference here: [Showdown Protocol](https://github.com/smogon/pokemon-showdown/blob/master/PROTOCOL.md)

* To run in development mode, `wails dev`
    1. just the front end can be run alone with `npm run dev` from the `frontend` directory
* To build Go code for front end use, `wails generate`
* To build an executable, `wails build`

### General Architecture

The websocket client sits in the Go backend. this allows it to be multi-threaded and more performant than if it were a front-end client.

the front-end should be display only. Anything that requires parsing should be handled in the backend. This helps keep the front-end focused on style and usability over code and maintainability. (the less Typescript the better! Lean in to what svelte and wails provide *hard*). Typescript was chosen to align with Go's stricter typing. It should never be difficult to tell what should go in our out of an IPC call to Go.

If the contents of a message need to go the front end, then should be prompted by emitting a wails event.

Sending messages is done through a direct function call.

#### Emitted Back End Events

* Popup events - create a popup
* private message events - go to a pm component
* challenge received events - go to a pm component?
* name taken events - create a popup event
* query response events - go to a pm component
* room events (all occur over the same event topic)
  * battle events
  * chat events
  * join events
  * leave events
  * notify events

* `loginSuccess`
  * to `Login.svelte`
  * indicates a successful user login upon receipt of `updateuser`
* `loginFail`
  * to `Login.svelte`
  * indicates user did not log in and the `Reason` message should be displayed
* `pm`
  * to `Client.svelte`
  * contains a pm to propagate to the appropriate Chat component
* `roomMsg`
  * to `Client.svelte`
  * contains a room message to propagate to the appropriate Battle/Chat room component  

#### UI Events

* `uiLogin`
  * from `Login.svelte` or `Register.svelte`
  * to `App.svelte`
  * changes the main app view from login page to client page
* `uiLogout`
  * from `Client.svelte`
  * to `App.svelte`
  * changes the main app view from client page to login page
* `uiRegister`
  * **not implemented yet**
  * from `Login.svelte`
  * to `App.svelte`
  * changes the mainapp view from the login page to the register page
* `uiPM#username`
  * from `Client.svelte`
  * to `Chat.svelte`
  * forwards a private message event to the chat component with the specified user
* `uiRoomMsg#roomId`
  * from `Client.svelte`
  * to `BattleChat.svelte`
  * sends a message to the in-room chat
