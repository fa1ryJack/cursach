<script setup>
import axios from "axios";
import { ref } from "vue";
import LikesCircles from "./LikesCircles.vue";

const userURL = ref("");
const loading = ref(false);
const loaded = ref(false);
const problem = ref(false);
const likesData = ref(null);

function handleSubmit(e) {
  e.preventDefault();
  loaded.value = false;
  loading.value = true;
  axios
    .get("http://localhost:8080/data?profile=" + userURL.value.trim())
    .then((res) => {
      likesData.value = res.data;
    })
    .catch((error) => {
      console.log("error!!!!", error);
      loading.value = false;
      loaded.value = false;
      problem.value = true;
    })
    .finally(() => {
      loading.value = false;
      loaded.value = true;
      problem.value = false;
    });
}
</script>

<template>
  <main>
    <form autocomplete="off">
      <label for="userURL">Paste user's link here:</label>
      <input
        placeholder="https://soundcloud.com/some-user-url"
        type="text"
        id="userURL"
        required
        v-model="userURL"
      />
      <button v-on:click="handleSubmit">Go!</button>
    </form>
    <div class="loader" v-if="loading"></div>
    <h2 class="error" v-if="problem">It seems an error has occured.</h2>
    <div class="stats-wrapper" v-if="loaded">
      <LikesCircles :likes-data="likesData" />
    </div>
  </main>
</template>

<style lang="scss" scoped>
@use "../assets/colors";

main,
.stats-wrapper {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.stats-wrapper {
  align-items: center;
  align-content: center;
  justify-content: center;
  width: 100%;
}

form {
  width: 100%;
  display: flex;
  align-items: center;
  align-content: center;
  justify-content: center;
  gap: 20px;

  input {
    width: 40%;
    padding: 6px 5px;
    border: 2px solid colors.$grey;
    border-radius: 5px;

    &:focus {
      border: 2px solid colors.$orange;
    }
  }

  button {
    padding: 10px 25px;
    cursor: pointer;
    color: white;
    background-color: colors.$black;
    border: 2px solid white;
    border-radius: 15px;

    &:hover {
      border-color: colors.$orange;
      background-color: colors.$orange;
    }
  }
}

@media (max-width: 740px) {
  form {
    flex-direction: column;
    align-items: center;
    align-content: center;

    input {
      width: 80%;
    }
  }
}

.loader {
  width: 50px;
  aspect-ratio: 1;
  border-radius: 50%;
  background: radial-gradient(farthest-side, colors.$orange 94%, #0000) top/8px
      8px no-repeat,
    conic-gradient(#0000 30%, colors.$orange);
  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);
  mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);
  animation: l13 1s infinite linear;
}
@keyframes l13 {
  100% {
    transform: rotate(1turn);
  }
}
</style>
