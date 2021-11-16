<template>
  <vs-card id="topology">
    <div slot="header" style="text-align: left">
       <vs-row vs-type="flex" vs-justify="space-between">
        <vs-col vs-w="3"> 
          <h3>Topology</h3>
        </vs-col>
        <vs-col vs-w="4" vs-type="flex" vs-justify="flex-end"> 
          <vs-button id="viewBt" size="small" :color="viewActiveOnly?'rgb(255, 130, 0)':'success'" icon-pack="fas" type="filled" icon="fa-eye" @click="toggleViewOption">
              {{viewActiveOnly?"Active only":"All paths"}}
            </vs-button>
          
        </vs-col>
      </vs-row>
    </div>
    <!-- <configuration v-show="showOption"/> -->
    <ECharts id="chart" ref="topo" :options="option" autoresize @click="handleClick" />

    <vs-prompt
      title="Set link properties"
      @cancel="closePrompt"
      @accept="acceptPrompt"
      @close="closePrompt"
      :active.sync="activePrompt"
    >
      <div class="link-prompt">
        <vs-row vs-align="center" vs-type="flex" vs-justify="center" vs-w="12">
          <vs-col vs-w="1">
            <span>Loss</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input v-model="tmpLoss" />
          </vs-col>
        </vs-row>

        <vs-row vs-align="center" vs-type="flex" vs-justify="center" vs-w="12">
          <vs-col vs-w="1">
            <span>Failed</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-switch v-model="tmpFailure" />
          </vs-col>
        </vs-row>
        <vs-row
          style="margin-top: 3px; margin-bottom: 3px"
          vs-align="center"
          vs-type="flex"
          vs-justify="center"
          vs-w="12"
        >
          <vs-col vs-w="1">
            <span>Delay</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input placeholder="100" v-model="tmpDelay" />
          </vs-col>
        </vs-row>
        <vs-row vs-align="center" vs-type="flex" vs-justify="center" vs-w="12">
          <vs-col vs-w="1">
            <span>Bandwidth</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input placeholder="100" v-model="tmpBandwidth" />
          </vs-col>
        </vs-row>
        <!-- <vs-row vs-align="center" vs-type="flex" vs-justify="center" vs-w="12">
          <vs-col vs-w="1">
            <span>Distance</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input placeholder="100" v-model="tmpDistance" />
          </vs-col>
        </vs-row> -->
      </div>
    </vs-prompt>
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
      viewActiveOnly: true,
      activeNodes: [],
      showOption: false,
      activePrompt: false,
      selectedLink: "",
      tmpFailure: 0,
      tmpLoss: 0,
      tmpDelay: 1,
      tmpBandwidth: 1,
      tmpDistance: 1000,
      linkStats: {},
      option: {
        tooltip: {
          trigger: "item",
          enterable: true,
          formatter: (item) => {
            if (item.name.indexOf(">") > 0) {
              // is a link
              var units = { loss: "%", bandwidth: " Gbps", delay: " us" };
              var link = this.linkStats[item.name];
              if (item.name.indexOf("GCC") > 0) {
                units.delay = " s";
                units.bandwidth = " bps";
              }
              return (
                "Loss: " +
                link.Loss +
                units.loss +
                "<br>Delay: " +
                link.Delay +
                units.delay +
                "<br>Bandwidth: " +
                link.Bandwidth +
                units.bandwidth
              );
            }
            return item.data.name;
          },
          // alwaysShowContent: true,
          // hideDelay:1000
        },
        series: [
          {
            type: "graph",
            layout: "none",
            zoom: 1.05,
            symbolSize: 40,
            lineStyle: {
              width: 2.2,
              color: "#555",
            },
            emphasis: {
              lineStyle: {
                width: 5,
                color: "#000",
              },
              label: {
                show: true,
              },
            },
            label: {
              show: true,
              fontSize: 12.5,
            },
            center: [600,130],
            nodes: [
              {
                name: "GCC",
                x: -350,
                y: 638,
                // symbolSize: 55,
                itemStyle: {
                  color: "purple",
                },
              },
              {
                name: "HMS",
                x: 80,
                y: 450,
                itemStyle: {}
              },

              {
                name: "STR",
                x: 200,
                y: -375,
                itemStyle: {}
              },

              {
                name: "PWR",
                x: 750,
                y: -600,
                itemStyle: {}
              },

              {
                name: "ECLSS",
                x: 1300,
                y: -375,
                itemStyle: {}
              },

              {
                name: "AGT",
                x: 1500,
                y: 150,
                itemStyle: {}
              },

              {
                name: "INT",
                x: 1300,
                y: 675,
                itemStyle: {}
              },
              {
                name: "EXT",
                x: 750,
                y: 900,
                itemStyle: {}
              },
              {
                name: "SW0",
                x: 750,
                y: 150,
                symbol: "rect",
                // symbolSize: 50,
                itemStyle: {
                  color: "#0079A3",
                },
              },
              {
                name: "SW1",
                x: 360,
                y: 325,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
              {
                name: "SW2",
                x: 425,
                y: -150,

                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
              {
                name: "SW3",
                x: 750,
                y: -300,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
              {
                name: "SW4",
                x: 1075,
                y: -150,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
              {
                name: "SW5",
                x: 1200,
                y: 150,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
              {
                name: "SW6",
                x: 1075,
                y: 450,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
              {
                name: "SW7",
                x: 750,
                y: 600,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
              },
            ],
            links: [
              {
                source: "HMS",
                target: "GCC",
                lineStyle: {
                  type: "dashed",
                  width: 2.5
                },
                emphasis: {
                  lineStyle: {
                    type: "dashed",
                  },
                },
              },

              {
                source: "HMS",
                target: "SW1",
                lineStyle: {}
              },
              {
                source: "HMS",
                target: "SW2",
                lineStyle: {}
              },
              {
                source: "HMS",
                target: "SW7",
                lineStyle: {}
              },

              {
                source: "STR",
                target: "SW2",
                lineStyle: {}
              },
              {
                source: "STR",
                target: "SW1",
                lineStyle: {}
              },
              {
                source: "STR",
                target: "SW3",
                lineStyle: {}
              },
              
              {
                source: "PWR",
                target: "SW3",
                lineStyle: {}
              },
              {
                source: "PWR",
                target: "SW2",
                lineStyle: {}
              },
              {
                source: "PWR",
                target: "SW4",
                lineStyle: {}
              },
              {
                source: "ECLSS",
                target: "SW4",
                lineStyle: {}
              },
              {
                source: "ECLSS",
                target: "SW3",
                lineStyle: {}
              },
              {
                source: "ECLSS",
                target: "SW5",
                lineStyle: {}
              },

              {
                source: "AGT",
                target: "SW5",
                lineStyle: {}
              },
              {
                source: "AGT",
                target: "SW4",
                lineStyle: {}
              },
              {
                source: "AGT",
                target: "SW6",
                lineStyle: {}
              },
              {
                source: "INT",
                target: "SW6",
                lineStyle: {}
              },
              {
                source: "INT",
                target: "SW5",
                lineStyle: {}
              },
              {
                source: "INT",
                target: "SW7",
                lineStyle: {}
              },

              {
                source: "EXT",
                target: "SW7",
                lineStyle: {}
              },
              {
                source: "EXT",
                target: "SW1",
                lineStyle: {}
              },
              {
                source: "EXT",
                target: "SW6",
                lineStyle: {}
              },

              {
                source: "SW1",
                target: "SW2",
                lineStyle: {}
              },
              {
                source: "SW2",
                target: "SW3",
                lineStyle: {}
              },
              {
                source: "SW3",
                target: "SW4",
                lineStyle: {}
              },
              {
                source: "SW4",
                target: "SW5",
                lineStyle: {}
              },
              {
                source: "SW5",
                target: "SW6",
                lineStyle: {}
              },
              {
                source: "SW6",
                target: "SW7",
                lineStyle: {}
              },
              {
                source: "SW7",
                target: "SW1",
                lineStyle: {}
              },
              {
                source: "SW1",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW2",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW3",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW4",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW5",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW6",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW7",
                target: "SW0",
                lineStyle: {}
              },
              {
                source: "SW8",
                target: "SW0",
                lineStyle: {}
              },
            ],
          },
          {
            z: -1,
            type: "graph",
            layout: "none",
            zoom:1.05,
            symbolSize: 45,
            center: [600,120],
            label: {
              show: true,
              fontSize: 12,
              color: "black",
              // align: "",
              fontFamily: "Menlo",
            },
            itemStyle: {
              // color:"black"
              color: "transparent",
            },
            nodes: [
              {
                name: "TX:0\nRX:0\n\nGCC",
                x: -350,
                y: 558,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nHMS",
                x: 80,
                y: 370,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSTR",
                x: 200,
                y: -455,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nPWR",
                x: 750,
                y: -680,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nECLSS",
                x: 1300,
                y: -455,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nAGT",
                x: 1500,
                y: 70,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nINT",
                x: 1300,
                y: 595,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nEXT",
                x: 750,
                y: 820,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW0",
                x: 750,
                y: 70,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW1",
                x: 360,
                y: 245,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW2",
                x: 425,
                y: -230,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW3",
                x: 750,
                y: -380,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW4",
                x: 1075,
                y: -230,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW5",
                x: 1200,
                y: 70,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW6",
                x: 1075,
                y: 370,
                label: {}
              },
              {
                name: "TX:0\nRX:0\n\nSW7",
                x: 750,
                y: 520,
                label: {}
              },  
            ],
          },
        ],
      },
    };
  },
  methods: {
    initLinkStatus() {
      this.linkStatus = {};
      for (var i = 0; i < this.option.series[0].links.length; i++) {
        var link = this.option.series[0].links[i];
        var name = link.source + " > " + link.target;

        this.linkStats[name] = {
          Loss: 0,
          Delay: 1,
          Bandwidth: 1,
          distance: 30,
        };
        if (name.indexOf("GCC") > 0) {
          this.linkStats[name] = {
            Loss: 0,
            Delay: 600,
            Bandwidth: 2000,
            distance: 200000000000,
          };
        }
      }
    },
    handleClick(item) {
      if (item.name.length <= 5) return;
      this.selectedLink = item.name;
      this.activePrompt = true;
    },
    closePrompt() {
      this.activePrompt = false;
    },
    acceptPrompt() {
      // window.console.log(this.tmpDelay, this.selectedLink);
      this.linkStats[this.selectedLink].Loss = this.tmpLoss;
      // this.linkStats[this.selectedLink].Delay = this.tmpDelay;
      this.linkStats[this.selectedLink].Delay = this.tmpDelay;
      this.linkStats[this.selectedLink].Bandwidth = this.tmpBandwidth;
      const params = new URLSearchParams()
      params.append('loss', this.tmpLoss)
      params.append('distance', this.tmpDistance)
      params.append('bandwidth', this.tmpBandwidth)
      this.$api.post(`/api/link/${this.selectedLink}`, params);

      this.activePrompt = false;
    },
    toggleViewOption() {
      this.viewActiveOnly=!this.viewActiveOnly
      if (!this.viewActiveOnly) {
        this.clearHighlights()
      } else {
        this.highLightActiveNodes()
      }
      
      
    },
    clearHighlights() {
      for (var ii=0;ii<this.option.series[0].nodes.length;ii++) {
        this.option.series[0].nodes[ii].itemStyle.opacity = 1
        this.option.series[1].nodes[ii].label.show = true
      }
      for (var j=0;j<this.option.series[0].links.length;j++) {
        this.option.series[0].links[j].lineStyle.width = 2.2
      }

      this.option = JSON.parse(JSON.stringify(this.option))
    },
    highLightActiveNodes() {
      if (this.activeNodes.length==0) return

      this.clearHighlights()

      for (var i=0;i<this.option.series[0].nodes.length;i++) {
        if (this.activeNodes.indexOf(i)<0) {
          this.option.series[0].nodes[i].itemStyle.opacity = 0.1
          this.option.series[1].nodes[i].label.show = false
        }
      }
      for (var j=0;j<this.option.series[0].links.length;j++) {
        var link = this.option.series[0].links[j]
        if (this.activeNodes.indexOf(link.source)<0 || this.activeNodes.indexOf(link.target)<0) {
          link.lineStyle.width = 0.1
        }
      }
    }
  },
  mounted() {
    window.topo = this;
    this.initLinkStatus();
    var nameIdxMap = {
      GCC: { idx: 0, name: "GCC" },
      HMS: { idx: 1, name: "HMS" },
      STR: { idx: 2, name: "STR" },
      PWR: { idx: 3, name: "PWR" },
      ECLSS: { idx: 4, name: "ECLSS" },
      AGT: { idx: 5, name: "AGT" },
      INT: { idx: 6, name: "INT" },
      EXT: { idx: 7, name: "EXT" },
      SW0: { idx: 8, name: "SW0" },
      SW1: { idx: 9, name: "SW1" },
      SW2: { idx: 10, name: "SW2" },
      SW3: { idx: 11, name: "SW3" },
      SW4: { idx: 12, name: "SW4" },
      SW5: { idx: 13, name: "SW5" },
      SW6: { idx: 14, name: "SW6" },
      SW7: { idx: 15, name: "SW7" },
    };
    this.$EventBus.$on("stats_comm", (stats) => {
      
      var tmpActiveNodes = []
      for (var i in stats) {
        var newStatsString = "TX:" +
          stats[i][0] +
          "\nRX:" +
          stats[i][1] +
          "\n\n" +
          nameIdxMap[i].name; 
        if (this.option.series[1].nodes[nameIdxMap[i].idx].name != newStatsString) {
          tmpActiveNodes.push(nameIdxMap[i].idx,nameIdxMap[i].name)
        }
        if (tmpActiveNodes.length>0) {
          this.activeNodes = tmpActiveNodes
        }
        this.option.series[1].nodes[nameIdxMap[i].idx].name = newStatsString
      }
    });
  },
  watch: {
    activeNodes: function() {
      if (this.viewActiveOnly)
        this.highLightActiveNodes()
    }
  }
};
</script>

<style scoped>
#topology {
  width: 100%;
  /* height: 830px; */
}
#chart {
  width: 100%;
  height: 550px;
}
.link-prompt {
  font-size: 1rem;
}
#viewBt {
  width:105px;
  font-size: .75rem; 
  font-weight: 600;
}
/* .btnx {
  font-size: 2rem;
  border-radius: 5px 0px 0px 5px;
}
.btn-drop {
   border-radius: 0px 5px 5px 0px;
} */
</style>