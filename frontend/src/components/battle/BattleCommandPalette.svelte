<!-- for inputting battle choices -->
<script lang="ts">
import type { BattleChoice, BattleRequest } from "../../util";
import { MakeBattleChoice } from "../../wailsjs/go/main/App";
import { battleRequests } from "../../store"

const enum CommandPaletteState {
    WaitForCommand,
    TeamOrder,
    SelectAction,
    SelectGimmickAction,
    SelectTarget,
    WaitForOpponent
}

export let roomName: string

let currentRequest: BattleRequest
let currentChoice: BattleChoice[]
let currentChoiceState: CommandPaletteState = CommandPaletteState.WaitForCommand
let currentChoiceIndex: number = 0

$: {
    currentRequest = $battleRequests[roomName]
    // fill current choice by the number of active pokemon in currentRequest
    currentChoiceState = CommandPaletteState.SelectAction
}

function sendChoice() {
    MakeBattleChoice(roomName, currentRequest.rqid, currentChoice)
}

function cancelChoice() {

}

function defaultChoice() {

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
        <p></p>
    {:else if currentChoiceState === CommandPaletteState.SelectAction}
        <p></p>
    {:else if currentChoiceState === CommandPaletteState.SelectGimmickAction}
        <p></p>
    {:else if currentChoiceState === CommandPaletteState.SelectTarget}
        <p></p>
    {:else if currentChoiceState === CommandPaletteState.WaitForOpponent}
        <p></p>
    {/if}
</main>