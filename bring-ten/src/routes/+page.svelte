<script>
	import GameTabs from '$lib/components/GameTabs.svelte';
	import { gsap } from 'gsap';
	import { onMount } from 'svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
	let rooms = $state({});

	rooms = data.rooms;
</script>

<div class="flex justify-center pt-6">
	<header>
		<title>BringTen</title>
	</header>
	<GameTabs {rooms} />

	<div class="box">Hey</div>
</div>
