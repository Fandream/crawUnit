<template>
  <div class="config-list">

    <!--preview results-->
    <el-dialog :visible.sync="dialogVisible"
               :title="'预览结果'"
               width="90%"
               :before-close="onDialogClose">
      <el-table class="table-header" :data="[{}]" :show-header="false">
        <el-table-column v-for="(f, index) in fields"
                         :key="f.name + '-' + index"
                         min-width="100px">
          <template>
            <el-input v-model="columnsDict[f.name]" size="mini" style="width: calc(100% - 15px)"></el-input>
            <a href="javascript:" style="margin-left: 2px;" @click="onDeleteField(index)">X</a>
            <!--<el-button size="mini" type="danger" icon="el-icon-delete" style="width:45px;margin-left:2px"></el-button>-->
          </template>
        </el-table-column>
      </el-table>
      <el-table :data="previewCrawlData"
                :show-header="false"
                border>
        <el-table-column v-for="(f, index) in fields"
                         :key="f.name + '-' + index"
                         :label="f.name"
                         min-width="100px">

          <template slot-scope="scope">
            {{getDisplayStr(scope.row[f.name])}}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
    <!--./preview results-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="spiderForm._id"
      @close="crawlConfirmDialogVisible = false"
    />
    <!--./crawl confirm dialog-->

    <!--tabs-->
    <el-tabs :active-name="activeTab" @tab-click="onTabClick">
      <!--Stages-->
      <el-tab-pane name="stages" :label="'Stages'">
        <!--config detail-->
        <el-row>
          <el-form label-width="150px" ref="form" :model="spiderForm.config">
          </el-form>
        </el-row>
        <!--./config detail-->

        <el-row>
          <div class="top-wrapper">
            <ul class="list">
              <li class="item">
                <label>{{'开始URL'}}: </label>
                <el-input
                  id="start-url"
                  v-model="spiderForm.config.start_url"
                  :placeholder="'开始URL'"
                  :class="startUrlClass"
                />
              </li>
              <li class="item">
                <label>{{'开始阶段'}}: </label>
                <el-select
                  id="start-stage"
                  v-model="spiderForm.config.start_stage"
                  :placeholder="'开始阶段'"
                  :class="startStageClass"
                >
                  <el-option
                    v-for="n in spiderForm.config.stages.map(s => s.name)"
                    :key="n"
                    :value="n"
                    :label="n"
                  />
                </el-select>
              </li>
              <li class="item">
                <label>{{'引擎'}}: </label>
                <el-select
                  v-model="spiderForm.config.engine"
                  :placeholder="'开始阶段'"
                  :class="startStageClass"
                  disabled
                >
                  <el-option
                    v-for="n in ['scrapy']"
                    :key="n"
                    :value="n"
                    :label="n"
                  />
                </el-select>
              </li>
              <li class="item">
                <label>{{'选择器类别'}}: </label>
                <div class="selector-type">
              <span class="selector-type-item" @click="onClickSelectorType('css')">
                <el-tag
                  :class="isCss ? 'active' : 'inactive'"
                  type="success"
                >
                  CSS
                </el-tag>
              </span>
                  <span class="selector-type-item" @click="onClickSelectorType('xpath')">
              <el-tag
                :class="isXpath ? 'active' : 'inactive'"
                type="primary"
              >
                XPath
              </el-tag>
            </span>
                </div>
              </li>
            </ul>
          </div>

          <div class="button-group-container">
            <div class="button-group">
              <el-button
                id="btn-run"
                size="small"
                type="danger"
                :disabled="isDisabled"
                icon="el-icon-video-play"
                @click="onCrawl"
              >
                {{'运行'}}
              </el-button>
              <el-button
                id="btn-convert"
                size="small"
                type="warning"
                :disabled="isDisabled"
                icon="el-icon-refresh-right"
                @click="onConvert"
              >
                {{'转化为自定义'}}
              </el-button>
              <!--              <el-button type="primary" @click="onExtractFields" v-loading="extractFieldsLoading">-->
              <!--                {{'ExtractFields'}}-->
              <!--              </el-button>-->
              <!--              <el-button type="warning" @click="onPreview" v-loading="previewLoading">{{'Preview'}}</el-button>-->
              <el-button
                id="btn-save"
                size="small"
                type="success"
                :disabled="saveLoading || isDisabled"
                @click="onSave"
                :icon="saveLoading ? 'el-icon-loading' : 'el-icon-check'"
              >
                {{'保存'}}
              </el-button>
            </div>
          </div>
        </el-row>

        <el-collapse
          :value="activeNames"
        >
          <el-collapse-item
            v-for="(stage, index) in spiderForm.config.stages"
            :key="index"
            :name="stage.name"
          >
            <template slot="title">
              <ul class="stage-list">
                <!--actions-->
                <li class="stage-item actions" style="min-width: 80px; flex-basis: 80px; justify-content: flex-end"
                    @click="$event.stopPropagation()">
                  <i class="action-item el-icon-copy-document" @click="onCopyStage(stage)"></i>
                  <i class="action-item el-icon-remove-outline" @click="onRemoveStage(stage)"></i>
                  <i class="action-item el-icon-circle-plus-outline" @click="onAddStage(stage)"></i>
                </li>
                <!--./actions-->

                <!--stage name-->
                <li class="stage-item" style="min-width: 240px" @click="$event.stopPropagation()">
                  <label>{{'阶段名称'}}: </label>
                  <div v-if="!stage.isEdit" @click="onEditStage(stage)" class="text-wrapper">
                    <span class="text">
                      {{stage.name}}
                    </span>
                    <i class="el-icon-edit-outline"></i>
                  </div>
                  <el-input
                    v-else
                    :ref="`stage-name-${stage.name}`"
                    class="edit-text"
                    v-model="stage.name"
                    :placeholder="'阶段名称'"
                    @focus="onStageNameFocus($event)"
                    @blur="stage.isEdit = false"
                  />
                </li>
                <!--./stage name-->

                <!--list-->
                <li class="stage-item" style="min-width: 240px">
                  <label>{{'列表'}}: </label>
                  <el-checkbox
                    style="text-align: left; flex-basis: 20px; margin-right: 5px"
                    :value="isList(stage)"
                    @change="onCheckIsList($event, stage)"
                  />
                  <el-popover v-model="stage.isListOpen" v-if="isList(stage)" placement="top" width="360">
                    <el-form label-width="120px">
                      <el-form-item :label="'选择器类别'">
                        <el-tag :class="!stage.list_xpath ? 'active' : 'inactive'" type="success"
                                @click="onSelectStageListType(stage, 'css')">CSS
                        </el-tag>
                        <el-tag :class="stage.list_xpath ? 'active' : 'inactive'" type="primary"
                                @click="onSelectStageListType(stage, 'xpath')">XPath
                        </el-tag>
                      </el-form-item>
                      <el-form-item :label="'选择器'" class="list-selector">
                        <el-input v-if="!stage.list_xpath" v-model="stage.list_css"/>
                        <el-input v-else v-model="stage.list_xpath"/>
                      </el-form-item>
                    </el-form>
                    <el-tag
                      v-if="!stage.list_xpath"
                      type="success"
                      slot="reference"
                      @click="onClickStageList($event, stage, 'css')"
                    >
                      <i v-if="!stage.isListOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      CSS
                    </el-tag>
                    <el-tag
                      v-else
                      type="primary"
                      slot="reference"
                      @click="onClickStageList($event, stage, 'xpath')"
                    >
                      <i v-if="!stage.isListOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      XPath
                    </el-tag>
                  </el-popover>
                </li>
                <!--./list-->

                <!--pagination-->
                <li class="stage-item" style="min-width: 240px">
                  <label>{{'分页'}}: </label>
                  <el-checkbox
                    style="text-align: left; flex-basis: 20px; margin-right: 5px"
                    :value="isPage(stage)"
                    @change="onCheckIsPage($event, stage)"
                    :disabled="!isList(stage)"
                  />
                  <el-popover v-model="stage.isPageOpen" v-if="isPage(stage)" placement="top" width="360">
                    <el-form label-width="120px">
                      <el-form-item :label="'Selector Type'">
                        <el-tag :class="!stage.page_xpath ? 'active' : 'inactive'" type="success"
                                @click="onSelectStagePageType(stage, 'css')">CSS
                        </el-tag>
                        <el-tag :class="stage.page_xpath ? 'active' : 'inactive'" type="primary"
                                @click="onSelectStagePageType(stage, 'xpath')">XPath
                        </el-tag>
                      </el-form-item>
                      <el-form-item :label="'Selector'" class="page-selector">
                        <el-input v-if="!stage.page_xpath" v-model="stage.page_css"/>
                        <el-input v-else v-model="stage.page_xpath"/>
                      </el-form-item>
                    </el-form>
                    <el-tag
                      v-if="!stage.page_xpath"
                      type="success"
                      slot="reference"
                      @click="onClickStagePage($event, stage, 'css')"
                    >
                      <i v-if="!stage.isPageOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      CSS
                    </el-tag>
                    <el-tag
                      v-else
                      type="primary"
                      slot="reference"
                      @click="onClickStagePage($event, stage, 'xpath')"
                    >
                      <i v-if="!stage.isPageOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      XPath
                    </el-tag>
                  </el-popover>
                </li>
                <!--./pagination-->

              </ul>
            </template>
            <fields-table-view
              type="list"
              title="列表页字段"
              :fields="stage.fields"
              :stage="stage"
              :stage-names="spiderForm.config.stages.map(s => s.name)"
            />
          </el-collapse-item>
        </el-collapse>
      </el-tab-pane>
      <!--./Stages-->

      <!--Graph-->
      <el-tab-pane name="process" :label="'Process'">
        <div id="process-chart"></div>
      </el-tab-pane>
      <!--./Graph-->

      <!--Setting-->
      <el-tab-pane name="settings" :label="'Settings'">
        <div class="actions" style="text-align: right;margin-bottom: 10px">
          <el-button type="success" size="small" :disabled="isDisabled" @click="onSave">
            {{'保存'}}
          </el-button>
        </div>
        <setting-fields-table-view
          type="list"
        />
      </el-tab-pane>
      <!--./Setting-->

      <!--Spiderfile-->
      <el-tab-pane name="spiderfile" label="Spiderfile">
        <div class="spiderfile-actions">
          <el-button
            type="primary"
            size="small"
            style="margin-right: 10px;"
            :disabled="isDisabled"
            @click="onSpiderfileSave"
          >
            <font-awesome-icon :icon="['fa', 'save']"/>
            {{'保存'}}
          </el-button>
        </div>
        <file-detail/>
      </el-tab-pane>
      <!--./Spiderfile-->
    </el-tabs>
    <!--./tabs-->
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import echarts from 'echarts'
import FieldsTableView from '../TableView/FieldsTableView'
import CrawlConfirmDialog from '../Common/CrawlConfirmDialog'

