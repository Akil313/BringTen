<script>
	import { roomState } from '$lib/rooms.svelte';
</script>

<div class="flex w-full justify-center">
	<div class="w-full max-w-xs">
		<form method="POST" action="?/join" class="mb-4 rounded bg-white px-8 pb-8 pt-6 shadow-md">
			<div class="mb-4">
				<label for="username" class="mb-2 block text-sm font-bold text-gray-700"> Username </label>
				<input
					placeholder="Username"
					type="text"
					class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
					id="username"
					name="name"
				/>
			</div>
			<div class="mb-4">
				<label for="rooms" class="mb-2 block text-sm font-bold text-gray-700"> Rooms </label>
				<select
					id="rooms"
					name="rooms"
					class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
				>
					<option selected>--Select Room--</option>
					{#each roomState.rooms as room}
						<option value={room}>{room}</option>
					{/each}
				</select>
			</div>
			<button
				class="focus:shadow-outline w-full rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
				type="submit"
			>
				Join Game Room
			</button>
			<div class="flex items-center justify-between"></div>
		</form>
		<div class="flex flex-col">
			{#each roomState.rooms as room}
				<span>{room}</span>
			{/each}
		</div>
	</div>
</div>
