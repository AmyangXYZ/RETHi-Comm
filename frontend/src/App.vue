<template>
  <div id="app">
    <Navbar/>
    <vs-row vs-w=12>
      <vs-col vs-offset=0.6 vs-w=5.4>
        <configuration/>
        <Comm/>
        <Console name="comm" height="168px"/>
      </vs-col>
      <vs-col vs-offset=0.4 vs-w=5>
        <Simulation v-if="mode=='Simulation'"/>
        <Fault/>
        
      </vs-col>
    </vs-row>
    
  </div>
</template>

<script>
import Comm from './components/Comm.vue'
import Configuration from './components/Configuration.vue';
import Console from "./components/Console.vue";
import Fault from './components/Fault.vue';
import Navbar from './components/Navbar.vue'
import Simulation from './components/Simulation.vue';

export default {
  name: 'App',
  components: {
    Comm,
    Console,
    Navbar,
    Simulation,
    Fault,
    Configuration,
  },
  data() {
    return {
      mode: "Simulation"
    }
  },
  mounted() {
    this.$EventBus.$on("mode", (m)=>{
      window.console.log(m)
      this.mode = m
    })
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
