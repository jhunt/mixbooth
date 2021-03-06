<template>
  <div id="wrap"><div id="app">
    <playlist :stream="stream" />
  </div></div>
</template>

<script>
import playlist from './components/playlist.vue'

export default {
  name: 'App',
  components: {playlist},
  computed: {
    stream() {
      return `${document.location.protocol}//${document.location.hostname}:8000/pirate-radio.opus`
    }
  }
}
</script>

<style>
body {
  margin: 0; padding: 0;
  background-color: #111;
  background-image: url(./assets/bg.jpg);
  background-size: cover;
  background-repeat: no-repeat;
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
}
#wrap {
  background-color: rgba(0,0,0, 0.5);
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin: 60px auto 0 auto;
  padding: 2em 1em;
  max-width: 600px;
  background-color: rgba(250, 250, 250, 0.7);
}
</style>
