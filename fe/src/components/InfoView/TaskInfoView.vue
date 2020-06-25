<template>
  <div class="info-view">
    <el-row>
      <el-form
        label-width="150px"
        :model="taskForm"
        ref="nodeForm"
        class="node-form"
        label-position="right"
      >
        <el-form-item :label="'任务ID'">
          <el-input
            v-model="taskForm._id"
            placeholder="任务ID"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'状态'">
          <status-tag :status="taskForm.status" />
          <el-badge
            v-if="taskForm.error_log_count > 0"
            :value="taskForm.error_log_count"
            style="margin-left:10px; cursor:pointer;"
          >
            <el-tag type="danger" @click="onClickLogWithErrors">
              <i class="el-icon-warning"></i>
              {{ "日志错误" }}
            </el-tag>
          </el-badge>
          <el-tag
            v-if="taskForm.status === 'finished' && taskForm.result_count === 0"
            type="danger"
            style="margin-left: 10px"
          >
            <i class="el-icon-warning"></i>
            {{ "空结果" }}
          </el-tag>
        </el-form-item>
        <el-form-item :label="'日志文件路径'">
          <el-input
            v-model="taskForm.log_path"
            placeholder="日志文件路径"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'参数'">
          <el-input
            v-model="taskForm.param"
            placeholder="参数"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'创建时间'">
          <el-input
            :value="getTime(taskForm.create_ts)"
            placeholder="创建时间"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'开始时间'">
          <el-input
            :value="getTime(taskForm.start_ts)"
            placeholder="开始时间"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'结束时间'">
          <el-input
            :value="getTime(taskForm.finish_ts)"
            placeholder="结束时间"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'等待时长(秒)'">
          <el-input
            :value="getWaitDuration(taskForm)"
            placeholder="等待时长(秒)"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'运行时长(秒)'">
          <el-input
            :value="getRuntimeDuration(taskForm)"
            placeholder="运行时长(秒)"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'总时长(秒)'">
          <el-input
            :value="getTotalDuration(taskForm)"
            placeholder="总时长(秒)"
            disabled
          ></el-input>
        </el-form-item>
        <el-form-item :label="'结果数'">
          <el-input
            v-model="taskForm.result_count"
            placeholder="结果数"
            disabled
          ></el-input>
        </el-form-item>
        <!--<el-form-item :label="'Average Results Count per Second'">-->
        <!--<el-input v-model="taskForm.avg_num_results" placeholder="Average Results Count per Second" disabled>-->
        <!--</el-input>-->
        <!--</el-form-item>-->
        <el-form-item :label="'错误信息'" v-if="taskForm.status === 'error'">
          <div class="error-message">
            {{ taskForm.error }}
          </div>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container">
      <el-button
        v-if="isRunning"
        size="small"
        type="danger"
        @click="onStop"
        icon="el-icon-video-pause"
      >
        {{ "停止" }}
      </el-button>
      <!--<el-button type="danger" @click="onRestart">Restart</el-button>-->
    </el-row>
  </div>
</template>

<script>
import { mapState } from "vuex";
import StatusTag from "../Status/StatusTag";
import dayjs from "dayjs";

export default {
  name: "NodeInfoView",
  components: { StatusTag },
  computed: {
    ...mapState("task", ["taskForm", "taskLog", "errorLogData"]),
    isRunning() {
      return ["pending", "running"].includes(this.taskForm.status);
    }
  },
  methods: {
    onRestart() {},
    onStop() {
      this.$store
        .dispatch("task/cancelTask", this.$route.params.id)
        .then(() => {
          this.$message.success(
            `Task "${this.$route.params.id}" has been sent signal to stop`
          );
        });
    },
    getTime(str) {
      if (!str || str.match("^0001")) return "NA";
      return dayjs(str).format("YYYY-MM-DD HH:mm:ss");
    },
    getWaitDuration(row) {
      if (!row.start_ts || row.start_ts.match("^0001")) return "NA";
      return dayjs(row.start_ts).diff(row.create_ts, "second");
    },
    getRuntimeDuration(row) {
      if (!row.finish_ts || row.finish_ts.match("^0001")) return "NA";
      return dayjs(row.finish_ts).diff(row.start_ts, "second");
    },
    getTotalDuration(row) {
      if (!row.finish_ts || row.finish_ts.match("^0001")) return "NA";
      return dayjs(row.finish_ts).diff(row.create_ts, "second");
    },
    onClickLogWithErrors() {
      this.$emit("click-log");
    }
  }
};
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

.el-tag {
  height: 36px;
  line-height: 36px;
}

.error-message {
  background-color: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
  border: 1px solid rgba(245, 108, 108, 0.2);
  border-radius: 4px;
  line-height: 18px;
  padding: 5px 10px;
}

.el-form-item {
  text-align: left;
}
</style>
