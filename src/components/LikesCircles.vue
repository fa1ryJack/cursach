<script setup>
import * as d3 from "d3";
const props = defineProps(["likesData"]);
</script>

<template>
  <div>
    <h3>{{ likesData.Name }}'s likes grouped by uploaders:</h3>
    <div ref="circles"></div>
  </div>
</template>

<!--   <iframe
    width="80%"
    height="300"
    src="https://w.soundcloud.com/player/?url=https%3A//api.soundcloud.com/tracks/81345014&color=%23ff5500&auto_play=false&hide_related=false&show_comments=true&show_user=true&show_reposts=false&show_teaser=true&visual=true"
  ></iframe> -->
