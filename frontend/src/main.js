import Vue from 'vue'

import 'normalize.css/normalize.css' // A modern alternative to CSS resets

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import locale from 'element-ui/lib/locale/lang/en' // lang i18n

import '@/styles/index.scss' // global css

import 'font-awesome/scss/font-awesome.scss'// FontAwesome
import { library } from '@fortawesome/fontawesome-svg-core'
import { fab } from '@fortawesome/free-brands-svg-icons'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { FontAwesomeIcon, FontAwesomeLayers, FontAwesomeLayersText } from '@fortawesome/vue-fontawesome'

import 'codemirror/lib/codemirror.js'
import { codemirror } from 'vue-codemirror-lite'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/darcula.css'

import App from './App'
import store from './store'
import router from './router'

import '@/icons' // icon

import request from './api/request'
import utils from './utils'

import iView from 'iview'; // 引入iview
import 'iview/dist/styles/iview.css';

// vue-tour
import VueTour from 'vue-tour'
import 'vue-tour/dist/vue-tour.css'
Vue.use(iView);

Vue.use(iView);

// code mirror
Vue.use(codemirror)

// element-ui
Vue.use(ElementUI, { locale })

// font-awesome
library.add(fab)
library.add(far)
library.add(fas)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('font-awesome-layers-text', FontAwesomeLayersText)

// vue-tour
Vue.use(VueTour)

Vue.config.productionTip = false

// inject request api
Vue.prototype.$request = request

// inject utils
Vue.prototype.$utils = utils

const app = new Vue({
  router,
  store,
  render: h => h(App)
})

async function Login () {
  const res = await store.dispatch('user/login')
  if (res.status === 200) {
    // success
    await store.dispatch('user/getInfo')
  } else {
    console.log(res)
  }
}
Login().then(i => {
  app.$mount('#app')
})
export default app
