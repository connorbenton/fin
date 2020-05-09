// actions.js
import moment from 'moment';

let state = {
  txData: {
    txSets: {},
  },
  customStart: '',
  customEnd: '',
};
// let state = {
//     webWorkerType: '',
//     isFetchTransactions: false,
//     fetchTransactionsItemDone: 0,
//     fetchTransactionsItemTotal: 0,
//     currName: '',
//     compareMatch: false,
//     trans1: {},
//     trans2: {},
//     compareAnswer: '',
//     catsToCompare: [],
//     finishedCats: [],
//     txData: {
//       txSets: {
//         last30TxSet: [],
//         thisMonthTxSet: [],
//         lastMonthTxSet: [],
//         lastSixMonthsTxSet: [],
//         thisYearTxSet: [],
//         lastYearTxSet: [],
//         fromBeginningTxSet: [],
//         customTxSet: [],
//       },
//       txTrees: {
//         last30TxTree: [],
//         thisMonthTxTree: [],
//         lastMonthTxTree: [],
//         lastSixMonthsTxTree: [],
//         thisYearTxTree: [],
//         lastYearTxTree: [],
//         fromBeginningTxTree: [],
//         customTxTree: [],
//       },
//     },

//     transactions: [],

//     categories: [],

//     accounts: [],

//     apiStateLoaded: false,

//     doneReloading: false,
//   }

