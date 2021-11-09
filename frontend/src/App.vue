<template>
  <div id="app">
    <!-- <img alt="Vue logo" src="./assets/logo.png"> -->
    <Navbar/>
    <vs-row vs-w=12>
      <vs-col vs-offset=0.4 vs-w=6.8>
        <Comm/>
      </vs-col>
    </vs-row>
    
  </div>
</template>

<script>
import Comm from './components/Comm.vue'

import Navbar from './components/Navbar.vue'

export default {
  name: 'App',
  components: {
    Comm,


    Navbar,

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
  margin-top: 30px;
}
</style>
