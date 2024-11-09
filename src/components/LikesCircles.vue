<script setup>
const props = defineProps(["data", "user"]);
</script>

<template>
  <div>
    <h3>{{ props.user }}'s likes grouped by uploaders:</h3>
  </div>
</template>

<!--   <iframe
    width="80%"
    height="300"
    src="https://w.soundcloud.com/player/?url=https%3A//api.soundcloud.com/tracks/81345014&color=%23ff5500&auto_play=false&hide_related=false&show_comments=true&show_user=true&show_reposts=false&show_teaser=true&visual=true"
  ></iframe> -->
