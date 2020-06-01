<template>
<div>
  <div class="container-fluid mt-4">
    <h1 class="h1">Posts Manager</h1>
    <v-alert :show="loading" variant="info">Loading...</v-alert>
    <v-row>
      <v-col>
        <table class="table table-striped">
          <thead>
            <tr>
              <th>ID</th>
              <th>Title</th>
              <th>Updated At</th>
              <th>&nbsp;</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="post in posts" :key="post.id">
              <td>{{ post.id }}</td>
              <td>{{ post.title }}</td>
              <td>{{ post.updated_at }}</td>
              <td class="text-right">
                <a href="#" @click.prevent="populatePostToEdit(post)">Edit</a> -
                <a href="#" @click.prevent="deletePost(post.id)">Delete</a>
              </td>
            </tr>
          </tbody>
        </table>
      </v-col>
      <v-col lg="3">
        <v-card :title="(model.id ? 'Edit Post ID#' + model.id : 'New Post')">
          <form @submit.prevent="savePost">
            <v-form label="Title">
              <v-text-field v-model="model.title" label="name"></v-text-field>
            </v-form>
            <v-form label="Body">
              <v-text-field rows="4" v-model="model.body" label="body"></v-text-field>
            </v-form>
            <div>
              <v-btn type="submit" variant="success">Save Post</v-btn>
            </div>
          </form>
        </v-card>
      </v-col>
    </v-row>
  </div>
  <div class="container-fluid mt-4">
    <h1 class="h1">Accounts Manager</h1>
    <v-alert :show="loading" variant="info">Loading...</v-alert>
    <v-row>
      <v-col>
        <table class="table table-striped">
          <thead>
            <tr>
              <th>ID</th>
              <th>Title</th>
              <th>Updated At</th>
              <th>&nbsp;</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="account in accounts" :key="account.id">
              <td>{{ account.id }}</td>
              <td>{{ account.title }}</td>
              <td>{{ account.updated_at }}</td>
              <td class="text-right">
                <a href="#" @click.prevent="populateAccountToEdit(account)">Edit</a> -
                <a href="#" @click.prevent="deleteAccount(account.id)">Delete</a>
              </td>
            </tr>
          </tbody>
        </table>
      </v-col>
      <v-col lg="3">
        <v-card :title="(model.id ? 'Edit Account ID#' + model.id : 'New Account')">
          <form @submit.prevent="saveAccount">
            <v-form label="Title">
              <v-text-field v-model="model.title" label="name"></v-text-field>
            </v-form>
            <v-form label="Body">
              <v-text-field rows="4" v-model="model.body" label="body"></v-text-field>
            </v-form>
            <div>
              <v-btn type="submit" variant="success">Save Account</v-btn>
            </div>
          </form>
        </v-card>
      </v-col>
    </v-row>
  </div>
  </div>
</template>

<script>
import api from '@/api'
export default {
  data () {
    return {
      loading: false,
      posts: [],
      accounts: [],
      model: {}
    }
  },
  async created () {
    this.refreshPosts()
    this.refreshAccounts()
  },
  methods: {
    async refreshAccounts () {
      this.loading = true
      this.accounts = await api.getAccounts()
      this.loading = false
    },
    async populateAccountToEdit (account) {
      this.model = Object.assign({}, account)
    },
    async saveAccount () {
      if (this.model.id) {
        await api.updateAccount(this.model.id, this.model)
      } else {
        await api.createAccount(this.model)
        console.log(JSON.stringify(this.model, 2, null))
      }
      this.model = {} // reset form
      await this.refreshAccounts()
    },
    async deleteAccount (id) {
      if (confirm('Are you sure you want to delete this account?')) {
        // if we are editing a account we deleted, remove it from the form
        if (this.model.id === id) {
          this.model = {}
        }
        await api.deleteAccount(id)
        await this.refreshAccounts()
      }
    },
    async refreshPosts () {
      this.loading = true
      this.posts = await api.getPosts()
      this.loading = false
    },
    async populatePostToEdit (post) {
      this.model = Object.assign({}, post)
    },
    async savePost () {
      if (this.model.id) {
        await api.updatePost(this.model.id, this.model)
      } else {
        await api.createPost(this.model)
      }
      this.model = {} // reset form
      await this.refreshPosts()
    },
    async deletePost (id) {
      if (confirm('Are you sure you want to delete this post?')) {
        // if we are editing a post we deleted, remove it from the form
        if (this.model.id === id) {
          this.model = {}
        }
        await api.deletePost(id)
        await this.refreshPosts()
      }
    }
  }
}
</script>