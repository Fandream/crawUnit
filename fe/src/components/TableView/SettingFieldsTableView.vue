<template>
  <div class="setting-list-table-view">
    <el-row>
      <el-table :data="list"
                class="table edit"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                :cell-style="getCellClassStyle"
      >
        <el-table-column class-name="action" width="80px" align="right">
          <template slot-scope="scope">
            <!--            <i class="action-item el-icon-copy-document" @click="onCopyField(scope.row)"></i>-->
            <i class="action-item el-icon-remove-outline" @click="onRemoveField(scope.row)"></i>
            <i class="action-item el-icon-circle-plus-outline" @click="onAddField(scope.row)"></i>
          </template>
        </el-table-column>
        <el-table-column :label="'名称'" width="240px">
          <template slot-scope="scope">
            <el-input
              v-model="scope.row.name"
              :placeholder="'名称'"
              suffix-icon="el-icon-edit"
              @change="onChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="'值'" width="auto" min-width="120px">
          <template slot-scope="scope">
            <el-input
              v-model="scope.row.value"
              :placeholder="'值'"
              suffix-icon="el-icon-edit"
              @change="onChange(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'SettingFieldsTableView',
  props: {
    type: {
      type: String,
      default: 'list'
    },
    title: {
      type: String,
      default: ''
    },
    stageNames: {
      type: Array,
      default () {
        return []
      }
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    list () {
      const list = []
      for (let name in this.spiderForm.config.settings) {
        if (this.spiderForm.config.settings.hasOwnProperty(name)) {
          const value = this.spiderForm.config.settings[name]
          list.push({ name, value })
        }
      }
      return list
    }
  },
  methods: {
    onChange (row) {
      if (this.list.filter(d => d.name === row.name).length > 1) {
        this.$message.error(`Duplicated field names for ${row.name}`)
      }
      this.$store.commit('spider/SET_SPIDER_FORM_CONFIG_SETTINGS', this.list)
    },
    onRemoveField (row) {
      const list = JSON.parse(JSON.stringify(this.list))
      for (let i = 0; i < list.length; i++) {
        if (row.name === list[i].name) {
          list.splice(i, 1)
        }
      }
      if (list.length === 0) {
        list.push({
          name: `VARIABLE_NAME_${Math.floor(new Date().getTime())}`,
          value: `VARIABLE_VALUE_${Math.floor(new Date().getTime())}`
        })
      }
      this.$store.commit('spider/SET_SPIDER_FORM_CONFIG_SETTINGS', list)
    },
    onAddField (row) {
      const list = JSON.parse(JSON.stringify(this.list))
      for (let i = 0; i < list.length; i++) {
        if (row.name === list[i].name) {
          const name = 'VARIABLE_NAME_' + Math.floor(new Date().getTime())
          const value = 'VARIABLE_VALUE_' + Math.floor(new Date().getTime())
          list.push({ name, value })
          break
        }
      }
      this.$store.commit('spider/SET_SPIDER_FORM_CONFIG_SETTINGS', list)
    },
    getCellClassStyle ({ row, columnIndex }) {
      if (columnIndex === 1) {
        // 字段名称
        if (!row.name) {
          return {
            'border': '1px solid red'
          }
        }
      } else if (columnIndex === 3) {
        // 选择器
        if (!row.css && !row.xpath) {
          return {
            'border': '1px solid red'
          }
        }
      }
    },
    onChangeNextStage (row) {
      this.list.forEach(f => {
        if (f.name !== row.name) {
          this.$set(f, 'next_stage', '')
        }
      })
    }
  },
  created () {
    if (this.list.length === 0) {
      this.$store.commit(
        'spider/SET_SPIDER_FORM_CONFIG_SETTING_ITEM',
        'VARIABLE_NAME_' + Math.floor(new Date().getTime()),
        'VARIABLE_VALUE_' + Math.floor(new Date().getTime())
      )
    }
  }
}
</script>

<style scoped>
  .el-table.edit >>> .el-table__body td {
    padding: 0;
  }

  .el-table.edit >>> .el-table__body td .cell {
    padding: 0;
    font-size: 12px;
  }

  .el-table.edit >>> .el-input__inner:hover {
    text-decoration: underline;
  }

  .el-table.edit >>> .el-input__inner {
    height: 36px;
    border: none;
    border-radius: 0;
    font-size: 12px;
  }

  .el-table.edit >>> .el-select .el-input .el-select__caret {
    line-height: 36px;
  }

  .el-table.edit >>> .button-selector-item {
    cursor: pointer;
    margin: 0 5px;
  }

  .el-table.edit >>> .el-tag.inactive {
    opacity: 0.5;
  }

  .el-table.edit >>> .action {
    background: none !important;
    border: none;
  }

  .el-table.edit >>> tr {
    border: none;
  }

  .el-table.edit >>> tr th {
    border-right: 1px solid rgb(220, 223, 230);
  }

  .el-table.edit >>> tr td:nth-child(2) {
    border-left: 1px solid rgb(220, 223, 230);
  }

  .el-table.edit >>> tr td {
    border-right: 1px solid rgb(220, 223, 230);
  }

  .el-table.edit::before {
    background: none;
  }

  .el-table.edit >>> .action-item {
    font-size: 14px;
    margin-right: 5px;
    cursor: pointer;
  }

  .el-table.edit >>> .action-item:last-child {
    margin-right: 10px;
  }

  .button-group-container {
    /*display: inline-block;*/
    /*width: 100%;*/
  }

  .button-group-container .title {
    float: left;
    line-height: 32px;
  }

  .button-group-container .button-group {
    float: right;
  }

  .action-button-group {
    display: flex;
    margin-left: 10px;
  }

  .action-button-group >>> .el-checkbox__label {
    font-size: 12px;
  }

  .el-table.edit >>> .el-select.disabled .el-input__inner {
    color: lightgrey;
  }
</style>
