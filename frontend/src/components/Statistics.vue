<template>
  <vs-card style="height: 390px">
    <div slot="header" style="text-align: left">
      <h3>Statistics of <span style="text-decoration: underline">{{selectedNode}}</span></h3>
    </div>
    <vs-row vs-align="center">
      <vs-col vs-offset="0.2" vs-w="2.5" id="summary" v-if="nodes_loaded">
        <span>
          Inbound: {{stats_summary[selectedNode].inbound}} pkts
        </span> <br>
        <span>
          Outbound: {{stats_summary[selectedNode].outbound}} pkts
        </span> <br>
        <span>
          Average delay: {{stats_summary[selectedNode].avg_delay}} us
        </span> <br>
        <span>
          Fault count: {{stats_summary[selectedNode].fault_cnt}}
        </span>
      </vs-col>
      <vs-col vs-w="9.3">
        <vs-tabs :value="2" >
          <vs-tab index="0" label="I/O">
            <ECharts
                id="bar-chart"
                ref="stats" 
                :options="optionIO"
                themes="macarons"
                autoresize
              />
          </vs-tab>
          <vs-tab index="1" label="Delay">
            <ECharts
                id="bar-chart"
                ref="stats" 
                :options="optionDelay"
                themes="macarons"
                autoresize
              />
          </vs-tab>
        </vs-tabs>
        
      </vs-col>
    </vs-row>
      
      
  </vs-card>
</template>

<script>
import ECharts from "vue-echarts/components/ECharts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/bar";
import "echarts/lib/chart/scatter";
import "echarts/lib/chart/graph";
import "echarts/lib/component/legend";
import "echarts/lib/component/dataZoom";
import "echarts/lib/component/tooltip";

export default {
  components: {
    ECharts,
  },
  data() {
    return {
      selectedNode: "HMS",
      stats_summary: {},
      traffic_history: {},
      nodes_loaded: false,
      nodes: [],
      refreshTimer: {},
      optionIO: {
        tooltip: {
          trigger: "axis",
        },
        grid: {
          // height: "70%",
          top: "10%",
          bottom: "20%",
          right: "5%",
          left: "5%",
        },
        legend: {
          data: ['Inbound', 'Outbound'] 
        },
        dataZoom: [
          {
            type: "inside",
            start: 0,
            end: 100,
          },
          {
            start: 0,
            end: 100,
            handleIcon:
              "M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z",
            handleSize: "80%",
            handleStyle: {
              color: "#fff",
              shadowBlur: 3,
              shadowColor: "rgba(0, 0, 0, 0.6)",
              shadowOffsetX: 2,
              shadowOffsetY: 2
            }
          }
        ],
        xAxis: {
          type: "time",
          data: [],
        },
        yAxis: {
          name: "Traffic",
          type: "value",
          boundaryGap: ["0%", "20%"],
          // min: 0,
        },
        series: [
          {
            name: 'Inbound',
            type: 'bar',
            // barGay: "30%",
            barCategoryGap: "40%",
            data:[],
            itemStyle: {
              color: "rgba(27,113,222,1)"
            }
          },
          {
            name: 'Outbound',
            type: 'bar',
            // barGay: "30%",
            barCategoryGap: "40%",
            data:[],
            itemStyle: {
              color: "rgba(243,171,71,1)"
            }
          },
        ],
      },
      optionDelay: {
        tooltip: {
          trigger: "axis",
        },
        grid: {
          // height: "20%",
          top: "10%",
          bottom: "25%",
          right: "5%",
          left: "10%",
        },
        dataZoom: [
          {
            type: "inside",
            start: 50,
            end: 100,
          },
          {
            start: 50,
            end: 100,
            handleIcon:
              "M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z",
            handleSize: "60%",
            handleStyle: {
              color: "#fff",
              shadowBlur: 3,
              shadowColor: "rgba(0, 0, 0, 0.6)",
              shadowOffsetX: 2,
              shadowOffsetY: 2
            }
          }
        ],
        xAxis: {
          type: "time",
          data: [],
        },
        yAxis: {
          name: "Delay",
          type: "value",
          // boundaryGap: ["0%", "10%"],
          // min: 0,
        },
        series: [
        ],
      },
    };
  },
  methods: {
    drawStatsIO() {
      this.$api.get("/api/stats/"+this.selectedNode+"/io")
      .then((res)=>{
        if (res.data.flag==0) return
        this.optionIO.series[0].data = []
        this.optionIO.series[1].data = []
        for (var i in res.data.data) {
          var stat = res.data.data[i]
          var diffRx = 0
          var diffTx = 0
          if (i>0) {
            diffRx = stat[1]-res.data.data[i-1][1]
            diffTx = stat[2]-res.data.data[i-1][2]
          }
          if (diffRx<0) diffRx = 0
          if (diffTx<0) diffTx = 0
          this.optionIO.series[0].data.push([stat[0], diffRx])
          this.optionIO.series[1].data.push([stat[0], diffTx])
          if (i==res.data.data.length-1) {
            if (this.stats_summary[this.selectedNode]==null) {
              this.stats_summary[this.selectedNode] = {inbound: stat[1], outbound:stat[2], avg_delay:0,fault_cnt:0}
            } else {
              this.stats_summary[this.selectedNode].inbound = stat[1]
              this.stats_summary[this.selectedNode].outbound = stat[2]
            }
            this.stats_summary = JSON.parse(JSON.stringify(this.stats_summary))
          }
        }
      })
    },
    drawStatsDelay() {
       this.$api.get("/api/stats/"+this.selectedNode+"/delay")
      .then((res)=>{
        this.optionDelay.series = []
        if (res.data.flag==0) return       
        if (this.selectedNode!="GCC") {
          this.optionDelay.yAxis.name = "Delay (us)"
        }
        var tmp = []
        for (var i=0;i<res.data.data.length;i++) {
          var entry = res.data.data[i]
          var line = {
            name: entry.source,
            type: 'line',
            data:entry.data,
            smooth:"true",
            symbol: "none",
            animation: false,
            sampling: "average"
          }
          tmp.push(line) 
        }
        this.optionDelay.series = tmp
      })
    }
  },
  mounted() {
    this.$EventBus.$on("topology", (topo)=>{
      this.nodes = topo.nodes
      for (var n in this.nodes) {
        this.stats_summary[this.nodes[n]] = {inbound:0, outbound:0, avg_delay:0, fault_cnt:0}
      }
      // window.console.log(this.stats_summary)
      this.nodes_loaded = true
      
    })
    
    this.drawStatsIO()
    this.drawStatsDelay()
    // this.refreshTimer = setInterval(this.drawStats, 10*1000)
    this.$EventBus.$on("selectedNode", (node)=>{
      // clearInterval(this.refreshTimer)
      this.selectedNode = node
      this.drawStatsIO()
      this.drawStatsDelay()
      // this.refreshTimer = setInterval(this.drawStats, 10*1000)
    })
  },
  created() {

  },
  watch: {
    // stats_summary: function() {
    // }
  },
};
</script>

<style scoped>
#summary {
  text-align: left;
  line-height: 2;
  font-size: 1rem;
  font-weight: 500;
}

#bar-chart {
  width: 100%;
  height: 280px;
}
</style>