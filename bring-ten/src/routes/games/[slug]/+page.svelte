<script>
	import { onMount } from 'svelte';
	import PlayingCard from '$lib/components/PlayingCard.svelte';
	import Scalable from 'scalable';
	import { publicApiURL } from '$lib/config/public.js';
	import gsap from 'gsap';
	import GameAnimation from '$lib/components/GameAnimation.svelte';
	import GameTabs from '$lib/components/GameTabs.svelte';

	/**
	 * Creates a new GameState object with default values
	 * @returns {import('$lib/types.js').GameState} A new GameState object with default empty values
	 */
	function createEmptyGameState() {
		return {
			name: '',
			position: -1,
			roomName: '',
			hostId: '',
			hand: [],
			validHand: [],
			deck: 0,
			currTurn: 0,
			dealer: -1,
			players: [{ pos: 0, id: '', name: '' }],
			team1Score: 0,
			team2Score: 0,
			trump: '',
			lift: [],
			playerBeg: false,
			playerStay: false,
			roundStart: false,
			gameStart: false,
			winner: ''
		};
	}

	/**
	 * Updates the gameState object from the values from the SSE
	 * @param {import('$lib/types.js').SSEState} state
	 */
	function updateGameState(state) {
		return {
			name: state.name,
			position: state.position,
			roomName: state.room_name,
			hostId: state.host_id,
			hand: state.hand,
			validHand: state.valid_hand,
			deck: state.deck,
			currTurn: state.curr_turn,
			dealer: state.dealer,
			players: state.players,
			team1Score: state.team_1_score,
			team2Score: state.team_2_score,
			trump: state.trump,
			lift: state.lift,
			playerBeg: state.player_beg,
			playerStay: state.player_stay,
			roundStart: state.round_start,
			gameStart: state.game_start,
			winner: state.winner
		};
	}

	/**
	 * This handles starting the game
	 * @param {MouseEvent} event
	 */
	async function startGameHandler(event) {
		const startUrl = `${publicApiURL}/rooms/${roomId}/start`;
		const resp = await fetch(startUrl, {
			method: 'POST',
			body: JSON.stringify({
				player_id: playerId
			})
		})
			.then((response) => response.json())
			.catch((error) => {
				console.error('Failed to join room:', error);
			});

		console.log(resp);
	}

	/**
	 * This handles starting the game
	 * @param {string} card
	 */
	function handleSelectCard(card) {
		if (card === selectedCard) {
			selectedCard = null;
			return;
		}
		selectedCard = card;
	}

	/**
	 * This handles starting the game
	 * @param {import("$lib/types.js").GameState} prev
	 * @param {import("$lib/types.js").GameState} curr
	 */
	function runAnimations(prev, curr) {
		const gameStarted = prev.gameStart === false && curr.gameStart === true;

		if (gameStarted) {
			gameAnimation.goAnim1();
		}
	}

	/**
	 * This handles starting the game
	 * @param {string} action
	 * @param {string} cardPlayed
	 */
	async function handleAction(action, cardPlayed = '') {
		if (!selectedCard && action == 'PLAY_CARD') {
			console.log('No Card Selected');
			return;
		}

		const playCardUrl = `${publicApiURL}/rooms/${roomId}/${playerId}/action`;
		const resp = await fetch(playCardUrl, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				action: action,
				card_played: cardPlayed
			})
		});
		return;
	}

	let { data } = $props();

	const roomId = data.slug;
	const playerId = data.playerId;
	const sseUrl = `${publicApiURL}/rooms/${roomId}/${playerId}/state`; // Your SSE endpoint

	let selectedCard = $state();

	/** @type {import('$lib/types.js').GameState} */
	let gameState = $state.raw(createEmptyGameState());

	/**
	 * @param {string[]} hand
	 * @param {string[]} validHand
	 * @return {Object<string, Boolean>} validHand
	 */
	function updatePlayerHand(hand, validHand) {
		/** @type {Object<string, Boolean>} handMap */
		const handMap = {};

		hand.forEach((card) => {
			handMap[card] = validHand.includes(card);
		});

		return handMap;
	}

	/** @type {Object<string, Boolean>} */
	const playerHand = $derived(updatePlayerHand(gameState.hand, gameState.validHand));

	let card_ghost_1 = 'w-20 aspect-[5/7] rounded border-2 bg-orange-400 opacity-40';
	let card_ghost_2 = 'w-20 aspect-[5/7] rounded border-2 bg-purple-400 opacity-40';
	const DECK_POS = 'left-[27%] top-[35%]';

	/** @type {HTMLElement} */
	let main_container;
	/** @type {HTMLElement} */
	let canvas;
	/** @type {{ destroy?: () => void }} */
	let scalableInstance;
	let ready = $state(false);

	/** @type GameAnimation */
	let gameAnimation;

	// Set up EventSource when component mounts
	onMount(() => {
		const startGameBtn = document.getElementById('start-game-btn');

		const tl2 = gsap.timeline({ paused: true });
		tl2.to('.shared-card-1', {
			top: '10%',
			left: '50%',
			duration: 1
		});

		scalableInstance = new Scalable(main_container, {
			align: 'center',
			verticalAlign: 'center',
			maxScale: 1.3
		});
		ready = true;

		if (roomId !== 'test') {
			console.log(sseUrl);
			const eventSource = new EventSource(sseUrl);

			eventSource.onmessage = function (event) {
				// Parse the SSE data
				const state = JSON.parse(event.data);
				console.log('Received new game state:', state);

				let prevGameState = gameState;
				// Update the store with the new game state
				gameState = { ...updateGameState(state) };

				runAnimations(prevGameState, gameState);
			};

			eventSource.onerror = function (error) {
				console.error('SSE connection error:', error);
				eventSource.close(); // Optionally handle reconnection here
			};

			// Clean up the EventSource when the component is destroyed
			return () => {
				scalableInstance?.destroy?.();
				eventSource.close();
			};
		}
	});
