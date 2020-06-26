<template>
  <div class="legend">
    <el-tag type="primary" size="small">
      <i class="el-icon-loading"></i>
      {{ "待定" }}
    </el-tag>
    <el-tag type="warning" size="small">
      <i class="el-icon-loading"></i>
      {{ "进行中" }}
    </el-tag>
    <el-tag type="success" size="small">
      <i class="el-icon-check"></i>
      {{ "已完成" }}
    </el-tag>
    <el-tag type="danger" size="small">
      <i class="el-icon-error"></i>
      {{ "错误" }}
    </el-tag>
    <el-tag type="info" size="small">
      <i class="el-icon-video-pause"></i>
      {{ "取消" }}
    </el-tag>
    <el-tag type="danger" size="small">
      <i class="el-icon-warning"></i>
      {{ "异常" }}
    </el-tag>
  </div>
</template>

<script>
export default {
  name: 'StatusLegend'
};
</script>

<style scoped>
.legend .el-tag {
  margin-right: 5px;
}
</style>
