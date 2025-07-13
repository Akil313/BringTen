import { redirect } from '@sveltejs/kit';
import { fail } from '@sveltejs/kit';
import * as db from '$lib/server/database'
import * as utils from '$lib/utils'
import { useRooms } from '$lib/rooms.svelte';
import { apiURL } from '$lib/config/private';

/** @type {import('./$types').PageServerLoad} */
export async function load({ params }) {
	const { fetchRooms } = useRooms();
	const roomList = await fetchRooms();
	return {
		rooms: roomList,
	};
}

/** @param {string} cookieString */
function decodeBase64Cookie(cookieString) {

	const match = cookieString.match(/playerInfo=([^;]+)/);

	if (!match) return null;

	const decodedValue = atob(match[1])
	return JSON.parse(decodedValue)
}

/**  @satisfies {import('./$types').Actions} */
export const actions = {
	join: async ({ cookies, request }) => {
		const data = await request.formData();
		console.log("Join request made", data)

		const roomId = data.get('rooms');
		const playerName = String(data.get('name'));

		let respData = null
		const url = `${apiURL}/rooms/${roomId}/join`
		try {
			let response = await fetch(url, {
				method: "POST",
				credentials: "include", // Ensures cookies are included in the request
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					"player_name": playerName
				})
			});

			if (!response.ok) {
				throw new Error(`HTTP error! Status: ${response.status}`);
			}

			// Parse the response JSON body
			respData = await response.json();

			console.log("Response Data:", respData);

		} catch (error) {
			console.error('Failed to join room:', error);
		}

		console.log("Join Response from Join Action:", respData);

		if (!playerName) {
			throw new Error("Player Name not set")
		}

		const playerId = respData.data.player_id

		cookies.set('player_id', playerId, { secure: false, path: '/' })
		cookies.set('player_name', playerName, { secure: false, path: '/' })

		redirect(303, `/games/${roomId}`)
	},
	create: async ({ cookies, request }) => {

		const data = await request.formData();

		const url = `${apiURL}/rooms`

		let respData = null
		try {
			let response = await fetch(url, {
				method: "POST",
				credentials: "include", // Ensures cookies are included in the request
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					"host_name": data.get("name"),
					"room_name": data.get("room_name"),
				})
			});

			if (!response.ok) {
				throw new Error(`HTTP error! Status: ${response.status}`);
			}

			// Parse the response JSON body
			respData = await response.json();

			console.log("Response Data:", respData);

		} catch (error) {
			console.error('Failed to create room:', error);
		}


		const roomId = respData.data.room_id ?? null;
		const roomName = respData.data.room_name ?? null;
		const playerName = respData.data.host_name ?? null

		if (!playerName) {
			throw new Error("Player Name not set")
		}

		const playerId = respData.data.host_id

		cookies.set('player_id', playerId, { secure: false, path: '/' })
		cookies.set('player_name', playerName, { secure: false, path: '/' })

		return redirect(303, `games/${roomId}`)
	}
}
