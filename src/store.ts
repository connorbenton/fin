import Vue from 'vue';
import Vuex from 'vuex';
import api from './api';
import moment from 'moment';

Vue.use(Vuex);

const workerActions = new Worker('./actions.ts', {type: 'module'});

// function buildTxData(state, key, start = 0, end = 0) {
//   return new Promise((resolve, reject) => {
//     try {
//       let setKey = key;
//       let treeKey = key.toString().substring(0, key.length - 3) + 'Tree';
//       // this.categories = this.$store.getters.getAllCategories;
//       // this.accounts = this.$store.getters.getAllAccounts;
//       // this.categories = await api.getCategories();
//       // this.accounts = await api.getAccounts();
//       let cats = state.categories.map(c => Object.assign({}, c));
//       let accs = state.accounts.map(a => Object.assign({}, a));

//       let a = moment();
//       switch (setKey) {
//         case 'last30TxSet':
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let b = moment(x.date, 'YYYY-MM-DD', true);
//             let c = a.diff(b, 'days');
//             return (c < 30);
//           });
//           break;
//         case 'thisMonthTxSet':
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let b = moment(x.date, 'YYYY-MM-DD', true);
//             return b.isSame(a, 'month');
//             // let c = a.diff(b, 'days');
//             // return (c < 30);
//           });
//           break;
//         case 'lastMonthTxSet':
//           a = a.subtract(1, 'month');
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let b = moment(x.date, 'YYYY-MM-DD', true);
//             return b.isSame(a, 'month');
//           });
//           break;
//         case 'lastSixMonthsTxSet':
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let b = moment(x.date, 'YYYY-MM-DD', true);
//             let c = a.diff(b, 'months');
//             return (c < 6);
//           });
//           break;
//         case 'thisYearTxSet':
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let b = moment(x.date, 'YYYY-MM-DD', true);
//             return b.isSame(a, 'year');
//           });
//           break;
//         case 'lastYearTxSet':
//           a = a.subtract(1, 'year');
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let b = moment(x.date, 'YYYY-MM-DD', true);
//             return b.isSame(a, 'year');
//           });
//           break;
//         case 'fromBeginningTxSet':
//           state.txData.txSets[setKey] = state.transactions;
//           break;
//         case 'customTxSet':
//           if (start === 0) state.txData.txSets[setKey] = state.txData.txSets.last30TxSet;
//           else {
//           let st = moment(start, 'YYYY-MM-DD', true);
//           let en = moment(end, 'YYYY-MM-DD', true);
//           state.txData.txSets[setKey] = state.transactions.filter(x => {
//             let c = moment(x.date, 'YYYY-MM-DD', true);
//             return !(c.isBefore(st) || c.isAfter(en));
//           });
//         }
//       }

//       // let loadcats = this.$store.getters.getAllCategories;
//       // let loadaccs = this.$store.getters.getAllAccounts;
//       // let cats = loadcats.slice(0);
//       // let accs = loadaccs.slice(0);
//       // let y = this.$store.getters.getAllCategories;
//       // this.transactions = await api.getTransactions();
//       // this.transactions = [...this.$store.getters.getAllTransactions];
//       let txSet = state.txData.txSets[setKey];
//       for (let i in txSet) {
//         txSet[i].accName = accs.find(
//           x => x.account_id === txSet[i].account_id
//         ).name;
//         let matchCat = cats.find(x => x.id === txSet[i].category);
//         txSet[i].catName = matchCat.subCategory;
//         if (typeof matchCat.count === "undefined") matchCat.count = 0;
//         matchCat.count = matchCat.count + 1;
//         if (typeof matchCat.total === "undefined") matchCat.total = 0;
//         matchCat.total =
//           matchCat.total + parseFloat(txSet[i].normalized_amount);

//         // matchCat.count = matchCat.count + 1;
//         // matchCat.total = matchCat.total + parseFloat(txSet[i].amount);
//       }

//       let txTree = {};

//       txTree.name = "Transactions by Category";
//       txTree.children = [];
//       txTree.value = 0;
//       for (let j in cats) {
//         // if (cats[j].excludeFromAnalysis || cats[j].topCategory === "Income")
//         if (cats[j].excludeFromAnalysis)
//           continue;
//         if (cats[j].subCategory === cats[j].topCategory) {
//           let newChild = {};
//           newChild.name = cats[j].topCategory;
//           newChild.children = [];
//           newChild.value = 0;
//           newChild.count = 0;
//           // newChild.dbID = '';

//           let children = cats.filter(
//             x => x.topCategory === cats[j].topCategory
//           );
//           for (let k in children) {
//             let subCatChildToPush = {};
//             if (children[k].subCategory === children[k].topCategory) {
//               subCatChildToPush.name = children[k].subCategory + ` (General)`;
//               newChild.dbID = children[k].id;
//             } else {
//               subCatChildToPush.name = children[k].subCategory;
//             }

