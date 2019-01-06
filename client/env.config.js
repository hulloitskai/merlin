const { NODE_ENV } = process.env;

const opts = {};
const setOpt = (key, val) => (opts[`VUE_APP_${key}`] = val);

switch (NODE_ENV) {
  case "development":
    setOpt("API_URL", "http://localhost:3000");
    break;
  default:
}

module.exports = { opts };
