<template>
  <div id="app">
    <!-- <img alt="Vue logo" src="./assets/louie.png"> -->
    <Header flash="Welcome to Your Vue.js App" username="riley"/>
    <router-view/>
  </div>
</template>

<script>
import Header from './partials/Header.vue'

export default {
  name: 'App',
  components: {
    Header
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
