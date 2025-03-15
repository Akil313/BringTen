<script>
	import '../app.css';
	let { children } = $props();
</script>

<div class="min-h-screen bg-red-200">
	{@render children()}
</div>
