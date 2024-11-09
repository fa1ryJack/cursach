<script setup>
const props = defineProps(["data", "user"]);
</script>

<template>
  <div>
    <h3>{{ props.user }}'s likes grouped by uploaders.</h3>
  </div>
</template>
