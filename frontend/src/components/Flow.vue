<template>
  <vs-card>
    <div slot="header" style="text-align: left; z-index: 20">
      <vs-row vs-type="flex" vs-justify="space-between">
        <vs-col vs-w="3">
          <h3>Flow Generator</h3>
        </vs-col>
        <vs-col vs-w="1.25" vs-type="flex" vs-justify="space-between">
          <vs-button
            size="small"
            color="primary"
            icon-pack="fas"
            type="relief"
            :disabled="started"
            icon="fa-play"
            @click="startSim"
          ></vs-button>
          <vs-button
            size="small"
            color="danger"
            icon-pack="fas"
            type="relief"
            :disabled="!started"
            icon="fa-stop"
            @click="stopSim"
          ></vs-button>
        </vs-col>
      </vs-row>
    </div>

    <vs-row vs-align="center">
      <vs-col vs-w="12">
        <vs-table :data="im" v-model="selectedFlows" multiple>
          <!-- <template slot="header">
            <vs-row>
              <vs-col vs-type="flex" vs-justify="center"> 
                <h3>
                  Interaction Matrix
                </h3>
              </vs-col>
            </vs-row>
          </template> -->
          <template slot="thead">
            <vs-th> Src \ Dst </vs-th>
            <vs-th v-for="(v,i) in im" :key="i">
              {{v.name}}
            </vs-th>
            <vs-th> Freq (Hz) </vs-th>
          </template>

          <template slot-scope="{ data }">
            <vs-tr :data="tr" :key="indextr" v-for="(tr, indextr) in data">
              <vs-td :data="data[indextr].name">
                {{ data[indextr].name }}
              </vs-td>

              <vs-td :key="i" v-for="(s, i) in tr.dst" :data="s">
                {{ s }}
                <template slot="edit">
                  <vs-select id="mode-select" v-model="tr.dst[i]">
                    <vs-select-item text="X" value="X" />
                    <vs-select-item text="-" value="-" />
                  </vs-select>
                </template>
              </vs-td>

              <vs-td :data="tr.freq">
                {{ tr.freq }}
                <template slot="edit">
                  <vs-input v-model="tr.freq" class="inputx" />
                </template>
              </vs-td>
            </vs-tr>
          </template>
        </vs-table>
      </vs-col>
    </vs-row>
  </vs-card>
</template>

<script>
export default {
  data() {
    return {
      started: false,
      selectedFlows: [],
      im: [               // 0   1   2   3   4   5   6   7   8   9  10
        {name:"GCC",   dst:["-","X","-","-","-","-","-","-","-","-","-"], freq:"0.1"},
        {name:"HMS",   dst:["X","-","X","-","-","-","-","-","-","-","X"], freq:"0.5"},
        {name:"STR",   dst:["-","X","-","X","-","X","-","X","X","-","-"], freq:"10"},
        {name:"SPL",   dst:["-","X","-","-","X","-","-","-","X","-","-"], freq:"10"},
        {name:"PWR",   dst:["-","X","X","-","-","X","X","X","X","-","-"], freq:"0.25"},
        {name:"ECLSS", dst:["-","X","-","-","X","-","-","-","X","-","-"], freq:"1"},
        {name:"AGT",   dst:["-","X","-","X","-","-","-","X","-","-","-"], freq:"5"},
        {name:"IE",    dst:["-","X","-","-","-","X","X","-","-","-","-"], freq:"0.25"},
        {name:"EXT",   dst:["-","X","-","-","-","-","X","-","-","-","-"], freq:"1"},
        {name:"DTB",   dst:["-","X","-","-","X","-","X","-","-","-","-"], freq:"1"},
        {name:"COORD", dst:["-","X","-","-","-","-","-","-","-","-","-"], freq:"0.1"},
      ]
    }
  },
  methods: {
    startSim() {
      var flowsCnt = 0;
      for (var i=0;i<this.selectedFlows.length;i++) {
        for (var j=0;j<11;j++) {
          if (this.selectedFlows[i].dst[j]=="X") {
            flowsCnt++
          }
        }
      }
      if (flowsCnt>0) {
        this.started = true
        this.$api.post("/api/flows", this.selectedFlows)
        this.$vs.notify({
          title:'Simulation',
          text:"Submitted "+ flowsCnt +" flows",
          // color:color
        })
      }
    },
    stopSim() {
      this.selectedFlows = []
      this.started = false
      this.$api.get("/api/flows/stop", this.selectedFlows)
      this.$vs.notify({
        title:'Simulation',
        text:"Stop all flows",
        color: "danger"
      })
    },
  },
  mounted() {
    this.$api.get("/api/flows/start_flag")
    .then((res)=>{
      // window.console.log(res.data)
      this.started = res.data=="0"?false:true
    })
  }
}
</script>

<style scoped>
.vs-table--content {
  overflow: hidden;
}

th,
td {
  text-align: center;
  padding-left: 10px;
  font-weight: 600;
  font-size: 0.85rem;
}

.vs-table--tbody-table .tr-values td {
  padding-top: 0px;
  padding-bottom: 0px;
}

.td-edit {
  text-decoration: none;
}
</style>