<script>
	import CreateGame from './CreateGame.svelte';
	import JoinGame from './JoinGame.svelte';

	let props = $props();

	const tabs = [
		{ id: 'join_game_tab', label: 'Join' },
		{ id: 'create_game_tab', label: 'Create' }
	];

	/** @param {string} activeTab */
	let activeTab = $state('join_game_tab');
</script>

<div class="flex h-[60vh] w-full flex-col">
	<div class="flex max-w-xl justify-around">
		{#each tabs as tab}
			<button
				class="px-x w-full cursor-pointer rounded-t-xl border-r-2 border-t-2 border-l-indigo-400 border-r-indigo-400 border-t-indigo-400 py-2 shadow-xl transition-colors duration-200 {activeTab ===
				tab.id
					? 'bg-white'
					: 'bg-gray-400 opacity-[30%]'}"
				id={tab.id}
				onclick={() => (activeTab = tab.id)}
			>
				{tab.label}
			</button>
		{/each}
	</div>
	<div class="h-full w-full">
		{#if activeTab === 'create_game_tab'}
			<CreateGame />
		{/if}
		{#if activeTab === 'join_game_tab'}
			<JoinGame rooms={props.rooms} />
		{/if}
	</div>
</div>
