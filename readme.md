# Pokemon Sliggooey

The goal of this project is to build a Pokemon showdown client that offers a better user experience than the default client.

By building in Go with Wails, we get the following benefits:

* Faster an simpler compiling
* Smaller Executables (showdown client for windows is 44 MB)
* Linux build target support
* No PHP

## Timeline

This projects seeks to be an improvement over the default showdown client. However, before it can surpass the default client, it must reach feature parity. The Roadmap below provides a general outline of the big milestones of this project.

Milestones will be as follows:

* Milestone 0: can play a full current gen random single battle initiated through a challenge
* Milestone 1: can play random battles
  * this means any format with random teams can be played, whether through direct challenges or through matchmaking on the ladder
  * fortunately. this step requires probably 95-99% support of the battle simulator protocol support, meaning future steps will not have to do much more work on supporting battles
* Milestone 2: can play any format and teambuilder
  * this is a big step. Any format must become playable and the teambuilder must be incorporated to help support that
  * Teambuidler includes team validation
  * Teambuilder includes team importing and exporting from pokepaste and json formats
  * teams must be packable to the WSS protocol packed format
* Milestone 3: Improvements over the default pokemon showdown client, which include (in no particular order)
  * Key bindings and shortcuts
  * Gamepad support
  * Damage Calculator integration
  * Automation in response to certain events
    * auto timer start
    * auto gl/hf/gg
    * custom automation event support
  * Appearance customization
  * Full range of player avatars available for selection
  * Team Builder Aid, such as:
    * Pikalytics integration (?)
    * EV/IV presets

## About

Sliggooey is a play on words of Sliggoo, the pokemon, and "GUI" (If you're not in the know, many people pronounce GUI as "Gooey" rather than "gee yoo eye").

## Development Details

This project is built in [Go 1.18](go.dev) with a [Svelte](svelte.dev)-Typescript front end with [Wails version 2](https://wails.io/)

Contributions are welcome and appreciated. Make an issue, a PR or hit me up on Twitter if you have an idea for a feature to add (or just want to help this get to version 1.0).

See the Pokemon Showdown websocket protocol and api reference here: [Showdown Protocol](https://github.com/smogon/pokemon-showdown/blob/master/PROTOCOL.md)

* To run in development mode, `wails dev` (**do this for first time setup too** - it will populate frontend/dist, frontend/wailsjs, and create the build/directory)
    1. just the front end can be run alone with `npm run dev` from the `frontend` directory
    2. While in `wails dev`, the front end is also accessible from `localhost:34115` in your browser
* To build Go code for front end use, `wails generate module`
* To build an executable, `wails build`

### Development Goals

1. try to unit test as much of the Go backend as possible
2. minimize TS code footprint (use svelte features and go backend as much as possible)
3. use Wails only for IPC (no events that don't cross IPC threshold)

### General Architecture

The websocket client sits in the Go backend. this allows it to be multi-threaded and more performant than if it were a front-end client.

the front-end should be display only. Anything that requires parsing should be handled in the backend. This helps keep the front-end focused on style and usability over code and maintainability. (the less Typescript the better! Lean in to what svelte and wails provide *hard*). Typescript was chosen to align with Go's stricter typing. It should never be difficult to tell what should go in our out of an IPC call to Go.

If the contents of a message need to go the front end, then should be prompted by emitting a wails event.

Sending messages is done through a direct function call.

#### Emitted Back End Events

Backend events are emitted to permanent components in the front end. If a component can be created or destroyed dynamically, a permanent front end component must receive the event and propagate it to the volatile component over a custom ui topic.

##### Backend->Frontend Event Topics

* `loginFail`
* `loginSuccess`
* `formats`
* `pm`
* `newRoom`
* `roomMsg`
* `popup`
* `challenged`
* `challengeEnd`

## Note Taking

Examples of how to get images of pokemon... from the showdown server

* Pokemon Sprites [Pokemon Spritesheet](https://play.pokemonshowdown.com/sprites/pokemonicons-sheet.png)
  * Have to do math based on pokedex number or something similar to get the right spot
  * 40 by 30 pixels in size
  * how to extract sprite from sprite sheet (12 sprites per row, 0 based index thanks to missingno #0 - sliggoo is #705. 59th row (30px per row 58 x 30 = 1740px), 10th column (40px per column 9 x 40 = 360px)

    ``` html
    <img style="width:40px; height:30px; background: url(https://play.pokemonshowdown.com/sprites/pokemonicons-sheet.png) transparent no-repeat scroll -360px -1740px;">
    ```

* Pokemon "models" (example Sliggoo) ![Garchomp gif](https://play.pokemonshowdown.com/sprites/ani/sliggoo.gif)
* Pokemon "models" from the back (example Sliggoo) ![Garchomp back gif](https://play.pokemonshowdown.com/sprites/ani-back/sliggoo.gif)
* Shiny ![Shiny Sliggoo](https://play.pokemonshowdown.com/sprites/ani-shiny/sliggoo.gif) ![Shiny Sliggoo](https://play.pokemonshowdown.com/sprites/ani-back-shiny/sliggoo.gif)
* Pokemon Sprites april fools ![afd](https://play.pokemonshowdown.com/sprites/afd/sliggoo.png)
* Type Labels ![Dragon](https://play.pokemonshowdown.com/sprites/types/Dragon.png)
* Generation 5 animated sprites ![wooloo](https://play.pokemonshowdown.com/sprites/gen5ani/wooloo.gif) (not available for all gen 6+ pokemon)
* Generation 5 still sprites ![Still Sliggoo](https://play.pokemonshowdown.com/sprites/gen5/sliggoo.png)
* Generation 4 and below have no back fill of new pokemon sprites
* Pokedex Image ![Sliggy](https://play.pokemonshowdown.com/sprites/dex/sliggoo.png)

## Side Notes

* the GUI will not render (dev or build) on Ubuntu 22.04 with the combination of Radeon Software for Linux 22.20 drivers and libwebkit2gtk-4.0-dev version 2.36.6-0ubuntu0.22.04.1
