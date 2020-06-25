<template>
  <div class="app-container schedule-list">
    <parameters-dialog
      :visible="isParametersVisible"
      :param="scheduleForm.param"
      @confirm="onParametersConfirm"
      @close="isParametersVisible = false"
    />

    <!--add popup-->
    <el-dialog
      :title="dialogTitle"
      :visible.sync="dialogVisible"
      width="640px"
      :before-close="onDialogClose"
    >
      <el-form
        label-width="180px"
        class="add-form"
        :model="scheduleForm"
        :inline-message="true"
        ref="scheduleForm"
        label-position="right"
      >
        <el-form-item :label="'定时任务名称'" prop="name" required>
          <el-input id="schedule-name" v-model="scheduleForm.name" :placeholder="'定时任务名称'"></el-input>
        </el-form-item>
        <el-form-item :label="'运行类型'" prop="run_type" required>
          <el-select id="run-type" v-model="scheduleForm.run_type" :placeholder="'运行类型'">
            <el-option value="all-nodes" :label="'所有节点'"/>
            <el-option value="selected-nodes" :label="'指定节点'"/>
            <el-option value="random" :label="'随机'"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="scheduleForm.run_type === 'selected-nodes'" :label="'节点'" prop="node_ids" required>
          <el-select id="node-ids" v-model="scheduleForm.node_ids" :placeholder="'节点'" multiple filterable>
            <el-option
              v-for="op in nodeList"
              :key="op._id"
              :value="op._id"
              :label="op.name"
              :disabled="op.status === 'offline'"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!isDisabledSpiderSchedule" :label="'爬虫'" prop="spider_id" required>
          <el-select
            id="spider-id"
            v-model="scheduleForm.spider_id"
            :placeholder="'爬虫'"
            filterable
            :disabled="isDisabledSpiderSchedule"
            @change="onSpiderChange"
          >
            <el-option
              v-for="op in spiderList"
              :key="op._id"
              :value="op._id"
              :label="`${op.display_name} (${op.name})`"
              :disabled="isDisabledSpider(op)"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-else :label="'爬虫'" required>
          <el-select
            :value="spiderId"
            :placeholder="'爬虫'"
            filterable
            :disabled="isDisabledSpiderSchedule"
          >
            <el-option
              v-for="op in spiderList"
              :key="op._id"
              :value="op._id"
              :label="`${op.display_name} (${op.name})`"
              :disabled="isDisabledSpider(op)"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="spiderForm.is_scrapy" :label="'Scrapy 爬虫'" prop="scrapy_spider" required
                      inline-message>
          <el-select v-model="scheduleForm.scrapy_spider" :placeholder="'Scrapy 爬虫'" :disabled="isLoading">
            <el-option
              v-for="s in spiderForm.spider_names"
              :key="s"
              :label="s"
              :value="s"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="spiderForm.is_scrapy" :label="'Scrapy 日志等级'" prop="scrapy_spider" required
                      inline-message>
          <el-select v-model="scheduleForm.scrapy_log_level" :placeholder="'Scrapy 日志等级'">
            <el-option value="INFO" label="INFO"/>
            <el-option value="DEBUG" label="DEBUG"/>
            <el-option value="WARN" label="WARN"/>
            <el-option value="ERROR" label="ERROR"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="'Cron'" prop="cron" required>
          <el-input
            class="cron"
            ref="cron"
            v-model="scheduleForm.cron"
            :placeholder="`${'[分] [时] [天] [月] [星期几]'}`"
            style="width: calc(100% - 100px)"
          >
          </el-input>
          <el-button
            class="cron-edit"
            type="primary"
            icon="el-icon-edit"
            style="width: 100px"
            @click="onShowCronDialog"
          >
            {{'编辑'}}
          </el-button>
        </el-form-item>
        <el-form-item :label="'执行命令'" prop="cmd">
          <el-input
            id="cmd"
            v-model="spiderForm.cmd"
            :placeholder="'执行命令'"
            disabled
          />
        </el-form-item>
        <el-form-item v-if="spiderForm.type === 'customized'" :label="'参数'" prop="param">
          <template v-if="spiderForm.is_scrapy">
            <el-input v-model="scheduleForm.param" :placeholder="'参数'" class="param-input"/>
            <el-button type="primary" icon="el-icon-edit" class="param-btn" @click="onOpenParameters"/>
          </template>
          <template v-else>
            <el-input v-model="scheduleForm.param" :placeholder="'参数'"></el-input>
          </template>
        </el-form-item>
        <el-form-item :label="'定时任务描述'" prop="description">
          <el-input id="schedule-description" v-model="scheduleForm.description" type="textarea"
                    :placeholder="'定时任务描述'"></el-input>
        </el-form-item>
      </el-form>
      <!--取消、保存-->
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="onCancel">{{'取消'}}</el-button>
        <el-button id="btn-submit" size="small" type="primary" @click="onAddSubmit" :disabled="isLoading">{{'提交'}}</el-button>
      </span>
    </el-dialog>
    <!--./add popup-->

    <!--view tasks popup-->
    <el-dialog
      :title="'Tasks'"
      :visible.sync="isViewTasksDialogVisible"
      width="calc(100% - 240px)"
      :before-close="() => this.isViewTasksDialogVisible = false"
    >
      <schedule-task-list ref="schedule-task-list"/>
    </el-dialog>
    <!--./view tasks popup-->

    <!--cron generation popup-->
    <el-dialog title="生成 Cron" :visible.sync="cronDialogVisible">
      <vue-cron-linux ref="vue-cron-linux" :data="scheduleForm.cron" :i18n="lang" @submit="onCronChange"/>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="cronDialogVisible = false">{{'取消'}}</el-button>
        <el-button size="small" type="primary" @click="onCronDialogSubmit">{{'确认'}}</el-button>
      </span>
    </el-dialog>
    <!--./cron generation popup-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="scheduleForm.spider_id"
      @close="() => crawlConfirmDialogVisible = false"
      @confirm="() => crawlConfirmDialogVisible = false"
    />
    <!--./crawl confirm dialog-->

    <el-card style="border-radius: 0" class="schedule-list">
      <!--filter-->
      <div class="filter">
        <div class="right">
          <el-button size="small" type="primary"
                     icon="el-icon-plus"
                     class="btn-add"
                     @click="onAdd">
            {{'添加定时任务'}}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table :data="filteredTableData"
                class="table" height="500"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border>
        <template v-for="col in columns">
          <el-table-column v-if="col.name === 'status'"
                           :key="col.name"
                           :property="col.name"
                           :label="col.label"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              <el-tooltip v-if="scope.row[col.name] === 'error'" :content="scope.row['message']" placement="top">
                <el-tag class="status-tag" type="danger">
                  {{scope.row[col.name] ? scope.row[col.name] : 'NA'}}
                </el-tag>
              </el-tooltip>
              <el-tag class="status-tag" v-else>
                {{scope.row[col.name] ? scope.row[col.name] : 'NA'}}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'run_type'" :key="col.name" :label="col.label"
                           :width="col.width">
            <template slot-scope="scope">
              <template v-if="scope.row.run_type === 'all-nodes'">{{'All Nodes'}}</template>
              <template v-else-if="scope.row.run_type === 'selected-nodes'">{{'Selected Nodes'}}</template>
              <template v-else-if="scope.row.run_type === 'random'">{{'Random'}}</template>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'node_names'" :key="col.name" :label="col.label"
                           :width="col.width">
            <template slot-scope="scope">
              {{scope.row.nodes.map(d => d.name).join(', ')}}
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'enable'" :key="col.name" :label="col.label" :width="col.width">
            <template slot-scope="scope">
              <el-switch
                v-model="scope.row.enabled"
                active-color="#13ce66"
                inactive-color="#ff4949"
                @change="onEnabledChange(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column v-else :key="col.name"
                           :property="col.name"
                           :label="col.label"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{scope.row[col.name]}}
            </template>
          </el-table-column>
        </template>
        <el-table-column :label="'操作'" class="actions" align="left" width="170" fixed="right">
          <template slot-scope="scope">
            <!--edit-->
            <el-tooltip :content="'编辑'" placement="top">
              <el-button type="warning" icon="el-icon-edit" size="mini" @click="onEdit(scope.row)"></el-button>
            </el-tooltip>
            <!--./edit-->

            <!--delete-->
            <el-tooltip :content="'删除'" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
            </el-tooltip>
            <!--./delete-->

            <!--view tasks-->
            <el-tooltip :content="'查看任务'" placement="top">
              <el-button type="primary" icon="el-icon-search" size="mini" @click="onViewTasks(scope.row)"></el-button>
            </el-tooltip>
            <!--./view tasks-->

            <!--run-->
            <el-tooltip :content="'运行'" placement="top">
              <el-button type="success" icon="fa fa-bug" size="mini" @click="onRun(scope.row)"></el-button>
            </el-tooltip>
            <!--./run-->
          </template>
        </el-table-column>
      </el-table>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import request from '../../api/request'
