// Augment process.env with options from env.config.js.
const { opts } = require("./env.config");
Object.assign(process.env, opts);

const { BASE_URL } = process.env;
module.exports = {
  baseUrl: BASE_URL,
  configureWebpack: {
    resolve: {
      alias: { vue$: "vue/dist/vue.runtime.esm.js" },
    },
  },
};
