<!-- for inputting battle choices -->
<script lang="ts">
import { BattleChoice, BattleRequest, tsPrint } from "../../util";
import { MakeBattleChoice } from "../../wailsjs/go/backend/App";
import { roomStates } from "../../store"

const enum CommandPaletteState {
    WaitForCommand,
    TeamOrder,
    SelectAction,
    SelectTarget,
    WaitForOpponent,
    WaitForNextRequest,
    ForceSwitch,
    GameOver
}

export let roomName: string

let currentRequest: BattleRequest
let currentChoice: BattleChoice[]
let currentChoiceState: CommandPaletteState = CommandPaletteState.WaitForCommand
let currentChoiceIndex: number = 0

$: {
    currentRequest = $roomStates[roomName]?.request as BattleRequest
    //tsPrint(`current request ${JSON.stringify(currentRequest)}`)  // if this component ever shits the bed, uncomment this line
    // fill current choice by the number of active pokemon in currentRequest
    if (currentRequest?.forceSwitch?.length) {
        currentChoiceState = CommandPaletteState.ForceSwitch
        currentChoice = Array(currentRequest?.forceSwitch.length).fill({})
        currentChoiceIndex = 0
    } else if (currentRequest?.wait) {
        currentChoiceState = CommandPaletteState.WaitForNextRequest
    } else {
        currentChoiceState = CommandPaletteState.SelectAction
        currentChoice = Array(currentRequest?.active?.length).fill({})
        currentChoiceIndex = 0
    }
}

$: {
    let gameIsActive = $roomStates[roomName]?.active
    if (!gameIsActive) {
        currentChoiceState = CommandPaletteState.GameOver
    }
}

$: {
    let newStateAfterUpdate = currentChoiceState
    if ($roomStates[roomName]?.gameType === "singles" && newStateAfterUpdate === CommandPaletteState.SelectTarget) {
        tsPrint(`sending single battle command in singles: ${JSON.stringify(currentChoice)}`)
        sendCommandChoice()
    }
}

$: {
    let choiceIndexAfterUpdate = currentChoiceIndex
    if (choiceIndexAfterUpdate >= currentChoice.length) {
        tsPrint(`sending choices because all targets selected ${JSON.stringify(currentChoice)}`)
        sendCommandChoice()
    }
}

function setMove(moveIndex: number) {
    return () => {
        currentChoice[currentChoiceIndex].Move = moveIndex
        currentChoice[currentChoiceIndex].Action = "move"
        currentChoiceState = CommandPaletteState.SelectTarget
    }
}

function setTarget(t: number) {
    return () => {
        currentChoice[currentChoiceIndex].Target = t
        currentChoiceIndex += 1
    }
}

function setSwitchPokemon(switchIndex: number) {
    return () => {
        currentChoice[currentChoiceIndex].Action = "switch"
        currentChoice[currentChoiceIndex].Target = switchIndex
        currentChoiceIndex += 1
    }
}

function sendChoice(choices: BattleChoice[]) {
    MakeBattleChoice(roomName, currentRequest.rqid, choices)
}

function sendCommandChoice() {
    sendChoice(currentChoice)
    currentChoiceState = CommandPaletteState.WaitForOpponent
}

function cancelChoice() {
    sendChoice([{Action: "undo"}])
    currentChoiceIndex = 0
    currentChoiceState = CommandPaletteState.SelectAction
}

function defaultChoice() {
    sendChoice([{Action: "default"}])
    currentChoiceState = CommandPaletteState.WaitForOpponent
}

function skipChoice() {
    sendChoice([{Action: "pass"}])
    currentChoiceState = CommandPaletteState.WaitForOpponent
}

function toggleGimmick() {

}

function startTimer() {

}

function killTimer() {

}

function forfeit() {

}
</script>

<main>
    <div>
        <!-- commands that do not require a state context -->
        <button on:click={skipChoice}>Skip</button>
        <button on:click={defaultChoice}>Default</button>
        <button>Timer</button>
        <button>Forfeit</button>
    </div>
    {#if currentChoiceState === CommandPaletteState.WaitForCommand}
        <p>Getting Ready...</p>
    {:else if currentChoiceState === CommandPaletteState.SelectAction && currentRequest !== null}
        <div>
            {#each currentRequest?.active[currentChoiceIndex]?.moves as move, index}
                <button disabled={move.disabled} on:click={setMove(index + 1)}>{move.move} {move.pp}/{move.maxpp}</button>
            {/each}
        </div>
        <!-- to do: change name that appears if game mod eis gen 8 and dynamax is allowed and dynamax is selected-->
        <input type="checkbox" name="gimmick" id="gimmick" on:click={toggleGimmick}>
        <div>
            {#each currentRequest?.side?.pokemon as p, index}
                <button disabled={p.active || p.condition === "0 fnt"} on:click={setSwitchPokemon(index + 1)}>{p.ident}</button>
            {/each}
        </div>
    {:else if currentChoiceState === CommandPaletteState.ForceSwitch && currentRequest !== null}
        {#each currentRequest?.side?.pokemon as p, index}
        <!-- todo do not allow switches to active or fainted pokemon -->
            <button on:click={setSwitchPokemon(index + 1)}>{p.ident}</button>
        {/each}
    {:else if currentChoiceState === CommandPaletteState.SelectTarget}
        <p></p>
    {:else if currentChoiceState === CommandPaletteState.WaitForOpponent}
        <p>Waiting for Opponent...</p>
        <button on:click={cancelChoice}>Redo Move</button>
    {/if}
</main>

<style>
    main {
        bottom: 0;
    }
</style>