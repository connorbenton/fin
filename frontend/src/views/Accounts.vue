<template>
  <v-container>
    <v-dialog v-model="sEdge" persistent max-width="800">
      <v-card>
        <v-toolbar flat dense>
          <v-spacer></v-spacer>
          <v-btn icon @click="closeSEdge()">
            <v-icon>close</v-icon>
          </v-btn>
        </v-toolbar>
        <iframe style="width:100%" height="700" :src="sEdgeURL" id="sEdgeRef" frameborder="0"></iframe>
      </v-card>
    </v-dialog>
    <v-dialog v-model="fetch" persistent width="60%" max-width="800">
      <v-card color="primary" dark>
        <v-card-text>
          {{dialogName}}
          <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-btn
      class="mt-2"
      :loading="loading3"
      :disabled="loading3"
      color="success"
      @click.native="fetchTransactions()"
    >Fetch Transactions for All Accounts</v-btn>

    <h1 v-if="USE_PLAID=='TRUE'" class="title mt-3">Plaid Connections</h1>
    <v-card
      v-if="USE_PLAID=='TRUE'"
      class="d-inline-block mx-auto my-3"
      max-width="1000"
      tile
      :key="redraw"
    >
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th>Name</th>
              <th>Last Transaction Fetch</th>
              <th :colspan="2">Status</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="(item, index1) in plaid(itemTokens)">
              <tr :key="`${index1}-${item.id}`">
                <td>{{item.institution}}</td>
                <td>{{ localeDate(item.last_downloaded_transactions) }}</td>
                <td v-if="item.needs_re_login == 1">
                  <v-btn
                    color="warning"
                    :loading="plaidRefresh"
                    dark
                    @click.native="startReLogin(item.item_id)"
                  >{{refreshText}}</v-btn>
                </td>
                <td v-else>
                  <v-icon color="success">check</v-icon>
                </td>
                <td v-if="$vuetify.breakpoint.smAndUp">
                  <v-tooltip right nudge-right="16">
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn
                        icon
                        v-bind="attrs"
                        v-on="on"
                        @click="toggleAccounts(index1, showPlaidAccounts)"
                      >
                        <v-icon small>build</v-icon>
                      </v-btn>
                    </template>
                    <span>Toggle Edit Account Names</span>
                  </v-tooltip>
                </td>
              </tr>
              <tr
                v-for="(account, j) in matchAccounts(item.item_id)"
                :key="j"
                v-show="showPlaidAccounts[index1]"
              >
                <td class="pl-10">
                  <v-btn icon @click="editAccountName(account)">
                    <v-icon small>create</v-icon>
                  </v-btn>
                  {{account.name}}
                </td>
                <td class="pl-10">{{account.type}}</td>
                <td class="pl-10" :colspan="2">{{formatBalance(account.balance, account.currency)}}</td>
              </tr>
            </template>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
    <v-btn v-if="USE_PLAID=='TRUE'" text class="d-block px-0">
      <plaid-link
        v-if="USE_PLAID=='TRUE'"
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
    <h1 v-if="USE_SALTEDGE=='TRUE'" class="title mt-3">Salt Edge Connections</h1>
    <v-card
      v-if="USE_SALTEDGE=='TRUE'"
      class="d-inline-block mx-auto my-3"
      max-width="1000"
      tile
      :key="redraw2"
    >
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th>Name</th>
              <th v-if="$vuetify.breakpoint.smAndUp">Last Transaction Fetch</th>
              <th v-else>Last Tx Fetch</th>
              <th>Last Refresh</th>
              <th :colspan="2">Next Refresh Available In</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="(item, index2) in saltEdge(itemTokens)">
              <tr :key="`${index2}-${item.id}`">
                <td>{{item.institution}}</td>
                <td
                  v-show="item.last_downloaded_transactions !== null"
                >{{localeDate(item.last_downloaded_transactions)}}</td>
                <td
                  v-show="item.last_downloaded_transactions == null"
                >{{"last transaction fetch: " + "Never"}}</td>
                <td v-if="item.interactive == 1">{{localeDate(item.last_refresh)}}</td>
                <td v-else>
                  <v-icon>remove</v-icon>
                </td>
                <td v-if="item.interactive == 1 && showRefresh">
                  <v-btn
                    :loading="loading4"
                    :disabled="loading4"
                    color="primary"
                    dark
                    @click.native="startRefreshInteractive(item.item_id)"
                  >
                    {{refreshText}}
                    <!-- Refresh Connection
                    <br />(needs Credentials)-->
                  </v-btn>
                </td>
                <td
                  v-else-if="item.interactive == 1"
                >{{ timeToRefresh(item.next_refresh_possible) }}</td>
                <td v-else>
                  <v-icon>remove</v-icon>
                </td>
                <td v-if="$vuetify.breakpoint.smAndUp">
                  <v-tooltip right nudge-right="16">
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn
                        icon
                        v-bind="attrs"
                        v-on="on"
                        @click="toggleAccounts(index2, showSaltEdgeAccounts)"
                      >
                        <v-icon small>build</v-icon>
                      </v-btn>
                    </template>
                    <span>Toggle Edit Account Names</span>
                  </v-tooltip>
                </td>
              </tr>
              <tr
                v-for="(account, j) in matchAccounts(item.item_id)"
                :key="j"
                v-show="showSaltEdgeAccounts[index2]"
              >
                <td class="pl-10">
                  <v-btn icon @click="editAccountName(account)">
                    <v-icon small>create</v-icon>
                  </v-btn>
                  {{account.name}}
                </td>
                <td class="pl-10">{{account.type}}</td>
                <td class="pl-10" :colspan="3">{{formatBalance(account.balance, account.currency)}}</td>
              </tr>
            </template>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
    <v-btn
      class="d-block"
      v-if="USE_SALTEDGE=='TRUE'"
      :loading="loading4"
      :disabled="loading4"
      @click.native="startCreateInteractive()"
    >Open Salt Edge To Add New Account</v-btn>
    <!-- <v-col style="max-width: 700px"> -->
    <h1 class="title mt-3">Import CSV from Mint.com</h1>
    <v-flex mt-4>
      <v-file-input
        prepend-icon="attach_file"
        accept=".csv"
        v-model="files"
        style="max-width: 400px"
        label="Choose CSV from mint.com to import"
      ></v-file-input>
      <v-btn
        :loading="loading2"
        :disabled="loading2"
        @click.native="importTransactions()"
      >Import CSV</v-btn>
    </v-flex>
    <!-- </v-col> -->
    <!-- <v-col class="px-6 py-2"> -->
    <v-row class="px-3 py-4">
      <v-btn color="warning" dark @click.native="resetDB()">Reset Database</v-btn>
    </v-row>
    <!-- </v-col> -->
    <!-- <v-col class="px-6 py-2"> -->
    <v-row class="px-3">
      <v-btn
        color="warning"
        dark
        @click.native="resetDBFull()"
      >Reset Database (Including Item Tokens)</v-btn>
    </v-row>
    <!-- </v-col> -->

    <v-dialog v-model="dialogChangeName" max-width="500px">
      <v-card>
        <v-card-title>
          <span class="headline">Edit Account Name</span>
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row>
              <v-text-field v-model="editedItem.name" label="Account name"></v-text-field>
            </v-row>
          </v-container>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="close">Cancel</v-btn>
          <v-btn text @click="save">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import PlaidLink from "vue-plaid-link";
