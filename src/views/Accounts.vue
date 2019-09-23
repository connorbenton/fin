<template>
    <v-container>
        <v-layout column>
            <h1 class="title my-3">My Accounts</h1>
                <v-flex mt-4>
                    <v-file-input v-model="files" label="Choose CSV from mint.com to import"></v-file-input>
                    <v-btn 
                    :disabled="dialog"
                    :loading="dialog"
                    @click="dialog = true">Import CSV</v-btn>
                </v-flex>
            <v-dialog
              v-model="dialog"
              persistent
              width="300"
            >
            <v-card
              color="primary"
              dark
            >
              <v-card-text>
                Importing CSV
                <v-progress-linear
                  indeterminate
                  color="white"
                  class="mb-0"
                ></v-progress-linear>
              </v-card-text>
            </v-card>
          </v-dialog>
        </v-layout>
            <section>
        <plaid-link
            env="sandbox"
            :publicKey="PLAID_PUBLIC_KEY"
            clientName="Test App"
            product="transactions"
            v-bind="{ onSuccess }">
            Open Plaid Link
        </plaid-link>
    </section>
   <v-card
    class="mx-auto"
    max-width="400"
    tile
   > 
    <v-list-item>
    </v-list-item>
  </v-card>
    </v-container>
</template>

<script>
import PlaidLink from 'vue-plaid-link'
import api from '@/api'
const Papa = require('papaparse');
  export default {
    data() {
        return {
            dialog: false,
            transactions: [],
            files:null,
            PLAID_PUBLIC_KEY: process.env.VUE_APP_PLAID_PUBLIC_KEY,
            connections: [],
            loading: false,
        }
    },
    created () {
      //console.log(this.PLAID_PUBLIC_KEY)
      this.refreshConnections()
    },
    components: { PlaidLink },
    methods: {
        onSuccess (token) {
            console.log(token)
        },
        async refreshConnections () {
          this.loading = true
          this.connections = await api.getSaltEdgeConnections()
          console.log(JSON.stringify(this.connections))
          this.loading = false
        },
    },
    watch: {
      dialog(val) {
        if (!val) return
        var self = this
        const that = this.files
        console.log(that)
        //console.log(that)
          Papa.parse(that, {
            header: true,
            transformHeader: function(h) {
              var f = h.replace(/\s/g, '')
              var i = f.charAt(0).toLowerCase() + f.slice(1)
              return i
            },
            //step: function (results, parser) {
            //  var upload = JSON.stringify(results.data, null, 2)
            //  if(results.data.amount) {
            //    api.createTransaction(results.data)
            //    console.log(results.data)
            //  }
            //},
            complete (results) {
              self.dialog = true
              that.doc = JSON.stringify(results.data, null, 2)
              //console.log(that.doc)
              //console.log(results.data)
              for (let i in results.data) {
                //console.log(results.data[i])
                //setTimeout(function timer(){
                  //console.log(results.data[i])
                  //api.createTransaction(results.data[i])
                //}, i*40)
              }
              api.createTransactionBulk(results.data)
              self.dialog = false
            },
            error (errors) {
              console.log('error', errors)
            }
          })
      }
    }
  }
</script>

<style scoped></style>