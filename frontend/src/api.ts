import axios from 'axios';
import store from './store';
import io from 'socket.io-client';
// import rateLimit from 'axios-rate-limit'

const client = axios.create({
//  maxRequests: 50,
//  perMilliseconds: 10,
  timeout: 60 * 4 * 1000, // wait 4 min for the long import calls
  // json: true
});

export default {
  async execute(method: any, resource: any, data: any = '') {
    // inject the accessToken for each request
    // let accessToken = await Vue.prototype.$auth.getAccessToken()
    return client({
      method,
      url: resource,
      data,
      // headers: {
      //  Authorization: `Bearer ${accessToken}`
      // }
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
    const ioClient = io();
    try {
    ioClient.on('compare', (compareSet: any, fn: any) => {
        if (compareSet.type === 'trans') {
          store.state.compareMatch = true;
          store.state.trans1 = compareSet.trans1;
          store.state.trans2 = compareSet.trans2;
          store.subscribe((mutation, state) => {
            if (mutation.type === 'answerGiven') {
              fn(mutation.payload);
              store.state.compareMatch = false;
            }
          });
        }
        if (compareSet.type === 'cats') {
          store.commit('newAssignCats', compareSet);
          store.subscribe((mutation, state) => {
            if (mutation.type === 'assignDone') {
              fn(mutation.payload);
            }
          });
        }
      });
    const res: any = await client.post(`/api/importTransactions`, data);
    //  .then( ioClient.on('compare', function(data) {
      // store.state.compareMatch = true;
      // store.state.trans1 = data.trans1;
      // store.state.trans2 = data.trans2;
      // console.log(data);
    //  })
    // ).then(function (response: any) {
    ioClient.off('compare');
      // console.log(response);
    return res;
    } catch (error) {
      // console.error(error);
    }
    // return this.execute('post', '/api/importTransactions', data)
  },
  updateTransaction(id: any, data: any) {
    // return this.execute('put', `/api/transactions/${id}`, data)
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
    // let ioClient = io("https://192.168.2.2:3000");
    const ioClient = io();
    //   reconnectionDelay: 1000,
    //   reconnection: true,
    //   reconnectionAttemps: 10,
    //   // transports: ['websocket'],
    //   agent: false,
    //   upgrade: false,
    //   rejectUnauthorized: false
    // });
    // ioClient.on("seq-num", (msg) => console.info(msg));
    // client.get(`/api/itemTokensFetchTransactions`)
    //  .then(x => x.request.response).then(store.state.isFetchTransactions = false).catch(error => error);

    store.commit('isFetch', true);
    // store.state.isFetchTransactions = true;
    // ioClient.on('check', function(data) {
    //   store.state.fetchTransactionsItemDone = data.curr;
    //   store.state.fetchTransactionsItemTotal = data.len;
    //   console.log(data.curr);
    // });
    // let config = {
      //   onDownloadProgress: progressEvent => {
      //     const dataChunk = progressEvent.currentTarget.response;
      //     store.state.fetchTransactionsItemDone = dataChunk.curr;
      //     store.state.fetchTransactionsItemTotal = dataChunk.len;
      //     console.log(store.state);
      //     console.log(dataChunk);
      //  }
    //  }
    try {
     ioClient.on('check', (data: any, fn: any) => {
      store.subscribe((mutation, state) => {
        if (mutation.type === 'newName') {
          fn();
        }
      });
      // store.state.fetchTransactionsItemDone = data.curr;
      // store.state.fetchTransactionsItemTotal = data.len;
      store.commit('newName', data.name);
      // store.state.currName = data.name;
      // console.log(data.curr);
     });
     const res: any = await client.get(`/api/itemTokensFetchTransactions`);
    // ).then(function (response) {
     ioClient.off('check');
     store.commit('isFetch', false);
      // store.state.isFetchTransactions = false;
     store.dispatch('getAll');
     return res;
      // console.log(response);
    } catch (error) {
      // console.error(error);
    }
    // }).catch(error => error);
    //  .then(x => x.request.response).then(store.state.isFetchTransactions = false).catch(error => error);

    // return this.execute('get', `/api/itemTokensFetchTransactions`)
    // return x.request.response
  },
  createAccount(data: any) {
    return this.execute('post', '/api/accounts', data);
  },
  updateAccount(id: any, data: any) {
    return this.execute('put', `/api/accounts/${id}`, data);
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
  getSaltEdgeConnections() {
    // setTimeout(() => {store.state.isFetchTransactions = false}, 2000);
    return this.execute('get', '/api/saltEdgeConnections');
  },
  refreshInteractive(id: any) {
    return this.execute('get', `/api/saltEdgeRefreshInteractive/${id}`);
  },
  createInteractive(id: any) {
    return this.execute('get', `/api/saltEdgeCreateInteractive/`);
  },
  resetDB() {
    return this.execute('get', `/api/resetDB/`);
  },
};
