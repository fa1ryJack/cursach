<script setup>
import axios from "axios";
import { ref, watch } from "vue";
import LikesCircles from "./LikesCircles.vue";

const userURL = ref("");
const loading = ref(false);
const loaded = ref(false);
const problem = ref(false);

function handleSubmit(e) {
  e.preventDefault();
  try {
    loading = true;
    axios
      .get("http://localhost:8080/data?profile=" + userURL.value.trim())
      .then((res) => {
        console.log(res.status);
      });
  } catch (error) {
    console.log(error);
    loading = false;
    loaded = false;
    problem = true;
  } finally {
    loading = false;
    loaded = true;
  }
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
    <h2 v-if="problem">It seems an error has occured.</h2>
    <div class="stats-wrapper" v-if="loaded">
      <LikesCircles />
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
</style>
