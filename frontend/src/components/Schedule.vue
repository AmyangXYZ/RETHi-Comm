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
          left:"3%",
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
          data:[1,2,3]
        },
        yAxis: {
          type: "category",
          name:"Ports",
          nameLocation:"center",
          nameGap: 25,
          nameTextStyle: {
            fontSize: 14,
            fontWeight:"600"
          },
          splitArea: {
            show: true
          },
          axisLabel: {
            fontSize:13,
          },
          data:[1,2,3]
        },
        visualMap: {
          type:"piecewise",
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
            data: [[1,1,1], [2,2,2]]
          }
        ]
      }
    }
  },
  mounted() {
    this.$EventBus.$on("selectedNode", (node)=>{
      this.selectedNode = node
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