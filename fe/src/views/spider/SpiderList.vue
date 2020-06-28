<template>
  <div class="app-container">
    <!--add dialog-->
    <el-dialog
      :title="'添加爬虫'"
      width="40%"
      :visible.sync="addDialogVisible"
      :before-close="onAddDialogClose"
    >
      <el-tabs :active-name="activeTabName">
        <!-- customized -->
        <el-tab-pane name="customized" :label="'自定义'">
          <el-form
            :model="spiderForm"
            ref="addCustomizedForm"
            inline-message
            label-width="120px"
          >
            <el-form-item :label="'爬虫名称'" prop="name" required>
              <el-input
                id="spider-name"
                v-model="spiderForm.name"
                :placeholder="'爬虫名称'"
              />
            </el-form-item>
            <el-form-item :label="'显示名称'" prop="display_name" required>
              <el-input
                id="display-name"
                v-model="spiderForm.display_name"
                :placeholder="'显示名称'"
              />
            </el-form-item>
            <el-form-item :label="'项目'" prop="project_id" required>
              <el-select
                v-model="spiderForm.project_id"
                :placeholder="'项目'"
                filterable
              >
                <el-option
                  v-for="p in projectList"
                  :key="p._id"
                  :value="p._id"
                  :label="p.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="'执行命令'" prop="cmd" required>
              <el-input
                id="cmd"
                v-model="spiderForm.cmd"
                :placeholder="'执行命令'"
                :disabled="spiderForm.is_scrapy"
              />
            </el-form-item>
            <el-form-item :label="'结果'" prop="col">
              <el-input
                id="col"
                v-model="spiderForm.col"
                :placeholder="'默认: ' + 'results_<spider_name>'"
              />
            </el-form-item>
            <el-form-item
              :label="'上传Zip文件'"
              label-width="120px"
              name="site"
            >
              <el-upload
                :action="$request.baseUrl + '/spiders'"
                :data="uploadForm"
                :headers="{ Authorization: token }"
                :on-success="onUploadSuccess"
                :file-list="fileList"
                :before-upload="beforeUpload"
              >
                <el-button
                  id="upload"
                  size="small"
                  type="primary"
                  icon="el-icon-upload"
                >
                  {{ "上传" }}
                </el-button>
              </el-upload>
            </el-form-item>
            <el-row>
              <el-col :span="8">
                <el-form-item :label="'是否为 Scrapy'" prop="is_scrapy">
                  <el-switch
                    v-model="spiderForm.is_scrapy"
                    active-color="#13ce66"
                    @change="onIsScrapy"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>

          <div class="actions">
            <el-button size="small" type="primary" @click="onAddCustomized">{{
              "添加"
            }}</el-button>
          </div>
        </el-tab-pane>
        <!-- configurable -->
      </el-tabs>
    </el-dialog>
    <!--./add dialog-->

    <!--running tasks dialog-->
    <el-dialog
      :visible.sync="isRunningTasksDialogVisible"
      :title="
        `${'最近任务'} (${'爬虫'}: ${activeSpider ? activeSpider.name : ''})`
      "
      width="920px"
    >
      <el-tabs v-model="activeSpiderTaskStatus">
        <el-tab-pane name="pending" :label="'待定'" />
        <el-tab-pane name="running" :label="'进行中'" />
        <el-tab-pane name="finished" :label="'已完成'" />
        <el-tab-pane name="error" :label="'错误'" />
        <el-tab-pane name="cancelled" :label="'已取消'" />
        <el-tab-pane name="abnormal" :label="'异常'" />
      </el-tabs>
      <el-table
        :data="activeNodeList"
        class="table"
        :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
        border
        default-expand-all
      >
        <el-table-column type="expand">
          <template slot-scope="scope">
            <h4 style="margin: 5px 10px">{{ "任务" }}</h4>
            <el-table
              :data="getTasksByNode(scope.row)"
              class="table"
              border
              style="margin: 5px 10px"
              max-height="240px"
              @row-click="onViewTask"
            >
              <el-table-column
                :label="'创建时间'"
                prop="create_ts"
                width="140px"
              />
              <el-table-column
                :label="'开始时间'"
                prop="start_ts"
                width="140px"
              />
              <el-table-column
                :label="'结束时间'"
                prop="finish_ts"
                width="140px"
              />
              <el-table-column :label="'参数'" prop="param" width="120px" />
              <el-table-column :label="'状态'" width="120px">
                <template slot-scope="scope">
                  <status-tag :status="scope.row.status" />
                </template>
              </el-table-column>
              <el-table-column
                :label="'结果数'"
                prop="result_count"
                width="80px"
              />
              <el-table-column :label="'操作'" width="auto">
                <template slot-scope="scope">
                  <el-button
                    v-if="['pending', 'running'].includes(scope.row.status)"
                    type="danger"
                    size="mini"
                    icon="el-icon-video-pause"
                    @click="onStop(scope.row, $event)"
                  />
                </template>
              </el-table-column>
            </el-table>
          </template>
        </el-table-column>
        <el-table-column :label="'节点'" width="150px" prop="name" />
        <el-table-column :label="'状态'" width="120px" prop="status">
          <template slot-scope="scope">
            <el-tag type="info" v-if="scope.row.status === 'offline'">{{
              "离线"
            }}</el-tag>
            <el-tag type="success" v-else-if="scope.row.status === 'online'">{{
              "在线"
            }}</el-tag>
            <el-tag type="danger" v-else>{{ "未知" }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="'描述'" width="auto" prop="description" />
      </el-table>
      <template slot="footer">
        <el-button
          type="primary"
          size="small"
          @click="isRunningTasksDialogVisible = false"
          >{{ "确定" }}</el-button
        >
      </template>
    </el-dialog>
    <!--./running tasks dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="activeSpiderId"
      :spiders="selectedSpiders"
      :multiple="isMultiple"
      @close="onCrawlConfirmDialogClose"
      @confirm="onCrawlConfirm"
    />
    <!--./crawl confirm dialog-->

    <!--copy dialog-->
    <copy-spider-dialog
      :visible="copyDialogVisible"
      :spider-id="activeSpiderId"
      @close="copyDialogVisible = false"
      @confirm="onCopyConfirm"
    />
    <!--./copy dialog-->

    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <el-form :inline="true">
            <!--            <el-form-item>-->
            <!--              <el-select clearable @change="onSpiderTypeChange" placeholder="爬虫类型" size="small" v-model="filter.type">-->
            <!--                <el-option v-for="item in types" :value="item.type" :key="item.type"-->
            <!--                           :label="item.type === 'customized'? '自定义':item.type "/>-->
            <!--              </el-select>-->
            <!--            </el-form-item>-->
            <el-form-item>
              <el-select
                v-model="filter.project_id"
                size="small"
                :placeholder="'项目'"
                @change="getList"
              >
                <el-option value="" :label="'所有项目'" />
                <el-option
                  v-for="p in projectList"
                  :key="p._id"
                  :value="p._id"
                  :label="p.name"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-input
                v-model="filter.keyword"
                size="small"
                :placeholder="'爬虫名称'"
                clearable
                @keyup.enter.native="onSearch"
              >
                <i slot="suffix" class="el-input__icon el-icon-search"></i>
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button
                size="small"
                type="success"
                class="btn refresh"
                @click="onRefresh"
              >
                {{ "搜索" }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="right">
          <el-button
            v-if="this.selectedSpiders.length"
            size="small"
            type="danger"
            icon="el-icon-video-play"
            class="btn add"
            @click="onCrawlSelectedSpiders"
            style="font-weight: bolder"
          >
            {{ "运行" }}
          </el-button>
          <el-button
            v-if="this.selectedSpiders.length"
            size="small"
            type="info"
            :icon="isStopLoading ? 'el-icon-loading' : 'el-icon-video-pause'"
            class="btn add"
            @click="onStopSelectedSpiders"
            style="font-weight: bolder"
          >
            {{ "停止" }}
          </el-button>
          <el-button
            v-if="this.selectedSpiders.length"
            size="small"
            type="danger"
            :icon="isRemoveLoading ? 'el-icon-loading' : 'el-icon-delete'"
            class="btn add"
            @click="onRemoveSelectedSpiders"
            style="font-weight: bolder"
          >
            {{ "删除" }}
          </el-button>
          <el-button
            size="small"
            type="success"
            icon="el-icon-plus"
            class="btn add"
            @click="onAdd"
            style="font-weight: bolder"
          >
            {{ "添加爬虫" }}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--legend-->
      <status-legend />
      <!--./legend-->

      <!--table list-->
      <el-table
        :data="spiderList"
        class="table"
        ref="table"
        :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
        border
        row-key="_id"
        @row-click="onRowClick"
        @sort-change="onSortChange"
        @selection-change="onSpiderSelect"
      >
        <el-table-column
          type="selection"
          width="45"
          align="center"
          reserve-selection
        />
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'type'"
            :key="col.name"
            :label="col.label"
            align="left"
            :width="col.width"
            :sortable="col.sortable"
          >
            <template slot-scope="scope">
              {{ scope.row.type }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_5_errors'"
            :key="col.name"
            :label="col.label"
            :width="col.width"
            :sortable="col.sortable"
            align="center"
          >
            <template slot-scope="scope">
              <div :style="{ color: scope.row[col.name] > 0 ? 'red' : '' }">
                {{ scope.row[col.name] }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'cmd'"
            :key="col.name"
            :label="col.label"
            :width="col.width"
            :sortable="col.sortable"
            align="left"
          >
            <template slot-scope="scope">
              <el-input v-model="scope.row[col.name]"></el-input>
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
            v-else-if="col.name === 'last_status'"
            :key="col.name"
            :label="col.label"
            align="left"
            :width="col.width"
            :sortable="col.sortable"
          >
            <template slot-scope="scope">
              <status-tag :status="scope.row.last_status" />
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="['is_scrapy', 'is_long_task'].includes(col.name)"
            :key="col.name"
            :label="col.label"
            align="left"
            :width="col.width"
            :sortable="col.sortable"
          >
            <template slot-scope="scope">
              <el-switch
                v-if="scope.row.type === 'customized'"
                v-model="scope.row[col.name]"
                active-color="#13ce66"
                disabled
              />
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'latest_tasks'"
            :key="col.name"
            :label="col.label"
            :width="col.width"
            :align="col.align"
            class-name="latest-tasks"
          >
            <template slot-scope="scope">
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'pending') > 0"
                type="primary"
                size="small"
              >
                <i class="el-icon-loading"></i>
                {{ getTaskCountByStatus(scope.row, "pending") }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'running') > 0"
                type="warning"
                size="small"
              >
                <i class="el-icon-loading"></i>
                {{ getTaskCountByStatus(scope.row, "running") }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'finished') > 0"
                type="success"
                size="small"
              >
                <i class="el-icon-check"></i>
                {{ getTaskCountByStatus(scope.row, "finished") }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'error') > 0"
                type="danger"
                size="small"
              >
                <i class="el-icon-error"></i>
                {{ getTaskCountByStatus(scope.row, "error") }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'cancelled') > 0"
                type="info"
                size="small"
              >
                <i class="el-icon-video-pause"></i>
                {{ getTaskCountByStatus(scope.row, "cancelled") }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'abnormal') > 0"
                type="danger"
                size="small"
              >
                <i class="el-icon-warning"></i>
                {{ getTaskCountByStatus(scope.row, "abnormal") }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :key="col.name"
            :property="col.name"
            :label="col.label"
            :sortable="col.sortable"
            :align="col.align || 'left'"
            :width="col.width"
          >
          </el-table-column>
        </template>
        <el-table-column
          :label="'操作'"
          align="left"
          fixed="right"
          min-width="220px"
        >
          <template slot-scope="scope">
            <el-tooltip :content="'查看'" placement="top">
              <el-button
                type="primary"
                icon="el-icon-search"
                size="mini"
                :disabled="isDisabled(scope.row)"
                @click="onView(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip :content="'删除'" placement="top">
              <el-button
                type="danger"
                icon="el-icon-delete"
                size="mini"
                :disabled="isDisabled(scope.row)"
                @click="onRemove(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip
              v-if="!isShowRun(scope.row)"
              :content="'No command line'"
              placement="top"
            >
              <el-button
                disabled
                type="success"
                icon="fa fa-bug"
                size="mini"
                @click="onCrawl(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip v-else :content="'运行'" placement="top">
              <el-button
                type="success"
                icon="fa fa-bug"
                size="mini"
                :disabled="isDisabled(scope.row)"
                @click="onCrawl(scope.row, $event)"
              />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          @current-change="onPageNumChange"
          @size-change="onPageSizeChange"
          :current-page.sync="pagination.pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pagination.pageSize"
          layout="sizes, prev, pager, next"
          :total="spiderTotal"
        >
        </el-pagination>
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import dayjs from 'dayjs';
import CrawlConfirmDialog from '../../components/Common/CrawlConfirmDialog';
import StatusTag from '../../components/Status/StatusTag';
import StatusLegend from '../../components/Status/StatusLegend';
import CopySpiderDialog from '../../components/Spider/CopySpiderDialog';

export default {
  name: 'SpiderList',
  components: {
    CopySpiderDialog,
    StatusLegend,
    CrawlConfirmDialog,
    StatusTag
  },
  data () {
    return {
      pagination: {
        pageNum: 1,
        pageSize: 10
      },
      importLoading: false,
      addConfigurableLoading: false,
      isEditMode: false,
      dialogVisible: false,
      addDialogVisible: false,
      crawlConfirmDialogVisible: false,
      isRunningTasksDialogVisible: false,
      activeSpiderId: undefined,
      activeSpider: undefined,
      filter: {
        project_id: '',
        keyword: '',
        type: 'all',
        owner_type: 'me'
      },
      sort: {
        sortKey: '',
        sortDirection: null
      },
      types: [],
      spiderFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      },
      fileList: [],
      activeTabName: 'customized',
      handle: undefined,
      activeSpiderTaskStatus: 'running',
      selectedSpiders: [],
      isStopLoading: false,
      isRemoveLoading: false,
      isMultiple: false,
      copyDialogVisible: false
    };
  },
  computed: {
    ...mapState('spider', [
      'importForm',
      'spiderList',
      'spiderForm',
      'spiderTotal',
      'templateList'
    ]),
    ...mapGetters('user', ['userInfo', 'token']),
    ...mapState('lang', ['lang']),
    ...mapState('project', ['projectList']),
    ...mapState('node', ['nodeList']),
    uploadForm () {
      return {
        name: this.spiderForm.name,
        display_name: this.spiderForm.display_name,
        col: this.spiderForm.col,
        cmd: this.spiderForm.cmd
      };
    },
    columns () {
      const columns = [];
      columns.push({
        name: 'display_name',
        label: '名称',
        width: '160',
        align: 'left',
        sortable: true
      });
      columns.push({
        name: 'type',
        label: '爬虫类型',
        width: '120',
        sortable: true
      });
      columns.push({ name: 'is_scrapy', label: '是否是Scrapy', width: '80' });
      columns.push({ name: 'latest_tasks', label: '最近任务', width: '80' });
      columns.push({
        name: 'last_status',
        label: '上次运行状态',
        width: '120'
      });
      columns.push({ name: 'last_run_ts', label: '上次运行', width: '140' });
      columns.push({ name: 'update_ts', label: '更新时间', width: '140' });
      columns.push({ name: 'create_ts', label: '创建时间', width: '140' });
      return columns;
    },
    activeNodeList () {
      return this.nodeList.filter(d => {
        return d.status === 'online';
      });
    },
    activeSpiderIds () {
      return this.selectedSpiders.map(d => d._id);
    }
  },
  methods: {
    onSpiderTypeChange (val) {
      this.filter.type = val;
      this.getList();
    },
    onPageSizeChange (val) {
      this.pagination.pageSize = val;
      this.getList();
    },
    onPageNumChange (val) {
      this.pagination.pageNum = val;
      this.getList();
    },
    onSearch () {
      this.getList();
    },
    onAdd () {
      let projectId = '000000000000000000000000';
      if (this.filter.project_id) {
        projectId = this.filter.project_id;
      }
      this.$store.commit('spider/SET_SPIDER_FORM', {
        project_id: projectId,
        template: this.templateList[0]
      });
      this.addDialogVisible = true;
    },
    onAddConfigurable () {
      this.$refs['addConfigurableForm'].validate(async res => {
        if (!res) return;

        let res2;
        try {
          res2 = await this.$store.dispatch('spider/addConfigSpider');
        } catch (e) {
          this.$message.error('Something wrong happened');
          return;
        }
        this.$router.push(`/spiders/${res2.data.data._id}`);
        this.getList();
      });
    },
    onAddCustomized () {
      this.$refs['addCustomizedForm'].validate(async res => {
        if (!res) return;
        let res2;
        try {
          res2 = await this.$store.dispatch('spider/addSpider');
        } catch (e) {
          this.$message.error('Something wrong happened');
          return;
        }
        this.$router.push(`/spiders/${res2.data.data._id}`);
        this.getList();
      });
    },
    onRefresh () {
      this.getList();
    },
    onSubmit () {
      const vm = this;
      const formName = 'spiderForm';
      this.$refs[formName].validate(valid => {
        if (valid) {
          if (this.isEditMode) {
            vm.$store.dispatch('spider/editSpider');
          } else {
            vm.$store.dispatch('spider/addSpider');
          }
          vm.dialogVisible = false;
        } else {
          return false;
        }
      });
    },
    onCancel () {
      this.$store.commit('spider/SET_SPIDER_FORM', {});
      this.dialogVisible = false;
    },
    onDialogClose () {
      this.$store.commit('spider/SET_SPIDER_FORM', {});
      this.dialogVisible = false;
    },
    onAddDialogClose () {
      this.addDialogVisible = false;
    },
    onEdit (row) {
      this.isEditMode = true;
      this.$store.commit('spider/SET_SPIDER_FORM', row);
      this.dialogVisible = true;
    },
    onRemove (row, ev) {
      ev.stopPropagation();
      this.$confirm('你确定要删除该爬虫?', '提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await this.$store.dispatch('spider/deleteSpider', row._id);
        this.$message({
          type: 'success',
          message: '成功删除'
        });
        await this.getList();
      });
    },
    onCrawl (row, ev) {
      ev.stopPropagation();
      this.crawlConfirmDialogVisible = true;
      this.activeSpiderId = row._id;
    },
    onCrawlConfirm () {
      setTimeout(() => {
        this.getList();
      }, 1000);
    },
    onCopy (row, ev) {
      ev.stopPropagation();
      this.copyDialogVisible = true;
      this.activeSpiderId = row._id;
    },
    onCopyConfirm () {
      setTimeout(() => {
        this.getList();
      }, 1000);
    },
    onView (row, ev) {
      ev.stopPropagation();
      this.$router.push('/spiders/' + row._id);
    },
    onImport () {
      this.$refs.importForm.validate(valid => {
        if (valid) {
          this.importLoading = true;
          // TODO: switch between github / gitlab / svn
          this.$store
            .dispatch('spider/importGithub')
            .then(response => {
              this.$message.success('Import repo successfully');
              this.getList();
            })
            .catch(response => {
              this.$message.error(response.data.error);
            })
            .finally(() => {
              this.dialogVisible = false;
              this.importLoading = false;
            });
        }
      });
    },
    openImportDialog () {
      this.dialogVisible = true;
    },
    isShowRun (row) {
      if (!this.isCustomized(row)) return true;
      return !!row.cmd;
    },
    isCustomized (row) {
      return row.type === 'customized';
    },
    fetchSiteSuggestions (keyword, callback) {
      this.$request
        .get('/sites', {
          keyword: keyword,
          page_num: 1,
          page_size: 100
        })
        .then(response => {
          const data = response.data.items.map(d => {
            d.value = d.name + ' | ' + d.domain;
            return d;
          });
          callback(data);
        });
    },
    onUploadSuccess (res) {
      // clear fileList
      this.fileList = [];

      // fetch spider list
      setTimeout(() => {
        this.getList();
      }, 500);

      // message
      this.$message.success('成功上传爬虫文件');

      // navigate to spider detail
      this.$router.push(`/spiders/${res.data._id}`);
    },
    beforeUpload (file) {
      return new Promise((resolve, reject) => {
        this.$refs['addCustomizedForm'].validate(res => {
          if (res) {
            resolve();
          } else {
            reject(new Error('form validation error'));
          }
        });
      });
    },
    getTime (str) {
      if (!str || str.match('^0001')) return 'NA';
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss');
    },
    onRowClick (row, column, event) {
      this.onView(row, event);
    },
    onSortChange ({ column, prop, order }) {
      this.sort.sortKey = order ? prop : '';
      this.sort.sortDirection = order;
      this.getList();
    },
    onClickTab (tab) {
      this.filter.type = tab.name;
      this.getList();
    },
    async getList () {
      let params = {
        page_num: this.pagination.pageNum,
        page_size: this.pagination.pageSize,
        sort_key: this.sort.sortKey,
        sort_direction: this.sort.sortDirection,
        keyword: this.filter.keyword,
        type: this.filter.type,
        project_id: this.filter.project_id,
        owner_type: this.filter.owner_type
      };
      await this.$store.dispatch('spider/getSpiderList', params);

      // 更新当前爬虫（任务列表）
      this.updateActiveSpider();
    },
    getTasksByStatus (row, status) {
      if (!row.latest_tasks) return [];
      return row.latest_tasks.filter(d => d.status === status);
    },
    getTaskCountByStatus (row, status) {
      return this.getTasksByStatus(row, status).length;
    },
    updateActiveSpider () {
      if (this.activeSpider) {
        for (let i = 0; i < this.spiderList.length; i++) {
          const spider = this.spiderList[i];
          if (this.activeSpider._id === spider._id) {
            this.activeSpider = spider;
          }
        }
      }
    },
    onViewRunningTasks (row, ev) {
      ev.stopPropagation();
      this.activeSpider = row;
      this.isRunningTasksDialogVisible = true;
    },
    getTasksByNode (row) {
      if (!this.activeSpider.latest_tasks) {
        return [];
      }
      return this.activeSpider.latest_tasks
        .filter(
          d => d.node_id === row._id && d.status === this.activeSpiderTaskStatus
        )
        .map(d => {
          d = JSON.parse(JSON.stringify(d));
          d.create_ts = d.create_ts.match('^0001')
            ? 'NA'
            : dayjs(d.create_ts).format('YYYY-MM-DD HH:mm:ss');
          d.start_ts = d.start_ts.match('^0001')
            ? 'NA'
            : dayjs(d.start_ts).format('YYYY-MM-DD HH:mm:ss');
          d.finish_ts = d.finish_ts.match('^0001')
            ? 'NA'
            : dayjs(d.finish_ts).format('YYYY-MM-DD HH:mm:ss');
          return d;
        });
    },
    onViewTask (row) {
      this.$router.push(`/tasks/${row._id}`);
    },
    async onStop (row, ev) {
      ev.stopPropagation();
      const res = await this.$store.dispatch('task/cancelTask', row._id);
      if (!res.data.error) {
        this.$message.success(`Task "${row._id}" has been sent signal to stop`);
        this.getList();
      }
    },
    onIsScrapy (value) {
      if (value) {
        this.spiderForm.cmd = 'scrapy crawl';
      }
    },
    onSpiderSelect (spiders) {
      this.selectedSpiders = spiders;
    },
    async onRemoveSelectedSpiders () {
      this.$confirm('您是否确认删除所选项？', '提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        this.isRemoveLoading = true;
        try {
          const res = await this.$request.delete('/spiders', {
            spider_ids: this.selectedSpiders.map(d => d._id)
          });
          if (!res.data.error) {
            this.$message.success('成功删除');
            this.$refs['table'].clearSelection();
            await this.getList();
          }
        } finally {
          this.isRemoveLoading = false;
        }
      });
    },
    async onStopSelectedSpiders () {
      this.$confirm('您是否确认停止所选项？', '提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        this.isStopLoading = true;
        try {
          const res = await this.$request.post('/spiders-cancel', {
            spider_ids: this.selectedSpiders.map(d => d._id)
          });
          if (!res.data.error) {
            this.$message.success('已经向所选任务发送取消任务信号');
            await this.getList();
          }
        } finally {
          this.isStopLoading = false;
        }
      });
    },
    onCrawlSelectedSpiders () {
      this.crawlConfirmDialogVisible = true;
      this.isMultiple = true;
    },
    onCrawlConfirmDialogClose () {
      this.crawlConfirmDialogVisible = false;
      this.isMultiple = false;
    },
    isDisabled (row) {
      return (
        row.is_public &&
        row.username !== this.userInfo.username &&
        this.userInfo.role !== 'admin'
      );
    }
  },
  async created () {
    // fetch project list
    await this.$store.dispatch('project/getProjectList');

    // project id
    if (this.$route.params.project_id) {
      this.filter.project_id = this.$route.params.project_id;
    }

    // fetch node list
    await this.$store.dispatch('node/getNodeList');

    // fetch spider list
    await this.getList();

    // fetch template list
    await this.$store.dispatch('spider/getTemplateList');

    // periodically fetch spider list
    this.handle = setInterval(() => {
      this.getList();
    }, 15000);
  },
  mounted () {
    const vm = this;
    this.$nextTick(() => {
      vm.$store.commit('spider/SET_SPIDER_FORM', this.spiderForm);
    });
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

  .filter-search {
    width: 240px;
  }

  .right {
    .btn {
      margin-left: 10px;
    }
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

.add-spider-wrapper {
  display: flex;
  justify-content: center;

  .add-spider-item {
    cursor: pointer;
    width: 180px;
    font-size: 18px;
    height: 120px;
    margin: 0 20px;
    flex-basis: 40%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .add-spider-item.primary {
    color: #409eff;
    background: rgba(64, 158, 255, 0.1);
    border: 1px solid rgba(64, 158, 255, 0.1);
  }

  .add-spider-item.success {
    color: #67c23a;
    background: rgba(103, 194, 58, 0.1);
    border: 1px solid rgba(103, 194, 58, 0.1);
  }

  .add-spider-item.info {
    color: #909399;
    background: #f4f4f5;
    border: 1px solid #e9e9eb;
  }
}

.el-autocomplete {
  width: 100%;
}
</style>

<style scoped>
.el-table >>> tr {
  cursor: pointer;
}

.actions {
  text-align: right;
}

.el-table >>> .latest-tasks .el-tag {
  margin: 3px 3px 0 0;
}
</style>
