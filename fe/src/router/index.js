import Vue from 'vue';
import Router from 'vue-router';

import store from '../store';
import request from '../api/request';
import stats from '../utils/stats';

import Layout from '../views/layout/Layout';
Vue.use(Router);
export const constantRouterMap = [
  { path: '/404', component: () => import('../views/404'), hidden: true },
  { path: '/', redirect: '/home' },

  // Crawlab Pages
  {
    path: '/home',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('../views/home/Home'),
        meta: {
          title: '主页'
        }
      }
    ]
  },
  {
    path: '/projects',
    component: Layout,
    meta: {
      title: '项目'
    },
    children: [
      {
        path: '',
        name: 'Project',
        component: () => import('../views/project/ProjectList'),
        meta: {
          title: '项目'
        }
      }
    ]
  },
  {
    path: '/spiders',
    component: Layout,
    meta: {
      title: '爬虫'
    },
    children: [
      {
        path: '',
        name: 'SpiderList',
        component: () => import('../views/spider/SpiderList'),
        meta: {
          title: '爬虫'
        }
      },
      {
        path: ':id',
        name: 'SpiderDetail',
        component: () => import('../views/spider/SpiderDetail'),
        meta: {
          title: '爬虫详情'
        },
        hidden: true
      }
    ]
  },
  {
    path: '/tasks',
    component: Layout,
    meta: {
      title: 'Task'
    },
    children: [
      {
        path: '',
        name: 'TaskList',
        component: () => import('../views/task/TaskList'),
        meta: {
          title: '任务'
        }
      },
      {
        path: ':id',
        name: 'TaskDetail',
        component: () => import('../views/task/TaskDetail'),
        meta: {
          title: '任务详情'
        },
        hidden: true
      }
    ]
  },
  {
    path: '/nodes',
    component: Layout,
    meta: {
      title: '节点'
    },
    children: [
      {
        path: '',
        name: 'NodeList',
        component: () => import('../views/node/NodeList'),
        meta: {
          title: '节点'
        }
      },
      {
        path: ':id',
        name: 'NodeDetail',
        component: () => import('../views/node/NodeDetail'),
        meta: {
          title: '节点详情'
        },
        hidden: true
      }
    ]
  },
  { path: '*', redirect: '/404', hidden: true }
];

const router = new Router({
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
});

router.beforeEach((to, from, next) => {
  if (to.meta && to.meta.title) {
    window.document.title = `CrawUnit - ${to.meta.title}`;
  } else {
    window.document.title = 'CrawUnit';
  }

  if (['/login', '/signup'].includes(to.path)) {
    next();
  } else {
    if (window.localStorage.getItem('token')) {
      next();
    } else {
      next('/login');
    }
  }
});

router.afterEach(async (to, from, next) => {
  if (to.path) {
    await store.dispatch('setting/getSetting');
    const res = await request.get('/version');
    const version = res.data.data;
    store.commit('version/SET_VERSION', version);
    sessionStorage.setItem('v', version);
    stats.sendPv(to.path);
  }
});

export default router;
