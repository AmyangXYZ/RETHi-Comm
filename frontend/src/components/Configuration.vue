<template>
<vs-card>
  <div slot="header" style="text-align:left;">
      <h3>Configurations</h3>
    </div>
  <vs-tabs :value="0" style="font-size:0.85rem">
    <vs-tab index="0" label="Links">
      <vs-table :data="configLinks" class="links">
        <template slot="thead">
          <vs-th> Link Type</vs-th>
          <vs-th> Bandwidth </vs-th>
          <vs-th> Speed </vs-th>
          <vs-th> Distance </vs-th>
          <vs-th> Expected Delay (100 bytes)</vs-th>
        </template>

        <template slot-scope="{ data }">
          <vs-tr :key="indextr" v-for="(tr, indextr) in data">
            <vs-td :data="data[indextr].link">
              {{ tr.link }}
            </vs-td>

            <vs-td :data="tr.bandwidth" v-if="tr.link[0] == 'G'">
              {{ tr.bandwidth }} Kbps
              <template slot="edit">
                <vs-input v-model="tr.bandwidth" class="inputx" />
              </template>
            </vs-td>

            <vs-td :data="tr.bandwidth" v-else>
              {{ tr.bandwidth }} Gbps
              <template slot="edit">
                <vs-input v-model="tr.bandwidth" class="inputx" />
              </template>
            </vs-td>

            <vs-td :data="tr.speed">
              {{ tr.speed }}
            </vs-td>

            <vs-td :data="tr.distance" v-if="tr.link[0] == 'G'">
              {{ wirelessDistance }} km
              <template slot="edit">
                <vs-row>
                  <vs-col vs-w="12" style="text-align: center; font-size: 0.95rem">
                    Range: 54500000 ~ 401300000 km
                  </vs-col>
                  <vs-col>
                    <vs-slider text-fixed="%" v-model="tr.distance" />
                  </vs-col>
                </vs-row>
              </template>
            </vs-td>

            <vs-td :data="tr.distance" v-else>
              {{ tr.distance }} m
              <template slot="edit">
                <vs-row>
                  <vs-col vs-w="12" style="text-align: center; font-size: 0.95rem">
                  </vs-col>
                  <vs-col>
                    <vs-slider :min=1 text-fixed="m" v-model="tr.distance" />
                  </vs-col>
                </vs-row>
              </template>
            </vs-td>

            <vs-td :data="tr.delay" v-if="tr.link[0] == 'G'">
              {{ wirelessDelay.toFixed(2) }} s
            </vs-td>

            <vs-td :data="tr.delay" v-else> {{ wiredDelay.toFixed(2) }} ns </vs-td>
          </vs-tr>
        </template>
      </vs-table>
    </vs-tab>
    <vs-tab index="1" label="Routing">
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          Shortest Path Routing        
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="routing"></vs-checkbox>
        </vs-col>
      </vs-row>

    </vs-tab>
    <vs-tab index="2" label="Scheduling">
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          MIMOMQ Prioirty Scheduling
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="mimomq"></vs-checkbox>
        </vs-col>
      </vs-row>
    </vs-tab>
    <vs-tab index="3" label="Reliability">
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          Auto Re-Reroute Upon Switch Failure
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="REROUTE_ENABLED"></vs-checkbox>
        </vs-col>
      </vs-row>
      <vs-row vs-type="flex" vs-justify="center">
        <vs-col vs-w="8">
          IEEE 802.1CB Frame Replication and Elimination for Reliability 
        </vs-col>
        <vs-col vs-w="3" style="font-size:0.5rem">
          <vs-checkbox size="small" v-model="FRER_ENABLED"></vs-checkbox>
        </vs-col>
      </vs-row>
    </vs-tab>

    <vs-tab index="4" label="Priorities">
      <vs-row  vs-type="flex" vs-justify="start" vs-align="center">
        <vs-col vs-w="3" :key="i" v-for="(s, i) in subsys">
          <vs-row vs-align="center">
            <vs-col vs-w="6">
              {{s.name}}
            </vs-col>
            <vs-col vs-w="4" style="font-size:0.5rem">
              <vs-input-number size="small" min="0" max="7" v-model="s.priority"/>
            </vs-col>
          </vs-row>
        </vs-col>
      </vs-row>
      <vs-row vs-type="flex" vs-justify="center" vs-align="center">
        <vs-col vs-w="1.1">
          <vs-button class="priority-bts" size="medium" color="danger" type="relief"  @click="submit">Submit</vs-button>
        </vs-col>
        <vs-col vs-offset="0.3" vs-w="1.1">
          <vs-button class="priority-bts" size="medium" color="primary" type="relief"  @click="reset">Reset</vs-button>
        </vs-col>
      </vs-row>
    </vs-tab>
  </vs-tabs>
