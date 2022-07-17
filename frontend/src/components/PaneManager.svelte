<!-- the Pane Manager allows for switching between Panes and manages the addition and removal of Panes -->
<script lang="ts">
import PaneTab from "./PaneTab.svelte";
import { panes } from "../store"
import { PaneType } from "../util";
import HomePane from "./HomePane.svelte";
import ChatPane from "./ChatPane.svelte";
import BattlePane from "./BattlePane.svelte";

</script>

<main>
    <div id="tabs">
        {#each $panes as p}
            <PaneTab info={p}/>
        {/each}
    </div>
    <div id="view">
        {#each $panes as p}
            {#if p.type === PaneType.HomePane}
                <HomePane info={p}/>
            {:else if p.type === PaneType.ChatPane}
                <ChatPane info={p}/>
            {:else if p.type === PaneType.BattlePane}
                <BattlePane info={p}/>
            {/if}
        {/each}
    </div>
</main>