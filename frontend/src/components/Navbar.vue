<template>
  <div>
    <vs-navbar v-model="activeItem" class="navbarx" >
      
      <div slot="title" id="title">
        <vs-navbar-title>
         <h4> MCVT - Communication Network</h4>
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
      <vs-navbar-item>
 
        <vs-row vs-align="center"
          vs-type="flex" vs-w="12">
          <vs-col vs-w="1.5">
            <h4>Mode:</h4>
          </vs-col>
          <vs-col vs-offset="2.5" vs-w="4">
            <vs-select id="mode-select" v-model="mode" >
              <vs-select-item :key="index" :value="item.value" :text="item.text" v-for="item,index in modes" />
            </vs-select>
         </vs-col>
        </vs-row>
      </vs-navbar-item>

      <vs-navbar-item id="uptime" index="2">
        Uptime: {{d}} day<span v-if="d>1">s</span> {{h.toString().padStart(2,'0')}}:{{m.toString().padStart(2,'0')}}:{{s.toString().padStart(2,'0')}}
      </vs-navbar-item>
    </vs-navbar>
  </div>
</template>
<script>
export default {
  data:()=>({
    modes:[{text:"Simulation", value:"Simulation"}, {text:"Deployment", value:"Deployment"}],
    mode: "Simulation",
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
      this.$api.get("/api/boottime")
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
    option() {

    }
  },
  
  mounted() {
    // this.getStarted()
    this.getUpTime()
  },
  watch: {
    mode:  function(val) {
      this.$EventBus.$emit("mode", val)
    }
  }
}
</script>

<style scoped>
.navbarx {
  height: 50px;
  margin-top: -30px;
  margin-bottom: 20px;
}
#title {
  margin-left: 40px;
}
#uptime {
  font-size: 1rem;
  font-weight: 600;
  margin-left: 20px;
  margin-right: 30px;
}
#button {
  margin-right: 15px;
}
#mode {
  font-size: 1.5rem;
}
#mode-select {
  width: 110px;
}

</style>