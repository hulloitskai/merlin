const { NODE_ENV } = process.env;

module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: ["plugin:vue/essential", "@vue/prettier"],
  rules: {
    "no-console": NODE_ENV === "production" ? "error" : "off",
    "no-debugger": NODE_ENV === "production" ? "error" : "off",
  },
  parserOptions: {
    parser: "babel-eslint",
  },
};
