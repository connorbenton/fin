import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';
import vuetify from './plugins/vuetify';
import { createNamespacedHelpers } from 'vuex';
import TransactionsTable from './components/TransactionsTable.vue';
import AccountsBar from './components/AccountsBar.vue';

Vue.config.productionTip = false;
Vue.component('TransactionsTable', TransactionsTable);
Vue.component('AccountsBar', AccountsBar);

export const vm = new Vue({
  router,
  store,
  vuetify,
  created() {
    this.$store.dispatch('getAll');
  },
  render: (h) => h(App),
}).$mount('#app');
