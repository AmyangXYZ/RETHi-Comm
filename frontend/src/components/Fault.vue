<template>
  <vs-card v-if="switches.length>0" id="fault">
    <div slot="header" style="text-align: left">
       <vs-row vs-type="flex" vs-justify="space-between">
        <vs-col vs-w="3"> 
          <h3>Fault Injector</h3>
        </vs-col>
        <vs-col vs-w="3" vs-type="flex" vs-justify="flex-end" > 
          <vs-button id="bt-apply" size="small" color="danger" type="relief"  @click="inject">Inject</vs-button>
          <vs-button id="bt-reset" size="small" color="primary" type="relief"  @click="reset">Reset</vs-button>
        </vs-col>
      </vs-row>
    </div>

    <vs-row vs-w="11" vs-align="center" vs-type="flex" vs-justify="space-around">
      <vs-col vs-w="2.5">
        <vs-row vs-w="12" vs-align="center" >
          <vs-col  vs-w="4" class="input-label">
            Type
          </vs-col>
          <vs-col vs-w="7">
              <vs-select class="selectbox" v-model="selectedType">
                <div :key="index" v-for="item,index in types">
                    <vs-select-item :value="index" :text="item" :key="index"/>
                </div>
              </vs-select>
          </vs-col>
        </vs-row>
      </vs-col>

      <vs-col vs-w="3">
        <vs-row vs-w="12" vs-align="center" >
          <vs-col vs-w="5" class="input-label">
            Affected
          </vs-col>
          <vs-col vs-offset="0.5" vs-w="6">
              <vs-select class="selectbox" v-model="affected">
                <vs-select-item :key="index" :value="index" :text="item" v-for="item,index in switches" />
              </vs-select>
          </vs-col>
        </vs-row>
      </vs-col>

      <vs-col vs-w="3">
        <vs-row vs-w="12" vs-align="center" >
          <vs-col  vs-w="8" class="input-label">
            Duration (s)
          </vs-col>
          <vs-col vs-offset="2.85" vs-w="0.5">
              <vs-input-number v-model="duration"/>
          </vs-col>
        </vs-row>
      </vs-col>


    </vs-row>
  </vs-card>
</template>

<script>
export default {
  data() {
    return {
      selectedType: 0,
      affected: 1,
      duration: 30,
      switches: [],
      types: ["Failure","Slow","Overflow","Flooding","Mis-routing"],
    }
  },
  methods: {
    inject() {
      const params = new URLSearchParams()
      params.append('type', this.types[this.selectedType])
      params.append('duration', this.duration)
      this.$api.post("/api/fault/switch/"+this.affected, params)
      let fault_log = "Inject "+this.types[this.selectedType]+" fault on SW"+this.affected+", duration: "+this.duration+" s"
      this.$vs.notify({
        title:'Fault Injector',
        text: fault_log,
        color: "danger"
      })
      this.$EventBus.$emit("fault-log",fault_log)
    },
    reset() {
      this.selectType = 0
      this.affected = 0
      this.duration = 0
      this.$api.get("/api/fault/clear")
      this.$vs.notify({
        title:'Fault Injector',
        text:"Reset and clear all faults",
        color: "primary"
      })
      this.$EventBus.$emit("fault-log","Reset and clear all faults")
    }
  },
  mounted() {
     this.$EventBus.$on("topology", (topo)=>{
      this.switches = []
      for (var i in topo.nodes) {
        if (topo.nodes[i].substr(0,2)=="SW")
          this.switches.push(topo.nodes[i])
      }
    })
  }
}
</script>

<style scoped>

.input-label {
  font-size: 0.95rem;
}
.selectbox {
  width:110px;  
}
#bt-apply {
  border-radius: 5px 0px 0px 5px;
  font-size: .75rem;
  font-weight: 600;
}
#bt-reset {
  border-radius: 0px 5px 5px 0px;
  font-size: .75rem;
  font-weight: 600;
}
</style>