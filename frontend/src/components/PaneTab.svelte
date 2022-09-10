<!-- Pane Tab is a tab that allows the user to selct the specific pane they want to see by clicking this -->
<script lang="ts">
    import { PaneType } from "../util"
    import { currentPaneStore, PaneInfo, panes } from "../store"

    export let info: PaneInfo

    function seeThisPane() {
        currentPaneStore.set(info.name)
    }

    function deleteThisPane() {
        panes.update((panes: PaneInfo[]) => {
            panes = panes.filter((p: PaneInfo) => p.name !== info.name)
            return panes
        })
    }
</script>

<main>
    <input type="button" value={info.name} on:click={seeThisPane}>
    {#if info.type !== PaneType.HomePane && info.type !== PaneType.ChatPane}
        <input type="button" value="&#10006" on:click={deleteThisPane}>
    {/if}
</main>