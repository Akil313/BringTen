<script>
	let username = $state('');
	let roomName = $state('');
</script>

<div class="flex w-full justify-center">
	<div class="w-full max-w-xs">
		<form action="?/create" method="POST" class="mb-4 rounded bg-white px-8 pb-8 pt-6 shadow-md">
			<div class="mb-4">
				<label for="username" class="mb-2 block text-sm font-bold text-gray-700"> Username </label>
				<input
					placeholder="Username"
					type="text"
					class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
					id="username"
					name="name"
					bind:value={username}
				/>
			</div>
			<div class="mb-4">
				<label for="room_name" class="mb-2 block text-sm font-bold text-gray-700">
					Room Name
				</label>
				<input
					placeholder="Room Name"
					type="text"
					class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
					id="room_name"
					name="room_name"
					bind:value={roomName}
				/>
			</div>
			<div class="flex items-center justify-between">
				<button
					class="focus:shadow-outline w-full rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
					type="submit"
				>
					Create Game Room
				</button>
			</div>
		</form>
	</div>
</div>
