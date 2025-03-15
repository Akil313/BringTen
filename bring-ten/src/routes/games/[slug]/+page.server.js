import * as db from '$lib/server/database';

/** 
 * @type {import('./$types').PageServerLoad} 
 * @param {import('@sveltejs/kit').ServerLoadEvent} event - The event object containing request data
 * @returns {Promise<{ slug: string , playerId: string | null, playerName: string | null }>} 
 * - The game slug, and player details if available.
 */
export async function load({ cookies, params }) {

	const playerId = cookies.get('player_id') ?? null;
	const playerName = cookies.get('player_name') ?? null;
	const slug = params.slug;

	console.log("Log player info from game room: ", playerId, playerName)

	if (!slug) {
		throw new Error("Slug is missing from URL parameter");
	}

	return {
		slug,
		playerId,
		playerName
	};
}