</script>

<div
	id="main-container"
	class={ready
		? 'visible-after-mount flex h-screen w-screen overflow-hidden bg-[#00131a] bg-[url(/imgs/cartographer.png)] text-lg'
		: 'hidden-before-mount'}
	bind:this={main_container}
>
	<div
		id="canvas"
		class="flex h-[720px] w-[1280px] bg-[#282f28] bg-[url($lib/images/table.png)] text-white"
		bind:this={canvas}
	>
		<div id="left-info-ctn" class="basis-1/6 border border-red-600 px-4 pt-4">
			<span class="absolute text-[1.5em]">{gameState.roomName}</span>
			<span class="absolute top-[10%] text-[0.8em]"
				># of Players: {gameState.players.length} / 4</span
			>
			{#if gameState?.gameStart === false && gameState.hostId === playerId}
				<button
					disabled={gameState.players.length !== 4}
					id="start-game-btn"
					onclick={startGameHandler}
					class="relative left-[50%] top-[45%] -translate-x-1/2 rounded-lg p-2 {gameState.players
						.length === 4
						? 'bg-blue-400'
						: 'bg-gray-500 opacity-[50%]'}
					"
				>
					Start Game
				</button>
			{/if}
		</div>
		<div id="play-field-ctn" class="relative basis-4/6 border border-blue-500">
			<!--
			<GameAnimation bind:this={gameAnimation} />
			-->
			<div id="deck" class="absolute left-[27%] top-[35%] z-[200] blur-[1px]">
				<PlayingCard
					cardString={'back'}
					selectCard={() => {}}
					isSelected={false}
					isPlayable={false}
				/>
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
							cardString={gameState.trump}
							selectCard={() => {}}
							isSelected={false}
							isPlayable={false}
						/>
					</div>
				</div>
			</div>
			<div class="absolute left-[10%] top-[18%]">
				<span class="text-[1.2em]">Team 1 Points</span>
				<p class="text-center text-[1.5em]">{gameState.team1Score}</p>
			</div>
			<div class="absolute left-[50%] top-[18%] -translate-x-1/2">
				<span class="text-[1.2em]">Player Turn</span>
				<p class="text-center text-[1.2em]">{gameState.players[gameState.currTurn].name}</p>
			</div>
			<div class="absolute left-[75%] top-[18%]">
				<span class="text-[1.2em]">Team 2 Points</span>
				<p class="text-center text-[1.5em]">{gameState.team2Score}</p>
			</div>
			<div class="absolute left-[47%] top-[28%] h-52 w-48">
				<div
					class="{gameState.lift?.[0] === undefined
						? 'z-[1]'
						: 'z-[6]'} absolute left-1/2 -translate-x-1/2"
				>
					{#if gameState.lift?.[0] === undefined}
						<div class={`${card_ghost_1}`}>
							<span class="absolute left-[10%] top-[10%]">Player 1</span>
						</div>
					{:else}
						<PlayingCard
							class=""
							cardString={gameState.lift[0]}
							selectCard={handleSelectCard}
							isSelected={false}
							isPlayable={false}
						/>
					{/if}
				</div>
				<div
					class="{gameState.lift?.[1] === undefined
						? 'z-[2]'
						: 'z-[7]'} absolute left-full top-1/2 -translate-x-full -translate-y-1/2"
				>
					{#if gameState.lift?.[1] === undefined}
						<div class={`${card_ghost_2}`}>
							<span class="absolute left-[10%] top-[30%] -translate-y-1/2">Player 2</span>
						</div>
					{:else}
						<PlayingCard
							class=""
							cardString={gameState.lift[1]}
							selectCard={handleSelectCard}
							isSelected={false}
							isPlayable={false}
						/>
					{/if}
				</div>
				<div
					class="{gameState.lift?.[2] === undefined
						? 'z-[3]'
						: 'z-[8]'} absolute left-1/2 top-full -translate-x-1/2 -translate-y-full"
				>
					{#if gameState.lift?.[2] === undefined}
						<div class={`${card_ghost_1}`}>
							<span class="absolute left-[10%] top-[60%]">Player 3</span>
						</div>
					{:else}
						<PlayingCard
							class=""
							cardString={gameState.lift[2]}
							selectCard={handleSelectCard}
							isSelected={false}
							isPlayable={false}
						/>
					{/if}
				</div>
				<div
					class="{gameState.lift?.[3] === undefined
						? 'z-[4]'
						: 'z-[9]'} absolute top-1/2 -translate-y-1/2"
				>
					{#if gameState.lift?.[3] === undefined}
						<div class={`${card_ghost_2}`}>
							<span class="absolute left-[10%] top-[70%] -translate-y-1/2"> Player 4 </span>
						</div>
					{:else}
						<PlayingCard
							class=""
							cardString={gameState.lift[3]}
							selectCard={handleSelectCard}
							isSelected={false}
							isPlayable={false}
						/>
					{/if}
				</div>
			</div>
			<div class="absolute left-[7%] top-[80%] flex">
				{#each gameState.hand as c, i (c + i)}
					<PlayingCard
						cardString={c}
						selectCard={handleSelectCard}
						isSelected={selectedCard === c}
						isValid={playerHand[c]}
						isPlayable={true}
					/>
				{/each}
			</div>
			{#if !gameState.roundStart}
				{#if gameState.position === gameState.currTurn && gameState.gameStart && gameState.playerBeg === false}
					<div class="absolute left-[50%] top-[62%] w-72 -translate-x-1/2">
						<button
							onclick={() => handleAction('STAY')}
							class="absolute rounded-lg border border-blue-500 bg-blue-400 px-6 py-2">Stay</button
						>
						<button
							onclick={() => handleAction('BEG')}
							class="absolute left-full -translate-x-full rounded-lg border border-blue-500 bg-blue-400 px-6 py-2"
							>Beg</button
						>
					</div>
				{/if}
				{#if gameState.position === gameState.dealer && gameState.gameStart && gameState.playerBeg === true}
					<div class="absolute left-[50%] top-[62%] w-72 -translate-x-1/2">
						<button
							onclick={() => handleAction('GO_AGAIN')}
							class="absolute w-32 rounded-lg border border-blue-500 bg-blue-400 py-2"
							>Go Again</button
						>
						<button
							onclick={() => handleAction('GIVE_ONE')}
							class="absolute left-full w-32 -translate-x-full rounded-lg border border-blue-500 bg-blue-400 py-2"
							>Give One</button
						>
					</div>
				{/if}
			{/if}
			{#if gameState.roundStart}
				<button
					type="button"
					aria-label="Button"
					disabled={gameState.currTurn !== gameState.position || gameState.roundStart !== true}
					onclick={gameState.currTurn === gameState.position
						? () => handleAction('PLAY_CARD', selectedCard)
						: () => {}}
					class="absolute left-[80%] top-[85%] w-32 rounded-lg {gameState.currTurn ===
						gameState.position && gameState.roundStart === true
						? 'bg-blue-400'
						: 'bg-gray-500 opacity-[50%]'} p-2"
				>
					Play Card
				</button>
			{/if}
		</div>
		<div class="flex grow-0 basis-1/6 flex-col">
			<span>position: {gameState.position}</span>
			<span>currTurn: {gameState.currTurn}</span>
			<span>playerBeg: {gameState.playerBeg}</span>
			<span>playerStay: {gameState.playerStay}</span>
			<span>roundStart: {gameState.roundStart}</span>
			<span>gameStart: {gameState.gameStart}</span>
			<span>winner: {gameState.winner}</span>
		</div>
	</div>
</div>

<style>
	.hidden-before-mount {
		visibility: hidden;
	}

	.visible-before-mount {
		visibility: visible;
		transition: none;
	}
</style>
