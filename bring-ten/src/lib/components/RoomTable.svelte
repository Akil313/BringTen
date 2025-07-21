<script>
	/**
	 * @typedef {Object} Props
	 * @property {(id: string) => void} joinRoom
	 * @property {import('../rooms.svelte').RoomList} rooms
	 */

	/** @type {Props} */
	let { joinRoom, rooms } = $props();
</script>

<div class="flex h-[60vh] flex-col bg-gray-200">
	<div class="grid grid-cols-5 px-4 py-2 text-center">
		<p>ID</p>
		<p>Room Name</p>
		<p>Host</p>
		<p>Num Players</p>
	</div>

	<div class="flex h-full grid-cols-5 flex-col gap-y-2 overflow-y-auto px-4 text-center">
		{#each Object.values(rooms) as r (r)}
			<div class="grid grid-cols-5 items-center rounded-md bg-white">
				{#each Object.values(r) as p (p)}
					<p class="justify-center">{p}</p>
				{/each}
				<button
					disabled={r.numPlayers == 4}
					class="focus:shadow-outline w-full rounded {r.numPlayers < 4
						? 'cursor-pointer bg-blue-500 hover:bg-blue-700'
						: 'cursor-default bg-gray-300 hover:bg-gray-400'}
							px-4 py-2 font-bold text-white focus:outline-none"
					onclick={r.numPlayers < 4 ? () => joinRoom(r.id) : () => {}}
				>
					Join
				</button>
			</div>
		{/each}
	</div>
</div>
