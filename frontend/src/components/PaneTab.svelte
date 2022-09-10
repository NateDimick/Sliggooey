<!-- Pane Tab is a tab that allows the user to selct the specific pane they want to see by clicking this -->
<script lang="ts">
    import { tsPrint } from "../util"
    import { currentPaneStore, PaneInfo, panes, roomChats } from "../store"

    export let info: PaneInfo

    function seeThisPane() {
        tsPrint(`Switching to pane ${info.name}`)
        currentPaneStore.set(info.name)
    }

    function deleteThisPane() {
        panes.update((panes: PaneInfo[]) => {
            panes = panes.filter((p: PaneInfo) => p.name !== info.name)
            return panes
        })
        roomChats.update((rms) => {
            delete rms[info.name]
            return rms
        })
    }
</script>

<main>
    <input type="button" value={info.name} on:click={seeThisPane}>
    {#if info.removable === true}
        <input type="button" value="&#10006" on:click={deleteThisPane}>
    {/if}
</main>