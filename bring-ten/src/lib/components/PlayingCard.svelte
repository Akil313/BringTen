<script>
	/**
	 * @typedef {Object} Props
	 * @property {string} cardString
	 */

	/** @type {Props}*/
	let { cardString } = $props();

	const cardStringArr = cardString.split('x');
	const value = cardStringArr[0];
	const suit = cardStringArr[1];
</script>

<div class="flex flex-col border border-green-500 px-2">
	<span>{value}</span>
	<span>{suit}</span>
</div>