import 'codemirror/lib/codemirror.js'
import 'codemirror/mode/yaml/yaml.js'
import FileDetail from '../File/FileDetail'
import SettingFieldsTableView from '../TableView/SettingFieldsTableView'

export default {
  name: 'ConfigList',
  components: {
    SettingFieldsTableView,
    FileDetail,
    CrawlConfirmDialog,
    FieldsTableView
  },
  watch: {
    activeTab () {
      setTimeout(() => {
        // 渲染流程图
        this.renderProcessChart()

        // 获取Spiderfile
        this.getSpiderfile()

        // 获取config
        this.$store.dispatch('spider/getSpiderData', this.spiderForm._id)
      }, 0)
    }
  },
  data () {
    return {
      crawlTypeList: [
        { value: 'list', label: '仅列表' },
        { value: 'detail', label: '仅详情页' },
        { value: 'list-detail', label: '列表＋详情页' }
      ],
      extractFieldsLoading: false,
      previewLoading: false,
      saveLoading: false,
      dialogVisible: false,
      crawlConfirmDialogVisible: false,
      columnsDict: {},
      fieldColumns: [
        { name: 'name', label: '名称' },
        { name: 'selector_type', label: '选择器类型' },
        { name: 'selector', label: '选择器' },
        { name: 'is_attr', label: '是否为属性' },
        { name: 'attr', label: '属性' },
        { name: 'next_stage', label: '下一阶段' }
      ],
      activeTab: 'stages',
      processChart: undefined,
      fileOptions: {
        mode: 'text/x-yaml',
        theme: 'darcula',
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        matchBrackets: true
      }
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm',
      'previewCrawlData'
    ]),
    ...mapGetters('user', [
      'userInfo'
    ]),
    isDisabled () {
      return this.spiderForm.is_public && this.spiderForm.username !== this.userInfo.username && this.userInfo.role !== 'admin'
    },
    fields () {
      if (this.spiderForm.crawl_type === 'list') {
        return this.spiderForm.fields
      } else if (this.spiderForm.crawl_type === 'detail') {
        return this.spiderForm.detail_fields
      } else if (this.spiderForm.crawl_type === 'list-detail') {
        return this.spiderForm.fields.concat(this.spiderForm.detail_fields)
      } else {
        return []
      }
    },
    isCss () {
      let i = 0
      this.spiderForm.config.stages.forEach(stage => {
        stage.fields.forEach(field => {
          if (!field.css) i++
        })
      })
      return i === 0
    },
    isXpath () {
      let i = 0
      this.spiderForm.config.stages.forEach(stage => {
        stage.fields.forEach(field => {
          if (!field.xpath) i++
        })
      })
      return i === 0
    },
    activeNames () {
      return this.spiderForm.config.stages.map(d => d.name)
    },
    startUrlClass () {
      if (!this.spiderForm.config.start_url) {
        return 'invalid'
      } else if (!this.spiderForm.config.start_url.match(/^https?:\/\/.+|^\/\/.+/i)) {
        return 'invalid'
      }

      return ''
    },
    startStageClass () {
      if (!this.spiderForm.config.start_stage) {
        return 'invalid'
      } else if (!this.activeNames.includes(this.spiderForm.config.start_stage)) {
        return 'invalid'
      }
      return ''
    },
    stageNodes () {
      const startStage = this.spiderForm.config.stages[this.spiderForm.config.start_stage]
      const nodes = []
      const allStageNames = new Set()

      let i = 0
      let currentStage = startStage
      while (currentStage) {
        // 加入节点信息
        nodes.push({
          x: i++,
          y: 0,
          itemStyle: {
            color: '#409EFF'
          },
          ...currentStage
        })

        // 记录该节点
        allStageNames.add(currentStage.name)

        // 设置当前阶段为下一阶段
        currentStage = this.getNextStage(currentStage)

        if (currentStage && allStageNames.has(currentStage.name)) {
          currentStage = undefined
        }
      }

      // 加入剩余阶段
      i = 0
      const restStages = this.spiderForm.config.stages
        .filter(stage => !allStageNames.has(stage.name))
      restStages.forEach(stage => {
        // 加入节点信息
        nodes.push({
          x: i++,
          y: 1,
          itemStyle: {
            color: '#F56C6C'
          },
          ...stage
        })
      })

      return nodes
    },
    stageEdges () {
      const edges = []
      const stages = this.spiderForm.config.stages
      stages.forEach(stage => {
        for (let i = 0; i < stage.fields.length; i++) {
          const field = stage.fields[i]
          if (field.next_stage) {
            edges.push({
              source: stage.name,
              target: field.next_stage
            })
          }
        }
      })
      return edges
    }
  },
  methods: {
    onSelectCrawlType (value) {
      this.spiderForm.crawl_type = value
    },
    async onSave () {
      this.saveLoading = true
      try {
        const res = await this.$store.dispatch('spider/postConfigSpiderConfig')
        if (!res.data.error) {
          this.$message.success('Spider info has been saved successfully')
          return true
        }
        return false
      } catch (e) {
        return false
      } finally {
        this.saveLoading = false
      }
    },
    onDialogClose () {
      this.dialogVisible = false
      this.fields.forEach(f => {
        f.name = this.columnsDict[f.name]
      })
    },
    onPreview () {
      this.$refs['form'].validate(res => {
        if (res) {
          this.onSave()
            .then(() => {
              this.previewLoading = true
              this.$store.dispatch('spider/getPreviewCrawlData')
                .then(() => {
                  this.fields.forEach(f => {
                    this.columnsDict[f.name] = f.name
                  })
                  this.dialogVisible = true
                })
                .catch(() => {
                  this.$message.error('Something wrong happened')
                })
                .finally(() => {
                  this.previewLoading = false
                })
            })
        }
      })
    },
    async onCrawl () {
      const res = await this.onSave()
      if (res) {
        this.crawlConfirmDialogVisible = true
      }
    },
    onExtractFields () {
      this.$refs['form'].validate(res => {
        if (res) {
          this.onSave()
            .then(() => {
              this.extractFieldsLoading = true
              this.$store.dispatch('spider/extractFields')
                .then(response => {
                  if (response.data.item_selector) {
                    this.$set(this.spiderForm, 'item_selector', response.data.item_selector)
                  }
                  if (response.data.item_selector_type) {
                    this.$set(this.spiderForm, 'item_selector_type', response.data.item_selector_type)
                  }

                  if (response.data.fields && response.data.fields.length) {
                    this.spiderForm.fields = response.data.fields
                  }

                  if (response.data.pagination_selector) {
                    this.spiderForm.pagination_selector = response.data.pagination_selector
                  }
                })
                .finally(() => {
                  this.extractFieldsLoading = false
                })
            })
        }
      })
    },
    onDeleteField (index) {
      this.fields.splice(index, 1)
    },
    getDisplayStr (value) {
      if (!value) return value
      value = value.trim()
      if (value.length > 20) return value.substr(0, 20) + '...'
      return value
    },
    onClickSelectorType (selectorType) {
      this.spiderForm.config.stages.forEach(stage => {
        // 列表
        if (selectorType === 'css') {
          if (stage.list_xpath) stage.list_xpath = ''
          if (!stage.list_css) stage.list_css = 'body'
        } else {
          if (stage.list_css) stage.list_css = ''
          if (!stage.list_xpath) stage.list_xpath = '//body'
        }

        // 分页
        if (selectorType === 'css') {
          if (stage.page_xpath) stage.page_xpath = ''
          if (!stage.page_css) stage.page_css = 'body'
        } else {
          if (stage.page_css) stage.page_css = ''
          if (!stage.page_xpath) stage.page_xpath = '//body'
        }

        // 字段
        stage.fields.forEach(field => {
          if (selectorType === 'css') {
            if (field.xpath) field.xpath = ''
            if (!field.css) field.css = 'body'
          } else {
            if (field.css) field.css = ''
            if (!field.xpath) field.xpath = '//body'
          }
        })
      })
    },
    onStageNameFocus (ev) {
      ev.stopPropagation()
    },
    onEditStage (stage) {
      this.$set(stage, 'isEdit', !stage.isEdit)
      setTimeout(() => {
        this.$refs[`stage-name-${stage.name}`][0].focus()
      }, 0)
    },
    onCopyStage (stage) {
      const stages = this.spiderForm.config.stages
      const ts = Math.floor(new Date().getTime()).toString()
      const newStage = JSON.parse(JSON.stringify(stage))
      newStage.name = `${stage.name}_copy_${ts}`
      for (let i = 0; i < stages.length; i++) {
        if (stage.name === stages[i].name) {
          stages.splice(i + 1, 0, newStage)
        }
      }
    },
    addStage (index) {
      const stages = this.spiderForm.config.stages
      const ts = Math.floor(new Date().getTime()).toString()
      const newStageName = `stage_${ts}`
      const newField = { name: `field_${ts}`, next_stage: '' }
      if (this.isCss) {
        newField['css'] = 'body'
      } else if (this.isXpath) {
        newField['xpath'] = '//body'
      } else {
        newField['xpath'] = '//body'
      }
      stages.splice(index + 1, 0, {
        name: newStageName,
        list_css: this.isCss ? 'body' : '',
        list_xpath: this.isXpath ? '//body' : '',
        page_css: '',
        page_xpath: '',
        fields: [newField]
      })
    },
    onRemoveStage (stage) {
      const stages = this.spiderForm.config.stages
      for (let i = 0; i < stages.length; i++) {
        if (stage.name === stages[i].name) {
          stages.splice(i, 1)
          break
        }
      }
      // 如果只剩一个stage，加入新的stage
      if (stages.length === 0) {
        this.addStage(0)
      }
      // 重置next_stage被设置为该stage的field
      stages.forEach(_stage => {
        _stage.fields.forEach(field => {
          if (field.next_stage === stage.name) {
            this.$set(field, 'next_stage', '')
          }
        })
      })
    },
    onAddStage (stage) {
      const stages = this.spiderForm.config.stages
      for (let i = 0; i < stages.length; i++) {
        if (stage.name === stages[i].name) {
          this.addStage(i)
          break
        }
      }
    },
    renderProcessChart () {
      const option = {
        title: {
          text: 'Stage Process'
        },
        series: [
          {
            animation: false,
            type: 'graph',
            // layout: 'force',
            layout: 'none',
            symbolSize: 50,
            roam: true,
            label: {
              normal: {
                show: true
              }
            },
            edgeSymbol: ['circle', 'arrow'],
            edgeSymbolSize: [4, 10],
            edgeLabel: {
              normal: {
                textStyle: {
                  fontSize: 20
                }
              }
            },
            focusOneNodeAdjacency: true,
            force: {
              initLayout: 'force',
              repulsion: 100,
              gravity: 0.00001,
              edgeLength: 200
            },
            // draggable: true,
            data: this.stageNodes,
            links: this.stageEdges,
            lineStyle: {
              normal: {
                opacity: 0.9,
                width: 2,
                curveness: 0
              }
            }
          }
        ],
        tooltip: {
          // formatter: '{b0}: {c0}<br />{b1}: {c1}',
          formatter: (params) => {
            if (!params.data.fields) return

            let str = ''
            str += `<label>${'Stage'}: ${params.name}</label><br>`
            str += `<ul style="list-style: none; padding: 0; margin: 0;">`
            // 列表
            if (params.data.list_css || params.data.list_xpath) {
              str += `<li><span style="display: inline-block;min-width: 50px;font-weight: bolder;text-align: right;margin-right: 3px">${'List'}: </span>${params.data.list_css || params.data.list_xpath}</li>`
            }
            if (params.data.page_css || params.data.page_xpath) {
              str += `<li><span style="display: inline-block;min-width: 50px;font-weight: bolder;text-align: right;margin-right: 3px">${'Pagination'}: </span>${params.data.page_css || params.data.page_xpath}</li>`
            }
            str += `</ul><br>`

            // 字段
            str += `<label>${'Fields'}: </label><br>`
            str += '<ul style="list-style: none; padding: 0; margin: 0;">'
            for (let i = 0; i < params.data.fields.length; i++) {
              const f = params.data.fields[i]
              str += `
<li>
<span style="display: inline-block; min-width: 50px; font-weight: bolder; text-align: right">${f.name}: </span>
${f.css || f.xpath} ${f.attr ? ('(' + f.attr + ')') : ''} ${f.next_stage ? (' --> ' + '<span style="font-weight:bolder">' + f.next_stage + '</span>') : ''}
</li>
`
            }
            str += '</ul>'
            return str
          }
        }
      }
      const el = document.querySelector('#process-chart')
      this.processChart = echarts.init(el)
      this.processChart.setOption(option)
      this.processChart.resize()
    },
    onTabClick (tab) {
      this.activeTab = tab.name
    },
    update () {
      if (this.activeTab !== 'stages') return

      // 手动显示tab下划线
      const elBar = document.querySelector('.el-tabs__active-bar')
      const elStages = document.querySelector('#tab-stages')
      const totalWidth = Number(getComputedStyle(elStages).width.replace('px', ''))
      const paddingRight = Number(getComputedStyle(elStages).paddingRight.replace('px', ''))
      elBar.setAttribute('style', 'width:' + (totalWidth - paddingRight) + 'px')
    },
    getSpiderfile () {
      this.$store.commit('file/SET_FILE_CONTENT', '')
      this.$store.commit('file/SET_CURRENT_PATH', 'Spiderfile')
      this.$store.dispatch('file/getFileContent', { path: 'Spiderfile' })
    },
    async onSpiderfileSave () {
      try {
        await this.$store.dispatch('spider/saveConfigSpiderSpiderfile')
        this.$message.success('Spiderfile saved successfully')
      } catch (e) {
        this.$message.error('Something wrong happened')
      }
    },
    isList (stage) {
      return !!stage.is_list
    },
    onCheckIsList (value, stage) {
      stage.is_list = value
      if (value) {

      } else {

      }
    },
    onClickStageList ($event, stage, type) {
      $event.stopPropagation()
    },
    onSelectStageListType (stage, type) {
      if (type === 'css') {
        if (!stage.list_css) stage.list_css = 'body'
        stage.list_xpath = ''
      } else {
        if (!stage.list_xpath) stage.list_xpath = '//body'
        stage.list_css = ''
      }
    },
    isPage (stage) {
      return !!stage.page_css || !!stage.page_xpath
    },
    onCheckIsPage (value, stage) {
      if (value) {
        if (!stage.page_css && !stage.page_xpath) {
          stage.page_xpath = '//body'
        }
      } else {
        stage.page_css = ''
        stage.page_xpath = ''
      }
    },
    onClickStagePage ($event, stage, type) {
      $event.stopPropagation()
    },
    onSelectStagePageType (stage, type) {
      if (type === 'css') {
        if (!stage.page_css) stage.page_css = 'body'
        stage.page_xpath = ''
      } else {
        if (!stage.page_xpath) stage.page_xpath = '//body'
        stage.page_css = ''
      }
    },
    getNextStageField (stage) {
      return stage.fields.filter(f => !!f.next_stage)[0]
    },
    getNextStage (stage) {
      const nextStageField = this.getNextStageField(stage)
      if (!nextStageField) return
      return this.spiderForm.config.stages[nextStageField.next_stage]
    },
    onConvert () {
      this.$confirm('Are you sure to convert this spider to customized spider?', 'Notification', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(async () => {
        this.spiderForm.type = 'customized'
        this.$store.dispatch('spider/editSpider')
          .then(res => {
            if (!res.data.error) {
              this.$store.commit('spider/SET_CONFIG_LIST_TS', +new Date())
              this.$message({
                type: 'success',
                message: 'Converted successfully'
              })
            } else {
              this.$message({
                type: 'error',
                message: 'Converted unsuccessfully'
              })
            }
            this.$store.dispatch('spider/getSpiderData', this.spiderForm._id)
          })
      })
    }
  }
}
</script>

