<template>
  <div id="app"><router-view /></div>
</template>

<style lang="scss">
@import url("https://rsms.me/inter/inter-ui.css");

body {
  // Normalize spacing in body and fonts.
  padding: 0;
  margin: 0;

  // prettier-ignore
  p, h1, h2, h3, h4, h5, h6 {
    margin: 0;
    padding: 0;
  }

  // Configure fonts.
  font-family: "Inter UI", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
    Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

// prettier-ignore
.flex { display: flex; }
// prettier-ignore
.grid { display: grid; }
</style>
