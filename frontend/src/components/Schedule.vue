<template>
  <vs-card >
    <div slot="header" style="text-align: left">
      <h3>GCL Schedule on <span style="text-decoration: underline">SW{{selectedSwitchID}}</span></h3>
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

export default {
  components: {
    ECharts,
  },
  data() {
    return {
      selectedSwitchID: 0,
      option: {
        grid: {
          top:"14%",
          bottom:"2%",
          left:"6%",
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
          },
          data:[],

        },
        yAxis: {
          type: "category",
          name:"Ports",
          inverse:"true",
          nameLocation:"center",
          nameGap: 70,
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
          type:"piecewise",
          // show:false,
          splitNumber: 8,
          precision:0,
          inRange: {
            color: ['darkred', 'orange', 'yellow'],
            // symbolSize: [30, 100]
          },
          pieces: [],
          min: 0,
          max: 8,
          calculable: true,
          orient: 'horizontal',
          left: 'right',
          top: '-5px',
          // left: "70%",
          textStyle: {
            fontSize:13,
          },
          formatter: function(v1) {
            return v1
          }
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
      // var colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc']
      // for (var a=0;a<8;a++) {
      //   this.option.visualMap.pieces.push(
      //     {value:a, color:colors[a]},
      //   )
      // }
      this.$api.get(`/api/switch/${this.selectedSwitchID}`).
      then((res)=>{
        if (!res.data.flag)
          return
        
        for (var t=0;t<res.data.data.GCL[0].length;t++) {
          this.option.xAxis.data.push(t)
        }
        for (var i=0;i<res.data.data.Neighbors.length;i++) {
          var gate = res.data.data.GCL[i]
          var neighbor = `to ${res.data.data.Neighbors[i]}`
 
          this.option.yAxis.data.push(neighbor)
          for (var j=0;j<gate.length;j++) {
            var w = gate[j]
            this.option.series[0].data.push({
              value:[j,neighbor,w.queue+1],
            })
          }
        }
      })
    
    }
  },
  mounted() {
    this.draw()
    this.$EventBus.$on("selectedNode", (node)=>{
      if (node.substr(0,2)=="SW") {
        this.selectedSwitchID = node.match(/[\d]+/)[0]
        this.draw()
      }
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