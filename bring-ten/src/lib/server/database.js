// My temporary in memory database. Shhhhh don't tell anyone
//
const db = new Map();

/** @param {string} roomId */
export async function getRoom(roomId) {

	await refreshRoomList();

	console.log("Checking for room with id:", roomId)
	const roomInfo = db.get(roomId)

	if (roomInfo === undefined) {
		throw new Error("No room with that name")
	}

	return roomInfo
}

export async function refreshRoomList() {

	const url = "http://localhost:8080/rooms"
	const roomsResp = await fetch(url, {
		method: "GET",
	})
		.then(response => response.json())
		.then(data => data)

	/** 
	 * @type {Array<{id: string, name: string, host: string, numPlayers: number}>} rooms
	 */
	const rooms = Object.values(roomsResp.data);

	rooms.forEach((room) => {
		db.set(room.id, room)
	});
}

export async function sendGetEvent() {

	const url = "http://localhost:8080/"
	await fetch(url, {
		method: "GET",
	})
		.then(response => response.json())
		.then(data => {
			console.log(data)
		})

}
