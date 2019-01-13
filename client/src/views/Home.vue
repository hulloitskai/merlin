<template>
  <div class="home flex">
    <div class="header flex">
      <div>
        <h1 class="title">merlin</h1>
        <p class="description">
          A system for accessing company finance data from EDGAR.
        </p>
      </div>
    </div>
    <div class="controls">
      <div class="grid">
        <div class="cik unit flex">
          <h4>CIK:</h4>
          <input v-model="cik" type="text" placeholder="1318605" />
        </div>
        <div class="acc-num unit flex">
          <h4>Accession Number:</h4>
          <input
            v-model="accNum"
            type="text"
            placeholder="0001564590-18-002956"
          />
        </div>
        <div class="buttons row flex">
          <div class="group">
            <button @click.prevent="loadSheets">Balance Sheet</button>
            <button @click.prevent="loadNotes">Financial Notes</button>
          </div>
        </div>
      </div>
    </div>
    <div class="viewer">
      <div v-if="apiData" class="group">
        <h3 class="title">API Response:</h3>
        <json-viewer :json="apiData" />
      </div>
      <div v-if="loading" class="spinner-container flex">
        <fulfilling-bouncing-circle-spinner
          animation-duration="infinite"
          :size="30"
          color="#ff1d5e"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { FulfillingBouncingCircleSpinner } from "epic-spinners";

import JSONViewer from "@/components/JSONViewer";
import APIClient from "@/utils/APIClient";

export default {
  data() {
    const api = new APIClient();
    return {
      cik: undefined,
      accNum: undefined,
      apiData: undefined,
      loading: false,
      api,
    };
  },
  methods: {
    validate() {
      const { cik, accNum } = this;
      if (!cik || !accNum) {
        alert("both cik and accession number fields must be filled out.");
        return false;
      }
      return true;
    },
    loadSheets() {
      const { validate, awaitData, api, cik, accNum } = this;
      if (!validate()) return;
      awaitData(api.getBalanceSheet(cik, accNum));
    },
    async loadNotes() {
      const { validate, awaitData, api, cik, accNum } = this;
      if (!validate()) return;
      awaitData(api.getFinancialNotes(cik, accNum));
    },
    async awaitData(promise) {
      this.apiData = undefined;
      this.loading = true;
      try {
        const data = await promise;
        this.apiData = data;
      } catch (err) {
        this.apiData = { error: `API request failed: ${err}` };
      } finally {
        this.loading = false;
      }
    },
  },
  components: {
    "json-viewer": JSONViewer,
    FulfillingBouncingCircleSpinner,
  },
};
</script>

<style lang="scss" scoped>
@import "@/styles/mixins.scss";

.home {
  flex-direction: column;
  align-items: stretch;

  color: rgb(80, 80, 80);
}

.header {
  padding: 1em;
  justify-content: center;

  .title {
    color: rgb(60, 60, 60);
    margin-bottom: 0.25em;
  }

  @include breakpoint(phablet) {
    margin-top: 1.25em;
  }
}

.controls {
  padding: 1em;

  > div {
    max-width: 32em;
    margin: auto;
    padding: 1em;
    grid-template: repeat(3, auto) / repeat(2, auto);

    border-radius: 0.25em;
    background-color: rgb(70, 70, 70);
    color: white;

    // prettier-ignore
    .row, .unit { grid-column: 1 / 3; }

    .unit {
      margin-bottom: 0.75em;
      flex-direction: column;
    }

    @include breakpoint(phablet) {
      padding: 1em 0.8em;

      .unit {
        grid-column: auto;
        grid-row: 1/ 3;

        margin: 0 0.5em 0.75em 0.5em;
      }
    }
  }

  .buttons {
    align-items: center;
    justify-content: center;

    button {
      $button-color: rgb(221, 221, 221);

      padding: 0.2em 0.5em;
      margin: 0.5em;

      outline: none;
      border: none;
      border-radius: 0.3em;
      background-color: $button-color;
      color: rgb(66, 66, 66);
      font-size: 11pt;
      cursor: pointer;

      transition: background 100ms ease-in-out;

      &:hover {
        background-color: lighten($button-color, 20%);
      }
    }
  }

  input {
    $input-color: rgb(73, 73, 73);

    margin-top: 0.3em;
    box-sizing: border-box;

    font-size: 11pt;
    font-weight: 500;

    color: $input-color;

    &::placeholder {
      color: lighten($input-color, 30%);
      font-weight: 400;
    }
  }
}

input {
  width: 100%;
  border-radius: 0.15em;
  padding: 0.2em 0.35em;
  box-sizing: border-box;

  border: none;
  outline: none;
}

.viewer {
  margin-top: 0.5em;
  padding: 1em;

  .group {
    margin: auto;
    max-width: 50em;
  }

  .spinner-container {
    justify-content: center;
    animation: fadein 1s;

    @keyframes fadein {
      // prettier-ignore
      from { opacity: 0; }
      // prettier-ignore
      to { opacity: 1; }
    }
  }

  // prettier-ignore
  .title { margin-bottom: 0.5em; }
}
</style>
