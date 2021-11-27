<template>
  <div id="app">
    <Navbar/>
    <Sidebar/>
    <vs-row vs-w=12>
      <vs-col vs-offset=0.8 vs-w=5.5  style="z-index:99">
        <Configuration/>
        <Topology/>
        
      </vs-col>
      <vs-col vs-offset=0.2 vs-w=5.1  style="z-index:99">
        <Simulation v-if="mode=='Simulation'"/>
        <Fault/>
        <Console name="comm" height="300px"/>
      </vs-col>
    </vs-row>
    <vs-row vs-w="12">
      <vs-col vs-offset=0.8 vs-w="10.8">
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
import Sidebar from './components/Sidebar.vue';

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
    Sidebar,
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
html, body {
  background-color: #EBECF1;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 80px;
}
</style>
