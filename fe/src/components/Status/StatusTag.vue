<template>
  <el-tag :type="type" class="status-tag">
    <i :class="icon"></i>
    {{ label }}
  </el-tag>
</template>

<script>
export default {
  name: 'StatusTag',
  props: {
    status: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      statusDict: {
        pending: { label: '待定', type: 'primary' },
        running: { label: '进行中', type: 'warning' },
        finished: { label: '已完成', type: 'success' },
        error: { label: '错误', type: 'danger' },
        cancelled: { label: '已取消', type: 'info' },
        abnormal: { label: '异常', type: 'danger' }
      }
    };
  },
  computed: {
    type () {
      const s = this.statusDict[this.status];
      if (s) {
        return s.type;
      }
      return '';
    },
    label () {
      const s = this.statusDict[this.status];
      if (s) {
        return s.label;
      }
      return 'NA';
    },
    icon () {
      if (this.status === 'finished') {
        return 'el-icon-check';
      } else if (this.status === 'pending') {
        return 'el-icon-loading';
      } else if (this.status === 'running') {
        return 'el-icon-loading';
      } else if (this.status === 'error') {
        return 'el-icon-error';
      } else if (this.status === 'cancelled') {
        return 'el-icon-video-pause';
      } else if (this.status === 'abnormal') {
        return 'el-icon-warning';
      } else {
        return 'el-icon-question';
      }
    }
  }
};
</script>

<style scoped></style>
