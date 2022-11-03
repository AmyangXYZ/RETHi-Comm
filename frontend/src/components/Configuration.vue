<template>
<vs-card>
  <div slot="header" style="text-align:left;">
      <h3>Configurations</h3>
    </div>
  <vs-tabs v-model="tabID" style="font-size:0.85rem">
    <vs-tab index="0" label="In-Habitat Links">
      <vs-table :data="configLinks" class="links">
        <template slot="thead">
          <vs-th> Link Type</vs-th>
          <vs-th> Bandwidth </vs-th>
          <vs-th> Speed </vs-th>
          <vs-th> Distance </vs-th>
          <vs-th> Expected Delay (100 bytes)</vs-th>
        </template>

        <template slot-scope="{ data }">
          <vs-tr :key="indextr" v-for="(tr, indextr) in data">
            <vs-td :data="data[indextr].link">
              {{ tr.link }}
            </vs-td>

            <vs-td :data="tr.bandwidth" v-if="tr.link[0] == 'G'">
              {{ tr.bandwidth }} Kbps
              <template slot="edit">
                <vs-input v-model="tr.bandwidth" class="inputx" />
              </template>
            </vs-td>

            <vs-td :data="tr.bandwidth" v-else>
              {{ tr.bandwidth }} Gbps
              <template slot="edit">
                <vs-input v-model="tr.bandwidth" class="inputx" />
              </template>
            </vs-td>

            <vs-td :data="tr.speed">
              {{ tr.speed }}
            </vs-td>

            <vs-td :data="tr.distance">
              {{ tr.distance }} m
              <template slot="edit">
                <vs-row>
                  <vs-col vs-w="12" style="text-align: center; font-size: 0.95rem">
                  </vs-col>
                  <vs-col>
                    <vs-slider :min=1 text-fixed="m" v-model="tr.distance" />
                  </vs-col>
                </vs-row>
              </template>
            </vs-td>

            <vs-td :data="tr.delay"> {{ wiredDelay.toFixed(2) }} ns </vs-td>
          </vs-tr>
        </template>
      </vs-table>
    </vs-tab>

    <vs-tab index="1" label="Ground-Habitat Links">
      <div id="chart-mars-distance"></div>
      <vs-row vs-align="center">
        <vs-col vs-offset="0.4" vs-w="5.3">
          <div id="chart-orbits"></div>
        </vs-col>
        <vs-col vs-offset="1" vs-w="5">
          <vs-row vs-align="center">
            <vs-col vs-offset="1" vs-w="4">
              <h3>Run Orbits</h3>
            </vs-col>
            <vs-col vs-offset="0.5" vs-w="2">
              <vs-button size="small" color="primary" :disabled="orbitsRunning" icon-pack="fas" type="relief" icon="fa-play" @click="orbitsRun"></vs-button>
            </vs-col>
            <vs-col vs-w="2">
              <vs-button size="small" color="success" :disabled="!orbitsRunning" icon-pack="fas" type="relief" icon="fa-stop" @click="orbitsReset"></vs-button>
            </vs-col>
          </vs-row>
          <vs-row vs-align="center" style="margin-top:20px">
            <vs-col vs-offset="-1.5">
              <h3>Ground-Habitat delay: {{currentMarsDelay}}s</h3>
            </vs-col>
            <vs-col vs-offset="2" vs-w="5" style="margin-top:10px">
                <vs-button class="priority-bts" size="medium" color="danger" type="relief"  @click="submitMarsDelay">Submit</vs-button>
            </vs-col>
          </vs-row>
        </vs-col>
      </vs-row>
    </vs-tab>

    <vs-tab index="2" label="Routing">
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          Shortest Path Routing        
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="routing"></vs-checkbox>
        </vs-col>
      </vs-row>

    </vs-tab>
    <vs-tab index="3" label="Scheduling">
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          MIMOMQ Prioirty Scheduling
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="mimomq"></vs-checkbox>
        </vs-col>
      </vs-row>
    </vs-tab>
    <vs-tab index="4" label="Reliability">
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          Auto Re-Reroute Upon Switch Failure
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="REROUTE_ENABLED"></vs-checkbox>
        </vs-col>
      </vs-row>
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          Duplication Elimination
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="DUP_ELI_ENABLED"></vs-checkbox>
        </vs-col>
      </vs-row>
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          IEEE 802.1CB Frame Replication and Elimination for Reliability 
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="FRER_ENABLED"></vs-checkbox>
        </vs-col>
      </vs-row>
    </vs-tab>

    <vs-tab index="5" label="Priorities">
      <vs-row  vs-type="flex" vs-justify="start" vs-align="center">
        <vs-col vs-w="3" :key="i" v-for="(s, i) in subsys">
          <vs-row vs-align="center">
            <vs-col vs-w="6">
              {{s.name}}
            </vs-col>
            <vs-col vs-w="4" style="font-size:0.5rem">
              <vs-input-number size="small" min="0" max="7" v-model="s.priority"/>
            </vs-col>
          </vs-row>
        </vs-col>
      </vs-row>
      <vs-row vs-type="flex" vs-justify="center" vs-align="center">
        <vs-col vs-w="1.1">
          <vs-button class="priority-bts" size="medium" color="danger" type="relief"  @click="submit">Submit</vs-button>
        </vs-col>
        <vs-col vs-offset="0.3" vs-w="1.1">
          <vs-button class="priority-bts" size="medium" color="primary" type="relief"  @click="reset">Reset</vs-button>
        </vs-col>
      </vs-row>
    </vs-tab>
  </vs-tabs>
