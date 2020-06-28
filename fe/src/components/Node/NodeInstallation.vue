<template>
  <div class="node-installation">
    <el-form class="search-form" inline>
      <el-form-item>
        <el-autocomplete
          class="search-box"
          size="small"
          clearable
          v-if="activeLang.executable_name === 'python'"
          v-model="depName"
          style="width: 240px"
          :placeholder="'搜索依赖'"
          :fetchSuggestions="fetchAllDepList"
          :minlength="2"
          @select="onSearch"
          @clear="onSearch"
        />
        <el-input
          v-else
          v-model="depName"
          style="width: 240px"
          :placeholder="'搜索依赖'"
        />
      </el-form-item>
      <el-form-item>
        <el-button size="small"
                   icon="el-icon-search"
                   type="success"
                   @click="onSearch"
        >
          {{'搜索'}}
        </el-button>
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="isShowInstalled" :label="'查看已安装'" @change="onIsShowInstalledChange"/>
      </el-form-item>
    </el-form>
    <el-tabs v-model="activeTab" @tab-click="onTabChange">
      <el-tab-pane v-for="lang in langList" :key="lang.name" :label="lang.name" :name="lang.executable_name"/>
    </el-tabs>
    <template v-if="activeLang.install_status === 'installed'">
      <template v-if="!['python', 'node'].includes(activeLang.executable_name)">
        <div class="install-wrapper">
          <el-button
            icon="el-icon-check"
            disabled
            type="success"
          >
            {{'已安装'}}
          </el-button>
        </div>
      </template>
      <template v-else>
        <el-table
          height="calc(100vh - 280px)"
          :data="computedDepList"
          :empty-text="depName ? 'No Data' : '请搜索依赖'"
          v-loading="loading"
          border
        >
          <el-table-column
            :label="'名称'"
            prop="name"
            width="180"
          />
          <el-table-column
            :label="!isShowInstalled ? '最新版本' : '版本'"
            prop="version"
            width="100"
          />
          <el-table-column
            v-if="!isShowInstalled"
            :label="'描述'"
            prop="description"
          />
          <el-table-column
            :label="'操作'"
          >
            <template slot-scope="scope">
              <el-button
                v-if="!scope.row.installed"
                :icon="getDepLoading(scope.row) ? 'el-icon-loading' : ''"
                :disabled="getDepLoading(scope.row)"
                size="mini"
                type="primary"
                @click="onClickInstallDep(scope.row)"
              >
                {{'下载'}}
              </el-button>
              <el-button
                v-else
                :icon="getDepLoading(scope.row) ? 'el-icon-loading' : ''"
                :disabled="getDepLoading(scope.row)"
                size="mini"
                type="danger"
                @click="onClickUninstallDep(scope.row)"
              >
                {{'卸载'}}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
    </template>
    <template v-else-if="activeLang.install_status === 'installing'">
      <div class="install-wrapper">
        <el-button
          icon="el-icon-loading"
          disabled
          type="warning"
        >
          {{'下载中'}}
        </el-button>
      </div>
    </template>
    <template v-else-if="activeLang.install_status === 'installing-other'">
      <div class="install-wrapper">
        <el-button
          loading="el-icon-close"
          disabled
          type="warning"
        >
          {{'其他语言正在安装'}}
        </el-button>
      </div>
    </template>
    <template v-else-if="activeLang.install_status === 'not-installed'">
      <div class="install-wrapper">
        <h4>{{'语言还未安装'}}</h4>
        <el-button
          icon="el-icon-check"
          type="primary"
          @click="onClickInstallLang"
        >
          {{'下载'}}
        </el-button>
      </div>
    </template>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'NodeInstallation',
  data () {
    return {
      activeTab: '',
      langList: [],
      depName: '',
      depList: [],
      loading: false,
      isShowInstalled: true,
      installedDepList: [],
      depLoadingDict: {},
      isLoadingInstallLang: false
    }
  },
  computed: {
    ...mapState('node', [
      'nodeForm'
    ]),
    activeLang () {
      for (let i = 0; i < this.langList.length; i++) {
        if (this.langList[i].executable_name === this.activeTab) {
          return this.langList[i]
        }
      }
      return {}
    },
    activeLangName () {
      return this.activeLang.executable_name
    },
    computedDepList () {
      if (this.isShowInstalled) {
        return this.installedDepList
      } else {
        return this.depList
      }
    }
  },
  methods: {
    async getDepList () {
      this.loading = true
      this.depList = []
      const res = await this.$request.get(`/nodes/${this.nodeForm._id}/deps`, {
        lang: this.activeLang.executable_name,
        dep_name: this.depName
      })
      this.loading = false
      this.depList = res.data.data

      if (this.activeLangName === 'python') {
        // 排序
        this.depList = this.depList.sort((a, b) => a.name > b.name ? 1 : -1)

        // 异步获取python附加信息
        this.depList.map(async dep => {
          const resp = await this.$request.get(`/system/deps/${this.activeLang.executable_name}/${dep.name}/json`)
          if (resp) {
            dep.version = resp.data.data.version
            dep.description = resp.data.data.description
          }
        })
      }
    },
    async getInstalledDepList () {
      if (this.activeLang.install_status !== 'installed') return
      if (!['Python', 'Node.js'].includes(this.activeLang.name)) return

      this.loading = true
      this.installedDepList = []
      const res = await this.$request.get(`/nodes/${this.nodeForm._id}/deps/installed`, {
        lang: this.activeLang.executable_name
      })
      this.loading = false
      this.installedDepList = res.data.data
    },
    async fetchAllDepList (queryString, callback) {
      const res = await this.$request.get(`/system/deps/${this.activeLang.executable_name}`, {
        dep_name: queryString
      })
      callback(res.data.data ? res.data.data.map(d => {
        return { value: d, label: d }
      }) : [])
    },
    onSearch () {
      this.isShowInstalled = false
      this.getDepList()
    },
    onIsShowInstalledChange (val) {
      if (val) {
        this.getInstalledDepList()
      } else {
        this.depName = ''
        this.depList = []
      }
    },
    async onClickInstallDep (dep) {
      const name = dep.name
      this.$set(this.depLoadingDict, name, true)
      const arr = this.$route.path.split('/')
      const id = arr[arr.length - 1]
      const data = await this.$request.post(`/nodes/${id}/deps/install`, {
        lang: this.activeLang.executable_name,
        dep_name: name
      })
      if (!data || data.error) {
        this.$notify.error({
          title: '安装依赖失败',
          message: '安装依赖失败: ' + name
        })
      } else {
        this.$notify.success({
          title: '安装依赖成功',
          message: '安装依赖成功: ' + name
        })
        dep.installed = true
      }
      this.$request.put('/actions', {
        type: 'install_dep'
      })
      this.$set(this.depLoadingDict, name, false)
    },
    async onClickUninstallDep (dep) {
      const name = dep.name
      this.$set(this.depLoadingDict, name, true)
      const arr = this.$route.path.split('/')
      const id = arr[arr.length - 1]
      const data = await this.$request.post(`/nodes/${id}/deps/uninstall`, {
        lang: this.activeLang.executable_name,
        dep_name: name
      })
      if (!data || data.error) {
        this.$notify.error({
          title: '卸载依赖失败',
          message: '卸载依赖失败: ' + name
        })
      } else {
        this.$notify.success({
          title: '卸载依赖成功',
          message: '卸载依赖成功: ' + name
        })
        dep.installed = false
      }
      this.$set(this.depLoadingDict, name, false)
    },
    getDepLoading (dep) {
      const name = dep.name
      if (this.depLoadingDict[name] === undefined) {
        return false
      }
      return this.depLoadingDict[name]
    },
    async onClickInstallLang () {
      this.isLoadingInstallLang = true
      const res = await this.$request.post(`/nodes/${this.nodeForm._id}/langs/install`, {
        lang: this.activeLang.executable_name
      })
      if (!res || res.error) {
        this.$notify.error({
          title: '安装语言失败',
          message: '安装语言失败: ' + this.activeLang.name
        })
      } else {
        this.$notify.success({
          title: '安装语言成功',
          message: '安装语言成功: ' + this.activeLang.name
        })
      }
      this.$request.put('/actions', {
        type: 'install_lang'
      })
      this.isLoadingInstallLang = false
    },
    onTabChange () {
      if (this.isShowInstalled) {
        this.getInstalledDepList()
      } else {
        this.depName = ''
        this.depList = []
      }
    }
  },
  async created () {
    const arr = this.$route.path.split('/')
    const id = arr[arr.length - 1]
    const res = await this.$request.get(`/nodes/${id}/langs`)
    this.langList = res.data.data
    this.activeTab = this.langList[0].executable_name || ''
    setTimeout(() => {
      this.getInstalledDepList()
    }, 100)
  }
}
</script>

<style scoped>
  .node-installation >>> .install-wrapper .el-button {
    min-width: 240px;
    font-weight: bolder;
    font-size: 18px
  }
</style>