//             subCatChildToPush.dbID = children[k].id;
//             let value = children[k].total;
//             let count = children[k].count;
//             subCatChildToPush.value =
//               typeof value === "undefined" ? 0 : -1 * value;
//             subCatChildToPush.count = typeof count === "undefined" ? 0 : count;
//             // subCatChildToPush.percent = "";
//             // if (newChild.name === 'Income' || child.value < 0) {
//             if (children[k].topCategory === 'Income' || subCatChildToPush.value < 0) {
//               subCatChildToPush.trueValue = subCatChildToPush.value;
//               subCatChildToPush.value = 0;
//             }
//             newChild.value = newChild.value + subCatChildToPush.value;
//             newChild.count = newChild.count + subCatChildToPush.count;
//             //  subCatChildToPush.count = children[k].count;
//             newChild.children.push(subCatChildToPush);
//           }
//           // newChild.children.map(obj => ({...obj, percent: (obj.value / newChild.value).toFixed(1)+"%"}));
//           newChild.children.forEach(function (child) {
//             // newChild.children[k].percent = '';
//             // if (newChild.name === 'Income' || child.value < 0) {
//               // child.trueValue = child.value;
//               // child.value = 0;
//             // }
//             child.percent = ((child.value / newChild.value) * 100).toFixed(1) + "%";
//             // let x = "y";
//           });
//           txTree.children.push(newChild);
//           txTree.value =
//             txTree.value + newChild.value;
//         }
//       }
//       let total = txTree.value;
//       // txTree.children = txTree.children.filter(child => child.count > 0 && child.value > 0);
//       // txTree.children = txTree.children.filter(child => child.value > -1);
//         // for (let index in txTree.children) {
//       txTree.children.forEach(child => {
//             if (child.name === 'Income' || child.value < 0) {
//               child.trueValue = child.value;
//               child.value = 0;
//             }
//         // let index = txTree.children.indexOf(child);
//         // let child = txTree.children[index]
//         child.percent = ((child.value / total) * 100).toFixed(1) + "%";
//         // if (child.value / total < 0.1) {
//         // if (child.count < 1 || child.value < 0) {
//           // txTree.children.splice(index, 1);
//           // return;
//         // }
//         // if (setKey == 'last30TxSet' && child.name === 'Health & Fitness')
//         // {
//         //   console.log('h');
//         // }
//       // }
//       });

