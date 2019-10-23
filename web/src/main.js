import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueRouter from 'vue-router'
import Buefy from 'buefy'
import 'buefy/dist/buefy.css'

import Login from "./components/Login.vue"
import Chat from "./components/Chat.vue"

Vue.use(Buefy)
Vue.use(VueAxios, axios)
Vue.use(VueRouter)

Vue.config.productionTip = false

const routes = [
  { path: '/', component: Login },
  { path: '/chat', component: Chat },
]

const router = new VueRouter({ routes })

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
