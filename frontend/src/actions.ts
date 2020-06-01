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

function buildTxData(stateObj: any, key: string, start: string = '', end: string = '') {
  return new Promise((resolve, reject) => {
    try {
      const setKey = key;
      const treeKey = key.toString().substring(0, key.length - 3) + 'Tree';
      // this.categories = this.$store.getters.getAllCategories;
      // this.accounts = this.$store.getters.getAllAccounts;
      // this.categories = await api.getCategories();
      // this.accounts = await api.getAccounts();
      const cats = stateObj.categories.map((c: any) => Object.assign({}, c));
      const accs = stateObj.accounts.map((ac: any) => Object.assign({}, ac));

      let a = moment();
      switch (setKey) {
        case 'last30TxSet':
          stateObj.txData.txMatrix.rangeDays[setKey] = 30;
          stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
            const b = moment(x.date, 'YYYY-MM-DD', true);
            const c = a.diff(b, 'days');
            return (c < 30);
          });
          break;
        case 'thisMonthTxSet':
          stateObj.txData.txMatrix.rangeDays[setKey] = moment().diff(moment().startOf('month'), 'days') + 1;
          stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
            const b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'month');
            // let c = a.diff(b, 'days');
            // return (c < 30);
          });
          break;
        case 'lastMonthTxSet':
          a = a.subtract(1, 'month');
          stateObj.txData.txMatrix.rangeDays[setKey] =
            a.clone().endOf('month').diff(a.clone().startOf('month'), 'days') + 1;
          stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
            const b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'month');
          });
          break;
        case 'lastSixMonthsTxSet':
          stateObj.txData.txMatrix.rangeDays[setKey] = moment().diff(moment().subtract(6, 'month'), 'days') + 1;
          stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
            const b = moment(x.date, 'YYYY-MM-DD', true);
            const c = a.diff(b, 'months');
            return (c < 6);
          });
          break;
        case 'thisYearTxSet':
          stateObj.txData.txMatrix.rangeDays[setKey] = moment().diff(moment().startOf('year'), 'days') + 1;
          stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
            const b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'year');
          });
          break;
        case 'lastYearTxSet':
          a = a.subtract(1, 'year');
          stateObj.txData.txMatrix.rangeDays[setKey] =
            a.clone().endOf('year').diff(a.clone().startOf('year'), 'days') + 1;
          stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
            const b = moment(x.date, 'YYYY-MM-DD', true);
            return b.isSame(a, 'year');
          });
          break;
        case 'fromBeginningTxSet':
          let oldestTrans = stateObj.transactions.map((t: any) => Object.assign({}, t));
          if (oldestTrans.length > 0) {
            oldestTrans = oldestTrans.reduce((y: any, z: any) => (moment(y.date) < moment(z.date)) ? y : z);
          }
          stateObj.txData.oldestDate = oldestTrans.date;
          stateObj.txData.txMatrix.rangeDays[setKey] = moment().diff(moment(oldestTrans.date), 'days') + 1;
          stateObj.txData.txSets[setKey] = stateObj.transactions;
          break;
        case 'customTxSet':
          if (start === '') {
            stateObj.txData.txSets[setKey] = stateObj.txData.txSets.last30TxSet;
            stateObj.txData.txMatrix.rangeDays[setKey] = 30;
          } else {
            const st = moment(start, 'YYYY-MM-DD', true);
            const en = moment(end, 'YYYY-MM-DD', true);
            stateObj.txData.txMatrix.rangeDays[setKey] = en.diff(st, 'days') + 1;
            stateObj.txData.txSets[setKey] = stateObj.transactions.filter((x: any) => {
              const c = moment(x.date, 'YYYY-MM-DD', true);
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
      const txSet = stateObj.txData.txSets[setKey];
      for (const set of txSet) {
        set.account_name = accs.find(
          (x: any) => x.account_id === set.account_id,
        ).name;
        const matchCat = cats.find((x: any) => x.id === set.category);
        set.category_name = matchCat.sub_category;
        if (matchCat.count === undefined) { matchCat.count = 0; }
        matchCat.count = matchCat.count + 1;
        if (matchCat.total === undefined) { matchCat.total = 0; }
        matchCat.total =
          matchCat.total + parseFloat(set.normalized_amount);

        // matchCat.count = matchCat.count + 1;
        // matchCat.total = matchCat.total + parseFloat(txSet[i].amount);
      }

      const txTree = {
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
      for (const j in cats) {
        // if (cats[j].exclude_from_analysis || cats[j].top_category === "Income")
        // if (cats[j].exclude_from_analysis)
        // continue;
        if (cats[j].sub_category === cats[j].top_category) {
          const newChild = {
            name: cats[j].top_category,
            children: [{}],
            value: 0,
            count: 0,
            trueCount: 0,
            trueValue: 0,
            dbID: 0,
          };
          // newChild.name = cats[j].top_category;
          // newChild.children = [];
          // newChild.value = 0;
          // newChild.count = 0;
          // newChild.trueCount = 0;
          // newChild.trueValue = 0;
          // newChild.dbID = '';

          const children = cats.filter(
            (x: any) => x.top_category === cats[j].top_category,
          );
          for (const child of children) {
            const subCatChildToPush = {
              name: '',
              dbID: 0,
              value: 0,
              count: 0,
              trueValue: 0,
              trueCount: 0,
            };
            if (child.sub_category === child.top_category) {
              subCatChildToPush.name = child.sub_category + ` (General)`;
              newChild.dbID = child.id;
            } else {
              subCatChildToPush.name = child.sub_category;
            }

            subCatChildToPush.dbID = child.id;
            const value = child.total;
            const count = child.count;
            subCatChildToPush.value =
              value === undefined ? 0 : -1 * value;
            subCatChildToPush.count = count === undefined ? 0 : count;
            subCatChildToPush.trueValue = subCatChildToPush.value;
            subCatChildToPush.trueCount = subCatChildToPush.count;
            // subCatChildToPush.percent = "";
            // if (newChild.name === 'Income' || child.value < 0) {
            if (child.exclude_from_analysis || child.top_category === 'Income' || subCatChildToPush.value < 0) {
              // subCatChildToPush.trueCount = subCatChildToPush.count;
              subCatChildToPush.count = 0;
              // subCatChildToPush.trueValue = subCatChildToPush.value;
              subCatChildToPush.value = 0;
            }
            newChild.value = newChild.value + subCatChildToPush.value;
            newChild.trueValue = newChild.trueValue + subCatChildToPush.trueValue;
            newChild.count = newChild.count + subCatChildToPush.count;
            newChild.trueCount = newChild.trueCount + subCatChildToPush.trueCount;
            //  subCatChildToPush.count = child.count;
            newChild.children.push(subCatChildToPush);
          }
          // newChild.children.map(obj => ({...obj, percent: (obj.value / newChild.value).toFixed(1)+"%"}));
          newChild.children.forEach((child: any) => {
            // newChild.child.percent = '';
            // if (newChild.name === 'Income' || child.value < 0) {
            // child.trueValue = child.value;
            // child.value = 0;
            // }
            child.percent = ((child.value / newChild.value) * 100).toFixed(1) + '%';
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
      const total = txTree.value;
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
        child.percent = ((child.value / total) * 100).toFixed(1) + '%';
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

      const treeKeyNoInvest = treeKey + 'NoInvest';
      // let txTreeNoInvest = {...txTree};
      // Object.assign(txTreeNoInvest, txTree);
      const txTreeNoInvest = JSON.parse(JSON.stringify(txTree));
      // txTreeNoInvest.name = "Transactions by Category (Financial Excluded)"
      if (txTreeNoInvest.value) {
        const fin = txTreeNoInvest.children[7];
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
          });
          txTreeNoInvest.children.forEach((cat: any) => {
            cat.percent = ((cat.value / txTreeNoInvest.value) * 100).toFixed(1) + '%';
          });
        }
      }
      stateObj.txData.txTrees[treeKeyNoInvest] = txTreeNoInvest;
      stateObj.txData.txTrees[treeKey] = txTree;
      resolve();
    } catch (error) {
      reject(error);
    }
  });
  // let a = this.transactionsTree;
  // let b = "ha";
}

function buildTxArray(type: string) {
  const PromiseArray = [];
  if (type === 'initial') {
    PromiseArray.push(buildTxData(state, 'last30TxSet'));
  } else if (type === 'afterInitial') {
    Object.keys(state.txData.txSets).forEach((s) => {
      if (s === 'last30TxSet') { return; }
      PromiseArray.push(buildTxData(state, s));
    });
  } else if (type === 'full') {
    Object.keys(state.txData.txSets).forEach((s) => {
      // if (s === 'last30TxSet') return;
      PromiseArray.push(buildTxData(state, s));
    });
  } else if (type === 'custom') {
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
  const PromiseArray = buildTxArray(e.data.webWorkerType);
  Promise.all(PromiseArray).then(() => {
    ctx.postMessage({ type: 'setTxData', payload: state.txData });
    if (e.data.webWorkerType === 'initial') {
      ctx.postMessage({ type: 'doneLoading', payload: true });
      const SecondPromiseArray = buildTxArray('afterInitial');
      Promise.all(SecondPromiseArray).then(() => {
        ctx.postMessage({ type: 'setTxData', payload: state.txData });
        ctx.postMessage({ type: 'setReloading', payload: true });
      });
    } else { ctx.postMessage({ type: 'setReloading', payload: true }); }
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
