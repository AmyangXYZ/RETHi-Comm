<template>
  <!-- <vs-card id="config"> -->

  <vs-table :data="settings" style="text-align: left">
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
                <vs-slider text-fixed="m" v-model="tr.distance" />
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
  <!-- </vs-card> -->
</template>

<script>
import { debounce } from "./debounce";
import axios from "axios";

export default {
  data() {
    return {
      settings: [
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
        54500000 + (this.settings[1].distance / 100) * (401300000 - 54500000)
      );
    },
    wiredDelay: function () {
      return (
        (this.settings[0].distance / (300000000 * 0.77) +
         800 / (this.settings[0].bandwidth * 1024 * 1024 * 1024)) *
        1000000000
      );
    },
    wirelessDelay: function () {
      return (
        this.wirelessDistance / 300000 +
        800 / (this.settings[1].bandwidth * 1024 * 1024)
      );
    },
  },
  watch: {
    wiredDelay: debounce(function () {
      axios.get(
        `http://localhost:8000/api/links?type=wired&distance=${this.settings[0].distance}&bandwidth=${this.settings[0].bandwidth}`
      );
    }, 200),
    wirelessDelay: debounce(function () {
      axios.get(
        `http://localhost:8000/api/links?type=wireless&distance=${this.wirelessDistance}&bandwidth=${this.settings[0].bandwidth}`
      );
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
#config {
  text-align: left;
  /* height: 238px; */
}
th {
  font-size: 1rem;
}
td {
  font-size: 1rem;
}
</style>