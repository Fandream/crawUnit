<template>
    <el-menu
      :show-timeout="200"
      :default-active="routeLevel1"
      :collapse="isCollapse"
      :background-color="variables.menuBg"
      :text-color="variables.menuText"
      :active-text-color="variables.menuActiveText"
      mode="horizontal"
    >
      crawUnit
      <sidebar-item
        v-for="route in routes"
        :class="route.path.replace('/', '')"
        :key="route.path"
        :item="route"
        :base-path="route.path"
      />
    </el-menu>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import variables from '@/styles/variables.scss'
import SidebarItem from './SidebarItem'

export default {
  components: { SidebarItem },
  computed: {
    ...mapState('user', [
      'adminPaths'
    ]),
    ...mapGetters([
      'sidebar'
    ]),
    routeLevel1 () {
      let pathArray = this.$route.path.split('/')
      return `/${pathArray[1]}`
    },
    routes () {
      return this.$router.options.routes.filter(d => {
        const role = this.$store.getters['user/userInfo'].role
        if (role === 'admin') return true
        return !this.adminPaths.includes(d.path)
      })
    },
    variables () {
      return variables
    },
    isCollapse () {
      return !this.sidebar.opened
    },
    version () {
      return this.$store.state.version.version || window.sessionStorage.getItem('v')
    }
  },
  data () {
    return {}
  },
  async created () {
  },
  mounted () {
  }
}
</script>

<style>
  #app .sidebar-container .el-menu {
    height: calc(100% - 50px);
  }

  .sidebar-container .sidebar-logo {
    height: 50px;
    display: flex;
    /*justify-content: center;*/
    align-items: center;
    padding-left: 20px;
    color: #fff;
    background: rgb(48, 65, 86);
    font-size: 24px;
    font-weight: 600;
    font-family: "Verdana", serif;
  }

  .sidebar-container .sidebar-logo.collapsed {
    padding-left: 8px;
  }

  .sidebar-container .sidebar-logo .version {
    margin-left: 5px;
    font-weight: normal;
    font-size: 11px;
  }
</style>
