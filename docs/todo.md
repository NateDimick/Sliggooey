# ToDo

last modified: 12/13/2022

this is Now slightly more organized list of actionable items that need to get done eventually, but in the near time horizon

## Blocking Milestone 0

Milestone 0: Can play current gen random battle, initiated through a challenge (either direction)

* finish handling battle messages, outside of team preview and a few other specific messages
* finish handling room messages
  * notify - back and front
  * chatTs - ?
  * : - ?
  * uhtml + **uhtmlchange** in the front

## Nice to haves

Potential Open source Contributor who is perusing this project, this list of items is perfect for you to start contributing if you feel compelled to!

* Websocket work
  * Replace Gorilla websocket with nyhoor websocket (RIP Gorilla toolkit)
  * handle disconnects and reconnects gracefully
* general backend work
  * Reduce code duplication in Message parsers
  * proper logging (replace goprint)
  * recovery on major goroutines (websocket listener, frontend event passer, error handler)
  * more unit testing!
  * organize code (perhaps sub modules)
* general frontend work
  * styling
  * hit enter button to submit text fields 9chat, login)
  * chat self
  * more unit testing!
  * refactor util.ts - it's getting too large
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

## Blocking Milestone 2

Milestone 2: Feature parity with Pokemon Showdown client

* Team Builder
  * Ui to build teams
  * Can Validate teams
  * Can Save Teams to local disk in either pokepaste or json format
* Team Loader
  * Can load a team from either pokepaste or json format
* Team Packer
  * Can convert loaded team to packed format for the wss protocol
* Battle Input for non-singles formats
