<template>
  <div class="environment-list">
    <el-row>
      <div class="button-group">
        <el-button
          size="small"
          type="primary"
          @click="addEnv"
          icon="el-icon-plus"
          >{{ "添加环境变量" }}</el-button
        >
        <el-button size="small" type="success" @click="save">{{
          "保存"
        }}</el-button>
      </div>
    </el-row>
    <el-row>
      <el-table :data="spiderForm.envs">
        <el-table-column :label="'变量'">
          <template slot-scope="scope">
            <el-input v-model="scope.row.name" :placeholder="'变量'"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="'值'">
          <template slot-scope="scope">
            <el-input v-model="scope.row.value" :placeholder="'值'"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="'操作'">
          <template slot-scope="scope">
            <el-button
              size="mini"
              icon="el-icon-delete"
              type="danger"
              @click="deleteEnv(scope.$index)"
            ></el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "EnvironmentList",
  computed: {
    ...mapState("spider", ["spiderForm"])
  },
  methods: {
    addEnv() {
      if (!this.spiderForm.envs) {
        this.$set(this.spiderForm, "envs", []);
      }
      this.spiderForm.envs.push({
        name: "",
        value: ""
      });
    },
    deleteEnv(index) {
      this.spiderForm.envs.splice(index, 1);
    },
    save() {
      this.$store
        .dispatch("spider/editSpider")
        .then(() => {
          this.$message.success("Spider info has been saved successfully");
        })
        .catch(error => {
          this.$message.error(error);
        });
    }
  }
};
</script>

<style scoped>
.button-group {
  width: 100%;
  text-align: right;
}

.el-table {
  min-height: 360px;
}
</style>
