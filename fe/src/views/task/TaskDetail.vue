<template>
  <div class="app-container">
    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="border-card">
      <el-tab-pane :label="'概览'" name="overview">
        <task-overview @click-log="activeTabName = 'log'" />
      </el-tab-pane>
      <!-- <el-tab-pane :label="'Log'" name="log">
        <log-view @search="getTaskLog(true)"/>
      </el-tab-pane> -->
      <el-tab-pane :label="'结果'" name="results">
        <general-table-view
          :data="taskResultsData"
          :columns="taskResultsColumns"
          :page-num="resultsPageNum"
          :page-size="resultsPageSize"
          :total="taskResultsTotalCount"
          @page-change="onResultsPageChange"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import TaskOverview from '../../components/Overview/TaskOverview';
import GeneralTableView from '../../components/TableView/GeneralTableView';
import LogView from '../../components/ScrollView/LogView';

export default {
  name: 'TaskDetail',
  components: {
    GeneralTableView,
    TaskOverview
  },
  data () {
    return {
      activeTabName: 'overview',
      handle: undefined
    };
  },
  computed: {
    ...mapState('task', [
      'taskForm',
      'taskResultsData',
      'taskResultsTotalCount',
      'taskLog',
      'logKeyword',
      'isLogAutoFetch',
      'currentLogIndex',
      'activeErrorLogItem'
    ]),
    ...mapGetters('task', ['taskResultsColumns', 'logData']),
    ...mapState('file', ['currentPath']),
    ...mapState('deploy', ['deployList']),
    resultsPageNum: {
      get () {
        return this.$store.state.task.resultsPageNum;
      },
      set (value) {
        this.$store.commit('task/SET_RESULTS_PAGE_NUM', value);
      }
    },
    resultsPageSize: {
      get () {
        return this.$store.state.task.resultsPageSize;
      },
      set (value) {
        this.$store.commit('task/SET_RESULTS_PAGE_SIZE', value);
      }
    },
    isLogAutoScroll: {
      get () {
        return this.$store.state.task.isLogAutoScroll;
      },
      set (value) {
        this.$store.commit('task/SET_IS_LOG_AUTO_SCROLL', value);
      }
    },
    isLogAutoFetch: {
      get () {
        return this.$store.state.task.isLogAutoFetch;
      },
      set (value) {
        this.$store.commit('task/SET_IS_LOG_AUTO_FETCH', value);
      }
    },
    isLogFetchLoading: {
      get () {
        return this.$store.state.task.isLogFetchLoading;
      },
      set (value) {
        this.$store.commit('task/SET_IS_LOG_FETCH_LOADING', value);
      }
    },
    currentLogIndex: {
      get () {
        return this.$store.state.task.currentLogIndex;
      },
      set (value) {
        this.$store.commit('task/SET_CURRENT_LOG_INDEX', value);
      }
    },
    logIndexMap () {
      const map = new Map();
      this.logData.forEach((d, index) => {
        map.set(d._id, index);
      });
      return map;
    },
    isRunning () {
      return ['pending', 'running'].includes(this.taskForm.status);
    }
  },
  methods: {
    onTabClick (tab) {},
    onSpiderChange (id) {
      this.$router.push(`/spiders/${id}`);
    },
    onResultsPageChange (payload) {
      const { pageNum, pageSize } = payload;
      this.resultsPageNum = pageNum;
      this.resultsPageSize = pageSize;
      this.$store.dispatch('task/getTaskResults', this.$route.params.id);
    },
    downloadCSV () {
      this.$store.dispatch('task/getTaskResultExcel', this.$route.params.id);
    },
    async getTaskLog (showLoading) {
      if (showLoading) {
        this.isLogFetchLoading = true;
      }
      await this.$store.dispatch('task/getTaskLog', {
        id: this.$route.params.id,
        keyword: this.logKeyword
      });
      this.currentLogIndex =
        this.logIndexMap.get(this.activeErrorLogItem.log_id) + 1 || 0;
      this.isLogFetchLoading = false;
      await this.$store.dispatch('task/getTaskErrorLog', this.$route.params.id);
    }
  },
  async created () {
    await this.$store.dispatch('task/getTaskData', this.$route.params.id);

    this.isLogAutoFetch = !!this.isRunning;
    this.isLogAutoScroll = !!this.isRunning;

    await this.$store.dispatch('task/getTaskResults', this.$route.params.id);

    this.getTaskLog();
    this.handle = setInterval(() => {
      if (this.isLogAutoFetch) {
        this.$store.dispatch('task/getTaskData', this.$route.params.id);
        this.$store.dispatch('task/getTaskResults', this.$route.params.id);
        this.getTaskLog();
      }
    }, 5000);
  },
  destroyed () {
    clearInterval(this.handle);
  }
};
</script>

<style scoped>
.selector {
  display: flex;
  align-items: center;
  position: absolute;
  right: 20px;
  /*float: right;*/
  z-index: 999;
  margin-top: -7px;
}

.selector .el-select {
  padding-left: 10px;
}

.log-view {
  margin: 20px;
  height: 640px;
}

.log-view pre {
  height: 100%;
  overflow-x: auto;
  overflow-y: auto;
}

.button-group {
  margin-bottom: 10px;
  text-align: right;
}
</style>
