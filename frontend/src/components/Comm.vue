<template>
  <vs-card id="topology">
    <div slot="header">
      <h3>Communication Network</h3>
    </div>
    <ECharts
      id="chart"
      ref="topo"
      :options="option"
      autoresize
      @click="handleClick"
    />

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
            <span>PDR</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input placeholder="100" v-model="tmpPDR" />
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
            <span>Latency</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input placeholder="100" v-model="tmpLatency" />
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
      </div>
    </vs-prompt>
    <Console name="comm" height="180px" />
  </vs-card>
</template>

<script>
import Console from "./Console.vue";
import ECharts from "vue-echarts/components/ECharts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/graph";
import "echarts/lib/component/tooltip";
import axios from 'axios'

export default {
  components: {
    ECharts,
    Console,
  },
  data() {
    return {
      activePrompt: false,
      selectedLink: "",
      tmpPDR: 100,
      tmpLatency: 10,
      tmpBandwidth: 10,
      linkStats: {},
      option: {
        tooltip: {
          trigger: "item",
          enterable: true,
          formatter: (item) => {
            if (item.name.indexOf(">") > 0) {
              // is a link
              var link = this.linkStats[item.name];
              return (
                "PDR:" +
                link.PDR +
                "<br>Latency: " +
                link.Latency +
                "<br>Bandwidth:" +
                link.Bandwidth
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
            symbolSize: 45,
            lineStyle: {
              width: 2.5,
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
              fontSize: 14,
            },
            center: [500, 120],
            nodes: [
              {
                name: "GCC",
                x: -500,
                y: 150,
                // symbolSize: 55,
                itemStyle: {
                  color: "purple",
                },
              },
              {
                name: "HMS",
                x: 0,
                y: 150,
                // symbolSize: 55,
                itemStyle: {
                  color: "darkred",
                },
              },

              {
                name: "AGT",
                x: 200,
                y: -375,
              },

              {
                name: "STR",
                x: 750,
                y: -600,
              },

              {
                name: "INV",
                x: 1300,
                y: -375,
              },

              {
                name: "PWR",
                x: 1500,
                y: 150,
              },

              {
                name: "ECLSS",
                x: 1300,
                y: 675,
              },
              {
                name: "INT",
                x: 750,
                y: 900,
              },
              {
                name: "EXT",
                x: 200,
                y: 675,
              },
              {
                name: "SW1",
                x: 300,
                y: 150,
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
              {
                name: "SW8",
                x: 425,
                y: 450,
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                },
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
            ],
            links: [
              {
                source: "HMS",
                target: "GCC",
                lineStyle: {
                  type: "dashed",
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
              },
              {
                source: "AGT",
                target: "SW2",
              },
              {
                source: "STR",
                target: "SW3",
              },
              {
                source: "INV",
                target: "SW4",
              },
              {
                source: "PWR",
                target: "SW5",
              },
              {
                source: "ECLSS",
                target: "SW6",
              },
              {
                source: "INT",
                target: "SW7",
              },
              {
                source: "EXT",
                target: "SW8",
              },
              {
                source: "SW1",
                target: "SW2",
              },
              {
                source: "SW2",
                target: "SW3",
              },
              {
                source: "SW3",
                target: "SW4",
              },
              {
                source: "SW4",
                target: "SW5",
              },
              {
                source: "SW5",
                target: "SW6",
              },
              {
                source: "SW6",
                target: "SW7",
              },
              {
                source: "SW7",
                target: "SW8",
              },
              {
                source: "SW8",
                target: "SW1",
              },
              {
                source: "SW1",
                target: "SW0",
              },
              {
                source: "SW2",
                target: "SW0",
              },
              {
                source: "SW3",
                target: "SW0",
              },
              {
                source: "SW4",
                target: "SW0",
              },
              {
                source: "SW5",
                target: "SW0",
              },
              {
                source: "SW6",
                target: "SW0",
              },
              {
                source: "SW7",
                target: "SW0",
              },
              {
                source: "SW8",
                target: "SW0",
              },
            ],
          },
          {
            z: -1,
            type: "graph",
            layout: "none",
            symbolSize: 45,
            center: [500, 120],
            label: {
              show: true,
              fontSize: 13,
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
                name: "TX:0\nRX:0\n\ngcc",
                x: -500,
                y: 70,
              },
              {
                name: "TX:0\nRX:0\n\nhms",
                x: 0,
                y: 70,
              },
              {
                name: "TX:0\nRX:0\n\nagt",
                x: 200,
                y: -455,
              },
              {
                name: "TX:0\nRX:0\n\ninv",
                x: 750,
                y: -680,
              },
              {
                name: "TX:0\nRX:0\n\nstr",
                x: 1300,
                y: -455,
              },
              {
                name: "TX:0\nRX:0\n\npwr",
                x: 1500,
                y: 70,
              },
              {
                name: "TX:0\nRX:0\n\neclss",
                x: 1300,
                y: 595,
              },
              {
                name: "TX:0\nRX:0\n\nint",
                x: 750,
                y: 820,
              },
              {
                name: "TX:0\nRX:0\n\next",
                x: 200,
                y: 595,
              },
              {
                name: "TX:0\nRX:0\n\nsw0",
                x: 300,
                y: 70,
              },
              {
                name: "TX:0\nRX:0\n\nsw1",
                x: 425,
                y: -230,
              },
              {
                name: "TX:0\nRX:0\n\nsw2",
                x: 750,
                y: -380,
              },
              {
                name: "TX:0\nRX:0\n\nsw3",
                x: 1200,
                y: 70,
              },
              {
                name: "TX:0\nRX:0\n\nsw4",
                x: 1075,
                y: 370,
              },
              {
                name: "TX:0\nRX:0\n\nsw5",
                x: 750,
                y: 520,
              },
              {
                name: "TX:0\nRX:0\n\nsw6",
                x: 1075,
                y: -230,
              },
              {
                name: "TX:0\nRX:0\n\nsw7",
                x: 425,
                y: 370,
              },
              {
                name: "TX:0\nRX:0\n\nsw8",
                x: 750,
                y: 70,
              },
            ],
          },
        ],
      },
    };
  },
  methods: {
    initLinkStat() {
      this.linkStats = {};
      for (var i = 0; i < this.option.series[0].links.length; i++) {
        var link = this.option.series[0].links[i];
        var name = link.source + " > " + link.target;

        this.linkStats[name] = {
          PDR: 100,
          Latency: 10,
          Bandwidth: 10,
        };
        if (name.indexOf("GCC") > 0) {
          this.linkStats[name] = {
            PDR: 100,
            Latency: 10000000,
            Bandwidth: 0.000001,
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
      window.console.log(this.tmpLatency)
      this.linkStats[this.selectedLink].PDR = this.tmpPDR;
      this.linkStats[this.selectedLink].Latency = this.tmpLatency;
      this.linkStats[this.selectedLink].Bandwidth = this.tmpBandwidth;
      axios.post("http://localhost:8000/api/link/ground", {
        pdr:this.tmpPDR,
        latency:this.tmpLatency,
        bandwidth: this.tmpBandwidth
      })
      
      this.activePrompt = false;
    },
  },
  mounted() {
    window.topo = this;
    this.initLinkStat();
    var nameIdxMap = {
      GCC: { id: 0, nick: "gcc" },
      HMS: { id: 1, nick: "hms" },
      AGT: { id: 2, nick: "agt" },
      INV: { id: 3, nick: "inv" },
      STR: { id: 4, nick: "str" },
      PWR: { id: 5, nick: "pwr" },
      ECLSS: { id: 6, nick: "eclss" },
      INT: { id: 7, nick: "int" },
      EXT: { id: 8, nick: "ext" },
      SW0: { id: 9, nick: "sw0" },
      SW1: { id: 10, nick: "sw1" },
      SW2: { id: 11, nick: "sw2" },
      SW3: { id: 12, nick: "sw3" },
      SW4: { id: 13, nick: "sw4" },
      SW5: { id: 14, nick: "sw5" },
      SW6: { id: 15, nick: "sw6" },
      SW7: { id: 16, nick: "sw7" },
      SW8: { id: 17, nick: "sw8" },
    };
    this.$EventBus.$on("stats_comm", (stats) => {
      for (var i in stats) {
          this.option.series[1].nodes[nameIdxMap[i].id].name = 
            "TX:" +
            stats[i][0] +
            "\nRX:" +
            stats[i][1] +
            "\n\n" +
            nameIdxMap[i].nick;
      }
    });
  },
};
</script>

<style scoped>
#topology {
  width: 100%;
  height: 830px;
}
#chart {
  width: 100%;
  height: 570px;
}
.link-prompt {
  font-size: 1rem;
}
</style>