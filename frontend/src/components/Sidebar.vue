<template lang="html">

  <div class="parentx">

    <vs-sidebar :reduce="reduce" :reduce-not-rebound="reduceNotRebound" :reduce-not-hover-expand="notExpand" 
      :hidden-background="hiddenBG" parent="body" 
      :default-index="cur_index"  color="danger" class="sidebarx" spacer v-model="active" >

      <div class="header-sidebar" slot="header">
        <vs-avatar src="https://www.purdue.edu/rethi/img/icons/Situational_Awareness_Icon.png"/>
      </div>
      <vs-sidebar-group open title="HMS" align="left">
        <vs-sidebar-item to="/" index="1" size="small" icon="dashboard" >
          Dashboard
        </vs-sidebar-item>
        <vs-sidebar-item to="/" index="2" icon="schedule" >
          Scheduler
        </vs-sidebar-item>
      </vs-sidebar-group>

    </vs-sidebar>
  </div>

</template>

<script>
export default {
  data() {
    return {
      cur_index: 1,
      active: true,
      paths: ["/dashboard", "/scheduler", "/analysis"],
      notExpand: false,
      hiddenBG: true,
      reduce: true,
      reduceNotRebound: true,
    };
  },
  methods: {
    biu() {
      this.cur_index = this.paths.indexOf(this.$route.path) + 1;
    },
  },
  created() {
    // this.biu();
  },
};
</script>

<style scoped>
.header-sidebar {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
}
</style>
