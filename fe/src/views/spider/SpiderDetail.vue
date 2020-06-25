<template>
  <div class="app-container">
    <!--selector-->
    <div class="selector">
      <label class="label">{{ "爬虫" }}: </label>
      <el-select
        id="spider-select"
        v-model="spiderForm._id"
        @change="onSpiderChange"
      >
        <el-option
          v-for="op in spiderList"
          :key="op._id"
          :value="op._id"
          :label="op.name"
        ></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="border-card">
      <el-tab-pane :label="'概览'" name="overview">
        <spider-overview />
      </el-tab-pane>
      <el-tab-pane v-if="isGit" :label="'Git'" name="git-settings">
        <git-settings />
      </el-tab-pane>
      <el-tab-pane v-if="isConfigurable" :label="'配置'" name="config">
        <config-list ref="config" @convert="onConvert" />
      </el-tab-pane>
      <el-tab-pane :label="'文件'" name="files">
        <file-list ref="file-list" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import { mapState } from "vuex";
import FileList from "../../components/File/FileList";
import SpiderOverview from "../../components/Overview/SpiderOverview";
import SpiderStats from "../../components/Stats/SpiderStats";
import ConfigList from "../../components/Config/ConfigList";
import SpiderSchedules from "./SpiderSchedules";
import SpiderScrapy from "../../components/Scrapy/SpiderScrapy";
import GitSettings from "../../components/Settings/GitSettings";

export default {
  name: "SpiderDetail",
  components: {
    GitSettings,
    SpiderScrapy,
    SpiderSchedules,
    ConfigList,
    SpiderStats,
    // EnvironmentList,
    FileList,
    SpiderOverview
  },
  watch: {
    configListTs() {
      this.onConvert();
    }
  },
  data() {
    return {
      activeTabName: "overview",
      redirectType: ""
    };
  },
  computed: {
    ...mapState("spider", ["spiderList", "spiderForm", "configListTs"]),
    ...mapState("file", ["currentPath"]),
    ...mapState("deploy", ["deployList"]),
    isCustomized() {
      return this.spiderForm.type === "customized";
    },
    isConfigurable() {
      return this.spiderForm.type === "configurable";
    },
    isScrapy() {
      return this.isCustomized && this.spiderForm.is_scrapy;
    },
    isGit() {
      return this.spiderForm.is_git;
    }
  },
  methods: {
    async onTabClick(tab) {
      if (this.activeTabName === "analytics") {
        setTimeout(() => {
          this.$refs["spider-stats"].update();
        }, 0);
      } else if (this.activeTabName === "config") {
        setTimeout(() => {
          this.$refs["config"].update();
        }, 0);
      } else if (this.activeTabName === "scrapy-settings") {
        await this.getScrapyData();
      } else if (this.activeTabName === "files") {
        await this.$store.dispatch("spider/getFileTree");
        if (this.currentPath) {
          await this.$store.dispatch("file/getFileContent", {
            path: this.currentPath
          });
        }
      }
    },
    onSpiderChange(id) {
      this.$router.push(`/spiders/${id}`);
    },
    async getScrapyData() {
      await Promise.all([
        this.$store.dispatch(
          "spider/getSpiderScrapySpiders",
          this.$route.params.id
        ),
        this.$store.dispatch(
          "spider/getSpiderScrapyItems",
          this.$route.params.id
        ),
        this.$store.dispatch(
          "spider/getSpiderScrapySettings",
          this.$route.params.id
        ),
        this.$store.dispatch(
          "spider/getSpiderScrapyPipelines",
          this.$route.params.id
        )
      ]);
    },
    async onClickScrapySpider(filepath) {
      this.activeTabName = "files";
      await this.$store.dispatch("spider/getFileTree");
      this.$refs["file-list"].clickSpider(filepath);
    },
    async onClickScrapyPipeline() {
      this.activeTabName = "files";
      await this.$store.dispatch("spider/getFileTree");
      this.$refs["file-list"].clickPipeline();
    },
    onConvert() {
      this.activeTabName = "overview";
    }
  },
  async created() {
    // get spider basic info
    await this.$store.dispatch("spider/getSpiderData", this.$route.params.id);

    // get spider file info
    await this.$store.dispatch("file/getFileList", this.spiderForm.src);

    // get spider tasks
    await this.$store.dispatch("spider/getTaskList", this.$route.params.id);

    // get spider list
    await this.$store.dispatch("spider/getSpiderList", { owner_type: "all" });
  }
};
</script>

<style scoped>
.selector {
  display: flex;
  align-items: center;
  position: absolute;
  right: 48px;
  /*float: right;*/
  z-index: 999;
  margin-top: 5px;
}

.selector .el-select {
  height: 30px;
  line-height: 30px;
  padding-left: 10px;
  width: 180px;
  border-radius: 0;
}

.selector .el-select >>> .el-input__icon,
.selector .el-select >>> .el-input__inner {
  border-radius: 0;
  height: 30px;
  line-height: 30px;
}

.label {
  text-align: right;
  width: 80px;
  color: #909399;
  font-weight: 100;
}
</style>
