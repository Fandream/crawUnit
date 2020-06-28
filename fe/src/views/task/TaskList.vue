<template>
  <div class="app-container">
    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <el-form
            class="filter-form"
            :model="filter"
            label-width="100px"
            label-position="right"
            inline
          >
            <el-form-item prop="node_id" :label="'节点'">
              <el-select
                v-model="filter.node_id"
                size="small"
                :placeholder="'节点'"
                @change="onFilterChange"
              >
                <el-option value="" :label="'全部'" />
                <el-option
                  v-for="node in nodeList"
                  :key="node._id"
                  :value="node._id"
                  :label="node.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="spider_id" :label="'爬虫'">
              <el-select
                v-model="filter.spider_id"
                size="small"
                :placeholder="'爬虫'"
                @change="onFilterChange"
                :disabled="isFilterSpiderDisabled"
              >
                <el-option value="" :label="'全部'" />
                <el-option
                  v-for="spider in spiderList"
                  :key="spider._id"
                  :value="spider._id"
                  :label="spider.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="status" :label="'状态'">
              <el-select
                v-model="filter.status"
                size="small"
                :placeholder="'状态'"
                @change="onFilterChange"
              >
                <el-option value="" :label="'全部'"></el-option>
                <el-option value="pending" :label="'待定'"></el-option>
                <el-option value="running" :label="'进行中'"></el-option>
                <el-option value="finished" :label="'已完成'"></el-option>
                <el-option value="error" :label="'错误'"></el-option>
                <el-option value="cancelled" :label="'取消'"></el-option>
                <el-option value="abnormal" :label="'异常'"></el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </div>
        <div class="right">
          <el-button
            class="btn-delete"
            @click="onRemoveMultipleTask"
            size="small"
            type="danger"
          >
            删除任务
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--legend-->
      <status-legend />
      <!--./legend-->

      <!--table list-->
      <el-table
        :data="filteredTableData"
        ref="table"
        class="table"
        :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
        border
        row-key="_id"
        @row-click="onRowClick"
        @selection-change="onSelectionChange"
      >
        <el-table-column
          type="selection"
          width="45"
          align="center"
          reserve-selection
        />
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'spider_name'"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
          >
            <template slot-scope="scope">
              <a
                href="javascript:"
                class="a-tag"
                @click="onClickSpider(scope.row)"
                >{{ scope.row[col.name] }}</a
              >
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name.match(/_ts$/)"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getTime(scope.row[col.name]) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'wait_duration'"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getWaitDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'runtime_duration'"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getRuntimeDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'total_duration'"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getTotalDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'node_name'"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              <a
                href="javascript:"
                class="a-tag"
                @click="onClickNode(scope.row)"
                >{{ scope.row[col.name] }}</a
              >
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'status'"
            :key="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              <status-tag :status="scope.row[col.name]" />
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :key="col.name"
            :property="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
          </el-table-column>
        </template>
        <el-table-column
          :label="'操作'"
          align="left"
          fixed="right"
          width="150px"
        >
          <template slot-scope="scope">
            <el-tooltip :content="'查看'" placement="top">
              <el-button
                type="primary"
                icon="el-icon-search"
                size="mini"
                @click="onView(scope.row)"
              ></el-button>
            </el-tooltip>
            <el-tooltip :content="'重新运行'" placement="top">
              <el-button
                type="warning"
                icon="el-icon-refresh"
                size="mini"
                @click="onRestart(scope.row, $event)"
              ></el-button>
            </el-tooltip>
            <el-tooltip :content="'删除'" placement="top">
              <el-button
                type="danger"
                icon="el-icon-delete"
                size="mini"
                @click="onRemove(scope.row, $event)"
              ></el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          @current-change="onPageChange"
          @size-change="onPageChange"
          :current-page.sync="pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pageSize"
          layout="sizes, prev, pager, next"
          :total="taskListTotalCount"
        >
        </el-pagination>
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import dayjs from 'dayjs';
import StatusTag from '../../components/Status/StatusTag';
import StatusLegend from '../../components/Status/StatusLegend';

