<template>
  <vs-card id="topology">
    <div slot="header" style="text-align: left">
      <vs-row vs-type="flex" vs-justify="space-between">
        <vs-col vs-w="3">
          <h3>Topology</h3>
        </vs-col>
        <vs-col vs-w="4" vs-type="flex" vs-justify="flex-end">
          <div v-show="!editMode">
            <vs-button
              class="buttons"
              size="small"
              :color="viewActiveOnly ? 'rgb(255, 130, 0)' : 'success'"
              icon-pack="fas"
              type="filled"
              icon="fa-eye"
              @click="toggleViewOption"
            >
              {{ viewActiveOnly ? "Active only" : "All paths" }}
            </vs-button>
            <vs-button
              class="buttons"
              style="margin-left: 10px"
              id="edit"
              size="small"
              color="primary"
              icon-pack="fas"
              type="filled"
              icon="fa-edit"
              @click="editEnable"
            >
              Edit
            </vs-button>
          </div>
          <div v-show="editMode">
            <vs-button
              class="buttons"
              size="small"
              color="success"
              icon-pack="fas"
              type="filled"
              icon="fa-check"
              @click="editEnable"
            >
              Apply
            </vs-button>
            <vs-button
              class="buttons"
              style="margin-left: 10px"
              size="small"
              color="danger"
              icon-pack="fas"
              type="filled"
              icon="fa-undo"
              @click="editEnable"
            >
              Cancel
            </vs-button>
          </div>
        </vs-col>
      </vs-row>

      <vs-row  vs-type="flex" vs-justify="center" v-show="editMode" style="margin-top:8px">
        <vs-col vs-offset="2" vs-w="2" >
          <vs-button
            class="buttons"
            size="small"
            color="success"
            icon-pack="fas"
            type="filled"
            icon="fa-plus"
            @click="addSwitch"
          >
            Add
          </vs-button>

          
        </vs-col>
        <vs-col vs-offset="1" vs-w="2">
          <vs-button
            class="buttons"
            size="small"
            color="primary"
            icon-pack="fas"
            type="filled"
            icon="fa-arrows-alt-h"
            @click="connect"
          >
            Connect
          </vs-button>
        </vs-col>
        <vs-col vs-offset="-0.5" vs-w="2">
          <vs-select
            class="conenct-select"
            v-model="connectHost0"
          >
            <vs-select-item :key="index" :value="item.value" :text="item.text" v-for="item,index in nodes" />
          </vs-select>
        </vs-col>
        <vs-col vs-offset="-0.2" vs-w="2">
          <vs-select
            class="conenct-select"
            v-model="connectHost1"
          >
            <vs-select-item :key="index" :value="item.value" :text="item.text" v-for="item,index in nodes" />
          </vs-select>
        </vs-col>
      </vs-row>
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
            <span>Loss</span>
          </vs-col>
          <vs-col vs-offset="2" vs-w="5">
            <vs-input v-model="tmpLoss" />
          </vs-col>
        </vs-row>

        <vs-row vs-align="center" vs-type="flex" vs-justify="center" vs-w="12">
          <vs-divider />
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
import "echarts/lib/component/graphic";

