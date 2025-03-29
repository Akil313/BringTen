<script>
	/**
	 * @typedef {Object} Props
	 * @property {(id: string) => void} joinRoom
	 * @property {import('../rooms.svelte').RoomList} rooms
	 */

	/** @type {Props} */
	let { joinRoom, rooms } = $props();
</script>

<div class="">
	<table class="w-full border border-red-300">
		<thead>
			<tr>
				<th>ID</th>
				<th>Room Name</th>
				<th>Host</th>
				<th>Num Players</th>
			</tr>
		</thead>
		<tbody>
			{#each Object.values(rooms) as r (r)}
				<tr>
					{#each Object.values(r) as p (p)}
						<td class=""> <div class="flex justify-center">{p}</div></td>
					{/each}
					<td>
						<button
							class="focus:shadow-outline w-full rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
							onclick={() => joinRoom(r.id)}
						>
							Join
						</button>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
