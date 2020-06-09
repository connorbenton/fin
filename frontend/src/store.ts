import Vue from 'vue';
import Vuex from 'vuex';
import api from './api';
import moment from 'moment';

Vue.use(Vuex);

// const workerActions = new Worker('./actions.ts', { type: 'module' });

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
    // txData: {
    //   oldestDate: null,
    //   txMatrix: {
    //     rangeDays: {},
    //     catData: {},
    //   },
    //   txSets: {
    //     last30TxSet: [],
    //     thisMonthTxSet: [],
    //     lastMonthTxSet: [],
    //     lastSixMonthsTxSet: [],
    //     thisYearTxSet: [],
    //     lastYearTxSet: [],
    //     fromBeginningTxSet: [],
    //     customTxSet: [],
    //   },
    //   txTrees: {
    //     last30TxTree: [],
    //     last30TxTreeNoInvest: [],
    //     thisMonthTxTree: [],
    //     thisMonthTxTreeNoInvest: [],
    //     lastMonthTxTree: [],
    //     lastMonthTxTreeNoInvest: [],
    //     lastSixMonthsTxTree: [],
    //     lastSixMonthsTxTreeNoInvest: [],
    //     thisYearTxTree: [],
    //     thisYearTxTreeNoInvest: [],
    //     lastYearTxTree: [],
    //     lastYearTxTreeNoInvest: [],
    //     fromBeginningTxTree: [],
    //     fromBeginningTxTreeNoInvest: [],
    //     customTxTree: [],
    //     customTxTreeNoInvest: [],
    //   },
    // },

    customStart: '',
    customEnd: '',

    itemTokens: [],

    analysisTrees: [] as any,

    transactions: [] as any,

    categories: [],

    accounts: [],

    apiStateLoaded: false,

    doneReloading: false,
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
  },
  mutations: {
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
      // const cats: any[] = state.categories;
      // const accs: any[] = state.accounts;
      // for (const trans of state.transactions) {
      //   // const trans: any = state.transactions[i];
      //   trans.category_name = cats.find(
      //     (x) => x.id === trans.category,
      //   ).sub_category;
      //   // state.transactions[i].category_name = cats.find(
      //   // x => x.id === state.transactions[i].category
      //   // ).sub_category;
      //   trans.account_name = accs.find(
      //     (x) => x.account_id === trans.account_id,
      //   ).name;
      // }
    },
    updateTransaction(state, transaction) {
      const transSet: any[] = state.transactions;
      const transToUpdate = transSet.find((x) => x.id === transaction.id);
      transToUpdate.category = transaction.category;
      transToUpdate.category_name = transaction.category_name;
    },
    updateAccount(state, account) {
      const accSet: any[] = state.accounts;
      const accToUpdate = accSet.find((x) => x.id === account.id);
      accToUpdate.ignore_transactions = account.ignore_transactions;
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
      //   var topArr = [];
      //   state.categories.forEach((obj) => {
      //     topArr.push(obj.top_category);
      //   });
      //   state.topCategories = [...new Set(topArr)];
    },
    updateAccounts(state, accounts) {
      state.accounts = accounts;
    },
    updateItemTokens(state, itemTokens) {
      state.itemTokens = itemTokens;
    },
    // analysisDataInitial(state) {
    //   // let a = moment();
    //   // state.txData.txSets.last30TxSet = state.transactions.filter(x => {
    //   //   let b = moment(x.date, 'YYYY-MM-DD', true);
    //   //   let c = a.diff(b, 'days');
    //   //   return (c < 30);
    //   // });
    //   buildTxData(state, 'last30TxSet');
    // },
    // analysisData(state) {
    //   let PromiseArray = [];
    //   Object.keys(state.txData.txSets).forEach(s => {
    //     if (s === 'last30TxSet') return;
    //     // PromiseArray.push(buildTxData(state, state.txData.txSets[s]));
    //     PromiseArray.push(buildTxData(state, s));
    //   });
    //   Promise.all(PromiseArray);
    // },
    doneLoading(state, value) {
      state.apiStateLoaded = value;
    },
    setReloading(state, value) {
      state.doneReloading = value;
    },
    // setTxData(state, value) {
    //   state.txData = value;
    // },
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
    async customFilter({ commit }, { startFromPage, endFromPage }) {
      const res = await api.customAnalyze({start: startFromPage, end: endFromPage});
      commit('updateCustomTree', res);
    //   commit('setCustomRange', { start: startFromPage, end: endFromPage });
    //   this.state.webWorkerType = 'custom';
    //   workerActions.postMessage(this.state);
    //   // await buildTxData(this.state, 'customTxSet', start, end);
      commit('setReloading', true);
    },
    // async reanalyze({ commit }) {
    //   this.state.webWorkerType = 'full';
    //   workerActions.postMessage(this.state);
    //   // commit('analysisDataInitial');
    //   // commit('analysisData');
    //   // commit('setReloading');
    //   // let res = await api.getAccounts();
    //   // commit('updateAccounts', res);
    // },
    async getAll({ commit }) {
      try {
        const [cats, accs, trans, toks, trees] =
          await Promise.all([api.getCategories(), api.getAccounts(), api.getTransactions(),
            api.getItemTokens(), api.getTrees()]);
        // Promise.all([
        //   api.getCategories(),
        //   api.getAccounts(),
        //   api.getTransactions(),
        //   api.getItemTokens(),
        // ]).then(([cats, accs, trans, toks]) => {
        // ]).then((values) => {
        // console.log(values);
        // this.state.customEnd = moment();
        // this.state.customEnd = moment().format('YYYY-MM-DD');
        // this.state.customStart = moment();
        // this.state.customStart = moment().subtract(29, 'days').format('YYYY-MM-DD');
        // this.state.customStart = this.state.customStart.format('YYYY-MM-DD');
        commit('updateCategories', cats);
        commit('updateAccounts', accs);
        commit('updateTransactions', trans);
        commit('updateItemTokens', toks);
        commit('updateTrees', trees);
        // this.state.webWorkerType = 'initial';
        // workerActions.postMessage(this.state);
        // commit('analysisDataInitial');
        commit('doneLoading', true);

        // Send this to worker as well
        // commit('analysisData');


        // }).then(() => {
        // this.state.webWorkerType = 'afterInitial';
        // workerActions.postMessage(this.state);
      } catch (e) {
        // console.error(e);
      }
      // }).catch((e) => console.error(e));
      // // let res2 = await api.getCategories();
      // commit('updateCategories', await api.getCategories());
      // // let res3 = await api.getAccounts();
      // commit('updateAccounts', await api.getAccounts());
      // // let res1 = await api.getTransactions();
      // commit('updateTransactions', await api.getTransactions());
      // commit('analysisDataInitial');
      // commit('doneLoading', true);
      // commit('analysisData');
    },

  },
});

// workerActions.onmessage = (e) => {
  // store.commit(e.data.type, e.data.payload);
// };

export default store;
