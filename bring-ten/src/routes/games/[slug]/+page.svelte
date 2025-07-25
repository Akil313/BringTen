<script>
	import { onMount } from 'svelte';
	import PlayingCard from '$lib/components/PlayingCard.svelte';
	// @ts-ignore
	import Scalable from 'scalable';
	import { publicApiURL } from '$lib/config/public.js';

	/** @typedef {Object} GameState
	 * @property {string} name
	 * @property {number} position
	 * @property {string} roomName
	 * @property {string[]} hand
	 * @property {string[]} validHand
	 * @property {number} deck
	 * @property {number} currTurn
	 * @property {number} dealer
	 * @property {Object[]} players
	 * @property {number} players[].pos
	 * @property {string} players[].id
	 * @property {string} players[].name
	 * @property {number} team1Score
	 * @property {number} team2Score
	 * @property {string} trump
	 * @property {string[]} lift
	 * @property {boolean} playerBeg
	 * @property {boolean} playerStay
	 * @property {boolean} roundStart
	 * @property {boolean} gameStart
	 * @property {string} winner
	 */

	/** @typedef {Object} SSEState
	 * @property {string} name
	 * @property {number} position
	 * @property {string} room_name
	 * @property {string[]} hand
	 * @property {string[]} valid_hand
	 * @property {number} deck
	 * @property {number} curr_turn
	 * @property {number} dealer
	 * @property {Object[]} players
	 * @property {number} players[].pos
	 * @property {string} players[].id
	 * @property {string} players[].name
	 * @property {number} team_1_score
	 * @property {number} team_2_score
	 * @property {string} trump
	 * @property {string[]} lift
	 * @property {boolean} player_beg
	 * @property {boolean} player_stay
	 * @property {boolean} round_start
	 * @property {boolean} game_start
	 * @property {string} winner
	 */

	/**
	 * Creates a new GameState object with default values
	 * @returns {GameState} A new GameState object with default empty values
	 */
	function createEmptyGameState() {
		if (roomId === 'test') {
			return {
				name: 'Lelouch',
				position: 1,
				roomName: 'The Test Empire',
				hand: ['5xS', '3xD', 'KxS', '10xH', '8xC', '4xS'],
				validHand: ['KxS', '5xS', '10xH'],
				deck: 24,
				currTurn: 1,
				dealer: 1,
				players: [
					{ pos: 0, id: '38hd8s', name: 'Shirly' },
					{ pos: 1, id: 'nn32d9', name: 'Lelouch' },
					{ pos: 2, id: 'l30dii', name: 'C.C' },
					{ pos: 3, id: 'ej2b35', name: 'Suzaku' }
				],
				team1Score: 0,
				team2Score: 0,
				trump: 'QxH',
				lift: ['JxC'],
				playerBeg: false,
				playerStay: false,
				roundStart: false,
				gameStart: false,
				winner: ''
			};
		}
		return {
			name: '',
			position: -1,
			roomName: '',
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
	 * @param {SSEState} state
	 */
	function updateGameState(state) {
		return {
			name: state.name,
			position: state.position,
			roomName: state.room_name,
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

	/** @type {GameState} */
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

	/** @type {HTMLElement} */
	let main_container;
	/** @type {HTMLElement} */
	let canvas;
	/** @type {{ destroy?: () => void }} */
	let scalableInstance;

	// Set up EventSource when component mounts
	onMount(() => {
		scalableInstance = new Scalable(main_container, {
			align: 'center',
			verticalAlign: 'center',
			maxScale: 1.3
		});
		if (roomId !== 'test') {
			console.log(sseUrl);
			const eventSource = new EventSource(sseUrl);

			eventSource.onmessage = function (event) {
				// Parse the SSE data
				const state = JSON.parse(event.data);
				console.log('Received new game state:', state);

				// Update the store with the new game state
				gameState = { ...updateGameState(state) };
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
	class="flex h-screen w-screen overflow-hidden bg-[#00131a] bg-[url(/imgs/cartographer.png)] text-lg"
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
			{#if gameState?.gameStart === false}
				<button
					disabled={gameState.players.length !== 4}
					onclick={startGameHandler}
					class="relative left-[50%] top-[45%] -translate-x-1/2 rounded-lg p-2 {gameState.players
						.length !== 4
						? 'bg-gray-500 opacity-[50%]'
						: 'bg-blue-300'}"
				>
					Start Game
				</button>
			{/if}
		</div>
		<div id="play-field-ctn" class="relative basis-4/6 border border-blue-500">
			<div class="absolute left-[10%] -translate-x-1/2">
				<span>Player Turn: {gameState.players[gameState.currTurn].name}</span>
			</div>
			<div class="relative left-[16%] top-[18%] inline-block -translate-x-1/2">
				<span class="text-[1.2em]">Team 1 Points</span>
				<p class="text-center text-[1.5em]">{gameState.team1Score}</p>
			</div>
			<div class="relative left-[68%] top-[18%] inline-block -translate-x-1/2">
				<span class="text-[1.2em]">Team 2 Points</span>
				<p class="text-center text-[1.5em]">{gameState.team2Score}</p>
			</div>
			<div class="absolute left-[30%] top-[35%] z-10 justify-between">
				<PlayingCard
					cardString={'back'}
					selectCard={handleSelectCard}
					isSelected={false}
					isPlayable={false}
				/>
			</div>
			<div id="deck" class="z-5 absolute left-[27%] top-[35%]">
				<!--
				<span
					class="absolute left-[10%] top-[10%] z-10 origin-center -translate-x-1/2 -translate-y-1/2 text-[1.2em] text-black"
					>{gameState.deck}
				</span>
				-->
				<PlayingCard
					class="z-5 opacity-100 blur-[1px]"
					cardString={'back'}
					selectCard={handleSelectCard}
					isSelected={false}
					isPlayable={false}
				/>
			</div>
			<div class="absolute left-[47%] top-[28%] h-52 w-48">
				<div class="z-1 absolute left-1/2 -translate-x-1/2">
					{#if gameState.lift?.[0] === undefined}
						<div class="flex h-28 w-20 rounded border-2 bg-orange-400 opacity-40">
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
				<div class="z-2 absolute left-full top-1/2 -translate-x-full -translate-y-1/2">
					{#if gameState.lift?.[1] === undefined}
						<div class="flex h-28 w-20 rounded border-2 bg-purple-400 opacity-40">
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
				<div class="z-3 absolute left-1/2 top-full -translate-x-1/2 -translate-y-full">
					{#if gameState.lift?.[2] === undefined}
						<div class="flex h-28 w-20 rounded border-2 bg-orange-400 opacity-40">
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
				<div class="z-4 absolute top-1/2 -translate-y-1/2">
					{#if gameState.lift?.[3] === undefined}
						<div class="flex h-28 w-20 rounded border-2 bg-purple-400 opacity-40">
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
				{#each gameState.hand as c (c)}
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
			<button
				type="button"
				aria-label="Button"
				disabled={gameState.currTurn !== gameState.position || gameState.roundStart !== true}
				onclick={gameState.currTurn === gameState.position
					? () => handleAction('PLAY_CARD', selectedCard)
					: () => {}}
				class="absolute left-[80%] top-[85%] w-32 rounded-lg {gameState.currTurn ===
					gameState.position && gameState.roundStart === true
					? 'bg-blue-300'
					: 'bg-gray-500 opacity-[50%]'} p-2"
			>
				Play Card
			</button>
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
