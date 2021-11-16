<template>
  <vs-card>
    <div slot="header" style="text-align:left">
      <h3>Statistics</h3>
    </div>
    <ECharts id="chart" ref="stats" :options="option" autoresize/>
  </vs-card>
</template>

<script>
import ECharts from "vue-echarts/components/ECharts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/graph";
import "echarts/lib/component/tooltip";

export default {
  components: {
    ECharts,
  },
  data() {
    return {
      option: {
        xAxis: {
          type: 'category',
          data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            type:"line",
            data:[1,2,3]
          }
        ]
      }
    }
  }
}
</script>

<style scoped>

</style>