<template>
  <div class="task-overview">
    <el-row class="action-wrapper">
      <el-button
        type="primary"
        size="small"
        icon="el-icon-position"
        @click="onNavigateToSpider"
      >
        {{'导航到爬虫'}}
      </el-button>
      <el-button
        type="warning"
        size="small"
        icon="el-icon-position"
        @click="onNavigateToNode"
      >
        {{'导航到节点'}}
      </el-button>
    </el-row>
    <el-row class="content">
      <el-col :span="12" style="padding-right: 20px;">
        <el-row class="task-info-overview-wrapper wrapper">
          <h4 class="title">{{'任务信息'}}</h4>
          <task-info-view @click-log="() => $emit('click-log')"/>
        </el-row>
        <el-row style="border-bottom:1px solid #e4e7ed;margin:0 0 20px 0;padding-bottom:20px;"/>
      </el-col>

      <el-col :span="12">
        <el-row class="task-info-spider-wrapper wrapper">
          <h4 class="title spider-title" @click="onNavigateToSpider">
            <i class="fa fa-search" style="margin-right: 5px"></i>
            {{'爬虫信息'}}</h4>
          <spider-info-view :is-view="true"/>
        </el-row>
        <el-row class="task-info-node-wrapper wrapper">
          <h4 class="title node-title" @click="onNavigateToNode">
            <i class="fa fa-search" style="margin-right: 5px"></i>
            {{'节点信息'}}</h4>
          <node-info-view :is-view="true"/>
        </el-row>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import SpiderInfoView from '../InfoView/SpiderInfoView'
import NodeInfoView from '../InfoView/NodeInfoView'
import TaskInfoView from '../InfoView/TaskInfoView'

export default {
  name: 'SpiderOverview',
  components: {
    NodeInfoView,
    SpiderInfoView,
    TaskInfoView
  },
  computed: {
    ...mapState('node', [
      'nodeForm'
    ]),
    ...mapState('spider', [
      'spiderForm'
    ])
  },
  methods: {
    onNavigateToSpider () {
      this.$router.push(`/spiders/${this.spiderForm._id}`)
    },
    onNavigateToNode () {
      this.$router.push(`/nodes/${this.nodeForm._id}`)
    }
  },
  created () {
  }
}
</script>

<style scoped>
  .title {
    margin: 10px 0 3px 0;
    text-align: center;
    display: inline-block;
  }

  .wrapper {
    text-align: center;
  }

  .spider-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }

  .node-title,
  .spider-title {
    cursor: pointer;
  }

  .node-title:hover,
  .spider-title:hover {
    text-decoration: underline;
  }

  .title > i {
    color: grey;
  }

  .content {
    margin-top: 10px;
  }

  .action-wrapper {
    text-align: right;
    padding-bottom: 10px;
    border-bottom: 1px solid #DCDFE6;
  }
</style>
