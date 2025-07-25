<script>
	import { cardImages } from '$lib/cards.svelte.js';

	/**
	 * @typedef {Object} Props
	 * @property {string} [ class ]
	 * @property {string} cardString
	 * @property {Function} selectCard
	 * @property {Boolean} isSelected
	 * @property {Boolean} [isValid]
	 * @property {Boolean} [isPlayable]
	 */

	/** @type {Props}*/
	let {
		class: customClass = '',
		isValid = false,
		isPlayable = true,
		selectCard,
		cardString,
		isSelected
	} = $props();

	const cardStringArr = cardString.split('x');
	const value = cardStringArr[0];
	const suit = cardStringArr[1];

	const testImg = '../images/cards/CLUB-1.svg';
	let cardImg = cardString ? cardImages[cardString] : cardImages['AxS'];

	function handleCardClick() {
		selectCard(cardString);
	}
</script>

{#if cardString}
	<button
		type="button"
		aria-label="Button"
		onclick={isPlayable && isValid ? handleCardClick : () => {}}
		class={`${customClass} relative w-20 transition-transform ${isSelected ? '-translate-y-4' : ''} ${isPlayable && !isValid ? 'opacity-50' : ''} ${isPlayable && isValid ? 'cursor-pointer hover:scale-110' : 'cursor-default'}`}
	>
		{#await cardImages[cardString]}
			<span>Card</span>
		{:then cardSvg}
			<svg
				class="h-full w-full"
				viewBox={cardString === 'back' ? '0 0 240 336' : '0 0 238.11073 332.5986'}
			>
				{@html cardSvg.default}
			</svg>
		{/await}
	</button>
{/if}
