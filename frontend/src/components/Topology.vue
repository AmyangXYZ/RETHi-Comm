<template>
  <vs-card id="topology">
    <div slot="header" style="text-align: left">
      <vs-row vs-type="flex" vs-justify="space-between">
        <vs-col vs-w="3">
          <h3>Topology</h3>
        </vs-col>
        <vs-col vs-w="9" vs-type="flex" vs-justify="flex-end">
          <div v-show="!editMode">
            
            <vs-row vs-w="12">
              <vs-col vs-offset="-0.5" vs-w="3.6">
                <vs-button
                  class="buttons"
                  style="width: 100px;height:30px;font-size:0.75rem"
                  size="small"
                  key="view"
                  :color="ANIMATION_ENABLED ? 'success' : 'grey'"
                
                  @click="pktTxAnimationToggle"
                >
                 Animation {{ ANIMATION_ENABLED ? "On" : "Off" }}
                </vs-button>
              </vs-col>
              <vs-col vs-offset="0.3" vs-w="2.5">
                <vs-button
                  class="buttons"
                  style="width:80px;height:30px;font-size:0.75rem"
                  size="small"
                  key="view"
                  :color="viewActiveOnly ? 'success' : 'grey'"

                  @click="toggleViewOption"
                >
                  {{ viewActiveOnly ? "Active only" : "All paths" }}
                </vs-button>
              </vs-col>

              <vs-col vs-offset="0.7" vs-w="3.1">
                <vs-select
                  class="connect-select"
                  v-model="selectedTopo"
                  width="90px"
                >
                  <vs-select-item
                    :key="index"
                    :value="index"
                    :text="item"
                    v-for="(item, index) in topoTags"
                  />
                </vs-select>
              </vs-col>
              <vs-col vs-offset="0.2" vs-w="1.5 ">
                <vs-button
                  style="margin-left: 5px"
                  class="buttons"
                  size="small"
                  color="primary"
                  icon-pack="fas"
                  type="filled"
                  icon="fa-edit"
                  icon-after
                  @click="editEnable"
                >
                  Edit
                </vs-button>
              </vs-col>
            </vs-row>
          </div>

          <div v-show="editMode">
            <vs-row vs-type="flex" vs-justify="flex-start">
              <vs-col vs-offset="-1" vs-w="3">
                <vs-input
                  style="width: 80px"
                  placeholder="Tag"
                  v-model="newTopoTag"
                />
              </vs-col>
              <vs-col vs-offset="1.8" vs-w="3">
                <vs-button
                  class="buttons"
                  size="small"
                  color="success"
                  icon-pack="fas"
                  type="filled"
                  icon="fa-save"
                  @click="editApply"
                >
                  Save
                </vs-button>
              </vs-col>
              <vs-col vs-offset="1" vs-w="3">
                <vs-button
                  class="buttons"
                  size="small"
                  color="danger"
                  icon-pack="fas"
                  type="filled"
                  icon="fa-undo"
                  @click="editReset"
                >
                  Reset
                </vs-button>
              </vs-col>
            </vs-row>
          </div>
        </vs-col>
      </vs-row>

      <div v-show="editMode">
        <vs-row vs-type="flex" vs-justify="center" style="margin-top: 8px">
          <vs-col vs-offset="-1" vs-w="1">
            <vs-input
              style="width: 80px"
              placeholder="SW#"
              v-model="newSwitchName"
            />
          </vs-col>
          <vs-col vs-offset="0.5" vs-w="0.5">
            <vs-button
              class="buttons"
              size="small"
              color="success"
              icon-pack="fas"
              type="filled"
              icon="fa-plus"
              @click="editAddSwitch"
            >
            </vs-button>
          </vs-col>
          <vs-col vs-offset="0.1" vs-w="1">
            <vs-button
              class="buttons"
              size="small"
              color="danger"
              icon-pack="fas"
              type="filled"
              icon="fa-minus"
              @click="editRemoveSwitch"
            >
            </vs-button>
          </vs-col>
          <vs-col vs-offset="1" vs-w="1">
            <vs-button
              class="buttons"
              size="small"
              color="primary"
              icon-pack="fas"
              type="filled"
              icon="fa-arrows-alt-h"
              @click="editConnect"
            >
            </vs-button>
          </vs-col>
          <vs-col vs-offset="-0.5" vs-w="1">
            <vs-select class="connect-select" v-model="connectHost0">
              <vs-select-item
                :key="index"
                :value="index"
                :text="item"
                v-for="(item, index) in nodes"
              />
            </vs-select>
          </vs-col>
          <vs-col vs-offset="0.4" vs-w="1">
            <vs-select class="connect-select" v-model="connectHost1">
              <vs-select-item
                :key="index"
                :value="index"
                :text="item"
                v-for="(item, index) in nodes"
              />
            </vs-select>
          </vs-col>
          <vs-col vs-offset="0.35" vs-w="1">
            <vs-button
              class="buttons"
              size="small"
              color="danger"
              icon-pack="fas"
              type="filled"
              icon="fa-cut"
              @click="editCut"
            >
            </vs-button>
          </vs-col>
        </vs-row>
      </div>
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
            <vs-switch v-model="tmpFailed" />
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
      nodes_position: [],
      ANIMATION_ENABLED:false,
      editMode: false,
      viewActiveOnly: false,
      activeNodes: [],
      showOption: false,
      selectedTopo: 0,
      topoTags: ["default"],
      newTopoTag: "",
      activePrompt: false,
      selectedLink: "",
      tmpFailed: false,
      tmpLoss: 0,
      tmpDelay: 1,
      tmpBandwidth: 1,
      tmpDistance: 1000,
      linkStats: {},
      newSwitchName: "",
      nodes: [], //
      packets:[],
      connectHost0: 0,
      connectHost1: 0,
      nodesIndexMap: {},
      flows: [],
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
        triggerOn: 'none',
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
          min:100,
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
            data: [],
            links: [],
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
            data: [],
            silent:true,
          },
          {
            z: 1,
            type: "graph",
            coordinateSystem: "cartesian2d",
            layout: "none",
            symbolSize: 8,
            data:[],
            animationDurationUpdate:500,
            // animationDuration:800
          }
        ],
      },
      option_backup: {},
    };
  },
  methods: {
    getTopologyTags() {
      this.$api.get("/api/topologies").then((res) => {
        if (res.data.flag == 0) {
          return;
        }
        this.topoTags = res.data.data;
        this.topoTags.sort(function (x, y) {
          return x == "default" ? -1 : y == "default" ? 1 : 0;
        });
        this.getTopology(this.topoTags[0]);
      });
    },
    postTopology(tag) {
      var nodes = [];
      var edges = [];

      for (var nn in this.option.series[0].data) {
        nodes.push({
          name: this.option.series[0].data[nn].name,
          value: this.option.series[0].data[nn].value,
        });
      }
      for (var ee in this.option.series[0].links) {
        edges.push([
          this.option.series[0].links[ee].source,
          this.option.series[0].links[ee].target,
        ]);
      }
      this.$api.post("/api/topology", {
        tag: tag,
        nodes: nodes,
        edges: edges,
      });
    },
    getTopology(tag) {
      this.$api.get("/api/topology?tag=" + tag).then((res) => {
        if (res.data.flag == 0) {
          if (tag == "default") setTimeout(this.getTopology, 200, tag);
          else return;
        }

        this.nodes = [];
        this.switchCnt = 0;
        for (var i = 0; i < res.data.data.nodes.length; i++) {
          var node = res.data.data.nodes[i];
          this.nodesIndexMap[node.name] = i
          node.itemStyle = {
            opacity: 1,
          };
          if (node.name.indexOf("GCC") != -1) {
            node.itemStyle.color = "purple";
          }
          if (node.name.indexOf("SW") != -1) {
            this.switchCnt++;
            node.symbol = "rect";
            node.itemStyle.color = "deepskyblue";
            // if (node.name == "SW0") node.itemStyle.color = "#0079A3";
          }
          this.nodes.push(node.name);
        }

        var edges = [];
        for (var j = 0; j < res.data.data.edges.length; j++) {
          var edge = res.data.data.edges[j];
          var link = {
            source: edge[0],
            target: edge[1],
            lineStyle: {},
          };
          if (edge[0] == "GCC" || edge[1] == "GCC") {
            link.lineStyle = {
              type: "dashed",
              width: 2.5,
            };
            link.emphasis = {
              lineStyle: {
                type: "dashed",
              },
            };
          }
          edges.push(link);
        }

        this.option.series[2].data = [] // clear pkt icons
        this.option.series[0].data = res.data.data.nodes;
        this.option.series[0].links = edges;

        this.$EventBus.$emit("topology", { nodes: this.nodes, edges: edges });

        this.initNodesStatistics();
        this.initLinkStatus();
      });
    },
    initNodesStatistics() {
      this.option.series[1].data = [];
      for (var i = 0; i < this.option.series[0].data.length; i++) {
        this.option.series[1].data.push({
          name: "TX:0\nRX:0\n\n" + this.option.series[0].data[i].name,
          value: [
            this.option.series[0].data[i].value[0],
            this.option.series[0].data[i].value[1] - 75,
          ],
          label: {
            show: true,
          },
        });
      }
    },
    monitorNodesStatistics() {
      this.$EventBus.$on("stats_comm", (stats) => {
        this.activeNodes = [];
        var tmpActiveNodes = [];
        for (var name in stats) {
          var newStatsString =
            "TX:" + stats[name][0] + "\nRX:" + stats[name][1] + "\n\n" + name;
          var idx = 0;
          for (var j = 0; j < this.nodes.length; j++) {
            if (this.nodes[j] == name) {
              idx = j;
              break;
            }
          }
          if (this.option.series[1].data[idx].name != newStatsString) {
            tmpActiveNodes.push(idx, name);
          }
          if (tmpActiveNodes.length > 0) {
            this.activeNodes = tmpActiveNodes;
          }
          this.option.series[1].data[idx].name = newStatsString;
        }
      });
    },
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
      if (item.dataType == "node") {
        this.$EventBus.$emit("selectedNode", item.data.name);
      } else if (item.dataType == "edge") {
        this.selectedLink = item.name;
        this.activePrompt = true;
      }
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
      this.linkStats[this.selectedLink].Failed = this.tmpFailed;
      const params = new URLSearchParams();
      params.append("loss", this.tmpLoss);
      params.append("distance", this.tmpDistance);
      params.append("bandwidth", this.tmpBandwidth);
      params.append("failed", this.tmpFailed);
      this.$api.post(`/api/link/${this.selectedLink}`, params);

      this.activePrompt = false;
    },
    toggleViewOption() {
      this.viewActiveOnly = !this.viewActiveOnly;
      if (!this.viewActiveOnly) {
        this.clearHighlights();
      } else {
        this.highlightActiveNodes();
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
    },
    highlightActiveNodes() {
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
      this.editMode = true;
      this.option_backup = JSON.stringify(this.option);
      this.addDraggableGraphicEle();
      this.option.tooltip = this.tooltipEdit;
    },
    editApply() {
      this.editMode = false;
      this.option.graphic = { elements: [] };
      this.postTopology(this.newTopoTag);
      this.topoTags.push(this.newTopoTag);
      this.selectedTopo = this.topoTags.length - 1;
      this.newTopoTag = "";

      // force update
      this.option = JSON.parse(JSON.stringify(this.option));
      this.option.tooltip = this.tooltipDefault;
    },
    editReset() {
      this.editMode = false;
      this.option.graphic = { elements: [] };
      // force update
      this.option = JSON.parse(this.option_backup);
      this.option.tooltip = this.tooltipDefault;
    },
    editAddSwitch() {
      this.nodes.push(this.newSwitchName);
      this.option.series[0].data.push({
        name: this.newSwitchName,
        value: [150, 100],
        symbol: "rect",
        itemStyle: {
          color: "deepskyblue",
          opacity: 1,
        },
      });

      this.option.series[1].data.push({
        name: "TX:0\nRX:0\n\n" + this.newSwitchName,
        value: [150, 25],
        label: {
          show: true,
        },
      });
      this.addDraggableGraphicEle();
    },
    editRemoveSwitch() {
      for (var i = 0; i < this.option.series[0].data.length; i++) {
        var n = this.option.series[0].data[i];
        if (n.name==this.newSwitchName) {
          this.option.series[0].data.splice(i, 1);
          this.option.series[1].data.splice(i, 1);
        }
      }
      for (var j=this.option.series[0].links.length-1;j>=0;j--) {
        var l = this.option.series[0].links[j];
        if (l.source==this.newSwitchName || l.target==this.newSwitchName) {
          this.option.series[0].links.splice(j, 1);
        }
      }
    },
    editConnect() {
      if (this.connectHost0 != this.connectHost1) {
        this.option.series[0].links.push({
          source: this.nodes[this.connectHost0],
          target: this.nodes[this.connectHost1],
        });
      }
    },
    editCut() {
      if (this.connectHost0 != this.connectHost1) {
        for (var i = 0; i < this.option.series[0].links.length; i++) {
          var l = this.option.series[0].links[i];
          if (
            (l.source == this.nodes[this.connectHost0] &&
              l.target == this.nodes[this.connectHost1]) ||
            (l.source == this.nodes[this.connectHost1] &&
              l.target == this.nodes[this.connectHost0])
          ) {
            this.option.series[0].links.splice(i, 1);
            break
          }
        }
      }
    },
    pktTxAnimationToggle() {
      this.packets = {}
      this.option.series[2].data = []
      this.ANIMATION_ENABLED = !this.ANIMATION_ENABLED 
      this.$api.get(`/api/animation/${this.ANIMATION_ENABLED}`)
    },
    pktTxAnimation(flow) {
      if (flow.finished == true) { 
        for (var i=0;i<this.option.series[2].data.length;i++) {
          if (this.option.series[2].data[i].name == flow.uid) {
            this.option.series[2].data.splice(i,1)
            return
          }
        }
      }
      if (this.packets[flow.uid]==null) {
        this.packets[flow.uid] = {name: flow.uid, value:JSON.parse(JSON.stringify(this.option.series[0].data[this.nodesIndexMap[flow.node]].value))}
        this.option.series[2].data.push(this.packets[flow.uid])
      } else {
        this.packets[flow.uid].value = JSON.parse(JSON.stringify(this.option.series[0].data[this.nodesIndexMap[flow.node]].value))
      }
    },
    pktTxAnimationUpdate() {
      for (var i=0;i<this.flows.length;i++) {
        var flow = this.flows[i]
        this.pktTxAnimation(flow)
      }
      // window.console.log(this.option.series[2].data.length)
      this.flows = []
    }
  },
  mounted() {
    window.topo = this;

    this.$api.get(`/api/animation`).
    then((res)=>{
      this.ANIMATION_ENABLED = Boolean(res.data)
    })

    this.option.tooltip = this.tooltipDefault;
    this.getTopologyTags();
    
    setTimeout(this.monitorNodesStatistics, 200);

    setInterval(this.pktTxAnimationUpdate, 200)

    this.$EventBus.$on("pkt_tx", (flow)=>{
      if (this.ANIMATION_ENABLED) {
        this.flows.push(flow)
        // this.pktTxAnimation(flow)
      }
    })
  },
  // created() {

  // },
  watch: {
    activeNodes: function () {
      if (this.viewActiveOnly) this.highlightActiveNodes();
    },
    selectedTopo: function () {
      this.getTopology(this.topoTags[this.selectedTopo]);
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
  font-size: 8rem;
  font-weight: 600;
}
.connect-select {
  width: 80px;
}
</style>