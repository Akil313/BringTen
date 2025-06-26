<script>
	import GameTabs from '$lib/components/GameTabs.svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
	let rooms = $state({});

	rooms = data.rooms;
</script>

<div class="flex justify-center">
	<header>
		<title>BringTen</title>
	</header>
	<GameTabs {rooms} />
</div>
