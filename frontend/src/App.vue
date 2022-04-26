<template>
  <div id="app">
    <Navbar/>
    <Sidebar/>
    <vs-row vs-w=12>
      <vs-col vs-offset=0.8 vs-w=5.2  style="z-index:99">
        <Configuration/>
        <Topology/>
        
      </vs-col>
      <vs-col vs-offset=0.2 vs-w=5.4  style="z-index:99">
        <Flow v-if="mode=='Simulation'"/>
        <Fault/>
        <Console name="comm" height="308px"/>
      </vs-col>
    </vs-row>
    <vs-row vs-w="12">
      <vs-col vs-offset=0.8 vs-w="6.8">
        <Statistics/>
        
      </vs-col>
      <vs-col vs-offset=0.2 vs-w="3.8">
        <Traffic/>
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
import Flow from './components/Flow.vue';
import Statistics from './components/Statistics.vue';
import Sidebar from './components/Sidebar.vue';
import Traffic from './components/Traffic.vue';

export default {
  name: 'App',
  components: {
    Topology,
    Console,
    Navbar,
    Flow,
    Fault,
    Configuration,
    Statistics,
    Sidebar,
    Traffic,
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