</vs-card>  
</template>

<script>
import { debounce } from "./debounce";

export default {
  data() {
    return {
      mimomq:true,
      routing:true,
      REROUTE_ENABLED:false,
      FRER_ENABLED: false,
      subsys: [
        { name:"GCC", priority:1},
        { name:"HMS", priority:1},
        { name:"STR", priority:1},
        { name:"SPL", priority:1},
        { name:"PWR", priority:1},
        { name:"ECLSS", priority:1},
        { name:"AGT", priority:1},
        { name:"IE", priority:1},
        { name:"EXT", priority:1},
        { name:"DTB", priority:1},
        { name:"COORD", priority:1},
      ],
      configLinks: [
        {
          link: "In-habitat",
          bandwidth: 1, // Gbps
          speed: ".77c",
          distance: 30,
        },
        {
          link: "Ground-habitat",
          bandwidth: 2, // Kbps
          speed: "c",
          distance: 1, // percentage
        },
      ],
    };
  },
  methods:{
    distanceFormatter() {
      return 10;
    },
    submit() {
      this.$api.post("/api/priorities", this.subsys)
    },
    reset() {
      for (var i=0;i<this.subsys.length;i++) {
        this.subsys[i].priority = 1
      }
    }
  },
  mounted () {
    this.$api.get(`/api/frer`).
    then((res)=>{
      this.FRER_ENABLED = Boolean(res.data)
    })
    this.$api.get(`/api/reroute`).
    then((res)=>{
      this.REROUTE_ENABLED = Boolean(res.data)
    })
  },
  computed: {
    wirelessDistance: function () {
      return Math.round(
        54500000 + (this.configLinks[1].distance / 100) * (401300000 - 54500000)
      );
    },
    wiredDelay: function () {
      return (
        (this.configLinks[0].distance / (300000000 * 0.77) +
         800 / (this.configLinks[0].bandwidth * 1024 * 1024 * 1024)) *
        1000000000
      );
    },
    wirelessDelay: function () {
      return (
        this.wirelessDistance / 300000 +
        800 / (this.configLinks[1].bandwidth * 1024 * 1024)
      );
    },
  },
  watch: {
    wiredDelay: debounce(function () {
      const params = new URLSearchParams()
      params.append('type', 'wired')
      params.append('distance', this.configLinks[0].distance)
      params.append('bandwidth', this.configLinks[0].bandwidth)
      this.$api.post('/api/links', params);
    }, 200),
    wirelessDelay: debounce(function () {
      const params = new URLSearchParams()
      params.append('type', 'wireless')
      params.append('distance', this.wirelessDistance*1000)
      params.append('bandwidth', this.configLinks[1].bandwidth)
      this.$api.post('/api/links', params);
    }, 200),
    FRER_ENABLED: debounce(function () {
      this.$api.get(`/api/frer/${this.FRER_ENABLED}`);
    }, 200),
    REROUTE_ENABLED: debounce(function () {
      this.$api.get(`/api/reroute/${this.REROUTE_ENABLED}`);
    }, 200),
  },
};
</script>

<style scoped>
.links {
  text-align: left;
}
.links th, td {
  font-size: 0.85rem;
}

.addr {
  text-align: left;
}
.addr td {
  font-size: 0.85rem;
}
.addr .vs-table--tbody-table .tr-values td {
  padding-top: 2px;
  padding-bottom: 2px;
}

.priority-bts {
  width: 100%;
}
</style>