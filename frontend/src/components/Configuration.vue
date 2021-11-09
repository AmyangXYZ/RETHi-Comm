<template>
  <vs-card id="config">
    <div slot="header">
      <h3>Settings</h3>
    </div>
    <h3>In-habitat (Wire):</h3>
    Bandwidth: 1Gbps<br />
    Speed: .77 c<br />
    Distance: 30 meter<br />
    One-way delay:  &lt; 1 us
    <h3>Earth-Mars (Wireless):</h3>
    Bandwidth: 18000bps (500 ~ 32000)<br />
    Speed: c<br />
    Distance: 54,500,000 ~ 401,300,000 km<br />
    One-way delay: 3~22 min
  </vs-card>
</template>

<script>
export default {};
</script>

<style scoped>
#config {
  text-align: left;
  /* height: 238px; */
}
</style>