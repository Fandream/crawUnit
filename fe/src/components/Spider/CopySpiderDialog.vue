<template>
  <el-dialog
    class="copy-spider-dialog"
    ref="form"
    :title="'复制爬虫'"
    :visible="visible"
    width="580px"
    :before-close="onClose"
  >
    <el-form
      label-width="160px"
      :model="form"
      ref="form"
    >
      <el-form-item
        :label="'新爬虫名称'"
        required
      >
        <el-input v-model="form.name" :placeholder="'新爬虫名称'"/>
      </el-form-item>
    </el-form>
    <template slot="footer">
      <el-button type="plain" size="small" @click="$emit('close')">{{'Cancel'}}</el-button>
      <el-button
        type="primary"
        size="small"
        :icon="isLoading ? 'el-icon-loading' : ''"
        :disabled="isLoading"
        @click="onConfirm"
      >
        {{'确认'}}
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'CopySpiderDialog',
  props: {
    spiderId: {
      type: String,
      default: ''
    },
    visible: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      form: {
        name: ''
      },
      isLoading: false
    }
  },
  methods: {
    onClose () {
      this.$emit('close')
    },
    onConfirm () {
      this.$refs['form'].validate(async valid => {
        if (!valid) return
        try {
          this.isLoading = true
          const res = await this.$request.post(`/spiders/${this.spiderId}/copy`, this.form)
          if (!res.data.error) {
            this.$message.success('已成功复制')
          }
          this.$emit('confirm')
          this.$emit('close')
        } finally {
          this.isLoading = false
        }
      })
    }
  }
}
</script>

<style scoped>

</style>
