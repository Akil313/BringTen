<script>
	import { onMount } from 'svelte';
	import PlayingCard from '$lib/components/PlayingCard.svelte';
	import Scalable from 'scalable';
	import { publicApiURL } from '$lib/config/public.js';
	import gsap from 'gsap';

	const DECK_POS = 'left-[27%] top-[35%]';

	/**
	 * @param {any} tl
	 * @param {number} player
	 * @returns {any}
	 */

	function shareSixCards(tl, player) {
		for (let x = 0; x < 6; x++) {
			tl.to(
				`.p${player}-shared-card-${x}`,
				{
					opacity: '100',
					duration: 0
				},
				0
			);
			if (player == 0) {
				tl.to(
					`.p${player}-shared-card-${x}`,
					{
						top: '10%',
						left: `${40 + x * 3}%`,
						duration: 0.5,
						ease: 'power2.out'
					},
					0
				);
			} else if (player == 1) {
				tl.to(
					`.p${player}-shared-card-${x}`,
					{
						top: `${30 + x * 3}%`,
						left: `${85}%`,
						rotation: 90,
						duration: 0.5,
						ease: 'power2.out'
					},
					0.5
				);
			} else if (player == 2) {
				tl.to(
					`.p${player}-shared-card-${x}`,
					{
						top: `59%`,
						left: `${40 + x * 3}%`,
						rotation: 180,
						duration: 0.5,
						ease: 'power2.out'
					},
					1.0
				);
			} else if (player == 3) {
				tl.to(
					`.p${player}-shared-card-${x}`,
					{
						top: `${30 + x * 3}%`,
						left: `${7.5}%`,
						rotation: -90,
						duration: 0.5,
						ease: 'power2.out'
					},
					1.5
				);
			}
		}

		return tl;
	}

	/**
	 * @param {any} tl
	 * @returns {any}
	 */
	function flipTrump(tl) {
		/** @type {NodeListOf<HTMLElement>} */
		const cards = document.querySelectorAll('.trump-card');

		cards.forEach((card) => {
			const front = card.querySelector('.card-face-front');
			const back = card.querySelector('.card-face-back');

			if (!front || !back) return;

			tl.to(card, {
				rotationY: '+=180',
				duration: 0.5,
				ease: 'power2.in'
			});
		});

		return tl;
	}

	/** @type {any} */
	let tl2;

	function setAnim1() {
		tl2 = shareSixCards(tl2, 0);
		tl2 = shareSixCards(tl2, 1);
		tl2 = shareSixCards(tl2, 2);
		tl2 = shareSixCards(tl2, 3);
		tl2 = flipTrump(tl2);

		for (let player = 0; player < 4; player++) {
			for (let c = 0; c < 6; c++) {
				if (player == 0) {
					tl2.to(
						`.p${player}-shared-card-${c}`,
						{
							top: '8%',
							opacity: '0',
							duration: 0.5,
							ease: 'power2.out'
						},
						3.0
					);
				} else if (player == 1) {
					tl2.to(
						`.p${player}-shared-card-${c}`,
						{
							left: `87%`,
							opacity: '0',
							rotation: 90,
							duration: 0.5,
							ease: 'power2.out'
						},
						'<'
					);
				} else if (player == 2) {
					tl2.to(
						`.p${player}-shared-card-${c}`,
						{
							top: `61%`,
							opacity: '0',
							rotation: 180,
							duration: 0.5,
							ease: 'power2.out'
						},
						'<'
					);
				} else if (player == 3) {
					tl2.to(
						`.p${player}-shared-card-${c}`,
						{
							left: `5.5%`,
							opacity: '0',
							rotation: -90,
							duration: 0.5,
							ease: 'power2.out'
						},
						'<'
					);
				}
			}
		}
	}

	export function goAnim1() {
		tl2.play();
	}

	onMount(() => {
		tl2 = gsap.timeline({ paused: true });
		// Initial styles
		gsap.set('.trump-card', {
			transformStyle: 'preserve-3d',
			perspective: 1000
		});

		gsap.set('.card-face', {
			transformStyle: 'preserve-3d',
			transformOrigin: '50% 50%',
			backfaceVisibility: 'hidden'
		});

		gsap.set('.card-face-back', {
			rotationY: 180
		});

		setAnim1();
	});
</script>

<div class="h-full w-full">
	<div id="deck" class="absolute left-[27%] top-[35%] z-[200] blur-[1px]">
		<PlayingCard cardString={'back'} selectCard={() => {}} isSelected={false} isPlayable={false} />
	</div>
	<div class="absolute left-[30%] top-[35%] z-[201] aspect-[5/7] w-20 justify-between">
		<div class="trump-card relative h-full w-full">
			<div class="card-face card-face-front backface-hidden absolute h-full w-full">
				<PlayingCard
					cardString={'back'}
					selectCard={() => {}}
					isSelected={false}
					isPlayable={false}
				/>
			</div>
			<div class="card-face card-face-back backface-hidden absolute h-full w-full">
				<PlayingCard
					cardString={'6xS'}
					selectCard={() => {}}
					isSelected={false}
					isPlayable={false}
				/>
			</div>
		</div>
	</div>
	<div>
		{#each { length: 4 }, player_num}
			{#each { length: 6 }, i}
				<div class={`opacity-0 p${player_num}-shared-card-${i} absolute z-[${i}] ${DECK_POS}`}>
					<PlayingCard
						class="w-16"
						cardString={'back'}
						selectCard={() => {}}
						isSelected={false}
						isPlayable={false}
					/>
				</div>
			{/each}
		{/each}
	</div>
</div>
