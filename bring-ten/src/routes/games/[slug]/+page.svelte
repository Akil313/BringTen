<script>
	import { onMount } from 'svelte';
	import PlayingCard from '$lib/components/PlayingCard.svelte';

	/** @typedef {Object} GameState
	 * @property {string} name
	 * @property {string} roomName
	 * @property {string[]} hand
	 * @property {string[]} validHand
	 * @property {string[]} deck
	 * @property {number} currTurn
	 * @property {number} dealer
	 * @property {string[]} players
	 * @property {number} team1Score
	 * @property {number} team2Score
	 * @property {string} trump
	 * @property {string[]} lift
	 * @property {boolean} playerBeg
	 * @property {boolean} playerStay
	 * @property {boolean} roundStart
	 * @property {string} winner
	 */

	/** @typedef {Object} SSEState
	 * @property {string} name
	 * @property {string} room_name
	 * @property {string[]} hand
	 * @property {string[]} valid_hand
	 * @property {string[]} deck
	 * @property {number} curr_turn
	 * @property {number} dealer
	 * @property {string[]} players
	 * @property {number} team_1_score
	 * @property {number} team_2_score
	 * @property {string} trump
	 * @property {string[]} lift
	 * @property {boolean} player_beg
	 * @property {boolean} player_stay
	 * @property {boolean} round_start
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
				roomName: 'The Test Empire',
				hand: ['5xS', '3xD', 'KxS', '10xH', '8xC', '4xS'],
				validHand: ['KxS', '5xS', '10xH'],
				deck: [],
				currTurn: 2,
				dealer: 1,
				players: ['38hd8s', 'nn32d9', 'l30dii', 'ej2b35'],
				team1Score: 0,
				team2Score: 0,
				trump: 'QxH',
				lift: ['JxC'],
				playerBeg: false,
				playerStay: false,
				roundStart: false,
				winner: ''
			};
		}
		return {
			name: '',
			roomName: '',
			hand: [],
			validHand: [],
			deck: [],
			currTurn: -1,
			dealer: -1,
			players: [],
			team1Score: 0,
			team2Score: 0,
			trump: '',
			lift: [],
			playerBeg: false,
			playerStay: false,
			roundStart: false,
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
			.then((data) => data)
			.catch((error) => {
				console.error('Failed to join room:', error);
			});
	}

	let { data } = $props();

	const roomId = data.slug;
	const playerId = data.playerId;
	const sseUrl = `http://localhost:8080/rooms/${roomId}/${playerId}/state`; // Your SSE endpoint

	/** @type {GameState} */
	let gameState = $state(createEmptyGameState());

	// Set up EventSource when component mounts
	onMount(() => {
		const eventSource = new EventSource(sseUrl);

		eventSource.onmessage = function (event) {
			// Parse the SSE data
			const state = JSON.parse(event.data);
			console.log('Received new game state:', state);

			// Update the store with the new game state
			gameState = updateGameState(state);
		};

		eventSource.onerror = function (error) {
			console.error('SSE connection error:', error);
			eventSource.close(); // Optionally handle reconnection here
		};

		// Clean up the EventSource when the component is destroyed
		return () => {
			eventSource.close();
		};
	});
</script>

<div class="px-8 pt-4">
	<h1 class="text-2xl">{gameState.roomName}</h1>
	<span># of Players: {gameState.players.length} / 4</span>
	<div class="min-w-xl flex justify-around">
		<span>{gameState.name}</span>
		<span>Player Turn: {gameState.currTurn}</span>
	</div>
	<div class="flex min-h-96 flex-col justify-around border border-blue-500">
		<span>Trump: {gameState.trump}</span>
		<span class="">Lift: {gameState.lift}</span>
		<div class="flex flex-col">
			<div class="flex">
				<span>Valid Hand: </span>
				<div class="flex gap-x-2">
					{#each gameState.validHand as c}
						<PlayingCard cardString={c} />
					{/each}
				</div>
			</div>
			<span class="">Hand: {gameState.hand}</span>
		</div>
	</div>
	<div class="flex flex-col">
		<span>Deck: {gameState.deck}</span>
		<span>playerBeg: {gameState.playerBeg}</span>
		<span>playerStay: {gameState.playerStay}</span>
		<span>roundStart: {gameState.roundStart}</span>
		<span>winner: {gameState.winner}</span>
	</div>

	<button onclick={startGameHandler} class="rounded-lg border border-blue-500 bg-blue-300 p-2">
		Start Game
	</button>
</div>
