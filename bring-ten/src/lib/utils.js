/** @param {number} strLength */
export const randomAlphaNumeric = strLength => {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789";
	let s = '';

	for (let i = 0; i < strLength; i++) {
		s += chars.charAt(Math.floor(Math.random() * chars.length));
	}

	return s
};
