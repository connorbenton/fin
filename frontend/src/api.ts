import axios from 'axios';
import store from './store';
import { EventEmitter } from 'events';

const client = axios.create({
  timeout: 60 * 4 * 1000, // wait 4 min for the long import calls
});

export default {
  async execute(method: any, resource: any, data: any = '') {
    return client({
      method,
      url: resource,
      data,
    }).then((req) => {
      return req.data;
    });
  },
  getTransactions() {
    return this.execute('get', '/api/transactions');
  },
  getTransaction(id: any) {
    return this.execute('get', `/api/transactions/${id}`);
  },
  createTransaction(data: any) {
    return this.execute('post', '/api/transactions', data);
  },
  async importTransactions(data: any) {
    const bus = new EventEmitter();
    let lock = false;
    const identifiedAccounts: any = [];
    const res = await this.execute('post', `/api/checkTransactions`, data);
    // console.log(res.trans)
    if (res.trans.transSets != null) {
      for (const compareSet of res.trans.transSets) {
        const matchingAccount = identifiedAccounts.find((x: any) => x.refAccountID === compareSet.trans2.account_id);
        if (matchingAccount == null) {
          store.state.trans1 = compareSet.trans1;
          store.state.trans2 = compareSet.trans2;
          store.state.compareMatch = true;
          let unsubscribeloop: any = null;
          lock = true;
          unsubscribeloop = await store.subscribe((mutation, state) => {
            if (mutation.type === 'answerGiven') {
              compareSet.isMatch = store.state.compareAnswer;
              if (compareSet.isMatch) {
                const matchIndex: any = {};
                matchIndex.importKey = compareSet.trans1.account_name;
                matchIndex.refAccountID = compareSet.trans2.account_id;
                matchIndex.refAccountName = compareSet.trans2.account_name;
                identifiedAccounts.push(matchIndex);
              }
              unsubscribeloop();
              lock = false;
              bus.emit('unlocked');
              store.state.compareMatch = false;
            }
          });
          if (lock) {
            await new Promise((resolve) => bus.once('unlocked', resolve));
          }
        }
      }
    }
    const catset: any = {};
    let catres: any = [];

    if (res.cats.catsToID != null) {
      catset.dbCats = store.state.categories;
      catset.compareCats = res.cats.catsToID;
      store.commit('newAssignCats', catset);

      let unsubscribe: any = null;
      lock = true;
      unsubscribe = await store.subscribe((mutation, state) => {
        if (mutation.type === 'assignDone') {
          catres = mutation.payload;
          unsubscribe();
          lock = false;
          bus.emit('unlocked');
        }
      });
      if (lock) {
        await new Promise((resolve) => bus.once('unlocked', resolve));
      }
    }

    const importData: any = {};
    importData.catres = catres;
    importData.identifiedAccounts = identifiedAccounts;
    importData.transactions = data;

    await this.execute('post', `/api/importTransactions`, importData);
    return;
  },

  updateTransaction(id: any, data: any) {
    return this.execute('put', `/api/transactions`, data);
  },
  deleteTransaction(id: any) {
    return this.execute('delete', `/api/transactions/${id}`);
  },
  getCategories() {
    return this.execute('get', '/api/categories');
  },
  getCategory(id: any) {
    return this.execute('get', `/api/categories/${id}`);
  },
  createCategory(data: any) {
    return this.execute('post', '/api/categories', data);
  },
  updateCategory(id: any, data: any) {
    return this.execute('put', `/api/categories/${id}`, data);
  },
  deleteCategory(id: any) {
    return this.execute('delete', `/api/categories/${id}`);
  },
  getAccounts() {
    return this.execute('get', '/api/accounts');
  },
  getAccount(id: any) {
    return this.execute('get', `/api/accounts/${id}`);
  },
  async fetchTransactions() {
    store.commit('isFetch', true);

    try {
      const res: any = await client.get(`/api/itemTokensFetchTransactions`);
      store.commit('isFetch', false);
      store.dispatch('getAll');
      return res;
    } catch (error) {
      // console.error(error);
    }

  },
  upsertAccountIgnore(data: any) {
    return this.execute('post', `/api/accountUpsertIgnore`, data);
  },
  upsertAccountName(data: any) {
    return this.execute('post', `/api/accountUpsertName`, data);
  },
  upsertTransaction(data: any) {
    return this.execute('post', `/api/transactionUpsert`, data);
  },
  deleteAccount(id: any) {
    return this.execute('delete', `/api/accounts/${id}`);
  },
  getItemTokens() {
    return this.execute('get', '/api/itemTokens');
  },
  getItemToken(id: any) {
    return this.execute('get', `/api/itemTokens/${id}`);
  },
  plaidCreateItemToken(data: any) {
    return this.execute('post', '/api/plaidItemTokens', data);
  },
  plaidGeneratePublicToken(data: any) {
    return this.execute('post', '/api/plaidGeneratePublicToken', data);
  },
  updateItemToken(id: any, data: any) {
    return this.execute('put', `/api/itemTokens/${id}`, data);
  },
  deleteItemToken(id: any) {
    return this.execute('delete', `/api/itemTokens/${id}`);
  },
  getSaltEdgeCategories() {
    return this.execute('get', '/api/saltEdge_Categories');
  },
  getPlaidCategories() {
    return this.execute('get', '/api/plaid_Categories');
  },
  refreshInteractive(id: any) {
    return this.execute('get', `/api/saltEdgeRefreshInteractive/${id}`);
  },
  createInteractive(id: any) {
    return this.execute('get', `/api/saltEdgeCreateInteractive/`);
  },
  resetDB() {
    return this.execute('get', `/api/resetDB`);
  },
  resetDBFull() {
    return this.execute('get', `/api/resetDBFull`);
  },
  getTrees() {
    return this.execute('get', `/api/analysisTrees`);
  },
  customAnalyze(data: any) {
    return this.execute('post', `/api/customTree`, data);
  },
};
