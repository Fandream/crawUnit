<template>
  <div class="app-container">
    <!--add popup-->
    <el-dialog
      :visible.sync="dialogVisible"
      width="640px"
      :before-close="onDialogClose">
      <el-form label-width="180px"
               class="add-form"
               :model="projectForm"
               :inline-message="true"
               ref="projectForm"
               label-position="right">
        <el-form-item :label="'项目名称'" prop="name" required>
          <el-input id="name" v-model="projectForm.name" :placeholder="'项目名称'"></el-input>
        </el-form-item>
<!--        <el-form-item :label="'项目描述'" prop="description">-->
<!--          <el-input-->
<!--            id="description"-->
<!--            type="textarea"-->
<!--            v-model="projectForm.description"-->
<!--            :placeholder="'项目描述'"-->
<!--          />-->
<!--        </el-form-item>-->
<!--        <el-form-item :label="'标签'" prop="tags">-->
<!--          <el-select-->
<!--            id="tags"-->
<!--            v-model="projectForm.tags"-->
<!--            :placeholder="'输入标签'"-->
<!--            allow-create-->
<!--            filterable-->
<!--            multiple-->
<!--          >-->
<!--          </el-select>-->
<!--        </el-form-item>-->
      </el-form>
      <!--取消、保存-->
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="onDialogClose">{{'取消'}}</el-button>
        <el-button id="btn-submit" size="small" type="primary" @click="onAddSubmit">{{'提交'}}</el-button>
      </span>
    </el-dialog>
    <!--./add popup-->

    <div class="action-wrapper">
      <div class="left">
<!--        <el-select-->
<!--          v-model="filter.tag"-->
<!--          size="small"-->
<!--          :placeholder="'选择标签'"-->
<!--          @change="onFilterChange"-->
<!--        >-->
<!--          <el-option value="" :label="'全部标签'"/>-->
<!--          <el-option-->
<!--            v-for="tag in projectTags"-->
<!--            :key="tag"-->
<!--            :label="tag"-->
<!--            :value="tag"-->
<!--          />-->
<!--        </el-select>-->
      </div>
      <div class="right">
        <el-button
          icon="el-icon-plus"
          type="primary"
          size="small"
          @click="onAdd"
        >
          {{'添加项目'}}
        </el-button>
      </div>
    </div>
    <div class="content">
      <div v-if="projectList.length === 0" class="empty-list">
        {{ '您没有创建项目，请点击 "添加项目" 按钮来创建一个新项目'}}
      </div>
      <ul v-else class="list">
        <li
          class="item"
          v-for="item in projectList.filter(d => d._id !== '000000000000000000000000')"
          :key="item._id"
          @click="onView(item)"
        >
          <el-card
            class="item-card"
          >
            <i v-if="!isNoProject(item)" class="btn-edit fa fa-edit" @click="onEdit(item)"></i>
            <i v-if="!isNoProject(item)" class="btn-close fa fa-trash-o" @click="onRemove(item)"></i>
            <el-row>
              <h4 class="title">{{ item.name }}</h4>
            </el-row>
            <el-row>
              <div style="display: flex; justify-content: space-between">
                <span class="spider-count">
                {{'爬虫数'}}: {{ item.spiders.length }}
                </span>
              </div>
            </el-row>
          </el-card>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'ProjectList',
  data () {
    return {
      defaultTags: [],
      dialogVisible: false,
      isClickAction: false,
      filter: {
        tag: ''
      }
    }
  },
  computed: {
    ...mapState('project', [
      'projectForm',
      'projectList',
      'projectTags'
    ])
  },
  methods: {
    onDialogClose () {
      this.dialogVisible = false
    },
    onFilterChange () {
      this.$store.dispatch('project/getProjectList', this.filter)
    },
    onAdd () {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('project/SET_PROJECT_FORM', { tags: [] })
    },
    onAddSubmit () {
      this.$refs.projectForm.validate(res => {
        if (res) {
          const form = JSON.parse(JSON.stringify(this.projectForm))
          if (this.isEdit) {
            this.$request.post(`/projects/${this.projectForm._id}`, form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('project/getProjectList')
              this.$message.success('项目已经保存')
            })
          } else {
            this.$request.put('/projects', form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('project/getProjectList')
              this.$message.success('项目已经添加')
            })
          }
        }
      })
    },
    onEdit (row) {
      this.isClickAction = true
      setTimeout(() => {
        this.isClickAction = false
      }, 100)

      this.$store.commit('project/SET_PROJECT_FORM', row)
      this.dialogVisible = true
      this.isEdit = true
    },
    onRemove (row) {
      this.isClickAction = true
      setTimeout(() => {
        this.isClickAction = false
      }, 100)

      this.$confirm('确认删除该项目?', '提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('project/removeProject', row._id)
          .then(() => {
            setTimeout(() => {
              this.$store.dispatch('project/getProjectList')
              this.$message.success('项目已经移除')
            }, 100)
          })
      }).catch(() => {
      })
    },
    onView (row) {
      if (this.isClickAction) return

      this.$router.push({
        name: 'SpiderList',
        params: {
          project_id: row._id
        }
      })
    },
    isNoProject (row) {
      return row._id === '000000000000000000000000'
    }
  },
  async created () {
    await this.$store.dispatch('project/getProjectList', this.filter)
    await this.$store.dispatch('project/getProjectTags')
  }
}
</script>

<style scoped>
  .action-wrapper {
    display: flex;
    justify-content: space-between;
    padding-bottom: 10px;
    border-bottom: 1px solid #EBEEF5;
  }

  .list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-wrap: wrap;
  }

  .list .item {
    width: 320px;
    margin: 10px;
  }

  .list .item .item-card {
    position: relative;
    cursor: pointer;
  }

  .list .item .item-card .title {
    margin: 10px 0 0 0;
  }

  .list .item .item-card .spider-count,
  .list .item .item-card .owner {
    font-size: 12px;
    color: grey;
    font-weight: bolder;
  }

  .list .item .item-card .description-wrapper {
    padding-bottom: 5px;
    margin-bottom: 0;
    border-bottom: 1px solid #EBEEF5;
  }

  .list .item .item-card .description {
    font-size: 12px;
    line-height: 16px;
    color: grey;
  }

  .list .item .item-card .tags {
    margin-bottom: -5px;
  }

  .list .item .item-card .tags .tag {
    margin: 0 5px 5px 0;
  }

  .list .item .item-card .el-row {
    margin-bottom: 5px;
  }

  .list .item .item-card .el-row:last-child {
    margin-bottom: 0;
  }

  .list .item .item-card .btn-edit {
    z-index: 1;
    color: grey;
    position: absolute;
    top: 11px;
    right: 40px;
  }

  .list .item .item-card .btn-close {
    z-index: 1;
    color: grey;
    position: absolute;
    top: 10px;
    right: 10px;
  }

  .empty-list {
    font-size: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    height: calc(100vh - 240px);
  }

</style>
