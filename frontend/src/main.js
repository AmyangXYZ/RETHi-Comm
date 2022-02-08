import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import router from './router'
import Vuesax from 'vuesax'
import 'vuesax/dist/vuesax.css'
import 'material-icons/iconfont/material-icons.css';
import '@fortawesome/fontawesome-free/css/all.css'
import '@fortawesome/fontawesome-free/js/all.js'

Vue.config.productionTip = false
Vue.prototype.$EventBus = new Vue()
Vue.use(Vuesax)

var axiosInstance = axios.create()
axiosInstance.defaults.baseURL = '/'
Vue.prototype.$api = axiosInstance;

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
