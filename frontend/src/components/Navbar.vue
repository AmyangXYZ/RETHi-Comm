<template>
  <div>
    <vs-navbar v-model="activeItem" class="nabarx">
      
      <div slot="title" id="title">
        <vs-navbar-title>
          MCVT - Communication Network
        </vs-navbar-title>
      </div>

      <!-- <vs-navbar-item id="button">
        <vs-button size="small" color="success" icon-pack="fas" :disabled="started" type="relief" icon="fa-play" @click="startAll"></vs-button>
      </vs-navbar-item>
      <vs-navbar-item id="button">
        <vs-button size="small" color="red" icon-pack="fas" type="relief" :disabled="!started" icon="fa-stop" @click="stopAll"></vs-button>
      </vs-navbar-item>
      <vs-navbar-item id="button">
        <vs-button size="small" color="primary" icon-pack="fas" type="relief" icon="fa-cog" @click="option"></vs-button>
      </vs-navbar-item> -->

      <vs-navbar-item id="uptime" index="2">
        Uptime: {{d}} day<span v-if="d>1">s</span> {{h.toString().padStart(2,'0')}}:{{m.toString().padStart(2,'0')}}:{{s.toString().padStart(2,'0')}}
      </vs-navbar-item>
    </vs-navbar>
  </div>
</template>
<script>
import axios from 'axios'
export default {
  data:()=>({
    popupActive: false,
    activeItem: 0,
    boottime:0,
    started: false,
    s: 0,
    m: 0,
    h: 0,
    d: 0,
  }),
  methods: {
    getUpTime() {
      axios.get("http://localhost:8000/api/boottime")
      .then((res)=>{
        this.boottime = parseInt(res.data)
        setInterval(()=>{
          var now = Date.now()
          now = Math.floor(now/1000)
          this.s = Math.floor((now-this.boottime)%60)
          this.m = Math.floor((now-this.boottime)/(60)%60)
          this.h = Math.floor((now-this.boottime)/(60*60)%24)
          this.d = Math.floor((now-this.boottime)/(60*60*24))
        },1000)
      })
    },
    getStarted() {
      axios.get("http://localhost:8000/api/started")
      .then((res)=>{
        var flag = parseInt(res.data)
        this.started = (flag=="1")?true:false
      })
    },
    startAll() {
      this.started = true
      axios.get("http://localhost:8000/api/start")
      this.$vs.notify({
        title:'Run',
        // text:'biubiubiu',
        color:"primary",
        time: "4000"
      })
    },
    stopAll() {
      this.$vs.notify({
        title:'Stop',
        // text:'biubiubiu',
        color:"primary",
        time: "4000"
      })
      this.started = false
      this.$EventBus.$emit("stopAll", 1)
      axios.get("http://localhost:8000/api/stop")
    },
    option() {

    }
  },
  
  mounted() {
    // this.getStarted()
    this.getUpTime()
  }
}
</script>

<style scoped>
.nabarx {
  height: 50px;
  margin-top: -30px;
  margin-bottom: 30px;
}
#title {
  margin-left: 40px;
}
#uptime {
  font-size: 1rem;
  font-weight: 600;
  margin-left: 10px;
  margin-right: 50px;
}
#button {
  margin-right: 15px;
}
</style>