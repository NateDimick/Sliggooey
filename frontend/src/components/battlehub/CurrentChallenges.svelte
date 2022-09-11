<script>
import ChatChallenge from "../chat/ChatChallenge.svelte";
import { coldChallenges } from "../../store"

</script>

<main>
    <h1>Active Challenges</h1>
    {#each $coldChallenges as c}
        <ChatChallenge challengeWith={c.With}/>
    {/each}
</main>