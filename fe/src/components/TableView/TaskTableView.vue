<template>
  <div class="task-table-view">
    <el-row class="title-wrapper">
      <h5 class="title">{{ title }}</h5>
      <el-button
        type="success"
        plain
        class="small-btn"
        size="mini"
        icon="fa fa-refresh"
        @click="onRefresh"
      ></el-button>
    </el-row>
    <el-table
      :data="taskList"
      border
      :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
      height="480px"
      @row-click="onClickTask"
    >
      <el-table-column property="node" :label="'节点'" width="120" align="left">
        <template slot-scope="scope">
          <a class="a-tag" @click="onClickNode(scope.row)">{{
            scope.row.node_name
          }}</a>
        </template>
      </el-table-column>
      <el-table-column
        property="spider_name"
        :label="'爬虫'"
        width="120"
        align="left"
      >
        <template slot-scope="scope">
          <a class="a-tag" @click="onClickSpider(scope.row)">{{
            scope.row.spider_name
          }}</a>
        </template>
      </el-table-column>
      <el-table-column property="param" :label="'参数'" width="120">
        <template slot-scope="scope">
          <span>{{ scope.row.param }}</span>
        </template>
      </el-table-column>
      <el-table-column
        property="result_count"
        :label="'结果数'"
        width="60"
        align="right"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.result_count }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="'状态'" align="left" width="100">
        <template slot-scope="scope">
          <status-tag :status="scope.row.status" />
        </template>
      </el-table-column>
      <!--<el-table-column property="create_ts" label="Create Time" width="auto" align="center"></el-table-column>-->
      <el-table-column
        property="create_ts"
        :label="'创建时间'"
        width="150"
        align="left"
      >
        <template slot-scope="scope">
          <a href="javascript:" class="a-tag" @click="onClickTask(scope.row)">
            {{ getTime(scope.row.create_ts).format("YYYY-MM-DD HH:mm:ss") }}
          </a>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { mapState } from "vuex";
import dayjs from "dayjs";
import StatusTag from "../Status/StatusTag";

export default {
  name: "TaskTableView",
  components: { StatusTag },
  data() {
    return {
      // setInterval handle
      handle: undefined,
      // 防抖
      clicked: false
    };
  },
  props: {
    title: String
  },
  computed: {
    ...mapState("spider", ["spiderForm"]),
    ...mapState("task", ["taskList"])
  },
  methods: {
    onClickSpider(row) {
      this.clicked = true;
      setTimeout(() => {
        this.clicked = false;
      }, 100);
      this.$router.push(`/spiders/${row.spider_id}`);
      if (this.$route.path.match(/\/nodes\//)) {
      }
    },
    onClickNode(row) {
      this.clicked = true;
      setTimeout(() => {
        this.clicked = false;
      }, 100);
      this.$router.push(`/nodes/${row.node_id}`);
      if (this.$route.path.match(/\/spiders\//)) {
      }
    },
    onClickTask(row) {
      if (!this.clicked) {
        this.$router.push(`/tasks/${row._id}`);
        if (this.$route.path.match(/\/nodes\//)) {
        } else if (this.$route.path.match(/\/spiders\//)) {
        }
      }
    },
    onRefresh() {
      if (this.$route.path.split("/")[1] === "spiders") {
        this.$store.dispatch("spider/getTaskList", this.$route.params.id);
      } else if (this.$route.path.split("/")[1] === "nodes") {
        this.$store.dispatch("node/getTaskList", this.$route.params.id);
      }
    },
    getTime(str) {
      return dayjs(str);
    }
  },
  mounted() {
    this.handle = setInterval(() => {
      this.onRefresh();
    }, 5000);
  },
  destroyed() {
    clearInterval(this.handle);
  }
};
</script>

<style scoped>
.task-table-view {
  margin-bottom: 10px;
}

.el-table .a-tag {
  text-decoration: underline;
}

.title {
  margin: 10px 0 3px 0;
  float: left;
}

.small-btn {
  float: right;
  width: 24px;
  margin: 0;
  padding: 5px;
}

.el-table >>> tr {
  cursor: pointer;
}
</style>
