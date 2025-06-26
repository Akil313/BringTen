/** 
 * @typedef {Object} Room
 * @property {string} id
 * @property {string} name
 * @property {string} host
 * @property {number} numPlayers
 */

/**
 * @typedef {Object<string, Room>} RoomList
 */

export function useRooms() {
	let rooms = $state({})

	$effect(() => {
		$inspect(rooms)
	})

	async function fetchRooms() {
		const serverURL = import.meta.env.VITE_API_URL;

		try {

			const response = await fetch(`${serverURL}/rooms`, {
				method: "GET"
			});
			const data = await response.json()
			rooms = data.data
			console.log("Fetched rooms:", rooms);
		} catch (error) {
			console.error('Failed to fetch rooms:', error);
			console.log('Defaulting to empty list');
			rooms = {};
		};
		return rooms
	}

	// Set the rooms directly from the SSR response or manual refresh
	/**
	 * Sets the rooms in the state.
	 * This function updates the reactive `rooms` state with new room data.
	 * 
	 * @param {Array<Object>} newRooms - The new list of rooms to update the state with.
	 * Each room should be an object with properties like `id`, `name`, `host`, etc.
	 */
	function setRooms(newRooms) {
		rooms = newRooms;
	}

	return {
		rooms,
		fetchRooms,
		setRooms
	}
}


