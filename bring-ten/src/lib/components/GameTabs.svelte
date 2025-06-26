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

<div class="flex max-w-[85%] flex-1 flex-col border border-blue-400">
	<span>Game Tabs</span>
	<div class="flex gap-x-2">
		{#each tabs as tab}
			<button
				class="px-x cursor-pointer py-2 transition-colors duration-200"
				class:bg-blue-400={activeTab == tab.id}
				id={tab.id}
				onclick={() => (activeTab = tab.id)}
			>
				{tab.label}
			</button>
		{/each}
	</div>
	<div class="my-2 h-1 bg-blue-200"></div>
	<div>
		{#if activeTab === 'create_game_tab'}
			<CreateGame />
		{/if}
		{#if activeTab === 'join_game_tab'}
			<JoinGame rooms={props.rooms} />
		{/if}
	</div>
</div>
