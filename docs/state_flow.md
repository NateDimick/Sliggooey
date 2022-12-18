# How Battle Rooms Handle State

As stated in the frontpage readme, one of the goals of this project is to write as little Typescript as possible (even though it is a pleasant language, much better than vanilla JS, Go is the language of choice for this project in doing anything complex). Svelte as a framework helps facilitate this goal by being primarily an HTML framework.

Svelte is also very good at dealing with state. Its stores are very powerful and reactive and make storing and displaying state a breeze. The conflict comes when that state needs to be updated. A Pokemon Showdown battle is constantly updating state. The state cannot simply accept each update from the server as the new state, because the server updates specific parts of the full state one piece at a time. Each individual update needs to placed in teh correct spot and not touch anything else.

The flow is that:

1. The showdown server sends an update
1. The application backend receives the message, and parses it into a more usable format without having any notion of what the current state is
1. The parsed update is sent to the frontend, which hold the state
1. the state must be updated

If the frontend holds the state and the update, then it makes sense for the frontend to reconcile thw two to produce the new state, right? Wrong! The classic blunder leads to two problems:

1. More Typescript! That is to be avoided. Remember the project goals.
1. Duplicated data structures, because wails module only generates models for structs that are arguments (or nested within arguments) or return types for bound functions! Because state updates are sent to the frontend as events, these data structures are not automatically ported to Typescript. This (duplicating data structures) is tedious to the extreme and also involves even more typescript (Boo!).

So, The approach that has been chosen for step 4 is to have the frontend ask tha backend to do the reconciliation by passing the state to it along with the update. Think of it like the frontend and backend both have halves of a locket, but only the backend can put them back together and then only the frontend can wear the locket.

This approach allows the structs that parsed server message data is placed into to be reused in state reconciliation, and because the returned room state is the return variable type, it's Typescript representation is automagically generated by Wails which helps us write less code (yay!). What a win-win.

There is one caveat that you will notice looking at the frontend - Javascript event listeners will fire whenever they are called, even if a previous call of the same event has not finished yet. This can lead to a race condition if many updates are fired to the frontend at once, and the frontend is waiting for the backend to give it the new state. This can lead to undesirable outcomes where the final update in a rapid-fire stream of updates is potentially the only update applied. The solution is to implement a mutex lock around the call to the backend state reconciliation function, which ensure sequential and atomic state updates. This is is caused by the IPC delay in making a call from teh frontend to the backend, which is obviously a tradeoff of passing state updates back and forth between the two ends, however the solution is not messy and the benefits, as listed above, make this implementation worth the extra safety measures.