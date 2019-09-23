import Vue from 'vue';
import Router from 'vue-router';
import Summary from './views/Summary.vue';

Vue.use(Router);

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'summary',
      component: Summary,
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('./views/Transactions.vue'),
    },
    {
      path: '/bills',
      name: 'bills',
      component: () => import('./views/Bills.vue'),
    },
    {
      path: '/investments',
      name: 'investments',
      component: () => import('./views/Investments.vue'),
    },
    {
      path: '/analysis',
      name: 'analysis',
      component: () => import('./views/Analysis.vue'),
    },
    {
      path: '/accounts',
      name: 'accounts',
      component: () => import('./views/Accounts.vue'),
    }
  ],
});

export default router;