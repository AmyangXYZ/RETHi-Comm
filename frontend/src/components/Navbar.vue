<template>
  <div>
    <vs-navbar id="nbar" color="linear-gradient(to right, #0A1E38 , #2A5788)" v-model="activeItem" class="navbarx" >
      
      <div slot="title" id="title">
        <vs-navbar-title id="ntitle">
          <div>MCVT</div><div id="subtitle">- Communication Network</div>
        </vs-navbar-title>
      </div>

      <vs-navbar-item>
        
        <vs-row vs-align="center"
          vs-type="flex" vs-w="12" style="margin-bottom:18px"> 
          <vs-col vs-w="1.5" style="color:white">
            <h4>Mode:</h4>
          </vs-col>
          <vs-col vs-offset="2.5" vs-w="4">
            <vs-select id="mode-select" v-model="mode" >
              <vs-select-item :key="index" :value="item.value" :text="item.text" v-for="item,index in modes" />
            </vs-select>
         </vs-col>
        </vs-row>
      </vs-navbar-item>

      <vs-navbar-item id="uptime" style="margin-bottom:18px; color:white">
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
  /* margin-top: -30px; */
  /* margin-bottom: 20px; */
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

#nbar {
  height: 98px;
  z-index: 99;
  position: absolute;
  margin-top: -80px;
  color: white;
}
#ntitle {
  font-size: 1.5rem;
  text-align: left;
  margin-bottom:18px;
  margin-left: 60px;
}

#ntitle > #subtitle {
  font-size: 1.3rem;
}

.nbar-items {
  color: white;
}

</style>