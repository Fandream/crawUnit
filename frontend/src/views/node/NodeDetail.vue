<template>
  <div class="app-container">
    <!--selector-->
    <div class="selector">
      <label class="label">{{'节点'}}: </label>
      <el-select size="small" v-model="nodeForm._id" @change="onNodeChange">
        <el-option v-for="op in nodeList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="border-card">
      <el-tab-pane :label="'概览'" name="overview">
        <node-overview></node-overview>
      </el-tab-pane>
      <el-tab-pane :label="'已部署爬虫'" name="spiders" v-if="false">
        {{'已部署爬虫'}}
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import NodeOverview from '../../components/Overview/NodeOverview'
import NodeInstallation from '../../components/Node/NodeInstallation'

export default {
  name: 'NodeDetail',
  components: {
    NodeOverview
  },
  data () {
    return {
      activeTabName: 'overview'
    }
  },
  computed: {
    ...mapState('node', [
      'nodeList',
      'nodeForm'
    ])
  },
  methods: {
    onTabClick (tab) {

    },
    onNodeChange (id) {
      this.$router.push(`/nodes/${id}`)
    }
  },
  created () {
    // get list of nodes
    this.$store.dispatch('node/getNodeList')

    // get node basic info
    this.$store.dispatch('node/getNodeData', this.$route.params.id)

    // get node task list
    this.$store.dispatch('node/getTaskList', this.$route.params.id)
  }
}
</script>

<style scoped>
  .selector {
    display: flex;
    align-items: center;
    position: absolute;
    right: 30px;
    margin-top: 4px;
    /*float: right;*/
    z-index: 999;
  }

  .selector .el-select {
    padding-left: 10px;
  }

  .label {
    color: #909399;
    font-weight: 100;
    width: 100px;
    text-align: right;
  }
</style>
<style lang="scss">
  .selector {
    .el-select {
      .el-input {
        .el-input_inner {
          height: 26px;
        }
      }
    }
  }
</style>
