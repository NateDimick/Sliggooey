# ToDo

last modified: 2/22/2023

this is Now slightly more organized list of actionable items that need to get done eventually, but in the near time horizon

## Blocking Milestone 0

Milestone 0: Can play current gen random battle, initiated through a challenge (either direction)

* Complete Battle Simulator UI functionality (not appearance, just functionality)
  * send messages to battle room chat for every event
  * display field/side state
  * gimmick switch (mega/dmax/tera)
  * consolidate state info from request into room state (pokemon items, abilities on player side)
  * bench sprites for pokemon like Arceus-Ghost
* finish handling room messages (may not need to happen until after milestone 0)
  * notify - back and front
  * chatTs - ?
  * : - ?
  * uhtml + **uhtmlchange** in the front

## Nice to haves

Potential Open source Contributor who is perusing this project, this list of items is perfect for you to start contributing if you feel compelled to!

* Websocket work
  * decide if there is a need to replace Gorilla websocket with [nhooyr websocket](https://github.com/nhooyr/websocket) (RIP Gorilla toolkit)
  * handle disconnects and reconnects gracefully
* general backend work
  * update to to more recent go version
  * Reduce code duplication in Message parsers
  * proper logging (replace goprint)
  * recovery on major goroutines (websocket listener, frontend event passer, error handler)
  * more unit testing!
  * organize code (perhaps sub modules)
* general frontend work
  * styling
  * hit enter button to submit text fields (chat, login)
  * chat self
  * more unit testing!
  * refactor util.ts - it's getting too large
  * handle errors for invalid choices (will be difficult to test as long as the ui prevents invalid choices from being made)
  * figure out how to eliminate relative pathing in frontend imports and instead just import from `src/store` or `wailsjs/...`
* general work
  * application build versioning
  * remember user login info
  * ci
    * determine if worth before implementing
    * gitlab actions, ideally
  * cd
    * again, determine if worth
    * only once app is available and distributable

## Blocking Milestone 1

Milestone 1: Can play any random battle format

* Frontend-Focused work
  * Battle Pane Manager
    * Battle Search

* Backend-Focused Work
  * support battle search

## Blocking Milestone 2

Milestone 2: Feature parity with Pokemon Showdown client

* Team Builder
  * Ui to build teams
  * Can Validate teams
  * Can Save Teams to local disk in either pokepaste or json format
* Team Loader
  * Can load a team from either pokepaste or json format
  * paste team in an immediately render the result on the page (look kinda like pokepaste)
  * edit paste and see updates
  * Save team locally
* Team Packer
  * Can convert loaded team to packed format for the wss protocol
* Battle Input for non-singles formats
