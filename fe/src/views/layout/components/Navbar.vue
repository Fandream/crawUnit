<template>
  <div class="navbar">
    <breadcrumb class="breadcrumb" />
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import Breadcrumb from "@/components/Breadcrumb";
import Hamburger from "@/components/Hamburger";

import showdown from "showdown";

export default {
  components: {
    Breadcrumb
  },
  data() {
    return {};
  }
};
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.navbar {
  height: 50px;
  line-height: 50px;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.12), 0 0 3px 0 rgba(0, 0, 0, 0.04);
}
</style>
<style scoped>
.navbar >>> .el-dialog__body {
  padding-top: 0;
}
</style>
