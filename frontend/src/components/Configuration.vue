<template>
<vs-card>
  <div slot="header" style="text-align:left;">
      <h3>Configurations</h3>
    </div>
  <vs-tabs :value="1" >
    <vs-tab index="0" label="IP Addr">
      <vs-table  :data="configAddrs" stripe class="addr">
        <template slot="thead">
          <vs-th> Subsys </vs-th>
          <vs-th> Local </vs-th>
          <vs-th> Remote </vs-th>
        </template>

        <template slot-scope="{ data }" >
          <vs-tr :key="indextr" v-for="(tr, indextr) in data">
            <vs-td :data="data[indextr].name" >
              {{ tr.name }}
            </vs-td>
            <vs-td :data="data[indextr].local">
              {{ tr.local }}
            </vs-td>
            <vs-td :data="data[indextr].remote">
              {{ tr.remote }}
            </vs-td>
          </vs-tr>
        </template>
      </vs-table>
    </vs-tab>
    <vs-tab index="1" label="Links">
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

    
    <vs-tab index="2" label="Routing">
    </vs-tab>
    <vs-tab index="3" label="Scheduling">
    </vs-tab>
  </vs-tabs>
</vs-card>  
</template>

<script>
import { debounce } from "./debounce";

export default {
  data() {
    return {
      configAddrs: [
        {name:"GCC", local: "0.0.0.0:10000", remote: "127.0.0.1:20000"},
        {name:"HMS", local: "0.0.0.0:10001", remote: "127.0.0.1:20001"},
        {name:"STR", local: "0.0.0.0:10002", remote: "127.0.0.1:20002"},
        {name:"PWR", local: "0.0.0.0:10003", remote: "127.0.0.1:20003"},
        {name:"ECLSS", local: "0.0.0.0:10004", remote: "127.0.0.1:20004"},
        {name:"AGT", local: "0.0.0.0:10005", remote: "127.0.0.1:20005"},
        {name:"INT", local: "0.0.0.0:10006", remote: "127.0.0.1:20006"},
        {name:"EXT", local: "0.0.0.0:10007", remote: "127.0.0.1:20007"},
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
          distance: 80, // percentage
        },
      ],
    };
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
      params.append('distance', this.configLinks[1].distance)
      params.append('bandwidth', this.configLinks[1].bandwidth)
      this.$api.post('/api/links', params);
    }, 200),
  },
  methods: {
    distanceFormatter() {
      return 10;
    },
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

</style>