export default {
  name: 'TaskList',
  components: { StatusLegend, StatusTag },
  data () {
    return {
      // setInterval handle
      handle: undefined,

      // determine if is edit mode
      isEditMode: false,

      // dialog visibility
      dialogVisible: false,

      // table columns
      columns: [
        { name: 'node_name', label: '节点', width: '120', align: 'center' },
        { name: 'spider_name', label: '爬虫', width: '120', align: 'center' },
        { name: 'status', label: '状态', width: '120', align: 'center' },
        // { name: 'create_ts', label: 'Create Time', width: '100' },
        { name: 'start_ts', label: '开始时间', width: '200', align: 'center' },
        { name: 'finish_ts', label: '结束时间', width: '200', align: 'center' },
        { name: 'wait_duration', label: '等待时长(秒)', width: '80', align: 'center' },
        { name: 'runtime_duration', label: '运行时长(秒)', width: '80', align: 'center' },
        { name: 'total_duration', label: '总时长(秒)', width: '80', align: 'center' },
        { name: 'result_count', label: '结果数', width: '80', align: 'center' }
        // { name: 'avg_num_results', label: 'Average Results Count per Second', width: '80' }
      ],

      multipleSelection: [],

      isFilterSpiderDisabled: false
    };
  },
  computed: {
    ...mapState('task', [
      'filter',
      'taskList',
      'taskListTotalCount',
      'taskForm'
    ]),
    ...mapState('spider', ['spiderList']),
    ...mapState('node', ['nodeList']),
    pageNum: {
      get () {
        return this.$store.state.task.pageNum;
      },
      set (value) {
        this.$store.commit('task/SET_PAGE_NUM', value);
      }
    },
    pageSize: {
      get () {
        return this.$store.state.task.pageSize;
      },
      set (value) {
        this.$store.commit('task/SET_PAGE_SIZE', value);
      }
    },
    filteredTableData () {
      return this.taskList
        .map(d => d)
        .sort((a, b) => (a.create_ts < b.create_ts ? 1 : -1))
        .filter(d => {
          // keyword
          if (!this.filter.keyword) return true;
          for (let i = 0; i < this.columns.length; i++) {
            const colName = this.columns[i].name;
            if (
              d[colName] &&
              d[colName]
                .toLowerCase()
                .indexOf(this.filter.keyword.toLowerCase()) > -1
            ) {
              return true;
            }
          }
          return false;
        });
    }
  },
  methods: {
    onSearch (value) {},
    onRefresh () {
      this.$store.dispatch('task/getTaskList');
    },
    onRemoveMultipleTask () {
      if (this.multipleSelection.length === 0) {
        this.$message({
          type: 'error',
          message: '请选择要删除的任务'
        });
        return;
      }
      this.$confirm('确定删除任务', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(() => {
          let ids = this.multipleSelection.map(item => item._id);
          this.$store.dispatch('task/deleteTaskMultiple', ids).then(resp => {
            if (resp.data.status === 'ok') {
              this.$message({
                type: 'success',
                message: '删除任务成功'
              });
              this.$store.dispatch('task/getTaskList');
              this.$refs['table'].clearSelection();
              return;
            }
            this.$message({
              type: 'error',
              message: resp.data.error
            });
          });
        })
        .catch(() => {});
    },
    onRemove (row, ev) {
      ev.stopPropagation();
      this.$confirm('您确定要删除该任务?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('task/deleteTask', row._id).then(() => {
          this.$message({
            type: 'success',
            message: '成功删除'
          });
        });
      });
    },
    onRestart (row, ev) {
      ev.stopPropagation();
      this.$confirm('确认重新运行该任务?', '提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('task/restartTask', row._id).then(() => {
          this.$message({
            type: 'success',
            message: '成功重新运行'
          });
        });
      });
    },
    onView (row) {
      this.$router.push(`/tasks/${row._id}`);
    },
    onClickSpider (row) {
      this.$router.push(`/spiders/${row.spider_id}`);
    },
    onClickNode (row) {
      this.$router.push(`/nodes/${row.node_id}`);
    },
    onPageChange () {
      setTimeout(() => {
        this.$store.dispatch('task/getTaskList');
      }, 0);
    },
    getTime (str) {
      if (str.match('^0001')) return 'NA';
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss');
    },
    getWaitDuration (row) {
      if (row.start_ts.match('^0001')) return 'NA';
      return dayjs(row.start_ts).diff(row.create_ts, 'second');
    },
    getRuntimeDuration (row) {
      if (row.finish_ts.match('^0001')) return 'NA';
      return dayjs(row.finish_ts).diff(row.start_ts, 'second');
    },
    getTotalDuration (row) {
      if (row.finish_ts.match('^0001')) return 'NA';
      return dayjs(row.finish_ts).diff(row.create_ts, 'second');
    },
    onRowClick (row, event, column) {
      if (column.label !== 'Action') {
        this.onView(row);
      }
    },
    onSelectionChange (val) {
      this.multipleSelection = val;
    },
    onFilterChange () {
      this.$store.dispatch('task/getTaskList');
    }
  },
  created () {
    this.$store.dispatch('task/getTaskList');
    this.$store.dispatch('spider/getSpiderList');
    this.$store.dispatch('node/getNodeList');
  },
  mounted () {
    this.handle = setInterval(() => {
      this.$store.dispatch('task/getTaskList');
    }, 5000);
  },
  destroyed () {
    clearInterval(this.handle);
  }
};
</script>

<style scoped lang="scss">
.el-dialog {
  .el-select {
    width: 100%;
  }
}

.filter {
  display: flex;
  justify-content: space-between;

  .left {
    .filter-select {
      width: 180px;
      margin-right: 10px;
    }
  }

  .filter-search {
    width: 240px;
  }

  .add {
  }
}

.table {
  margin-top: 8px;
  border-radius: 5px;

  .el-button {
    padding: 7px;
  }
}

.delete-confirm {
  background-color: red;
}

.el-table .a-tag {
  text-decoration: underline;
}

.pagination {
  margin-top: 10px;
  text-align: right;
}
</style>

<style scoped>
.el-table >>> tr {
  cursor: pointer;
}

.el-table >>> .el-badge .el-badge__content {
  font-size: 7px;
}
</style>