function buildTxData(state: any, key: string, start: string = '', end: string = '') {
  return new Promise((resolve, reject) => {
    try {
      let setKey = key;
      let treeKey = key.toString().substring(0, key.length - 3) + 'Tree';
      // this.categories = this.$store.getters.getAllCategories;
      // this.accounts = this.$store.getters.getAllAccounts;
      // this.categories = await api.getCategories();
      // this.accounts = await api.getAccounts();
      let cats = state.categories.map((c: any) => Object.assign({}, c));
      let accs = state.accounts.map((a: any) => Object.assign({}, a));

      let a = moment();
      switch (setKey) {
        case 'last30TxSet':
          state.txData.txMatrix.rangeDays[setKey] = 30;
          state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
            let b = moment(x.date, 'YYYY-MM-DD', true);
            let c = a.diff(b, 'days');
            return (c < 30);
          });
          break;
        case 'thisMonthTxSet':
          state.txData.txMatrix.rangeDays[setKey] = moment().diff(moment().startOf('month'), 'days') + 1;
          state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
            let b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'month');
            // let c = a.diff(b, 'days');
            // return (c < 30);
          });
          break;
        case 'lastMonthTxSet':
          a = a.subtract(1, 'month');
          state.txData.txMatrix.rangeDays[setKey] = a.clone().endOf('month').diff(a.clone().startOf('month'), 'days') + 1;
          state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
            let b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'month');
          });
          break;
        case 'lastSixMonthsTxSet':
          state.txData.txMatrix.rangeDays[setKey] = moment().diff(moment().subtract(6, 'month'), 'days') + 1;
          state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
            let b = moment(x.date, 'YYYY-MM-DD', true);
            let c = a.diff(b, 'months');
            return (c < 6);
          });
          break;
        case 'thisYearTxSet':
          state.txData.txMatrix.rangeDays[setKey] = moment().diff(moment().startOf('year'), 'days') + 1;
          state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
            let b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'year');
          });
          break;
        case 'lastYearTxSet':
          a = a.subtract(1, 'year');
          state.txData.txMatrix.rangeDays[setKey] = a.clone().endOf('year').diff(a.clone().startOf('year'), 'days') + 1;
          state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
            let b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'year');
          });
          break;
        case 'fromBeginningTxSet':
          let oldestTrans = state.transactions.map((a: any) => Object.assign({}, a));
          if (oldestTrans.length > 0) oldestTrans = oldestTrans.reduce((a: any, b: any) => (moment(a.date) < moment(b.date)) ? a : b);
          state.txData.oldestDate = oldestTrans.date;
          state.txData.txMatrix.rangeDays[setKey] = moment().diff(moment(oldestTrans.date), 'days') + 1;
          state.txData.txSets[setKey] = state.transactions;
          break;
        case 'customTxSet':
          if (start === '') {
            state.txData.txSets[setKey] = state.txData.txSets.last30TxSet;
            state.txData.txMatrix.rangeDays[setKey] = 30;
          }
          else {
            const st = moment(start, 'YYYY-MM-DD', true);
            const en = moment(end, 'YYYY-MM-DD', true);
            state.txData.txMatrix.rangeDays[setKey] = en.diff(st, 'days') + 1;
            state.txData.txSets[setKey] = state.transactions.filter((x: any) => {
              let c = moment(x.date, 'YYYY-MM-DD', true);
              return !(c.isBefore(st) || c.isAfter(en));
            });
          }
      }

      // let loadcats = this.$store.getters.getAllCategories;
      // let loadaccs = this.$store.getters.getAllAccounts;
      // let cats = loadcats.slice(0);
      // let accs = loadaccs.slice(0);
      // let y = this.$store.getters.getAllCategories;
      // this.transactions = await api.getTransactions();
      // this.transactions = [...this.$store.getters.getAllTransactions];
      let txSet = state.txData.txSets[setKey];
      for (let i in txSet) {
        txSet[i].accName = accs.find(
          (x: any) => x.account_id === txSet[i].account_id
        ).name;
        let matchCat = cats.find((x: any) => x.id === txSet[i].category);
        txSet[i].catName = matchCat.subCategory;
        if (typeof matchCat.count === "undefined") matchCat.count = 0;
        matchCat.count = matchCat.count + 1;
        if (typeof matchCat.total === "undefined") matchCat.total = 0;
        matchCat.total =
          matchCat.total + parseFloat(txSet[i].normalized_amount);

        // matchCat.count = matchCat.count + 1;
        // matchCat.total = matchCat.total + parseFloat(txSet[i].amount);
      }

      let txTree = {
        name: 'Transactions by Category',
        children: [{}],
        value: 0,
        count: 0,
        trueValue: 0,
        trueCount: 0,
      };

      // txTree.name = "Transactions by Category";
      // txTree.children = [];
      // txTree.value = 0;
      // txTree.count = 0;
      // txTree.trueValue = 0;
      // txTree.trueCount = 0;
      for (let j in cats) {
        // if (cats[j].excludeFromAnalysis || cats[j].topCategory === "Income")
        // if (cats[j].excludeFromAnalysis)
        // continue;
        if (cats[j].subCategory === cats[j].topCategory) {
          let newChild = {
            name: cats[j].topCategory,
            children: [{}],
            value: 0,
            count: 0,
            trueCount: 0,
            trueValue: 0,
            dbID: 0,
          };
          // newChild.name = cats[j].topCategory;
          // newChild.children = [];
          // newChild.value = 0;
          // newChild.count = 0;
          // newChild.trueCount = 0;
          // newChild.trueValue = 0;
          // newChild.dbID = '';

          let children = cats.filter(
            (x: any) => x.topCategory === cats[j].topCategory
          );
          for (let k in children) {
            let subCatChildToPush = {
              name: '',
              dbID: 0,
              value: 0,
              count: 0,
              trueValue: 0,
              trueCount: 0,
            };
            if (children[k].subCategory === children[k].topCategory) {
              subCatChildToPush.name = children[k].subCategory + ` (General)`;
              newChild.dbID = children[k].id;
            } else {
              subCatChildToPush.name = children[k].subCategory;
            }

            subCatChildToPush.dbID = children[k].id;
            let value = children[k].total;
            let count = children[k].count;
            subCatChildToPush.value =
              typeof value === "undefined" ? 0 : -1 * value;
            subCatChildToPush.count = typeof count === "undefined" ? 0 : count;
            subCatChildToPush.trueValue = subCatChildToPush.value;
            subCatChildToPush.trueCount = subCatChildToPush.count;
            // subCatChildToPush.percent = "";
            // if (newChild.name === 'Income' || child.value < 0) {
            if (children[k].excludeFromAnalysis || children[k].topCategory === 'Income' || subCatChildToPush.value < 0) {
              // subCatChildToPush.trueCount = subCatChildToPush.count;
              subCatChildToPush.count = 0;
              // subCatChildToPush.trueValue = subCatChildToPush.value;
              subCatChildToPush.value = 0;
            }
            newChild.value = newChild.value + subCatChildToPush.value;
            newChild.trueValue = newChild.trueValue + subCatChildToPush.trueValue;
            newChild.count = newChild.count + subCatChildToPush.count;
            newChild.trueCount = newChild.trueCount + subCatChildToPush.trueCount;
            //  subCatChildToPush.count = children[k].count;
            newChild.children.push(subCatChildToPush);
          }
          // newChild.children.map(obj => ({...obj, percent: (obj.value / newChild.value).toFixed(1)+"%"}));
          newChild.children.forEach(function (child: any) {
            // newChild.children[k].percent = '';
            // if (newChild.name === 'Income' || child.value < 0) {
            // child.trueValue = child.value;
            // child.value = 0;
            // }
            child.percent = ((child.value / newChild.value) * 100).toFixed(1) + "%";
            // let x = "y";
          });
          txTree.children.push(newChild);
          txTree.value =
            txTree.value + newChild.value;
          txTree.count =
            txTree.count + newChild.count;
          txTree.trueValue = txTree.trueValue + newChild.trueValue;
          txTree.trueCount = txTree.trueCount + newChild.trueCount;
        }
      }
      let total = txTree.value;
      // txTree.children = txTree.children.filter(child => child.count > 0 && child.value > 0);
      // txTree.children = txTree.children.filter(child => child.value > -1);
      // for (let index in txTree.children) {
      txTree.children.forEach((child: any) => {
        // if (child.name === 'Income' || child.value < 0) {
        //   // child.trueCount = child.count;
        //   child.count = 0;
        //   // child.trueValue = child.value;
        //   child.value = 0;
        // }
        // let index = txTree.children.indexOf(child);
        // let child = txTree.children[index]
        child.percent = ((child.value / total) * 100).toFixed(1) + "%";
        // if (child.value / total < 0.1) {
        // if (child.count < 1 || child.value < 0) {
        // txTree.children.splice(index, 1);
        // return;
        // }
        // if (setKey == 'last30TxSet' && child.name === 'Health & Fitness')
        // {
        //   console.log('h');
        // }
        // }
      });

      let treeKeyNoInvest = treeKey + 'NoInvest';
      // let txTreeNoInvest = {...txTree};
      // Object.assign(txTreeNoInvest, txTree);
      let txTreeNoInvest = JSON.parse(JSON.stringify(txTree))
      // txTreeNoInvest.name = "Transactions by Category (Financial Excluded)"
      if (txTreeNoInvest.value) {
        let fin = txTreeNoInvest.children[5];
        if (fin.value) {
          txTreeNoInvest.value = txTreeNoInvest.value - fin.value;
          txTreeNoInvest.count = txTreeNoInvest.count - fin.count;
          fin.trueValue = fin.value;
          fin.value = 0;
          fin.children.forEach((subcat: any) => {
            subcat.trueCount = subcat.count;
            subcat.count = 0;
            subcat.trueValue = subcat.value;
            subcat.value = 0;
            subcat.percent = 0;
          })
          txTreeNoInvest.children.forEach((cat: any) => {
            cat.percent = ((cat.value / txTreeNoInvest.value) * 100).toFixed(1) + "%";
          })
        }
      }
      state.txData.txTrees[treeKeyNoInvest] = txTreeNoInvest;
      state.txData.txTrees[treeKey] = txTree;
      resolve();
    }
    catch (error) {
      reject(error);
    }
  })
  // let a = this.transactionsTree;
  // let b = "ha";
}