import api from "@/api";
import moment from "moment";
const Papa = require("papaparse");

export default {
  data() {
    return {
      // console: null,
      apiStateLoaded: false,
      showRefresh: false,
      fetch: false,
      sEdge: false,
      sEdgeURL: "",
      dialogName: "Fetching Transactions",
      dialog: false,
      dialog2: false,
      dialogChangeName: false,
      vh: null,
      transactions: [],
      showPlaidAccounts: [],
      showSaltEdgeAccounts: [],
      editedIndex: -1,
      editedItem: {
        name: ""
      },
      files: null,
      environment:
        process.env.VUE_APP_PLAID_ENVIRONMENT || window._env_.PLAID_ENVIRONMENT,
      PLAID_PUBLIC_KEY:
        process.env.VUE_APP_PLAID_PUBLIC_KEY || window._env_.PLAID_PUBLIC_KEY,
      USE_PLAID: process.env.VUE_APP_USE_PLAID || window._env_.USE_PLAID,
      USE_SALTEDGE:
        process.env.VUE_APP_USE_SALTEDGE || window._env_.USE_SALTEDGE,
      updateToken: null,
      itemTokens: [],
      accounts: [],
      plaidRefresh: false,
      loading: false,
      loading2: false,
      loading3: false,
      loading4: false,
      refreshUrl: null,
      reactive: true,
      redraw: 0,
      redraw2: 1,
      redraw3: 2,
    };
  },

  created() {
    // console.log(this.PLAID_PUBLIC_KEY);
    // this.console = window.console;

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
  mounted() {
      // let recaptchaScript = document.createElement('script')
      // document.head.appendChild(recaptchaScript)
      let vm = this;
  },
  beforeDestroy() {
    this.unsub();
  },
  components: { PlaidLink },
  computed: {
    refreshText() {
      if (this.$vuetify.breakpoint.smAndUp) {
        return "Refresh Connection";
      } else {
        return "Refresh";
      }
    }
  },
  methods: {
    closeSEdge() {
      this.loading4 = false;
      this.sEdge = false;
      document.getElementById("sEdgeRef").src = "about:blank";
      this.fetchRefreshData();
    },
    editAccountName(account) {
      this.editedIndex = this.accounts.indexOf(account);
      this.editedItem = Object.assign({}, account);
      this.dialogChangeName = true;
    },
    close() {
      this.dialogChangeName = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      });
    },

    async save() {
      const itemToUpdate = this.editedItem;
      Object.assign(this.accounts[this.editedIndex], this.editedItem);
      this.close();
      // console.log(this.editedItem);
      await this.$store.commit("updateAccountName", itemToUpdate);
      await api.upsertAccountName(itemToUpdate);
      this.fetch = true;
      await this.$store.dispatch("getTransactions");
      this.fetch = false;
    },
    matchAccounts(id) {
      return this.accounts.filter(acc => acc.item_id === id);
    },
    toggleAccounts(index, array) {
      const tempVal = array[index] == undefined ? true : !array[index];
      this.$set(array, index, tempVal);
      // console.log(this.showPlaidAccounts);
    },
    localeDate(date) {
      // console.log(date)
      if (date == null) {
        return "Never";
      } else {
        if (this.$vuetify.breakpoint.smAndUp) {
          return new Date(date).toLocaleString();
        } else {
          return moment(date).format("MMM D hh:mm");
        }
      }
    },
    timeToRefresh(time) {
      let diffmin = moment().diff(moment(time), "minutes");
      let diffsec = moment().diff(moment(time), "seconds");
      if (diffsec > 0) this.showRefresh = true;
      if (diffsec > -60) return diffsec * -1 + " seconds";
      else return diffmin * -1 + " minutes";
    },
    formatBalance: function(bal, code) {
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: code
      }).format(bal);
    },
    async startCreateInteractive() {
      var vm = this;
      vm.loading4 = true;
      vm.refreshUrl = await api.createInteractive();
      vm.sEdgeURL = vm.refreshUrl;
      vm.sEdge = true;

      await new Promise(function(resolve, reject) {
        (function waitForSE() {
          try {
            if (
              document
                .getElementById("sEdgeRef")
                .contentWindow.location.host.indexOf(document.domain) != -1
            )
              return resolve();
          } catch {}
          setTimeout(waitForSE, 30);
        })();
      });

      vm.loading4 = false;
      vm.sEdge = false;
      document.getElementById("sEdgeRef").src = "about:blank";
      vm.fetchRefreshData();
    },
    async resetDB() {
      this.dialogName = "Resetting Database";
      this.fetch = true;
      await api.resetDB();
      this.fetch = false;
      this.dialogName = "Fetching Transactions";
      this.refreshData();
    },
    // async resetToken() {
    //   this.dialogName = "Resetting Database";
    //   this.fetch = true;
    //   await api.resetToken();
    //   this.fetch = false;
    //   this.dialogName = "Fetching Transactions";
    //   this.refreshData();
    // },
    async resetDBFull() {
      if (confirm("are you sure?")) {
        this.dialogName = "Resetting Database";
        this.fetch = true;
        await api.resetDBFull();
        this.fetch = false;
        this.dialogName = "Fetching Transactions";
        this.refreshData();
      }
    },
    async startRefreshInteractive(id) {
      var vm = this;

      vm.refreshUrl = await api.refreshInteractive(id);
      if (vm.refreshUrl == "") {
        alert("Try again in a minute");
        return;
      }
      vm.loading4 = true;
      vm.sEdgeURL = vm.refreshUrl;
      vm.sEdge = true;

      await new Promise(function(resolve, reject) {
        (function waitForSE() {
          try {
            if (
              document
                .getElementById("sEdgeRef")
                .contentWindow.location.host.indexOf(document.domain) != -1
            )
              return resolve();
          } catch {}
          setTimeout(waitForSE, 30);
        })();
      });

      vm.loading4 = false;
      vm.sEdge = false;
      document.getElementById("sEdgeRef").src = "about:blank";
      vm.fetchRefreshData();
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
      await this.fetchTransactions();
      this.refreshData();
    },
    async startReLogin(id) {
      this.plaidRefresh = true;
      let ItemToUpload = {
        item_id: id
      };
      let tok = await api.plaidGeneratePublicToken(ItemToUpload);
      this.updateToken = tok.public_token;

      await this.$nextTick();
      this.plaidRefresh = false;
      await this.$refs.plaidLinkRef.onScriptLoaded();
      this.$refs.plaidLinkRef.handleOnClick();
      await this.fetchTransactions();
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
      this.redraw3 += 1;
    },
    async fetchRefreshData() {
      this.dialogName =
        "Checking for new SaltEdge connections to add and fetching transactions";
      this.fetch = true;
      let res = await api.fetchTransactions();
      this.itemTokens = await api.getItemTokens();
      this.accounts = await api.getAccounts();
      this.fetch = false;
      this.dialogName = "Fetching Transactions";
      this.redraw -= 1;
      this.redraw2 += 1;
      this.redraw3 += 1;
    },
    async refreshData() {
      this.itemTokens = await api.getItemTokens();
      this.accounts = await api.getAccounts();
      this.redraw -= 1;
      this.redraw2 += 1;
      this.redraw3 += 1;
    },
    async fetchTransactions() {
      this.fetch = true;
      await api.fetchTransactions();
      this.fetch = false;
      this.redraw -= 1;
      this.redraw2 += 1;
      this.redraw3 += 1;
    },
    async importTransactions() {
      if (this.files == null) return;
      this.dialogName = "Importing Data from CSV";
      this.fetch = true;

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
      this.files = null;
      this.$store.dispatch("getAll");
      this.fetch = false;
      this.dialogName = "Fetching Transactions";
    }
  }
};
</script>