export default {
  components: {
    ECharts,
  },
  data() {
    return {
      editMode: false,
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
      switchCnt: 7,
      nodes: [
        {text:"GCC", value:0},
        {text:"HMS", value:1},
        {text:"STR", value:2},
        {text:"PWR", value:3},
        {text:"ECLSS", value:5},
        {text:"AGT", value:5},
        {text:"INT", value:6},
        {text:"EXT", value:7},
        {text:"SW0", value:8},
        {text:"SW1", value:9},
        {text:"SW2", value:10},
        {text:"SW3", value:11},
        {text:"SW4", value:12},
        {text:"SW5", value:13},
        {text:"SW6", value:14},
        {text:"SW7", value:15},
      ],
      connectHost0: "",
      connectHost1: "",
      tooltipDefault: {
        trigger: "item",
        enterable: true,
        formatter: (item) => {
          if (item.dataType == "edge") {
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
          return null;
        },
      },
      tooltipEdit: {
        formatter: (item) => {
          if (item.dataType == "node")
            return "X: " + item.value[0] + "<br>Y: " + item.value[1];
          else return null;
        },
      },
      option: {
        tooltip: {},
        grid: {
          right: "1%",
          left: "1%",
          top: "1%",
          bottom: "1%",
        },
        xAxis: {
          type: "value",
          position: "top",
          // min:-500,
          max: 2050,
          interval: 50,
          axisTick: {
            show: false,
          },
          axisLabel: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              width: 1,
              opacity: 0.5,
            },
          },
        },
        yAxis: {
          type: "value",
          inverse: true,
          min: -100,
          max: 1700,
          interval: 50,
          axisTick: {
            show: false,
          },
          axisLabel: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              width: 1,
              opacity: 0.5,
            },
          },
        },
        graphic: {
          elements: [],
        },
        series: [
          {
            type: "graph",
            coordinateSystem: "cartesian2d",
            layout: "none",
            animation: false,
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
            data: [
              {
                name: "GCC",
                value: [150, 1300],
                itemStyle: {
                  color: "purple",
                  opacity: 1,
                },
              },
              {
                name: "HMS",
                value: [500, 1150],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "STR",
                value: [650, 325],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "PWR",
                value: [1200, 75],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "ECLSS",
                value: [1725, 325],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "AGT",
                value: [1950, 850],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "INT",
                value: [1725, 1400],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "EXT",
                value: [1200, 1600],
                itemStyle: {
                  opacity: 1,
                },
              },
              {
                name: "SW0",
                value: [1200, 850],
                symbol: "rect",
                itemStyle: {
                  color: "#0079A3",
                  opacity: 1,
                },
              },
              {
                name: "SW1",
                value: [850, 1000],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
              {
                name: "SW2",
                value: [900, 525],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
              {
                name: "SW3",
                value: [1200, 375],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
              {
                name: "SW4",
                value: [1525, 525],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
              {
                name: "SW5",
                value: [1650, 850],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
              {
                name: "SW6",
                value: [1525, 1200],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
              {
                name: "SW7",
                value: [1200, 1300],
                symbol: "rect",
                itemStyle: {
                  color: "deepskyblue",
                  opacity: 1,
                },
              },
            ],
            links: [
              {
                source: "HMS",
                target: "GCC",
                lineStyle: {
                  type: "dashed",
                  width: 2.5,
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
                lineStyle: {},
              },
              {
                source: "HMS",
                target: "SW2",
                lineStyle: {},
              },
              {
                source: "HMS",
                target: "SW7",
                lineStyle: {},
              },

              {
                source: "STR",
                target: "SW2",
                lineStyle: {},
              },
              {
                source: "STR",
                target: "SW1",
                lineStyle: {},
              },
              {
                source: "STR",
                target: "SW3",
                lineStyle: {},
              },

              {
                source: "PWR",
                target: "SW3",
                lineStyle: {},
              },
              {
                source: "PWR",
                target: "SW2",
                lineStyle: {},
              },
              {
                source: "PWR",
                target: "SW4",
                lineStyle: {},
              },
              {
                source: "ECLSS",
                target: "SW4",
                lineStyle: {},
              },
              {
                source: "ECLSS",
                target: "SW3",
                lineStyle: {},
              },
              {
                source: "ECLSS",
                target: "SW5",
                lineStyle: {},
              },

              {
                source: "AGT",
                target: "SW5",
                lineStyle: {},
              },
              {
                source: "AGT",
                target: "SW4",
                lineStyle: {},
              },
              {
                source: "AGT",
                target: "SW6",
                lineStyle: {},
              },
              {
                source: "INT",
                target: "SW6",
                lineStyle: {},
              },
              {
                source: "INT",
                target: "SW5",
                lineStyle: {},
              },
              {
                source: "INT",
                target: "SW7",
                lineStyle: {},
              },

              {
                source: "EXT",
                target: "SW7",
                lineStyle: {},
              },
              {
                source: "EXT",
                target: "SW1",
                lineStyle: {},
              },
              {
                source: "EXT",
                target: "SW6",
                lineStyle: {},
              },

              {
                source: "SW1",
                target: "SW2",
                lineStyle: {},
              },
              {
                source: "SW2",
                target: "SW3",
                lineStyle: {},
              },
              {
                source: "SW3",
                target: "SW4",
                lineStyle: {},
              },
              {
                source: "SW4",
                target: "SW5",
                lineStyle: {},
              },
              {
                source: "SW5",
                target: "SW6",
                lineStyle: {},
              },
              {
                source: "SW6",
                target: "SW7",
                lineStyle: {},
              },
              {
                source: "SW7",
                target: "SW1",
                lineStyle: {},
              },
              {
                source: "SW1",
                target: "SW0",
                lineStyle: {},
              },
              {
                source: "SW2",
                target: "SW0",
                lineStyle: {},
              },
              {
                source: "SW3",
                target: "SW0",
                lineStyle: {},
              },
              {
                source: "SW4",
                target: "SW0",
                lineStyle: {},
              },
              {
                source: "SW5",
                target: "SW0",
                lineStyle: {},
              },
              {
                source: "SW6",
                target: "SW0",
                lineStyle: {},
              },
              {
                source: "SW7",
                target: "SW0",
                lineStyle: {},
              },
            ],
          },
          {
            z: 1,
            type: "graph",
            coordinateSystem: "cartesian2d",
            layout: "none",
            symbolSize: 45,
            label: {
              show: true,
              fontSize: 12,
              color: "black",
              fontFamily: "Menlo",
            },
            itemStyle: {
              color: "transparent",
            },
            animation: false,
            data: [
              {
                name: "TX:0\nRX:0\n\nGCC",
                value: [150, 1225],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:29\n\nHMS",
                value: [500, 1075],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nSTR",
                value: [650, 250],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:28\n\nPWR",
                value: [1200, 0],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:28\n\nECLSS",
                value: [1725, 250],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:114\nRX:0\n\nAGT",
                value: [1950, 775],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nINT",
                value: [1725, 1325],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nEXT",
                value: [1200, 1525],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:29\nRX:29\n\nSW0",
                value: [1200, 780],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:29\nRX:29\n\nSW1",
                value: [850, 925],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nSW2",
                value: [900, 450],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:1\nRX:29\n\nSW3",
                value: [1200, 300],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:114\nRX:114\n\nSW4",
                value: [1525, 450],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nSW5",
                value: [1650, 775],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nSW6",
                value: [1525, 1125],
                label: {
                  show: true,
                },
              },
              {
                name: "TX:0\nRX:0\n\nSW7",
                value: [1200, 1225],
                label: {
                  show: true,
                },
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
      const params = new URLSearchParams();
      params.append("loss", this.tmpLoss);
      params.append("distance", this.tmpDistance);
      params.append("bandwidth", this.tmpBandwidth);
      this.$api.post(`/api/link/${this.selectedLink}`, params);

      this.activePrompt = false;
    },
    toggleViewOption() {
      this.viewActiveOnly = !this.viewActiveOnly;
      if (!this.viewActiveOnly) {
        this.clearHighlights();
      } else {
        this.highLightActiveNodes();
      }
    },
    clearHighlights() {
      for (var ii = 0; ii < this.option.series[0].data.length; ii++) {
        this.option.series[0].data[ii].itemStyle.opacity = 1;
        this.option.series[1].data[ii].label.show = true;
      }
      for (var j = 0; j < this.option.series[0].links.length; j++) {
        this.option.series[0].links[j].lineStyle.width = 2.2;
      }

      // this.option = JSON.parse(JSON.stringify(this.option))
    },
    highLightActiveNodes() {
      if (this.activeNodes.length == 0) return;
      this.clearHighlights();
      for (var i = 0; i < this.option.series[0].data.length; i++) {
        if (this.activeNodes.indexOf(i) < 0) {
          this.option.series[0].data[i].itemStyle.opacity = 0.1;
          this.option.series[1].data[i].label.show = false;
        }
      }
      for (var j = 0; j < this.option.series[0].links.length; j++) {
        var link = this.option.series[0].links[j];
        if (
          this.activeNodes.indexOf(link.source) < 0 ||
          this.activeNodes.indexOf(link.target) < 0
        ) {
          link.lineStyle.width = 0.1;
        }
      }
    },
    addDraggableGraphicEle() {
      const topoChart = this.$refs.topo;
      this.option.graphic.elements = [];
      for (var i = 0; i < this.option.series[0].data.length; i++) {
        this.option.graphic.elements.push({
          type: "circle",
          position: topoChart.convertToPixel(
            "grid",
            this.option.series[0].data[i].value
          ),
          shape: {
            r: 20,
          },
          z: 200,
          info: i,
          invisible: true,
          draggable: true,
          ondrag: function (item) {
            window.topo.onDrag(item);
          },
          onmousemove: function (item) {
            window.topo.onMove(item);
          },
          onmouseout: function () {
            window.topo.onMoveOut();
          },
        });
      }
    },
    onDrag(item) {
      const topoChart = this.$refs.topo;
      var nodeIdx = parseInt(item.target.info);
      var pos = topoChart.convertFromPixel("grid", [
        item.offsetX,
        item.offsetY,
      ]);
      pos[0] = Math.floor(pos[0]) - (Math.floor(pos[0]) % 25);
      pos[1] = Math.floor(pos[1]) - (Math.floor(pos[1]) % 25);
      this.option.series[0].data[nodeIdx].value = pos;
      this.option.series[1].data[nodeIdx].value = [pos[0], pos[1] - 75];
      this.addDraggableGraphicEle();
      // window.console.log(this.option.series[0].data[nodeIdx].value, [item.offsetX,item.offsetY])
    },
    onMove(item) {
      const topoChart = this.$refs.topo;
      var nodeIdx = parseInt(item.target.info);
      topoChart.dispatchAction({
        type: "showTip",
        seriesIndex: 0,
        dataIndex: nodeIdx,
      });
    },
    onMoveOut() {
      const topoChart = this.$refs.topo;
      topoChart.dispatchAction({
        type: "hideTip",
      });
    },
    editEnable() {
      this.editMode = !this.editMode;
      if (this.editMode) {
        this.addDraggableGraphicEle();
        this.option.tooltip = this.tooltipEdit;
      } else {
        this.option.graphic = { elements: [] };
        // force update
        this.option = JSON.parse(JSON.stringify(this.option));
        this.option.tooltip = this.tooltipDefault;
      }
    },
    addSwitch() {
      var name = "SW" + ++this.switchCnt;
      this.nodes.push(
        {text:name, value: this.nodes.length}
      )
      this.option.series[0].data.push({
        name: name,
        value: [150, 100],
        symbol: "rect",
        itemStyle: {
          color: "deepskyblue",
          opacity: 1,
        },
      });
      this.option.series[1].data.push({
        name: "TX:0\nRX:0\n\n" + name,
        value: [150, 25],
        label: {
          show: true,
        },
      });
      this.addDraggableGraphicEle();
    },
    connect() {
      this.option.series[0].links.push({
        source:this.connectHost0,
        target:this.connectHost1
      })
    }
  },
  mounted() {
    window.topo = this;
    this.initLinkStatus();
    this.option.tooltip = this.tooltipDefault;
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
      var tmpActiveNodes = [];
      for (var i in stats) {
        var newStatsString =
          "TX:" +
          stats[i][0] +
          "\nRX:" +
          stats[i][1] +
          "\n\n" +
          nameIdxMap[i].name;
        if (
          this.option.series[1].data[nameIdxMap[i].idx].name != newStatsString
        ) {
          tmpActiveNodes.push(nameIdxMap[i].idx, nameIdxMap[i].name);
        }
        if (tmpActiveNodes.length > 0) {
          this.activeNodes = tmpActiveNodes;
        }
        this.option.series[1].data[nameIdxMap[i].idx].name = newStatsString;
      }
    });
  },
  watch: {
    activeNodes: function () {
      if (this.viewActiveOnly) this.highLightActiveNodes();
    },
  },
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
.buttons {
  /* width:105px; */
  font-size: 0.75rem;
  font-weight: 600;
}
.conenct-select {
  width: 100px;
}
</style>