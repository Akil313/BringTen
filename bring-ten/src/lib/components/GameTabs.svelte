<script>
	import CreateGame from './CreateGame.svelte';
	import JoinGame from './JoinGame.svelte';

	const tabs = [
		{ id: 'join', label: 'Join' },
		{ id: 'create', label: 'Create' }
	];

	/** @param {string} activeTab */
	let activeTab = $state('join');
</script>

<div class="flex flex-col border border-blue-400">
	<span>Game Tabs</span>
	<div class="flex gap-x-2">
		{#each tabs as tab}
			<button
				class="px-x cursor-pointer py-2 transition-colors duration-200"
				class:bg-blue-400={activeTab == tab.id}
				onclick={() => (activeTab = tab.id)}
			>
				{tab.label}
			</button>
		{/each}
	</div>
	<div class="my-2 h-1 bg-blue-200"></div>
	<div>
		{#if activeTab === 'create'}
			<CreateGame />
		{/if}
		{#if activeTab === 'join'}
			<JoinGame />
		{/if}
	</div>
</div>
