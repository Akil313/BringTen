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
let rooms = $state({})

fetch("http://localhost:8080/rooms", {
	method: "GET"
})
	.then(response => response.json())
	.then(data => {
		rooms = data.data
	})
	.catch((error) => {
		console.error('Failed to fetch rooms:', error);
		rooms = {};
	});

export const roomState = {
	get rooms() {
		return rooms;
	}
}
