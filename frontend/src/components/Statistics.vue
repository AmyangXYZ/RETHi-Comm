<template>
  <vs-card style="height: 420px">
    <div slot="header" style="text-align: left">
      <h3>Statistics</h3>
    </div>
    <vs-row vs-w="12" vs-align="center">
      <vs-col vs-w="5">
        <ECharts
          id="pie-chart"
          ref="stats"
          :options="optionPie"
          themes="macarons"
          autoresize
        />
        <ECharts
          id="line-chart"
          ref="stats"
          :options="optionLine"
          themes="macarons"
          autoresize
        />
        
      </vs-col>
      <vs-col vs-w="3">
        
      </vs-col>
      <vs-col vs-offset="0" vs-w="3">
        <ECharts
          id="pie-chart"
          ref="stats"
          :options="optionBar"
          themes="macarons"
          autoresize
        />
      </vs-col>
    </vs-row>
  </vs-card>
</template>

<script>
import ECharts from "vue-echarts/components/ECharts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/pie";
import "echarts/lib/chart/graph";
import "echarts/lib/component/title";
import "echarts/lib/component/tooltip";

export default {
  components: {
    ECharts,
  },
  data() {
    return {
      mdoe: ["by-sender", "by-receiver"],
      trafficAmountHistory: [],
      trafficPerSubsysHistory: {},
      trafficDiffPerSubsysHistory: {},
      selectedTs: 0,
      optionLine: {
        tooltip: {
          trigger: "axis",
          formatter: (item) => {
            this.selectedTs = item[0].data[0]
            return new Date(item[0].data[0]-5000).toLocaleString('en-US', {hour12: false}).substr(12, 8)+"~"+item[0].axisValueLabel+"<br>"+item[0].data[1]+" packets"
          }
        },
        grid: {
          height: "40%",
          top: "20%",
          bottom: "10%",
        },
        xAxis: {
          type: "time",
          data: [],
        },
        yAxis: {
          name: "Packets",
          type: "value",
          boundaryGap: ["50%", "100%"],
          // min: 0,
        },
        series: [
          {
            lineStyle: {
              // color: "rgba(90,11,192,1)",
              width:3
            },
            symbolSize: 4,
            data: [],
            type: "line",
            smooth: true,
            // animation: false,
          },
        ],
      },
      optionPie: {
        color: [
          "#5470c6",
          "#91cc75",
          "#fac858",
          "#ee6666",
          "#73c0de",
          "#3ba272",
          "#fc8452",
          "#9a60b4",
          "#ea7ccc",
        ],
        title: {
          text: "Traffic Summary",
          subtext: "",
          left: "center",
          textStyle: {
            fontSize: 16,
            fontWeight: 500,
          },
        },
        tooltip: {
          trigger: "item",
        },
        legend: {
          orient: "vertical",
          left: "left",
        },
        series: [
          {
            type: "pie",
            top: "20%",
            bottom: "0",
            radius: "100%",
            startAngle: 225,
            stillShowZeroSum: true,
            label: {
              position:"inside",
              fontSize: 14,
              formatter: (item)=> {
                if(item.data.value>0)
                  return `${item.data.name}`
                return ""
              }
            },
            data: [
              { value: 1, name: "GCC" },
              { value: 1, name: "HMS" },
              { value: 1, name: "STR" },
              { value: 1, name: "PWR" },
              { value: 1, name: "ECLSS" },
              { value: 1, name: "AGT" },
              { value: 1, name: "INT" },
              { value: 1, name: "EXT" },
            ],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.7)",
              },
            },
          },
        ],
      },
      optionBar: {},
    };
  },
  methods: {
    updateLineChart() {
      this.optionLine.series[0].data.push([+new Date(),0])
      this.$EventBus.$on("stats_comm", (stats) => {
        var trafficAmount = 0;
        var timestamp = +new Date();
        this.trafficDiffPerSubsysHistory[timestamp] = {}
        for (var n in stats) {
          if (stats[n].indexOf("SW") == -1) {
            trafficAmount += stats[n][1]; // rx amount

            

            if (this.trafficPerSubsysHistory[n] == null) {
              this.trafficPerSubsysHistory[n] = [
                { ts: timestamp, traffic: stats[n][1] },
              ];
            } else {
              for (var i in this.optionPie.series[0].data) {
                if (this.optionPie.series[0].data[i].name == n) {
                  
                  this.optionPie.series[0].data[i].value =
                    stats[n][1] -
                    this.trafficPerSubsysHistory[n][
                      this.trafficPerSubsysHistory[n].length - 1
                    ].traffic;
                  this.trafficDiffPerSubsysHistory[timestamp][n] = this.optionPie.series[0].data[i].value
                  break;
                }
              }
              
              this.trafficPerSubsysHistory[n].push({
                ts: timestamp,
                traffic: stats[n][1],
              });
            }
          }
        }
        if (this.trafficAmountHistory.length > 0) {
          this.optionLine.series[0].data.push([
            timestamp,
            trafficAmount -
              this.trafficAmountHistory[this.trafficAmountHistory.length - 1]
                .traffic,
          ]);
        }
        if (this.optionLine.series[0].data.length > 50) {
          this.optionLine.series[0].data.shift();
        }
        this.trafficAmountHistory.push({
          ts: timestamp,
          traffic: trafficAmount,
        });
      });
    },
  },
  mounted() {
    this.updateLineChart();
  },
  watch: {
    selectedTs: function() {
      for (var s in this.trafficDiffPerSubsysHistory[this.selectedTs]) {
        for (var i in this.optionPie.series[0].data) {
          if (this.optionPie.series[0].data[i].name == s) {
            
            this.optionPie.series[0].data[i].value = this.trafficDiffPerSubsysHistory[this.selectedTs][s]
            window.console.log(s,this.optionPie.series[0].data[i].value)
            break
          }
        }
      }
    },
  },
};
</script>

<style scoped>
#pie-chart {
  width: 100%;
  height: 200px
}
#line-chart {
  width: 100%;
  height: 200px
}
</style>