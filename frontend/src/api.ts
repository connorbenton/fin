import axios from 'axios';
import store from './store';
// import io from 'socket.io-client';
// import rateLimit from 'axios-rate-limit'
import { EventEmitter } from 'events';

// class Socket {
//   public ee: EventEmitter;
//   public ws: WebSocket;
//   constructor(wsurl: any, ee = new EventEmitter()) {
//       const ws = new WebSocket(wsurl);
//       this.ee = ee;
//       this.ws = ws;
//       ws.onmessage = this.message.bind(this);
//       ws.onopen = this.open.bind(this);
//       ws.onclose = this.close.bind(this);
//   }
//   public on(this: Socket, name: any, fn: any) {
//       this.ee.on(name, fn);
//   }
//   public off(name: any, fn: any) {
//       this.ee.removeListener(name, fn);
//   }
//   public emit(name: any, data: any) {
//       const message = JSON.stringify({name, data});
//       this.ws.send(message);
//   }
//   public message(e: any) {
//       try {
//           const msgData = JSON.parse(e.data);
//           this.ee.emit(msgData.name, msgData.data);
//       } catch (err) {
//           const error = {
//               message: err,
//           };
//           console.log(err);
//           this.ee.emit(error.message);
//       }
//   }
//   public open() {
//       this.ee.emit('connected');
//   }
//   public close() {
//       this.ee.emit('disconnected');
//   }
// }

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

    // console.log(catres);
    // console.log(identifiedAccounts);

    const importData: any = {};
    importData.catres = catres;
    importData.identifiedAccounts = identifiedAccounts;
    importData.transactions = data;

    // const resImport = await this.execute('post', `/api/importTransactions`, data);
    await this.execute('post', `/api/importTransactions`, importData);
    return;
  },
  // async importTransactions(data: any) {
  //   // const socket = new Socket('ws://localhost/ws');
  //   // const ioClient = io();
  //   // socket.on('connected', () => { console.log('Connected'); });
  //   // const socket = new Socket('ws://fintrack-go:6060/ws');
  //   try {
  //   // ioClient.on('compare', (compareSet: any, fn: any) => {
  //   socket.on('compare', (compareSet: any, fn: any) => {
  //       if (compareSet.type === 'trans') {
  //         store.state.comparematch = true;
  //         store.state.trans1 = compareset.trans1;
  //         store.state.trans2 = compareset.trans2;
  //         store.subscribe((mutation, state) => {
  //           if (mutation.type === 'answergiven') {
  //             fn(mutation.payload);
  //             store.state.comparematch = false;
  //           }
  //         });
  //       }
  //       if (compareSet.type === 'cats') {
  //         store.commit('newAssignCats', compareSet);
  //         store.subscribe((mutation, state) => {
  //           if (mutation.type === 'assignDone') {
  //             fn(mutation.payload);
  //           }
  //         });
  //       }
  //     });
  //   const res: any = await client.post(`/api/importTransactions`, data);
  //   //  .then( ioClient.on('compare', function(data) {
  //     // store.state.compareMatch = true;
  //     // store.state.trans1 = data.trans1;
  //     // store.state.trans2 = data.trans2;
  //     // console.log(data);
  //   //  })
  //   // ).then(function (response: any) {
  //   // ioClient.off('compare');
  //   socket.off('compare');
  //     // console.log(response);
  //   return res;
  //   } catch (error) {
  //     // console.error(error);
  //   }
  //   // return this.execute('post', '/api/importTransactions', data)
  // },
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
    // const ioClient = io();
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
      //  ioClient.on('check', (data: any, fn: any) => {
      //   store.subscribe((mutation, state) => {
      //     if (mutation.type === 'newName') {
      //       fn();
      //     }
      //   });
      // store.state.fetchTransactionsItemDone = data.curr;
      // store.state.fetchTransactionsItemTotal = data.len;
      // store.commit('newName', data.name);
      // store.state.currName = data.name;
      // console.log(data.curr);
      //  });
      const res: any = await client.get(`/api/itemTokensFetchTransactions`);
      // ).then(function (response) {
      //  ioClient.off('check');
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
  // getSaltEdgeConnections() {
  //   // setTimeout(() => {store.state.isFetchTransactions = false}, 2000);
  //   return this.execute('get', '/api/saltEdgeConnections');
  // },
  refreshInteractive(id: any) {
    return this.execute('get', `/api/saltEdgeRefreshInteractive/${id}`);
  },
  createInteractive(id: any) {
    return this.execute('get', `/api/saltEdgeCreateInteractive/`);
  },
  resetDB() {
    return this.execute('get', `/api/resetDB`);
  },
  // resetToken() {
  //   return this.execute('get', `/api/resetToken`);
  // },
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
