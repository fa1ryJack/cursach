<script setup>
import axios from 'axios';
import { ref } from 'vue';

const userURL = ref('')

function handleSubmit(e){
    e.preventDefault()
    axios.get("http://localhost:8080/data?profile=" + userURL)
    .then((data) => {
        console.log(data.status)
    })
}

</script>

<template>
    <form>
        <label for="userURL">Paste user's link here:</label>
        <input type="text" id="userURL" v-model="userURL">
        <button v-on:click="handleSubmit">Go!</button>
    </form>
</template>