//       let treeKeyNoInvest = treeKey + 'NoInvest';
//       // let txTreeNoInvest = {...txTree};
//       // Object.assign(txTreeNoInvest, txTree);
//       let txTreeNoInvest = JSON.parse(JSON.stringify(txTree))
//       // txTreeNoInvest.name = "Transactions by Category (Financial Excluded)"
//       if (txTreeNoInvest.value) {
//         let fin = txTreeNoInvest.children[5];
//         if (fin.value) {
//           txTreeNoInvest.value = txTreeNoInvest.value - fin.value;
//           fin.trueValue = fin.value;
//           fin.value = 0;
//           fin.children.forEach(subcat => {
//             subcat.trueValue = subcat.value;
//             subcat.value = 0;
//             subcat.percent = 0;
//           })
//           txTreeNoInvest.children.forEach(cat => {
//             cat.percent = ((cat.value / txTreeNoInvest.value) * 100).toFixed(1) + "%";
//           })
//         }
//       }
//       state.txData.txTrees[treeKeyNoInvest] = txTreeNoInvest;
//       state.txData.txTrees[treeKey] = txTree;
//       resolve();
//     }
//     catch (error) {
//       reject(error);
//     }
//   })
//   // let a = this.transactionsTree;
//   // let b = "ha";
// }

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
    txData: {
      oldestDate: null,
      txMatrix: {
        rangeDays: {},
        catData: {},
      },
      txSets: {
        last30TxSet: [],
        thisMonthTxSet: [],
        lastMonthTxSet: [],
        lastSixMonthsTxSet: [],
        thisYearTxSet: [],
        lastYearTxSet: [],
        fromBeginningTxSet: [],
        customTxSet: [],
      },
      txTrees: {
        last30TxTree: [],
        last30TxTreeNoInvest: [],
        thisMonthTxTree: [],
        thisMonthTxTreeNoInvest: [],
        lastMonthTxTree: [],
        lastMonthTxTreeNoInvest: [],
        lastSixMonthsTxTree: [],
        lastSixMonthsTxTreeNoInvest: [],
        thisYearTxTree: [],
        thisYearTxTreeNoInvest: [],
        lastYearTxTree: [],
        lastYearTxTreeNoInvest: [],
        fromBeginningTxTree: [],
        fromBeginningTxTreeNoInvest: [],
        customTxTree: [],
        customTxTreeNoInvest: [],
      },
    },

    customStart: '',
    customEnd: '',

    itemTokens: [],

    transactions: [],

    categories: [],

    accounts: [],

    apiStateLoaded: false,
    
    doneReloading: false,
  },
  getters: {
    ItemsDone: state => { return state.fetchTransactionsItemDone; },
    ItemsTotal: state => { return state.fetchTransactionsItemTotal; },
    getName: state => { return state.currName; },
    getAllTransactions: state => { return state.transactions },
    getAllCategories: state => { return state.categories },
    getAllAccounts: state => { return state.accounts },
    getAllItemTokens: state => { return state.itemTokens },
    getTrans1: state => { return state.trans1; },
    getTrans2: state => { return state.trans2; }
  },
  mutations: {
    incrementItem(state) {
      state.fetchTransactionsItemDone++
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
      state.transactions = transactions;
      let cats: any[] = state.categories;
      let accs: any[] = state.accounts;
      for (let i in state.transactions) {
        let trans: any = state.transactions[i];
        trans.catName = cats.find(
          x => x.id === trans.category
        ).subCategory;
        // state.transactions[i].catName = cats.find(
          // x => x.id === state.transactions[i].category
        // ).subCategory;
        trans.accName = accs.find(
          x => x.account_id === trans.account_id
        ).name;
      }
    },
    updateTransaction(state, transaction) {
      let transSet: any[] = state.transactions;
      let transToUpdate = transSet.find(x => x.id === transaction.id);
      transToUpdate.category = transaction.category;
      transToUpdate.catName = transaction.catName;
    },
    updateCategories(state, categories) {
      state.categories = categories;
    //   var topArr = [];
    //   state.categories.forEach((obj) => {
    //     topArr.push(obj.topCategory);
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
    setTxData(state, value) {
      state.txData = value;
    },
    setCustomRange(state, value) {
      state.customStart = value.start;
      state.customEnd = value.end;
    },

  },
  actions: {
    async getTransactions({ commit }) {
      let res = await api.getTransactions();
      commit('updateTransactions', res);
    },
    async getCategories({ commit }) {
      let res = await api.getCategories();
      commit('updateCategories', res);
    },
    async getAccounts({ commit }) {
      let res = await api.getAccounts();
      commit('updateAccounts', res);
    },
    async customFilter({ commit }, {startFromPage, endFromPage}) {
      commit('setCustomRange', {start:startFromPage, end:endFromPage});
      this.state.webWorkerType = 'custom';
      workerActions.postMessage(this.state);
      // await buildTxData(this.state, 'customTxSet', start, end);
      // commit('setReloading', true);
    },
    async reanalyze({ commit }) {
      this.state.webWorkerType = 'full';
      workerActions.postMessage(this.state);
      // commit('analysisDataInitial');
      // commit('analysisData');
      // commit('setReloading');
      // let res = await api.getAccounts();
      // commit('updateAccounts', res);
    },
    async getAll({ commit }) {
      Promise.all([
        api.getCategories(),
        api.getAccounts(),
        api.getTransactions(),
        api.getItemTokens()
      ]).then(([cats, accs, trans, toks]) => {
        // ]).then((values) => {
        // console.log(values);
        // this.state.customEnd = moment();
        this.state.customEnd = moment().format('YYYY-MM-DD');
        // this.state.customStart = moment();
        this.state.customStart = moment().subtract(29, 'days').format('YYYY-MM-DD');
        // this.state.customStart = this.state.customStart.format('YYYY-MM-DD');
        commit('updateCategories', cats);
        commit('updateAccounts', accs);
        commit('updateTransactions', trans);
        commit('updateItemTokens', toks);
        this.state.webWorkerType = 'initial';
        workerActions.postMessage(this.state);
        // commit('analysisDataInitial');
        // commit('doneLoading', true);

        //Send this to worker as well
        // commit('analysisData');


      // }).then(() => {
          // this.state.webWorkerType = 'afterInitial';
          // workerActions.postMessage(this.state);
      }).catch(e => console.error(e));
      // // let res2 = await api.getCategories();
      // commit('updateCategories', await api.getCategories());
      // // let res3 = await api.getAccounts();
      // commit('updateAccounts', await api.getAccounts());
      // // let res1 = await api.getTransactions();
      // commit('updateTransactions', await api.getTransactions());
      // commit('analysisDataInitial');
      // commit('doneLoading', true);
      // commit('analysisData');
    }

  },
});

workerActions.onmessage = e => {
  store.commit(e.data.type, e.data.payload);
};

export default store;