<template>
  <div>
      <h1 class="h1" align="center">CrawUnit</h1>
    <div class="app-container">
      <el-row>
        <ul class="metric-list">
          <li class="metric-item" v-for="m in metrics" @click="onClickMetric(m)" :key="m.name">
            <div class="metric-icon" :class="m.color">
              <i :class="m.icon"></i>
            </div>
            <div class="metric-content" :class="m.color">
              <div class="metric-label">
                {{m.label}}
              </div>
              <div class="metric-number">
                {{overviewStats[m.name]}}
              </div>
            </div>
          </li>
        </ul>
      </el-row>
    </div>
    <div class="app-container">
      <el-row>
        <ul class="metric-list">
          <li class="metric-item" v-for="m in metrics2" @click="onClickMetric(m)" :key="m.name">
            <div class="metric-icon" :class="m.color">
              <!--            <font-awesome-icon :icon="m.icon"/>-->
              <i :class="m.icon"></i>
            </div>
            <div class="metric-content" :class="m.color">
              <div class="metric-label">
                {{m.label}}
              </div>
              <div class="metric-number">
                {{overviewStats[m.name]}}
              </div>
            </div>
          </li>
        </ul>
      </el-row>
    </div>
  </div>
</template>

<script>
import request from '../../api/request'
import echarts from 'echarts'
export default {
  name: 'Home',
  data () {
    return {
      echarts: {},
      overviewStats: {},
      dailyTasks: [],
      metrics: [
        { name: 'task_count', label: '总任务数', icon: 'fa fa-check', color: 'blue', path: 'tasks' },
        { name: 'spider_count', label: '爬虫', icon: 'fa fa-bug', color: 'green', path: 'spiders' }
      ],
      metrics2: [
        { name: 'active_node_count', label: '在线节点', icon: 'fa fa-server', color: 'red', path: 'nodes' },
        { name: 'project_count', label: '项目', icon: 'fa fa-code-fork', color: 'grey', path: 'projects' }
      ]
    }
  },
  methods: {
    initEchartsDailyTasks () {
      const option = {
        xAxis: {
          type: 'category',
          data: this.dailyTasks.map(d => d.date)
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          data: this.dailyTasks.map(d => d.task_count),
          type: 'line',
          areaStyle: {},
          smooth: true
        }],
        tooltip: {
          trigger: 'axis',
          show: true
        }
      }
      this.echarts.dailyTasks = echarts.init(this.$el.querySelector('#echarts-daily-tasks'))
      this.echarts.dailyTasks.setOption(option)
    },
    onClickMetric (m) {
      this.$router.push(`/${m.path}`)
    }
  },
  created () {
    request.get('/stats/home')
      .then(response => {
        // overview stats
        this.overviewStats = response.data.data.overview

        // daily tasks
        this.dailyTasks = response.data.data.daily
        this.initEchartsDailyTasks()
      })
  },
  mounted () {
    // this.$ba.trackPageview('/')
  }
}
</script>
<style scoped>
  .h1{
    font-size: 64px;
    width: 200px;
    position: relative;
    color: #1482f0;
    left: 50%;
    transform: translateX(-60%);

  }
  .h1::after{
    content:'';
    width: 50px;
    height:50px;
    position: absolute;
    left:50%;
    display: block;
    background: url('/favicon.ico') no-repeat;
    background-size: 100%;
    left: -80px;
    top: 22px;
  }
</style>
<style scoped lang="scss">

  .metric-list {
    margin-top: 0;
    padding-left: 0;
    list-style: none;
    display: flex;
    font-size: 16px;

    .metric-item:last-child .metric-card {
      margin-right: 0;
    }

    .metric-item:hover {
      transform: scale(1.05);
      transition: transform 0.5s ease;
    }

    .metric-item {
      flex-basis: 50%;
      height: 100px;
      display: flex;
      color: white;
      cursor: pointer;
      transform: scale(1);
      transition: transform 0.5s ease;

      .metric-icon {
        display: inline-flex;
        width: 64px;
        align-items: center;
        justify-content: center;
        border-top-left-radius: 5px;
        border-bottom-left-radius: 5px;
        font-size: 24px;

        svg {
          width: 24px;
        }
      }

      .metric-content {
        display: flex;
        width: calc(100% - 80px);
        align-items: center;
        opacity: 0.85;
        font-size: 14px;
        padding-left: 15px;
        border-top-right-radius: 5px;
        border-bottom-right-radius: 5px;

        .metric-number{
          font-weight: bolder;
          /*margin-bottom: 5px;*/
          margin: 0 auto;
          font-size: 50px;
        }
        .metric-label {
          font-weight: bolder;
          /*margin-bottom: 5px;*/
          margin: 0 auto;
          font-size: 40px;
        }
      }
      .metric-icon.blue,
      .metric-content.blue {
        background: #409eff;
      }

      .metric-icon.green,
      .metric-content.green {
        background: #67c23a;
      }

      .metric-icon.red,
      .metric-content.red {
        background: #f56c6c;
      }

      .metric-icon.orange,
      .metric-content.orange {
        background: #E6A23C;
      }

      .metric-icon.grey,
      .metric-content.grey {
        background: #97a8be;
      }
    }
  }

  .title {
    padding: 0;
    margin: 0;
  }

  #echarts-daily-tasks {
    height: 360px;
    width: 100%;
  }

  .el-card {
    /*border: 1px solid lightgrey;*/
  }

  .svg-inline--fa {
    width: 100%;
    height: 100%;
  }
</style>
