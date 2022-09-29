<template>
  <vs-card >
    <div slot="header" style="text-align: left">
      <h3>GCL Schedule on <span style="text-decoration: underline">{{selectedNode}}</span></h3>
    </div>
    <ECharts id="sch" autoresize :options="option"/>
  </vs-card>
</template>

<script>
import ECharts from "vue-echarts/components/ECharts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/bar";
import "echarts/lib/chart/scatter";
import "echarts/lib/chart/heatmap";
import "echarts/lib/component/visualMap";
import "echarts/lib/component/legend";
import "echarts/lib/component/dataZoom";
import "echarts/lib/component/tooltip";

const scheduleL2 = require("../assets/schedule_level2.json");

export default {
  components: {
    ECharts,
  },
  data() {
    return {
      selectedNode: "",
      option: {
        grid: {
          top:"13%",
          bottom:"2%",
          left:"4%",
          right:"2%",
        },
        xAxis: {
          position: "top",
          name:"Time",
          nameLocation:"center",
          nameGap: 25,
          nameTextStyle: {
            fontSize: 14,
            fontWeight:"600"
          },
          type: "category",
          splitArea: {
            show: true
          },
          axisLabel: {
            fontSize:13,
            formatter:(item)=>{
              return item*800
            }
          },
          data:[],

        },
        yAxis: {
          type: "category",
          name:"Ports",
          nameLocation:"center",
          nameGap: 45,
          nameTextStyle: {
            fontSize: 14,
            fontWeight:"600"
          },
          splitArea: {
            show: true
          },
          axisLabel: {
            fontSize:14,
          },
          data:[]
        },
        visualMap: {
          // type:"piecewise",
          show:false,
          pieces: [],
          min: 0,
          max: 10,
          calculable: true,
          orient: 'horizontal',
          left: 'right',
          top: '-5px',
          textStyle: {
            fontSize:13,
          },
        },
        series: [
          {
            type: "heatmap",
            emphasis: {
              itemStyle: {
                borderWidth:0,
                shadowBlur: 10,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            },
            data: [],
            animation:false
          }
        ]
      }
    }
  },
  methods: {
    draw() {
      this.option.xAxis.data = []
      this.option.yAxis.data = []
      this.option.series[0].data = []

      for (var i=0;i<=5000000/800;i++) {
        this.option.xAxis.data.push(i)
      }

      for (var link in scheduleL2) {
        var sw = link.match(/\d+/g)[0]
        if (this.selectedNode.match(/\d+/g)[0]==sw) {
          var port = `to ${link.match(/\d+/g)[1]}`
          this.option.yAxis.data.push(port)
          var windows = scheduleL2[link]
          for (var j=0;j<windows.length;j++) {
            var cell = windows[j]
            for (var t=cell[1]/800;t<cell[2]/800;t++) {
              this.option.series[0].data.push({
                value: [t,port,cell[0]+2],
                itemStyle: {},
                label: {}
              })
            }
          }
        }
      }
    }
  },
  mounted() {
    this.$EventBus.$on("selectedNode", (node)=>{
      this.selectedNode = node
      this.draw()
    })
  }
}
</script>

<style scoped>
#sch {
  width:100%;
  height: 300px;
}
</style>