import VueCronLinux from '../../components/Cron'
import {
  mapState
} from 'vuex'
import ParametersDialog from '../../components/Common/ParametersDialog'
import ScheduleTaskList from '../../components/Schedule/ScheduleTaskList'
import CrawlConfirmDialog from '../../components/Common/CrawlConfirmDialog'

export default {
  name: 'ScheduleList',
  components: {
    CrawlConfirmDialog,
    ScheduleTaskList,
    VueCronLinux,
    ParametersDialog
  },
  data () {
    return {
      columns: [
        { name: 'name', label: '名称', width: '150px' },
        { name: 'cron', label: 'Cron', width: '120px' },
        { name: 'run_type', label: '运行类型', width: '120px' },
        { name: 'node_names', label: '节点', width: '150px' },
        { name: 'spider_name', label: '爬虫', width: '150px' },
        { name: 'scrapy_spider', label: 'Scrapy爬虫', width: '150px' },
        { name: 'param', label: '参数', width: '150px' },
        { name: 'description', label: '描述', width: '200px' },
        { name: 'enable', label: '启用/禁用', width: '120px' },
        { name: 'username', label: '所有者', width: '100px' }
        // { name: 'status', label: 'Status', width: '100px' }
      ],
      isEdit: false,
      dialogTitle: '',
      dialogVisible: false,
      cronDialogVisible: false,
      expression: '',
      spiderList: [],
      nodeList: [],
      isShowCron: false,
      isLoading: false,
      isParametersVisible: false,
      isViewTasksDialogVisible: false,
      crawlConfirmDialogVisible: false
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    ...mapState('schedule', [
      'scheduleList',
      'scheduleForm'
    ]),
    lang () {
      const lang = this.$store.state.lang.lang || window.localStorage.getItem('lang')
      if (!lang) return 'cn'
      if (lang === 'zh') return 'cn'
      return 'en'
    },
    filteredTableData () {
      return this.scheduleList
    },
    spider () {
      for (let i = 0; i < this.spiderList.length; i++) {
        if (this.spiderList[i]._id === this.scheduleForm.spider_id) {
          return this.spiderList[i]
        }
      }
      return {}
    },
    isDisabledSpiderSchedule () {
      return false
    }
  },
  methods: {
    onDialogClose () {
      this.dialogVisible = false
    },
    onCancel () {
      this.dialogVisible = false
    },
    onAdd () {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', { node_ids: [] })
    },
    onAddSubmit () {
      this.$refs.scheduleForm.validate(res => {
        if (res) {
          const form = JSON.parse(JSON.stringify(this.scheduleForm))
          form.cron = '0 ' + this.scheduleForm.cron
          if (this.isEdit) {
            request.post(`/schedules/${this.scheduleForm._id}`, form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('schedule/getScheduleList')
              this.$message.success('已保存定时任务')
            })
          } else {
            request.put('/schedules', form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('schedule/getScheduleList')
              this.$message.success('已添加定时任务')
            })
          }
        }
      })
    },
    isShowRun (row) {
    },
    async onEdit (row) {
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.dialogVisible = true
      this.isEdit = true

      this.isLoading = true
      if (!this.scheduleForm.scrapy_log_level) {
        this.$set(this.scheduleForm, 'scrapy_log_level', 'INFO')
      }
      await this.$store.dispatch('spider/getSpiderData', row.spider_id)
      if (this.spiderForm.is_scrapy) {
        await this.$store.dispatch('spider/getSpiderScrapySpiders', row.spider_id)
        if (!this.scheduleForm.scrapy_spider) {
          if (this.spiderForm.spider_names && this.spiderForm.spider_names.length > 0) {
            this.$set(this.scheduleForm, 'scrapy_spider', this.spiderForm.spider_names[0])
          }
        }
      }
      this.isLoading = false
    },
    onRemove (row) {
      this.$confirm('确定删除定时任务?', '提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('schedule/removeSchedule', row._id)
          .then(() => {
            setTimeout(() => {
              this.$store.dispatch('schedule/getScheduleList')
              this.$message.success('已删除定时任务')
            }, 100)
          })
      }).catch(() => {
      })
    },
    isDisabledSpider (spider) {
      if (spider.type === 'customized') {
        return !spider.cmd
      } else {
        return false
      }
    },
    async onEnabledChange (row) {
      let res
      if (row.enabled) {
        res = await this.$store.dispatch('schedule/enableSchedule', row._id)
      } else {
        res = await this.$store.dispatch('schedule/disableSchedule', row._id)
      }
      if (!res || res.data.error) {
        this.$message.error(`${row.enabled ? 'Enabling' : 'Disabling'} the schedule unsuccessful`)
      } else {
        this.$message.success(`${row.enabled ? 'Enabling' : 'Disabling'} the schedule successful`)
      }
    },
    onCronChange (value) {
      this.$set(this.scheduleForm, 'cron', value)
    },
    onCronDialogSubmit () {
      const valid = this.$refs['vue-cron-linux'].submit()
      if (valid) {
        this.cronDialogVisible = false
      }
    },
    onOpenParameters () {
      this.isParametersVisible = true
    },
    onParametersConfirm (value) {
      this.scheduleForm.param = value
      this.isParametersVisible = false
    },
    async onSpiderChange (spiderId) {
      await this.$store.dispatch('spider/getSpiderData', spiderId)
      if (this.spiderForm.type === 'customized' && this.spiderForm.is_scrapy) {
        this.isLoading = true
        await this.$store.dispatch('spider/getSpiderScrapySpiders', spiderId)
        this.isLoading = false
        this.$set(this.scheduleForm, 'scrapy_spider', this.spiderForm.spider_names[0])
        this.$set(this.scheduleForm, 'scrapy_log_level', 'INFO')
      }
    },
    onShowCronDialog () {
      this.cronDialogVisible = true
    },
    async onViewTasks (row) {
      this.isViewTasksDialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      setTimeout(() => {
        this.$refs['schedule-task-list'].update()
      }, 100)
    },
    async onRun (row) {
      this.crawlConfirmDialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
    }
  },
  created () {
    this.$store.dispatch('schedule/getScheduleList')

    // 节点列表
    request.get('/nodes', {}).then(response => {
      this.nodeList = response.data.data.map(d => {
        d.systemInfo = {
          os: '',
          arch: '',
          num_cpu: '',
          executables: []
        }
        return d
      })
    })

    // 爬虫列表
    request.get('/spiders', { owner_type: 'all' })
      .then(response => {
        this.spiderList = response.data.data.list || []
      })
  }
}
</script>

<style scoped>
  .filter .right {
    text-align: right;
  }

  .table {
    min-height: 360px;
    margin-top: 10px;
  }

  .table .el-button {
    width: 28px;
    height: 28px;
    padding: 0;
  }

  .status-tag {
    cursor: pointer;
  }

  .schedule-list >>> .param-input {
    width: calc(100% - 56px);
  }

  .schedule-list >>> .param-input .el-input__inner {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    border-right: none;
  }

  .schedule-list >>> .param-btn {
    width: 56px;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }

  .cron {
    width: calc(100% - 100px);
  }

  .cron >>> .el-input__inner {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    border-right: none;
  }

  .cron-edit {
    width: 100px;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
</style>
