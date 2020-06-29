<template>
  <div class="info-view">
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="spiderForm._id"
      @close="crawlConfirmDialogVisible = false"
    />

    <el-row>
      <el-form label-width="150px"
               :model="spiderForm"
               ref="spiderForm"
               class="spider-form"
               label-position="right">
        <el-form-item :label="'爬虫ID'">
          <el-input v-model="spiderForm._id" :placeholder="'爬虫ID'" disabled></el-input>
        </el-form-item>
        <el-form-item :label="'爬虫名称'">
          <el-input v-model="spiderForm.display_name" :placeholder="'爬虫名称'" :disabled="isView || isPublic"/>
        </el-form-item>
        <el-form-item :label="'项目'" prop="project_id" required>
          <el-select
            v-model="spiderForm.project_id"
            :placeholder="'项目'"
            filterable
            :disabled="isView || isPublic"
          >
            <el-option
              v-for="p in projectList"
              :key="p._id"
              :value="p._id"
              :label="p.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="'代码目录'">
          <el-input v-model="spiderForm.src" :placeholder="'代码目录'" disabled></el-input>
        </el-form-item>
        <template v-if="spiderForm.type === 'customized'">
          <el-form-item :label="'执行命令'" prop="cmd" required :inline-message="true">
            <el-input
              v-model="spiderForm.cmd"
              :placeholder="'执行命令'"
              :disabled="isView || spiderForm.is_scrapy || isPublic"
            />
          </el-form-item>
        </template>
        <el-form-item :label="'结果集'" prop="col">
          <el-input
            v-model="spiderForm.col"
            :placeholder="'默认: ' + 'results_<spider_name>'"
            :disabled="isView || isPublic"
          />
        </el-form-item>
        <el-row>
        <el-col :span="6">
          <el-form-item v-if="spiderForm.type === 'customized' && !isView" :label="'是否为 Scrapy'" prop="is_scrapy">
            <el-switch
              v-model="spiderForm.is_scrapy"
              active-color="#13ce66"
              :disabled="isView || isPublic"
              @change="onIsScrapyChange"
            />
          </el-form-item>
        </el-col>
        </el-row>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button size="small" v-if="isShowRun && !isPublic" type="danger" @click="onCrawl"
                 icon="el-icon-video-play" style="margin-right: 10px">
        {{'运行'}}
      </el-button>
      <el-upload
        v-if="spiderForm.type === 'customized'"
        :action="$request.baseUrl + `/spiders/${spiderForm._id}/upload`"
        :headers="{Authorization:token}"
        :on-progress="() => this.uploadLoading = true"
        :on-error="onUploadError"
        :on-success="onUploadSuccess"
        :file-list="fileList"
        style="display:inline-block;margin-right:10px"
      >
        <el-button v-if="!isPublic" size="small" type="primary" icon="el-icon-upload" v-loading="uploadLoading">
          {{'上传'}}
        </el-button>
      </el-upload>
      <el-button v-if="!isPublic" size="small" type="success" @click="onSave" icon="el-icon-check">
        {{'保存'}}
      </el-button>
    </el-row>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import CrawlConfirmDialog from '../Common/CrawlConfirmDialog'

export default {
  name: 'SpiderInfoView',
  components: { CrawlConfirmDialog },
  props: {
    isView: {
      default: false,
      type: Boolean
    }
  },
  data () {
    const cronValidator = (rule, value, callback) => {
      let patArr = []
      for (let i = 0; i < 6; i++) {
        patArr.push('[/*,0-9]+')
      }
      const pat = '^' + patArr.join(' ') + '$'
      if (this.spiderForm.cron_enabled) {
        if (!value) {
          callback(new Error('cron cannot be empty'))
        } else if (!value.match(pat)) {
          callback(new Error('cron format is invalid'))
        }
      }
      callback()
    }
    const dedupValidator = (rule, value, callback) => {
      if (!this.spiderForm.is_dedup) {
        return callback()
      } else {
        if (value) {
          return callback()
        } else {
          return callback(new Error('dedup field cannot be empty'))
        }
      }
    }
    return {
      uploadLoading: false,
      fileList: [],
      crawlConfirmDialogVisible: false,
      cmdRule: [
        { message: 'Execute Command should not be empty', required: true }
      ],
      cronRules: [
        { validator: cronValidator, trigger: 'blur' }
      ],
      dedupRules: [
        { validator: dedupValidator, trigger: 'blur' }
      ]
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    ...mapGetters('user', [
      'userInfo',
      'token'
    ]),
    ...mapState('project', [
      'projectList'
    ]),
    isConfigurable () {
      return this.spiderForm.type === 'configurable'
    },
    isShowRun () {
      if (this.spiderForm.type === 'customized') {
        return !!this.spiderForm.cmd
      } else {
        return true
      }
    },
    isPublic () {
      return this.spiderForm.is_public && this.spiderForm.username !== this.userInfo.username && this.userInfo.role !== 'admin'
    }
  },
  methods: {
    onCrawl () {
      this.crawlConfirmDialogVisible = true
    },
    onSave () {
      this.$refs['spiderForm'].validate(async valid => {
        if (!valid) return
        const res = await this.$store.dispatch('spider/editSpider')
        if (!res.data.error) {
          this.$message.success('Spider info has been saved successfully')
        }
        await this.$store.dispatch('spider/getSpiderData', this.$route.params.id)
        if (this.spiderForm.is_scrapy) {
          await this.$store.dispatch('spider/getSpiderScrapySpiders', this.$route.params.id)
        }
      })
    },
    fetchSiteSuggestions (keyword, callback) {
      this.$request.get('/sites', {
        keyword: keyword,
        page_num: 1,
        page_size: 100
      }).then(response => {
        const data = response.data.items.map(d => {
          d.value = `${d.name} | ${d.domain}`
          return d
        })
        callback(data)
      })
    },
    onSiteSelect (item) {
      this.spiderForm.site = item._id
    },
    onUploadSuccess () {
      this.$store.dispatch('spider/getFileTree')

      this.uploadLoading = false

      this.$message.success('Uploaded spider files successfully')
    },
    onUploadError () {
      this.uploadLoading = false
    },
    onIsScrapyChange (value) {
      if (value) {
        this.spiderForm.cmd = 'scrapy crawl'
      }
    },
    onIsDedupChange (value) {
      if (value && !this.spiderForm.dedup_method) {
        this.spiderForm.dedup_method = 'overwrite'
      }
    }
  },
  async created () {
    // fetch project list
    await this.$store.dispatch('project/getProjectList')

    // 兼容项目ID
    if (!this.spiderForm.project_id) {
      this.$set(this.spiderForm, 'project_id', '000000000000000000000000')
    }
  }
}
</script>

<style scoped>
  .spider-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }

  .el-autocomplete {
    width: 100%;
  }

  .info-view >>> .el-upload-list {
    display: none;
  }
</style>
