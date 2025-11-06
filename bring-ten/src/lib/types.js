/** @typedef {Object} GameState
 * @property {string} name
 * @property {number} position
 * @property {string} roomName
 * @property {string} hostId
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
 * @property {string} host_id
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

export { };