</vs-card>  
</template>

<script>
import { debounce } from "./debounce";
import * as echarts from "echarts";
import mars from "./mars.json"

export default {
  data() {
    return {
      selectedDate: "2022-11-01",
      orbitsTimer: {},
      orbitsIdx:0,
      orbitsRunning:false,
      currentMarsDelay: 1,
      chartMarsDistance: {},
      optionMarsDistance: {
        grid:{
          top:"12%",
          left:"10%",
          right:"10%",
          bottom:"20%"
        },
        tooltip: {
          trigger: "axis"
        },
        xAxis:{
          type:"time",
          axisPointer: {
            handle: {
              margin: 35,
              size:23,
              show:true,
            }
          }
        },
        yAxis:{
          name:"Delay (s)",
          type:"value"
        },
        axisPointer: {
          type:"line",
          link: {
            xAxisIndex: "all"
          },
          
          value: "2022-11-04",
          lineStyle: {
            type:"dashed",
            width:1.5,
            color:'#888'
          },
        },
        series:[{
          name:"Delay",
          type:"line",
          symbol:"none",
          data: [],
          animation:false,
        }]
      },
      chartOrbits:{},
      optionOrbits:{
        grid:{
          top:"3%",
          bottom:"2%",
          left:"2%",
          right:"2%"
        },
        tooltip: {
          trigger: "axis"
        },
        xAxis:{
          type:"value",
          axisLabel:{
            show:false
          }
        },
        yAxis:{
          type:"value",
          axisLabel:{
            show:false
          }
        },
        series:[
          {
            name:"mars-orbit",
            type:"line",
            symbol: "none",
            // lineStyle: {type:"dashed"},
            color:"grey",
            data: [],
            animation:false,
          },
          {
            name:"earth-orbit",
            type:"line",
            color:"grey",
            // lineStyle: {type:"dashed"},
            symbol: "none",
            data: [],
            animation:false,
          },
          {
            name:"mars",
            type:"scatter",
            symbolSize:12,
            color:"red",
            data:[],
            animationDurationUpdate:400,
          },
          {
            name:"earth",
            color:"#5470c6",
            type:"scatter",
            symbolSize:12*1.5,
            data:[],
            animationDurationUpdate:400
          }
        ]
      },
      dateIdxMap:{},
      tabID: 0,
      mimomq:true,
      routing:true,
      REROUTE_ENABLED:false,
      FRER_ENABLED: false,
      DUP_ELI_ENABLED: false,
      subsys: [
        { name:"GCC", priority:1},
        { name:"HMS", priority:1},
        { name:"STR", priority:1},
        { name:"SPL", priority:1},
        { name:"PWR", priority:1},
        { name:"ECLSS", priority:1},
        { name:"AGT", priority:1},
        { name:"IE", priority:1},
        { name:"EXT", priority:1},
        { name:"DTB", priority:1},
        { name:"COORD", priority:1},
      ],
      configLinks: [
        {
          link: "In-habitat",
          bandwidth: 1, // Gbps
          speed: ".77c",
          distance: 30,
        },
        // {
        //   link: "Ground-habitat",
        //   bandwidth: 2, // Kbps
        //   speed: "c",
        //   distance: 1, // percentage
        // },
      ],
    };
  },
  methods:{
    drawDistance() {
      this.optionMarsDistance.series[0].data = []
      for (let i=0;i<mars.length;i++) {
        let point = mars[i]
        this.dateIdxMap[point.date] = i
        this.optionMarsDistance.series[0].data.push([point.date,point.delay])
      }
      this.loadPlanetsPos()
      this.chartMarsDistance.setOption(this.optionMarsDistance)
      this.chartMarsDistance.on('highlight', (item) => {
        let date = this.optionMarsDistance.series[0].data[item.batch[0].dataIndex][0]
        this.selectedDate = date
        // console.log(date)
      });
    },
    drawOrbits() {
      for (let i=0;i<2;i++) {
        this.optionOrbits.series[i].data = []
      }
      
      for (let i=0;i<mars.length;i++) {
        let point = mars[i]
        this.optionOrbits.series[0].data.push([point.mars[0],point.mars[1]])
        this.optionOrbits.series[1].data.push([point.earth[0],point.earth[1]])
      }
      this.chartOrbits.setOption(this.optionOrbits)
    },
    loadPlanetsPos() {
      let point = mars[this.dateIdxMap[this.selectedDate]]
      this.optionOrbits.series[2].data = [[point.mars[0],point.mars[1]]]
      this.optionOrbits.series[3].data = [[point.earth[0],point.earth[1]]]
    
      this.chartOrbits.setOption(this.optionOrbits)
    },
    orbitsPtrInc() {
      this.selectedDate = mars[this.orbitsIdx].date
      this.orbitsIdx++
    },
    orbitsRun() {
      this.orbitsRunning = true
      this.optionOrbits.series[2].animationDurationUpdate = 1500
      this.optionOrbits.series[3].animationDurationUpdate = 1500
      this.chartMarsDistance.setOption(this.optionMarsDistance)
      this.orbitsTimer = setInterval(() => {
        this.orbitsPtrInc()
      }, 250);
    },
    orbitsReset() {
      this.orbitsRunning = false
      this.optionOrbits.series[2].animationDurationUpdate = 400
      this.optionOrbits.series[3].animationDurationUpdate = 400
      this.selectedDate = "2022-11-01"
      this.chartMarsDistance.setOption(this.optionMarsDistance)
      clearInterval(this.orbitsTimer)
    },
    submitMarsDelay() {
      this.$api.get("/api/mars/"+this.currentMarsDelay)
    },
    distanceFormatter() {
      return 10;
    },
    submit() {
      this.$api.post("/api/priorities", this.subsys)
    },
    reset() {
      for (var i=0;i<this.subsys.length;i++) {
        this.subsys[i].priority = 1
      }
    }
  },
  mounted () {
    this.$api.get(`/api/frer`).
      then((res)=>{
        this.FRER_ENABLED = Boolean(res.data)
      })
      this.$api.get(`/api/reroute`).
      then((res)=>{
        this.REROUTE_ENABLED = Boolean(res.data)
    })
  },
  computed: {
    wiredDelay: function () {
      return (
        (this.configLinks[0].distance / (300000000 * 0.77) +
         800 / (this.configLinks[0].bandwidth * 1024 * 1024 * 1024)) *
        1000000000
      );
    },
  },
  watch: {
    wiredDelay: debounce(function () {
      const params = new URLSearchParams()
      params.append('type', 'wired')
      params.append('distance', this.configLinks[0].distance)
      params.append('bandwidth', this.configLinks[0].bandwidth)
      this.$api.post('/api/links', params);
    }, 200),
    tabID: debounce(function (id) {
      if (id==1) {
        this.chartMarsDistance = echarts.init(document.getElementById("chart-mars-distance"))
        this.chartOrbits = echarts.init(document.getElementById("chart-orbits"))

        this.drawDistance()
        this.drawOrbits()
        this.currentMarsDelay = parseInt(mars[this.dateIdxMap[this.selectedDate]].delay)
      }
    }, 100),
    selectedDate: debounce(function () {
      this.loadPlanetsPos()
      this.currentMarsDelay = parseInt(mars[this.dateIdxMap[this.selectedDate]].delay)
      this.optionMarsDistance.xAxis.axisPointer.value = this.selectedDate
      this.chartMarsDistance.setOption(this.optionMarsDistance)
    }, 50),
    // wirelessDelay: debounce(function () {
    //   const params = new URLSearchParams()
    //   params.append('type', 'wireless')
    //   params.append('distance', this.wirelessDistance*1000)
    //   params.append('bandwidth', this.configLinks[1].bandwidth)
    //   this.$api.post('/api/links', params);
    // }, 200),
    FRER_ENABLED: debounce(function () {
      this.$api.get(`/api/frer/${this.FRER_ENABLED}`);
    }, 200),
    DUP_ELI_ENABLED: debounce(function () {
      this.$api.get(`/api/dupeli/${this.DUP_ELI_ENABLED}`);
    }, 200),
    REROUTE_ENABLED: debounce(function () {
      this.$api.get(`/api/reroute/${this.REROUTE_ENABLED}`);
    }, 200),
  },
};
</script>

<style scoped>

.links {
  text-align: left;
}
.links th, td {
  font-size: 0.85rem;
}

#chart-mars-distance {
  width:100%;
  height: 240px;
}

#chart-orbits {
  width:100%;
  height: 300px;
}

.priority-bts {
  width: 100%;
}
</style>