<template>
  <div class="json-viewer">
    <pre><code v-html="highlightedJSON" /></pre>
  </div>
</template>

<script>
import highlight from "json-format-highlight";

export default {
  props: {
    json: { default: () => null },
  },
  computed: {
    highlightedJSON() {
      return highlight(this.json);
    },
  },
};
</script>

<style lang="scss" scoped>
.json-viewer {
  padding: 1em;
  overflow: scroll;

  border-radius: 0.25em;
  background-color: rgb(245, 245, 245);

  // prettier-ignore
  pre { margin: 0; }
}
</style>
