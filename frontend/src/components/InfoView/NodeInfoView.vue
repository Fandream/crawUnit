<template>
  <div class="info-view">
    <el-row>
      <el-form label-width="150px"
               :model="nodeForm"
               ref="nodeForm"
               class="node-form"
               label-position="right">
        <el-form-item :label="'节点名称'">
          <el-input v-model="nodeForm.name" :placeholder="'节点名称'" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="'节点IP'" prop="ip" required>
          <el-input v-model="nodeForm.ip" :placeholder="'节点IP'" disabled></el-input>
        </el-form-item>
        <el-form-item :label="'节点MAC'" prop="ip" required>
          <el-input v-model="nodeForm.mac" :placeholder="'节点MAC'" disabled></el-input>
        </el-form-item>
        <el-form-item :label="'描述'">
          <el-input type="textarea" v-model="nodeForm.description" :placeholder="'描述'" :disabled="isView">
          </el-input>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button size="small" type="success" @click="onSave">{{'保存'}}</el-button>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'NodeInfoView',
  props: {
    isView: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    ...mapState('node', [
      'nodeForm'
    ])
  },
  methods: {
    onSave () {
      this.$refs.nodeForm.validate(valid => {
        if (valid) {
          this.$store.dispatch('node/editNode')
            .then(() => {
              this.$message.success('Node info has been saved successfully')
            })
        }
      })
    }
  }
}
</script>

<style scoped>
  .node-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }
</style>
