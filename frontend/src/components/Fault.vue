<template>
  <vs-card id="fault">
    <div slot="header" style="text-align: left">
       <vs-row vs-type="flex" vs-justify="space-between">
        <vs-col vs-w="3"> 
          <h3>Fault Injector</h3>
        </vs-col>
        <vs-col vs-w="3" vs-type="flex" vs-justify="flex-end" > 
          <vs-button id="bt-apply" size="small" color="danger" type="relief"  @click="apply">Apply</vs-button>
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
                <vs-select-item :key="index" :value="item.value" :text="item.text" v-for="item,index in types" />
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
                <vs-select-item :key="index" :value="item.value" :text="item.text" v-for="item,index in components" />
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
      affected: 0,
      duration: 30,
      components: [ // switches only now
        {text:"SW0", value:0},
        {text:"SW1", value:1},
        {text:"SW2", value:2},
        {text:"SW3", value:3},
        {text:"SW4", value:4},
        {text:"SW5", value:5},
        {text:"SW6", value:6},
        {text:"SW7", value:7},
      ],
      types: [
        {text:"Failure", value:0},
        {text:"Long Delay", value:1},
        {text:"Mis-routing", value:2},
        {text:"Congestion", value:3},
      ],
    }
  },
  methods: {
    apply() {
      const params = new URLSearchParams()
      params.append('type', this.types[this.selectedType].text)
      params.append('duration', this.duration)
      this.$api.post("/api/fault/switch/"+this.affected, params)
    },
    reset() {
      this.selectType = 0
      this.affected = 0
      this.duration = 0
    }
  }
}
</script>

<style scoped>
#fault {
  /* height: 300px */
}
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