<style scoped>

  .button-group-container {
    margin-top: 10px;
    /*border-bottom: 1px dashed #dcdfe6;*/
    padding-bottom: 20px;
  }

  .button-group {
    text-align: right;
  }

  .list-fields-container {
    margin-top: 20px;
    /*border-bottom: 1px dashed #dcdfe6;*/
    padding-bottom: 20px;
  }

  .detail-fields-container {
    margin-top: 20px;
  }

  .title {
    color: #606266;
    font-size: 14px;
  }

  .el-table.table-header >>> td {
    padding: 0;
  }

  .el-table.table-header >>> .cell {
    padding: 0;
  }

  .el-table.table-header >>> .el-input .el-input__inner {
    border-radius: 0;
  }

  .selector-type-item {
    margin: 0 5px;
    cursor: pointer;
    font-weight: bolder;
  }

  .el-tag {
    margin-right: 5px;
    font-weight: bolder;
    cursor: pointer;
  }

  .el-tag.inactive {
    opacity: 0.5;
  }

  .stage-list {
    width: 100%;
    /*width: calc(80px + 320px);*/
    display: flex;
    flex-wrap: wrap;
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .stage-list .stage-item {
    /*flex-basis: 320px;*/
    min-width: 120px;
    display: flex;
    align-items: center;
  }

  .stage-list .stage-item label {
    flex-basis: 90px;
    margin-right: 10px;
    justify-self: flex-end;
    text-align: right;
  }

  .stage-list .stage-item .el-input {
    flex-basis: calc(100% - 90px);
    height: 32px;
  }

  .stage-list .stage-item .el-input .el-input__inner {
    height: 32px;
    inline-size: 32px;
  }

  .stage-list .stage-item .action-item {
    cursor: pointer;
    width: 13px;
    margin-right: 5px;
  }

  .stage-list .stage-item .action-item:last-child {
    margin-right: 10px;
  }

  .stage-list .stage-item .text-wrapper {
    display: flex;
    align-items: center;
    max-width: calc(100% - 90px - 10px);
  }

  .stage-list .stage-item .text-wrapper .text {
    text-overflow: ellipsis;
    overflow: hidden;
  }

  .stage-list .stage-item .text-wrapper .text:hover {
    text-decoration: underline;
  }

  .stage-list .stage-item .text-wrapper i {
    margin-left: 5px;
  }

  .stage-list .stage-item >>> .edit-text {
    height: 32px;
    line-height: 32px;
  }

  .stage-list .stage-item >>> .edit-text .el-input__inner {
    height: 32px;
    line-height: 32px;
  }

  .top-wrapper {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .top-wrapper .list {
    list-style: none;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    padding: 0;
  }

  .top-wrapper .list .item {
    margin-bottom: 10px;
    display: flex;
    align-items: center;
  }

  .top-wrapper .list .item label {
    width: 100px;
    text-align: right;
    margin-right: 10px;
    font-size: 12px;
  }

  .top-wrapper .list .item label + * {
    width: 240px;
  }

  .invalid >>> .el-input__inner {
    border: 1px solid red !important;
  }

  #process-chart {
    width: 100%;
    height: 480px;
  }

  .config-list >>> .file-content {
    height: calc(100vh - 280px);
  }

  .spiderfile-actions {
    margin-bottom: 5px;
    text-align: right;
  }
</style>
