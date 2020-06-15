<template>
  <v-content :key="apiStateLoaded">
    <v-col alignSelf="center" v-if="$vuetify.breakpoint.mdAndDown">

      <v-row align="center" justify="center">
      <v-col class="pt-0">
            <!-- <v-btn
                text
                @click="accountsButton = true"
                >Accounts</v-btn
            >
            <v-btn
                text
                @click="accountsButton = false"
                >Transactions Table</v-btn
            > -->
        <v-tabs centered v-model="tab" :background-color="$vuetify.theme.dark?'accent':''">
        <v-tabs-slider></v-tabs-slider>
        <v-tab href="#tab-1">Accounts</v-tab>
        <v-tab href="#tab-2">Transaction Table</v-tab>
        </v-tabs>

      </v-col>
      </v-row>
      <v-row align="center" justify="center" v-if="$vuetify.breakpoint.mdAndDown">
        <v-tabs-items v-model="tab" 
        >
        <v-tab-item value="tab-1">
          <!-- <AccountsBar v-show="accountsButton"></AccountsBar> -->
          <AccountsBar></AccountsBar>
          </v-tab-item>
          <v-tab-item value="tab-2">
          <!-- <TransactionsTable v-bind:transactionsToDisplay="filteredTrans" v-show="!accountsButton"></TransactionsTable> -->
          <TransactionsTable v-bind:transactionsToDisplay="filteredTrans"></TransactionsTable>
          </v-tab-item>
        </v-tabs-items>
      </v-row>
    </v-col>
    <!-- <v-col v-else alignSelf="start"> -->
      <v-row align="start" style="maxWidth: 2400px" v-else>
        <v-col cols="3" class style="min-width: 340px">
          <AccountsBar></AccountsBar>
        </v-col>
        <v-col class="ml-8" cols="8"> 
          <TransactionsTable v-bind:transactionsToDisplay="filteredTrans"></TransactionsTable>
        </v-col>
      </v-row>
    <!-- </v-col> -->
  </v-content>
</template>

<style>
.v-expansion-panel-content__wrap {
  padding: 0px;
}
</style>

<script>
import TransactionsTable from "../components/TransactionsTable.vue";
import AccountsBar from "../components/AccountsBar.vue";
import Currency from "currency.js";
import api from "@/api";
export default {
  data() {
    return {
      // accountsButton: true,
      apiStateLoaded: "",
      transactions: [],
      accounts: [],
      categories: [],
      itemTokens: [],
    };
  },
  computed: {
    filteredTrans() {
      return this.$store.getters.getFilteredTrans;
    },
  },
  async created() {
    //Checks if API has loaded into Vuex store, then loads data
    this.apiStateLoaded = this.$store.state.apiStateLoaded;
    if (this.apiStateLoaded) {
      this.importFromStore();
    }
    this.unsub = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "doneLoading") {
        if (mutation.payload) {
          this.importFromStore();
        }
        this.apiStateLoaded = mutation.payload;
      }
    });
  },
  beforeDestroy() {
    this.unsub();
  },
  methods: {
    //Initializes data from Vuex
    importFromStore() {
      this.categories = this.$store.getters.getAllCategories;
      this.accounts = this.$store.getters.getAllAccounts;
      this.transactions = this.$store.getters.getAllTransactions;
      Object.freeze(this.transactions);
      this.itemTokens = this.$store.getters.getAllItemTokens;
    },
  }
};
</script>