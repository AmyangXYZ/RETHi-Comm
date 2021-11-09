<template>
  <div id="app">
    <!-- <img alt="Vue logo" src="./assets/logo.png"> -->
    <Navbar/>
    <vs-row vs-w=12>
      <vs-col vs-offset=0.4 vs-w=6.4>
        <Comm/>
      </vs-col>
      <vs-col vs-offset=0.2 vs-w=4.6>
        <Configuration/>
        <Console name="comm" height="500px"/>
      </vs-col>
    </vs-row>
    
  </div>
</template>

<script>
import Comm from './components/Comm.vue'
import Configuration from './components/Configuration.vue';
import Console from "./components/Console.vue";
import Navbar from './components/Navbar.vue'

export default {
  name: 'App',
  components: {
    Comm,
    Console,
    Navbar,
    Configuration,

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
