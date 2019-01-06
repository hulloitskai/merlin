/**
 * @param {string} key
 * @returns {string}
 */
const getOpt = key => process.env[`VUE_APP_${key}`];

// Set API_URL.
let API_URL = getOpt("API_URL");
if (!API_URL) API_URL = "/api";

export { getOpt, API_URL };
