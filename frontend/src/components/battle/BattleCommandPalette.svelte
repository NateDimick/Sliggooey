<!-- for inputting battle choices -->
<script lang="ts">
import { BattleChoice, BattleRequest, tsPrint } from "../../util";
import { MakeBattleChoice } from "../../wailsjs/go/main/App";
import { roomStates } from "../../store"

const enum CommandPaletteState {
    WaitForCommand,
    TeamOrder,
    SelectAction,
    SelectTarget,
    WaitForOpponent
}

export let roomName: string

let currentRequest: BattleRequest
let currentChoice: BattleChoice[]
let currentChoiceState: CommandPaletteState = CommandPaletteState.WaitForCommand
let currentChoiceIndex: number = 0

$: {
    currentRequest = $roomStates[roomName]?.request
    tsPrint(JSON.stringify(currentRequest))
    // fill current choice by the number of active pokemon in currentRequest
    currentChoiceState = CommandPaletteState.SelectAction
    currentChoice = Array(currentRequest?.active?.length).fill({})
    currentChoiceIndex = 0
}

function setMove(moveIndex: number) {
    return () => {
        currentChoice[currentChoiceIndex].Move = moveIndex
        currentChoiceState = CommandPaletteState.SelectTarget
    }
}

function setTarget(t: number) {
    return () => {
        currentChoice[currentChoiceIndex].Target = t
        currentChoiceIndex += 1
    }
}

function sendChoice(choices: BattleChoice[]) {
    MakeBattleChoice(roomName, currentRequest.rqid, choices)
}

function sendCommandChoice() {
    sendChoice(currentChoice)
}

function cancelChoice() {
    sendChoice([{Action: "undo"}])
    currentChoiceIndex = 0
    currentChoiceState = CommandPaletteState.SelectAction
}

function defaultChoice() {
    sendChoice([{Action: "default"}])
}

function skipChoice() {
    sendChoice([{Action: "pass"}])
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
    </div>
    {#if currentChoiceState === CommandPaletteState.WaitForCommand}
        <p>Getting Ready...</p>
    {:else if currentChoiceState === CommandPaletteState.SelectAction}
        {#if currentRequest != undefined}
            {#each currentRequest?.active[currentChoiceIndex]?.moves as move, index}
                <button on:click={setMove(index + 1)}>{move.move}</button>
            {/each}
            <!-- to do: change name that appears if game mod eis gen 8 and dynamax is allowed and dynamx is selected-->
            <input type="checkbox" name="gimmick" id="gimmick" on:click={toggleGimmick}>
        {/if}
    {:else if currentChoiceState === CommandPaletteState.SelectTarget}
        <p></p>
    {:else if currentChoiceState === CommandPaletteState.WaitForOpponent}
        <p>Waiting for Opponent...</p>
        <button on:click={cancelChoice}>Redo Move</button>
    {/if}
</main>