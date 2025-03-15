/** @type {Array.<string>} */
let rooms = $state([])

fetch("http://localhost:8080/rooms", {
	method: "GET"
})
	.then(response => response.json())
	.then(data => {
		rooms.push(...Object.keys(data.data))
	})
	.catch((error) => {
		console.error('Failed to fetch rooms:', error);
		rooms = [];
	});

export const roomState = {
	get rooms() {
		return rooms;
	}
}
