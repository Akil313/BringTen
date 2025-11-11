<script>
	import '../app.css';
	let { children } = $props();
</script>

<div class="min-h-screen bg-[#2E3347]">
	{@render children()}
</div>
