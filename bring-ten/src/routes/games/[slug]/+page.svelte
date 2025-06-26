<script>
	import { onMount } from 'svelte';
	import PlayingCard from '$lib/components/PlayingCard.svelte';

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
		const startUrl = `http://localhost:8080/rooms/${roomId}/start`;
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

		const playCardUrl = `http://localhost:8080/rooms/${roomId}/${playerId}/action`;
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
	const sseUrl = `http://localhost:8080/rooms/${roomId}/${playerId}/state`; // Your SSE endpoint

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
	$inspect(gameState, playerHand);

	// Set up EventSource when component mounts
	onMount(() => {
		if (roomId !== 'test') {
			const eventSource = new EventSource(sseUrl);

			eventSource.onmessage = function (event) {
				// Parse the SSE data
				const state = JSON.parse(event.data);
				console.log('Received new game state:', state);

				// Update the store with the new game state
				gameState = { ...updateGameState(state) };
				//playerHand = updatePlayerHand(gameState.hand, gameState.validHand);
				//console.log(playerHand);
			};

			eventSource.onerror = function (error) {
				console.error('SSE connection error:', error);
				eventSource.close(); // Optionally handle reconnection here
			};

			// Clean up the EventSource when the component is destroyed
			return () => {
				eventSource.close();
			};
		}
	});
</script>

<div class="flex h-screen flex-col px-4">
	<h1 class="text-2xl">{gameState.roomName}</h1>
	<span># of Players: {gameState.players.length} / 4</span>
	<div class="flex justify-around">
		<span>{gameState.name}</span>
		<span>Player Turn: {gameState.players[gameState.currTurn].name}</span>
	</div>
	<div class="flex flex-col justify-around border border-blue-500">
		<div>
			<span>Trump:</span>
			<PlayingCard
				cardString={gameState.trump}
				selectCard={handleSelectCard}
				isSelected={false}
				isPlayable={false}
			/>
		</div>
		<div class="flex flex-initial border border-green-600">
			<PlayingCard
				cardString={'back'}
				selectCard={handleSelectCard}
				isSelected={false}
				isPlayable={false}
			/>
			<span class="translate-8 absolute">{gameState.deck}</span>
		</div>
		<div class="flex-col">
			<span class="">Lift:</span>
			{#each gameState.lift as c}
				<PlayingCard cardString={c} selectCard={() => {}} isSelected={false} isPlayable={false} />
			{/each}
		</div>
		<div class="flex flex-col">
			<span>Hand: </span>
			<div class="">
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
		</div>
		{#if !gameState.roundStart}
			<span>{gameState.position === gameState.currTurn}</span>
			<span>{gameState.gameStart}</span>
			<span>{!gameState.playerBeg}</span>
			{#if gameState.position === gameState.currTurn && gameState.gameStart && gameState.playerBeg === false}
				<div class="flex w-full flex-row justify-around">
					<button
						onclick={() => handleAction('STAY')}
						class="rounded-lg border border-blue-500 bg-blue-300 p-2">Stay</button
					>
					<button
						onclick={() => handleAction('BEG')}
						class="rounded-lg border border-blue-500 bg-blue-300 p-2">Beg</button
					>
				</div>
			{/if}
			{#if gameState.position === gameState.dealer && gameState.gameStart && gameState.playerBeg === true}
				<div class="flex w-full flex-row justify-around">
					<button
						onclick={() => handleAction('GO_AGAIN')}
						class="rounded-lg border border-blue-500 bg-blue-300 p-2">Go Again</button
					>
					<button
						onclick={() => handleAction('GIVE_ONE')}
						class="rounded-lg border border-blue-500 bg-blue-300 p-2">Give 1</button
					>
				</div>
			{/if}
		{/if}
		<button
			type="button"
			aria-label="Button"
			disabled={gameState.currTurn !== gameState.position}
			onclick={gameState.currTurn === gameState.position
				? () => handleAction('PLAY_CARD', selectedCard)
				: () => {}}
			class="rounded-lg border border-blue-500 {gameState.currTurn === gameState.position
				? 'bg-blue-300'
				: 'bg-gray-300'} p-2"
		>
			Submit Card
		</button>
	</div>
	<div class="flex flex-col">
		<span>position: {gameState.position}</span>
		<span>currTurn: {gameState.currTurn}</span>
		<span>playerBeg: {gameState.playerBeg}</span>
		<span>playerStay: {gameState.playerStay}</span>
		<span>roundStart: {gameState.roundStart}</span>
		<span>gameStart: {gameState.gameStart}</span>
		<span>winner: {gameState.winner}</span>
		<span>players: {JSON.stringify(gameState.players)}</span>
	</div>

	<button
		disabled={gameState.players.length !== 4}
		onclick={startGameHandler}
		class="rounded-lg border border-blue-500 p-2 {gameState.players.length !== 4
			? 'bg-gray-200'
			: 'bg-blue-300'}"
	>
		Start Game
	</button>
</div>
