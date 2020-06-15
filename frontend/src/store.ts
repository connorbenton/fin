import Vue from 'vue';
import Vuex from 'vuex';
import api from './api';

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    webWorkerType: '',
    isFetchTransactions: false,
    fetchTransactionsItemDone: 0,
    fetchTransactionsItemTotal: 0,
    currName: '',
    compareMatch: false,
    trans1: {},
    trans2: {},
    compareAnswer: '',
    catsToCompare: [],
    finishedCats: [],

    filteredTrans: [] as any,

    customStart: '',
    customEnd: '',

    itemTokens: [],

    analysisTrees: [] as any,

    transactions: [] as any,

    categories: [],

    accounts: [],

    apiStateLoaded: false,

    doneReloading: false,

    isDark: null,
  },
  getters: {
    ItemsDone: (state) => state.fetchTransactionsItemDone,
    ItemsTotal: (state) => state.fetchTransactionsItemTotal,
    getName: (state) => state.currName,
    getAllTransactions: (state) => state.transactions,
    getAllCategories: (state) => state.categories,
    getAllAccounts: (state) => state.accounts,
    getAllItemTokens: (state) => state.itemTokens,
    getAllTrees: (state) => state.analysisTrees,
    getTrans1: (state) => state.trans1,
    getTrans2: (state) => state.trans2,
    getDark: (state) => state.isDark,
    getFilteredTrans: (state) => state.filteredTrans,
  },
  mutations: {
    updateFilteredTrans(state, payload) {
      state.filteredTrans = payload;
      // console.log(payload);
    },
    incrementItem(state) {
      state.fetchTransactionsItemDone++;
    },
    newName(state, payload) {
      state.currName = payload;
    },
    isFetch(state, payload) {
      state.isFetchTransactions = payload;
    },
    answerGiven(state, payload) {
      state.compareAnswer = payload;
    },
    newAssignCats(state, payload) {
      state.catsToCompare = payload;
    },
    assignDone(state, payload) {
      state.finishedCats = payload;
    },
    updateTransactions(state, transactions) {
      Object.freeze(transactions);
      state.transactions = transactions;
      state.filteredTrans = transactions;
    },
    updateTransaction(state, transactions) {
      const transSet: any[] = state.transactions;
      transactions.forEach((transaction: any) => {
      const transToUpdate = transSet.find((x) => x.id === transaction.id);
      transToUpdate.category = transaction.category;
      transToUpdate.category_name = transaction.category_name;
      });
    },
    updateAccountIgnore(state, account) {
      const accSet: any[] = state.accounts;
      const accToUpdate = accSet.find((x) => x.id === account.id);
      accToUpdate.ignore_transactions = account.ignore_transactions;
    },
    updateAccountName(state, account) {
      const accSet: any[] = state.accounts;
      const accToUpdate = accSet.find((x) => x.id === account.id);
      accToUpdate.name = account.name;
    },
    updateTrees(state, trees) {
      state.analysisTrees = trees;
    },
    updateCustomTree(state, tree) {
      const index = state.analysisTrees.findIndex((x: any) => x.name === 'custom');
      state.analysisTrees[index] = tree;
      // console.log(treeSet);
      // console.log(state.analysisTrees);
    },
    updateCategories(state, categories) {
      state.categories = categories;
    },
    updateAccounts(state, accounts) {
      state.accounts = accounts;
    },
    updateItemTokens(state, itemTokens) {
      state.itemTokens = itemTokens;
    },
    doneLoading(state, value) {
      state.apiStateLoaded = value;
    },
    setReloading(state, value) {
      state.doneReloading = value;
    },
    setCustomRange(state, value) {
      state.customStart = value.start;
      state.customEnd = value.end;
    },

  },
  actions: {
    async getTransactions({ commit }) {
      const res = await api.getTransactions();
      commit('updateTransactions', res);
    },
    async getCategories({ commit }) {
      const res = await api.getCategories();
      commit('updateCategories', res);
    },
    async getAccounts({ commit }) {
      const res = await api.getAccounts();
      commit('updateAccounts', res);
    },
    async getTrees({ commit }) {
      const res = await api.getTrees();
      commit('updateTrees', res);
      commit('setReloading', true);
    },
    async customFilter({ commit }, { startFromPage, endFromPage }) {
      const res = await api.customAnalyze({start: startFromPage, end: endFromPage});
      commit('updateCustomTree', res);
      commit('setReloading', true);
    },
    async getAll({ commit }) {
      try {
        const [cats, accs, trans, toks, trees] =
          await Promise.all([api.getCategories(), api.getAccounts(), api.getTransactions(),
            api.getItemTokens(), api.getTrees()]);

        commit('updateCategories', cats);
        commit('updateAccounts', accs);
        commit('updateTransactions', trans);
        commit('updateItemTokens', toks);
        commit('updateTrees', trees);
        commit('doneLoading', true);
      } catch (e) {
        // console.error(e);
      }
    },

  },
});

export default store;
