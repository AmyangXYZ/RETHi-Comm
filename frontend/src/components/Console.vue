<template>
  <vs-card id="console">
    <div slot="header"  style="text-align:left">
      <h3>Logs</h3>
    </div>
    <vs-tabs style="font-size:0.85rem">
      <vs-tab label="Packets">
        <textarea :style="cssVars" autofocus id="logs" ref="logs" :value="log" @change="v => log = v" disabled />
      </vs-tab>
      <vs-tab label="Faults">
        <textarea :style="cssVars" autofocus id="logs" ref="logs" :value="logFaults" @change="v => logFaults = v" disabled />
      </vs-tab>
    </vs-tabs>
  </vs-card>
</template>

<script>
const WSLOG_HEARTBEAT = -1
const WSLOG_MSG = 0
const WSLOG_STAT = 1
const WSLOG_PKT_TX = 2

export default {
  data() {
    return {
      dropped:false,
      ws: {},
      cnt: 0,
      log: "",
      logFaults: "",
      newLogs: "",
    };
  },
  props:[
    "name",
    "height"
  ],
  computed: {
    cssVars () {
      return {
        '--height': this.height
      }
    }
  },
  methods: {
    startWS() {
      this.ws = new WebSocket("ws://localhost:8000/ws/"+this.name)
      this.ws.onopen = () => {
        this.dropped = false
        this.newLogs = "[+] Console connected"
      }
      this.ws.onclose = () => {
        if(!this.dropped)
          this.newLogs+="\n[!] Connection dropped, trying to reconnect..."
        else
          this.newLogs+="."
        this.dropped = true
        this.$nextTick(() => {
          this.$refs.logs.scrollTop = this.$refs.logs.scrollHeight
        })
        
        setTimeout(this.startWS, 1000)
      }
      this.ws.onmessage = (evt) => {    
        // window.console.log(evt.data)
        if(this.cnt>2000 && this.cnt%2000==1) this.log = ""
        this.cnt++
        
        var entry = JSON.parse(evt.data)
        
        if(entry.type == WSLOG_HEARTBEAT) {
          // heartbeat 
        }
        if(entry.type == WSLOG_MSG) {
          // this.log+="\n[*] "+entry.msg+""
          var date = new Date()
          var options = { hour12: false };
        
          this.newLogs+="\n["+date.toLocaleString('en-US', options).replace(",","")+"] "+entry.msg

        }
        if(entry.type == WSLOG_STAT) {
          this.$EventBus.$emit("stats_"+this.name, entry["stats_"+this.name])
        }
        if(entry.type == WSLOG_PKT_TX) {
          this.$EventBus.$emit("pkt_tx", entry["pkt_tx"])
        }
      }
    },
    writeLogs() {
      this.log+=this.newLogs
      this.$nextTick(() => {
        this.$refs.logs.scrollTop = this.$refs.logs.scrollHeight
      })
      this.newLogs = ""
    }
  },
  mounted() {
    this.startWS()
    setInterval(() => {
      this.writeLogs()
    }, 160);
    this.$EventBus.$on("fault-log",(fault_log)=>{
      var date = new Date()
      var options = { hour12: false };
      this.logFaults += "["+date.toLocaleString('en-US', options).replace(",","")+"] "+fault_log+"\n"
    })
  },
};
</script>

<style scoped>

#logs {
  margin-top: 1px;
  width: 100%;
  height: var(--height);
  font-size: 0.8rem;
  line-height: 1.3;
  border-radius: 6px;
  padding: 3px;
  box-sizing: border-box;
  resize: none;
  outline: none;
  text-transform: none;
  text-decoration: none;
}

textarea {
  border: none
}
textarea:disabled {
  background: white;
}
</style>