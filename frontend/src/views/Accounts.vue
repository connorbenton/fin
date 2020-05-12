<template>
  <v-container>
    <v-dialog v-model="fetch" persistent width="60%" max-width="800">
      <v-card color="primary" dark>
        <v-card-text>
          {{dialogName}}
          <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-btn
      :loading="loading3"
      :disabled="loading3"
      color="success"
      @click.native="fetchTransactions()"
    >Fetch Transactions for All Accounts</v-btn>

    <h1 class="title mt-3">Plaid Connections</h1>
    <v-card class="d-inline-block mx-auto my-3" max-width="1000" tile :key="redraw">
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th>Name</th>
              <th>Last Transaction Fetch</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index1) in plaid(itemTokens)" :key="index1">
              <td>{{item.institution}}</td>
              <td>{{ localeDate(item.lastDownloadedTransactions) }}</td>
              <td v-if="item.needsReLogin == 1">
                <v-btn
                  color="warning"
                  :loading="plaidRefresh"
                  dark
                  @click.native="startReLogin(item.item_id)"
                >
                  Refresh Connection
                  <br />(needs New Login)
                </v-btn>
              </td>
              <td v-else>
                <v-icon color="success">mdi-check</v-icon>
              </td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
    <v-btn text class="d-block px-0">
      <plaid-link
        :env="environment"
        :publicKey="PLAID_PUBLIC_KEY"
        clientName="Test App"
        product="transactions"
        :token="updateToken"
        ref="plaidLinkRef"
        v-bind="{ onSuccess }"
      >
        <template slot="button" slot-scope="props">
          <v-btn @click="props.onClick">Open Plaid Link to Add New Account</v-btn>
        </template>
      </plaid-link>
    </v-btn>
    <h1 class="title mt-3">Salt Edge Connections</h1>
    <v-card class="d-inline-block mx-auto my-3" max-width="1000" tile :key="redraw2">
      <!-- <v-dialog v-model="dialog2" persistent>
                <v-card>
                <div id="wrapper" style="position:relative">
                  <iframe width=100%, :height="vh" :src="refreshUrl"></iframe>
                </div>
                <v-card-actions>
                <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="dialog2 = false">Close</v-btn>
                </v-card-actions>
                </v-card>
      </v-dialog>-->
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th>Name</th>
              <th>Last Transaction Fetch</th>
              <th>Last Refresh</th>
              <th>Next Refresh Available In</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index2) in saltEdge(itemTokens)" :key="index2">
              <td>{{item.institution}}</td>
              <td
                v-show="item.lastDownloadedTransactions !== null"
              >{{new Date(item.lastDownloadedTransactions).toLocaleString()}}</td>
              <td
                v-show="item.lastDownloadedTransactions == null"
              >{{"last transaction fetch: " + "Never"}}</td>
              <td v-if="item.interactive == 1">{{new Date(item.lastRefresh).toLocaleString()}}</td>
              <td v-else>
                <v-icon>mdi-minus</v-icon>
              </td>
              <td v-if="item.interactive == 1 && showRefresh">
                <v-btn
                  :loading="loading4"
                  :disabled="loading4"
                  color="primary"
                  dark
                  @click.native="startRefreshInteractive(item.item_id)"
                >
                  Refresh Connection
                  <br />(needs Credentials)
                </v-btn>
              </td>
              <td v-else-if="item.interactive == 1">{{ timeToRefresh(item.nextRefreshPossible) }}</td>
              <td v-else>
                <v-icon>mdi-minus</v-icon>
              </td>
              <!-- <td
                v-if="item.interactive == 1"
              >
              {{new Date(item.nextRefreshPossible).toLocaleString()}}</td>
              {{ timeToRefresh(item.nextRefreshPossible) }}</td>
              <td v-else><v-icon>mdi-minus</v-icon></td>-->
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
    <v-btn
      class="d-block"
      :loading="loading4"
      :disabled="loading4"
      @click.native="startCreateInteractive()"
    >Open Salt Edge To Add New Account</v-btn>
    <v-layout column>
      <h1 class="title mt-3">Import CSV from Mint.com</h1>
      <v-flex mt-4>
        <v-file-input accept=".csv" v-model="files" label="Choose CSV from mint.com to import"></v-file-input>
        <v-btn
          :loading="loading2"
          :disabled="loading2"
          @click.native="importTransactions()"
        >Import CSV</v-btn>
      </v-flex>
      <!-- <v-dialog v-model="dialog" persistent width="300">
        <v-card color="primary" dark>
          <v-card-text>
            Importing CSV
            <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>-->
    </v-layout>
    <v-btn color="warning" dark class="my-8" @click.native="resetDB()">Reset Database</v-btn>
  </v-container>
