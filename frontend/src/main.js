import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import Vuesax from 'vuesax'
import 'vuesax/dist/vuesax.css'
import 'material-icons/iconfont/material-icons.css';
import '@fortawesome/fontawesome-free/css/all.css'
import '@fortawesome/fontawesome-free/js/all.js'

Vue.config.productionTip = false
Vue.prototype.$EventBus = new Vue()
Vue.use(Vuesax)

var axiosInstance = axios.create()
axiosInstance.defaults.baseURL = 'http://localhost:8000/'
Vue.prototype.$api = axiosInstance;

new Vue({
  render: h => h(App),
}).$mount('#app')
