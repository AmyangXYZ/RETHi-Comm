<template>
  <div id="app">
    <Navbar/>
    <vs-row vs-w=12>
      <vs-col vs-offset=0.6 vs-w=5.4>
        <Configuration/>
        <Topology/>
        <Console name="comm" height="168px"/>
      </vs-col>
      <vs-col vs-offset=0.4 vs-w=5>
        <Simulation v-if="mode=='Simulation'"/>
        <Fault/>
        <Statistics/>
      </vs-col>
    </vs-row>
    
  </div>
</template>

<script>
import Topology from './components/Topology.vue'
import Configuration from './components/Configuration.vue';
import Console from "./components/Console.vue";
import Fault from './components/Fault.vue';
import Navbar from './components/Navbar.vue'
import Simulation from './components/Simulation.vue';
import Statistics from './components/Statistics.vue';

export default {
  name: 'App',
  components: {
    Topology,
    Console,
    Navbar,
    Simulation,
    Fault,
    Configuration,
    Statistics,
  },
  data() {
    return {
      mode: "Simulation"
    }
  },
  mounted() {
    this.$EventBus.$on("mode", (m)=>{
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
