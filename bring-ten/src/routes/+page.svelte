<script>
	import GameTabs from '$lib/components/GameTabs.svelte';
	import Header from '$lib/components/Header.svelte';
	import { gsap } from 'gsap';
	import { onMount } from 'svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
	let rooms = $state({});

	rooms = data.rooms;
</script>

<div class="mx-60 flex flex-col items-center gap-y-4 pt-10">
	<Header />
	<GameTabs {rooms} />
</div>