</template>

<script>
import PlaidLink from "vue-plaid-link";
import api from "@/api";
import moment from "moment";
// import VueFriendlyIframe from 'vue-friendly-iframe'
const Papa = require("papaparse");
// const waitFor = delay => new Promise(resolve => setTimeout(resolve, delay));

export default {
  data() {
    return {
      apiStateLoaded: false,
      showRefresh: false,
      fetch: false,
      dialogName: "Fetching Transactions",
      dialog: false,
      dialog2: false,
      vh: null,
      transactions: [],
      files: null,
      environment: process.env.VUE_APP_ENVIRONMENT || window._env_.VUE_APP_ENVIRONMENT,
      PLAID_PUBLIC_KEY: process.env.VUE_APP_PLAID_PUBLIC_KEY || window._env_.VUE_APP_PLAID_PUBLIC_KEY,
      updateToken: null,
      itemTokens: [],
      accounts: [],
      //connections: [],
      plaidRefresh: false,
      // loader: null,
      // loader2: null,
      // loader3: null,
      loading: false,
      loading2: false,
      loading3: false,
      loading4: false,
      refreshUrl: null,
      reactive: true,
      redraw: 0,
      redraw2: 1
    };
  },
  // computed: {
  //   dialogName: {
  //     get () {
  //       return this.$store.getters.getName
  //     }
  //   },
  //   fetch: {
  //     get () {
  //       return this.$store.state.isFetchTransactions
  //     }
  //   }
  // },
  created() {
    console.log(this.PLAID_PUBLIC_KEY);
    this.apiStateLoaded = this.$store.state.apiStateLoaded;
    if (this.apiStateLoaded) {
      this.initialData();
    }
    this.unsub = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "isFetch") {
        this.fetch = mutation.payload;
      }
      if (mutation.type === "newName") {
        this.dialogName = mutation.payload;
      }
      if (mutation.type === "doneLoading") {
        if (mutation.payload) {
          this.initialData();
        }
        this.apiStateLoaded = mutation.payload;
      }
    });

    // this.refreshData();
  },
  beforeDestroy() {
    this.unsub();
  },
  components: { PlaidLink },
  methods: {
    localeDate(date) {
      if (date == null) {
        return "Never";
      } else return new Date(date).toLocaleString();
    },
    timeToRefresh(time) {
      let diffmin = moment().diff(moment(time), "minutes");
      // let diffmin = moment().diff(moment("2020-05-07T15:33:49.000Z"), 'minutes');
      let diffsec = moment().diff(moment(time), "seconds");
      if (diffsec > 0) this.showRefresh = true;
      if (diffsec > -60) return diffsec * -1 + " seconds";
      else return diffmin * -1 + " minutes";
    },
    async startCreateInteractive() {
      var vm = this;
      vm.loading4 = true;
      vm.refreshUrl = await api.createInteractive();
      // console.log(vm.refreshUrl)
      // vm.vh = window.parent.innerHeight / 1.5
      // vm.vw = window.parent.innerWidth / 1.5
      // window.open("https://vuetifyjs.com", '_blank', 'toolbar=0,location=0,menubar=0,height=800,width=700')
      let win = window.open(
        vm.refreshUrl,
        "_blank",
        "toolbar=0,location=0,menubar=0,height=800,width=700"
      );
      let interval = setInterval(() => {
        // if (win.location.href === process.env.VUE_APP_BASE_URL) {
        try {
          if (win.document.domain === document.domain) {
            clearInterval(interval);
            vm.fetchRefreshData();
            win.close();
          }
        } catch (e) {
          if (win.closed) {
            clearInterval(interval);
            vm.fetchRefreshData();
            return;
          }
        }
      }, 500);
      vm.loading4 = false;
      // vm.dialog2 = true
    },
    async resetDB() {
      this.dialogName = "Resetting Database";
      this.fetch = true;
      await api.resetDB();
      this.fetch = false;
      this.dialogName = "Fetching Transactions";
      this.refreshData();
    },
    async startRefreshInteractive(id) {
      var vm = this;
      vm.loading4 = true;
      vm.refreshUrl = await api.refreshInteractive(id);
      // console.log(vm.refreshUrl)
      // vm.vh = window.parent.innerHeight / 1.5
      // vm.vw = window.parent.innerWidth / 1.5
      // window.open("https://vuetifyjs.com", '_blank', 'toolbar=0,location=0,menubar=0,height=800,width=700')
      let win = window.open(
        vm.refreshUrl,
        "_blank",
        "toolbar=0,location=0,menubar=0,height=800,width=700"
      );
      let interval = setInterval(() => {
        // if (win.location.href === process.env.VUE_APP_BASE_URL) {
        try {
          if (win.document.domain === document.domain) {
            clearInterval(interval);
            vm.fetchRefreshData();
            win.close();
          }
        } catch (e) {
          if (win.closed) {
            clearInterval(interval);
            vm.fetchRefreshData();
            return;
          }
        }
      }, 500);
      vm.loading4 = false;
      // vm.dialog2 = true
    },
    async onSuccess(token, metadata) {
      // console.log(token)
      // console.log(metadata)
      let TokenToUpload = {
        token: token,
        name: metadata.institution.name
      };
      this.dialogName =
        "Upserting account data for " + metadata.institution.name;
      this.fetch = true;
      await api.plaidCreateItemToken(TokenToUpload);
      this.fetch = false;
      this.dialogName = "Fetching Transactions";
      this.refreshData();
    },
    async startReLogin(id) {
      this.plaidRefresh = true;
      let ItemToUpload = {
        item_id: id
      };
      let tok = await api.plaidGeneratePublicToken(ItemToUpload);
      // console.log(this.updateToken);
      this.updateToken = tok.public_token;
      // let btn = this.$refs.plaidBtn;
      // console.log(this.PLAID_PUBLIC_KEY);
      // console.log(this.updateToken);
      // await waitFor(2000);
      // await this.$nextTick();
      await this.$nextTick();
      this.plaidRefresh = false;
      await this.$refs.plaidLinkRef.onScriptLoaded();
      this.$refs.plaidLinkRef.handleOnClick();
      // btn.click();
    },
    saltEdge: function(itemTokens) {
      return itemTokens.filter(itemToken => itemToken.provider == "SaltEdge");
    },
    plaid: function(itemTokens) {
      return itemTokens.filter(itemToken => itemToken.provider == "Plaid");
    },
    async initialData() {
      this.itemTokens = this.$store.getters.getAllItemTokens;
      this.accounts = this.$store.getters.getAllAccounts;
      this.redraw -= 1;
      this.redraw2 += 1;
    },
    async fetchRefreshData() {
      this.dialogName = "Checking for new SaltEdge connections to add";
      this.fetch = true;
      // await waitFor(2000);
      let res = await api.getSaltEdgeConnections();
      this.itemTokens = res.resTokens;
      this.accounts = res.resAccounts;
      this.fetch = false;
      this.dialogName = "Fetching Transactions";
      this.redraw -= 1;
      this.redraw2 += 1;
    },
    async refreshData() {
      this.itemTokens = await api.getItemTokens();
      this.accounts = await api.getAccounts();
      // this.itemTokens = await api.getItemTokens();
      // // console.log("itemTokens: ", (this.itemTokens))
      // this.accounts = await api.getAccounts();
      // // console.log("accounts: ", this.accounts)
      this.redraw -= 1;
      this.redraw2 += 1;
      // this.redraw = Math.random()
      //   .toString(36)
      //   .substring(7);
      // this.redraw2 = Math.random()
      //   .toString(36)
      //   .substring(7);
      // this.$forceUpdate()
    },
    async fetchTransactions() {
      this.fetch = true;
      await api.fetchTransactions();
      this.fetch = false;
      this.redraw -= 1;
      this.redraw2 += 1;
    },
    async importTransactions() {
      this.dialogName = "Importing Data from CSV";
      this.fetch = true;
      // if (!val) return;
      // var self = this;
      // const that = this.files;
      // console.log(that)
      //console.log(that)
      const parseFile = rawFile => {
        return new Promise(resolve => {
          Papa.parse(rawFile, {
            header: true,
            transformHeader: function(h) {
              var f = h.replace(/\s/g, "");
              var i = f.charAt(0).toLowerCase() + f.slice(1);
              return i;
            },
            complete: results => {
              resolve(results.data);
            }
          });
        });
      };
      let parsedData = await parseFile(this.files);
      await api.importTransactions(parsedData);
      this.$store.dispatch("getAll");
      this.fetch = false;
      this.dialogName = "Fetching Transactions";

      // Papa.parse(this.files, {
      //   header: true,
      //   transformHeader: function(h) {
      //     var f = h.replace(/\s/g, "");
      //     var i = f.charAt(0).toLowerCase() + f.slice(1);
      //     return i;
      //   },
      //   //step: function (results, parser) {
      //   //  var upload = JSON.stringify(results.data, null, 2)
      //   //  if(results.data.amount) {
      //   //    api.createTransaction(results.data)
      //   //    console.log(results.data)
      //   //  }
      //   //},
      //   async complete(results) {
      // self.dialog = true;
      // that.doc = JSON.stringify(results.data, null, 2);
      // //console.log(that.doc)
      // //console.log(results.data)
      // for (let i in results.data) {
      //   //console.log(results.data[i])
      //   //setTimeout(function timer(){
      //   //console.log(results.data[i])
      //   //api.createTransaction(results.data[i])
      //   //}, i*40)
      // }
      // await api.importTransactions(results.data);
      // that.$store.dispatch("getAll");
      // // await api.importTransactions(that.doc)
      // self.dialog = false;
      //   },
      //   error(errors) {
      //     console.log("error", errors);
      //   }
      // });
    }
    //async updateSaltEdgeConnections () {
    //  this.loading = true
    //  this.connections = await api.getSaltEdgeConnections()
    //  this.connections = this.connections.data
    //  console.log('Salt Edge Data')
    //  console.log(JSON.stringify(this.connections))
    //  this.loading = false
    //},
  },
  watch: {
    // async loader3() {
    //   // console.log(this)
    //   var vm = this;
    //   const l = this.loader3;
    //   this[l] = !this[l];
    //   if (l) {
    //     // console.log(vm)
    //     vm.refreshUrl = await api.refreshInteractive(vm.item.item_id);
    //     // console.log(vm.refreshUrl)
    //   }
    //   vm.dialog2 = true;
    //   this[l] = false;
    //   this.loader3 = null;
    // },
    // async loader2() {
    //   const l = this.loader2;
    //   this[l] = !this[l];
    //   if (l) {
    //     await api.fetchTransactions();
    //     // let fetchStatus = await api.fetchTransactions()
    //     // console.log(fetchStatus)
    //     await this.refreshData();
    //   }
    //   this[l] = false;
    //   this.loader2 = null;
    // },
    // async loader() {
    //   const l = this.loader;
    //   this[l] = !this[l];
    //   //if (l) api.getSaltEdgeConnections()
    //   //if (l) setTimeout(() => (this[l] = false), 3000)

    //   if (l) {
    //     // this.$store.state.isFetchTransactions = true
    //     await api.getSaltEdgeConnections();
    //     // this.$store.state.isFetchTransactions = false
    //     await this.refreshData();
    //   }
    //   //this.connections = await api.getSaltEdgeConnections()
    //   //this.connections = this.connections.data
    //   //console.log('Salt Edge Data')
    //   //console.log(JSON.stringify(this.connections))
    //   this[l] = false;
    //   this.loader = null;
    // },
    dialog(val) {}
  }
};
</script>

<style scoped>
</style>