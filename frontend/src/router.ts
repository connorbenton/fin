import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

const router = new Router({
  mode: 'history',
  // base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      redirect: '/transactions',
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('./views/Transactions.vue'),
    },
    {
      path: '/analysis',
      name: 'analysis',
      component: () => import('./views/Analysis.vue'),
    },
    {
      path: '/databasego',
      name: 'databasego',
      component: () => import('./views/DatabaseGo.vue'),
    },
    {
      path: '/database',
      name: 'database',
      component: () => import('./views/Database.vue'),
    },
    {
      path: '/accounts',
      name: 'accounts',
      component: () => import('./views/Accounts.vue'),
    },
  ],
});

export default router;
