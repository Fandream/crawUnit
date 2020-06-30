import Vue from 'vue'
import Vuex from 'vuex'
import app from './modules/app'
import user from './modules/user'
import tagsView from './modules/tagsView'
import dialogView from './modules/dialogView'
import node from './modules/node'
import spider from './modules/spider'
import task from './modules/task'
import file from './modules/file'
import project from './modules/project'
import getters from './getters'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    app,
    user,
    tagsView,
    dialogView,
    node,
    spider,
    task,
    file,
    project
  },
  getters
})

export default store
