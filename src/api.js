import Vue from 'vue'
import axios from 'axios'
//import rateLimit from 'axios-rate-limit'

const client = axios.create({
//  maxRequests: 50,
//  perMilliseconds: 10,
  json: true
})

export default {
  async execute (method, resource, data) {
    // inject the accessToken for each request
    //let accessToken = await Vue.prototype.$auth.getAccessToken()
    return client({
      method,
      url: resource,
      data,
      //headers: {
      //  Authorization: `Bearer ${accessToken}`
      //}
    }).then(req => {
      return req.data
    })
  },
  getTransactions () {
    return this.execute('get', '/api/transactions')
  },
  getTransaction (id) {
    return this.execute('get', `/api/transactions/${id}`)
  },
  createTransaction (data) {
    return this.execute('post', '/api/transactions', data)
  },
  createTransactionBulk (data) {
    return this.execute('post', '/api/transactionsBulk', data)
  },
  updateTransaction (id, data) {
    return this.execute('put', `/api/transactions/${id}`, data)
  },
  deleteTransaction (id) {
    return this.execute('delete', `/api/transactions/${id}`)
  },
  getCategories () {
    return this.execute('get', '/api/categories')
  },
  getCategory (id) {
    return this.execute('get', `/api/categories/${id}`)
  },
  createCategory (data) {
    return this.execute('post', '/api/categories', data)
  },
  updateCategory (id, data) {
    return this.execute('put', `/api/categories/${id}`, data)
  },
  deleteCategory (id) {
    return this.execute('delete', `/api/categories/${id}`)
  },
  getAccounts () {
    return this.execute('get', '/api/accounts')
  },
  getAccount (id) {
    return this.execute('get', `/api/accounts/${id}`)
  },
  createAccount (data) {
    return this.execute('post', '/api/accounts', data)
  },
  updateAccount (id, data) {
    return this.execute('put', `/api/accounts/${id}`, data)
  },
  deleteAccount (id) {
    return this.execute('delete', `/api/accounts/${id}`)
  },
  getProviderTokens () {
    return this.execute('get', '/api/providerTokens')
  },
  getProviderToken (id) {
    return this.execute('get', `/api/providerTokens/${id}`)
  },
  createProviderToken (data) {
    return this.execute('post', '/api/providerTokens', data)
  },
  updateProviderToken (id, data) {
    return this.execute('put', `/api/providerTokens/${id}`, data)
  },
  deleteProviderToken (id) {
    return this.execute('delete', `/api/providerTokens/${id}`)
  },
  getSaltEdgeConnections () {
    return this.execute('get', '/api/saltEdgeConnections')
  },
}