function buildTxArray(type: string) {
  let PromiseArray = [];
  if (type === 'initial') {
    PromiseArray.push(buildTxData(state, 'last30TxSet'));
  }
  else if (type === 'afterInitial') {
    Object.keys(state.txData.txSets).forEach(s => {
      if (s === 'last30TxSet') return;
      PromiseArray.push(buildTxData(state, s));
    });
  }
  else if (type === 'full') {
    Object.keys(state.txData.txSets).forEach(s => {
      // if (s === 'last30TxSet') return;
      PromiseArray.push(buildTxData(state, s));
    });
  }
  else if (type === 'custom') {
    PromiseArray.push(buildTxData(state, 'customTxSet', state.customStart, state.customEnd));
  }
  return PromiseArray;
  // Promise.all(PromiseArray).then(() => {
  // self.postMessage({ type: 'SET_TX_DATA', payload: state.txData})
  // });

}

const ctx: Worker = self as any;
ctx.onmessage = (e: any) => {
  state = e.data;
  let PromiseArray = buildTxArray(e.data.webWorkerType);
  Promise.all(PromiseArray).then(() => {
    ctx.postMessage({ type: 'setTxData', payload: state.txData });
    if (e.data.webWorkerType === 'initial') {
      ctx.postMessage({ type: 'doneLoading', payload: true });
      let SecondPromiseArray = buildTxArray('afterInitial');
      Promise.all(SecondPromiseArray).then(() => {
        ctx.postMessage({ type: 'setTxData', payload: state.txData });
        ctx.postMessage({ type: 'setReloading', payload: true });
      })
    }
    else ctx.postMessage({ type: 'setReloading', payload: true });
  });
  // Perform the calculation
  // We can trigger any mutations from here!
  //   self.postMessage({ type: 'SET_WORKING', payload: true });

  //   const primes = calculatePrimes(400, 1000000000);
  //   self.postMessage({ type: 'SET_ITEMS', payload: primes });
  // self.postMessage({ type: 'SET_TX_DATA', payload: state.txData})

  // We can trigger any mutations from here!
  // Set the loading state back to false
  // self.postMessage({ type: 'SET_WORKING', payload: false });
};