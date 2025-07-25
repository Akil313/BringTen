<script>
	import { useRooms } from '$lib/rooms.svelte';
	import RoomTable from './RoomTable.svelte';

	let roomId = $state('');

	/** @type {import('../../routes/$types').PageData} */
	let props = $props();

	/** @param {CustomEvent<string>} event */
	function handleRoomJoin(event) {
		console.log(event.detail);
	}

	/** @type { import('../../lib/rooms.svelte').RoomList } */
	const testRoomList = {};

	for (let i = 0; i < 20; i++) {
		const id = `test-room-${i}`;
		testRoomList[id] = {
			id: id,
			name: `Test Room ${i}`,
			host: `Nigga #${i}`,
			numPlayers: Math.floor(Math.random() * 4) + 1
		};
	}
</script>

<div class="flex w-full justify-center">
	<div class=" max-w-full grow">
		<form method="POST" action="?/join" class="mb-4 rounded bg-white px-8 pb-8 pt-6 shadow-md">
			<div class="mb-4">
				<label for="username" class="mb-2 block text-sm font-bold text-gray-700"> Username </label>
				<input
					placeholder="Username"
					type="text"
					class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
					id="join_game_username"
					name="name"
					required
				/>
			</div>
			<div class="mb-4">
				<input type="hidden" name="rooms" bind:value={roomId} />
				<RoomTable
					rooms={props.rooms}
					joinRoom={(id) => {
						roomId = id;
					}}
				/>
			</div>
			<div class="flex items-center justify-between"></div>
		</form>
	</div>
</div>
