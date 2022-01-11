<template>
  <vs-card id="console">
    <div slot="header"  style="text-align:left">
      <h3>Console</h3>
    </div>
    <textarea :style="cssVars" autofocus id="logs" ref="logs" :value="log" @change="v => log = v" disabled />
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
      log: ""
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
      var loc = window.location;
      "ws://"+loc.host+"/ws"
      this.ws = new WebSocket("ws://127.0.0.1:8000/ws/"+this.name)
      // this.ws = new WebSocket("ws://"+loc.host+"/ws")
      this.ws.onopen = () => {
        this.dropped = false
        this.log = "[+] Console connected"
        this.$nextTick(() => {
          this.$refs.logs.scrollTop = this.$refs.logs.scrollHeight
        })

      }
      this.ws.onclose = () => {
        if(!this.dropped)
          this.log+="\n[!] Connection dropped, trying to reconnect..."
        else
          this.log+="."
        this.dropped = true
        this.$nextTick(() => {
          this.$refs.logs.scrollTop = this.$refs.logs.scrollHeight
        })
        
        setTimeout(this.startWS, 1000)
      }
      this.ws.onmessage = (evt) => {    
        // window.console.log(evt.data)
        if(this.cnt>1000 && this.cnt%1000==1) this.log = ""
        this.cnt++
        
        var entry = JSON.parse(evt.data)
        
        if(entry.type == WSLOG_HEARTBEAT) {
          // heartbeat 
        }
        if(entry.type == WSLOG_MSG) {
          // this.log+="\n[*] "+entry.msg+""
          var date = new Date()
          var options = { hour12: false };
        
          this.log+="\n["+date.toLocaleString('en-US', options).replace(",","")+"] "+entry.msg
          this.$nextTick(() => {
            this.$refs.logs.scrollTop = this.$refs.logs.scrollHeight
          })

        }
        if(entry.type == WSLOG_STAT) {
          this.$EventBus.$emit("stats_"+this.name, entry["stats_"+this.name])
        }
        if(entry.type == WSLOG_PKT_TX) {
          this.$EventBus.$emit("pkt_tx", entry["pkt_tx"])
          window.console.log(entry["pkt_tx"])
        }
      }
    }
  },
  mounted() {
    this.startWS()
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