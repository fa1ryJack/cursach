<script setup></script>

<template>
  <div class="greetings">
    <h1 class="orange">Welcome!</h1>
    <h3>
      Here you can get likes grouped by uploaders from
      <a target="_blank" href="https://soundcloud.com/you/"
        >Soundcloud's user</a
      >
      by pasting their link in the box below.
    </h3>
  </div>
</template>

<style lang="scss" scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  position: relative;
  top: -10px;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}
</style>
