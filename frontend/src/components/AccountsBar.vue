<template>
  <v-content>
    <!-- <v-row align="center" justify="center"> -->
    <!-- <v-col class="flex-grow-0 flex-shrink-1 mx-auto"> -->
    <v-col align="center" justify="center" class="mx-auto">
      <v-expansion-panels
        v-model="panel"
        style="minWidth: 340px; maxWidth: 340px;"
        accordion
        multiple
      >
        <v-expansion-panel class="ma-0" v-for="(type, i) in accountTypes" :key="i">
          <v-card outlined>
            <v-expansion-panel-header class="px-4 py-2">
              <v-col align="center" class="pa-0">
                <v-row justify="center" class="subtitle my-2">{{type}}</v-row>
                <v-row no-gutters align="center" class="pr-2" style="flex-wrap: wrap">
                  <v-col
                    v-for="(val, v) in AmountMethod(accounts, type, 'sum')"
                    :key="v"
                    :class="{ 'green--text': type !== 'Credit' }"
                    class="ma-2"
                  >{{val.amount}}</v-col>
                </v-row>
              </v-col>
            </v-expansion-panel-header>
          </v-card>
          <v-expansion-panel-content class="nopad">
            <v-list class="pa-0">
              <v-list-item-group v-model="model[i]" multiple color="indigo">
                <v-list-item
                  v-for="(acct, index) in printAccounts[i]"
                  :key="index"
                  class="pa-0"
                  @click="showTransactionsForAccount($event, acct)"
                >
                  <v-col class="ma-0 pa-0">
                    <v-row no-gutters align="center" class="px-4" style="flex-wrap: nowrap">
                      <v-col cols="8" class="flex-grow-1 flex-shrink-0 pa-0">
                        <v-list-item-content class="text-left pa-0">
                          <v-list-item-title class="body-2">{{acct.name}}</v-list-item-title>
                          <v-list-item-subtitle class="caption">{{acct.institution}}</v-list-item-subtitle>
                        </v-list-item-content>
                      </v-col>
                      <v-col cols="4" class="flex-grow-0 flex-shrink-1 pa-0">
                        <v-list-item-content class="text-right pa-0">
                          <v-list-item-title
                            class="body-2"
                          >{{formatBalance(acct.balance, acct.currency)}}</v-list-item-title>
                          <v-list-item-subtitle
                            class="caption"
                          >{{timeSince(acct.updated_at, acct.item_id)}}</v-list-item-subtitle>
                        </v-list-item-content>
                      </v-col>
                    </v-row>
                    <v-row
                      v-if="showAcc"
                      no-gutters
                      align="center"
                      class="px-4"
                      style="flex-wrap: nowrap"
                    >
                      <v-switch
                        dense
                        inset
                        class="ma-0 pa-0"
                        v-model="acct.ignore_transactions"
                        :label="`Account is Hidden from Sum: ${acct.ignore_transactions.toString()}`"
                        @click.stop="toggleAccountShow(acct)"
                      ></v-switch>
                    </v-row>
                  </v-col>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
      <v-row class="ma-4" align="center" justify="center">
        <v-switch v-model="showAcc" label="Choose Accounts to Include in Sum"></v-switch>
      </v-row>
    </v-col>
    <!-- </v-row> -->
  </v-content>
</template>

<style>
.v-expansion-panel-content__wrap {
  padding: 0px;
}
</style>

<script>
import TransactionsTable from "../components/TransactionsTable.vue";
import Currency from "currency.js";
import api from "@/api";
export default {
  data() {
    return {
      accountTypes: ["Cash", "Credit", "Investment"],
      panel: [0, 1, 2],
      transactions: this.$store.getters.getAllTransactions,
      accounts: this.$store.getters.getAllAccounts,
      itemTokens: this.$store.getters.getAllItemTokens,
      model: [],
      cashAccountTypes: [
        "account",
        "bonus",
        "debit_card",
        "ewallet",
        "savings",
        "card",
        "depository"
      ],
      creditAccountTypes: ["credit", "credit_card"],
      investmentAccountTypes: ["investment"],
      printedAccounts: [],
      showAcc: false
    };
  },
  computed: {
    //Sets list of accounts to print in the accounts pane
    printAccounts() {
      for (let index in this.accountTypes) {
        this.printedAccounts[index] = this.AmountMethod(
          this.accounts,
          this.accountTypes[index],
          "account"
        );
        // console.log(this.printedAccounts);
      }
      return this.printedAccounts;
    }
  },
  methods: {
    //Filters transactions based on accounts selected in accounts pane
    filteredItems() {
      if (this.model.every(element => element.length === 0))
        return this.transactions;
      let accountsToFilter = [];
      for (let i in this.model) {
        for (let j in this.model[i]) {
          accountsToFilter.push(
            this.printedAccounts[i][this.model[i][j]].account_id
          );
        }
      }
      let filteredTrans = this.transactions.filter(x => {
        return accountsToFilter.indexOf(x.account_id) !== -1;
      });
      return filteredTrans;
    },
    toggleAccountShow(acct) {
      try {
        acct.ignore_transactions = !acct.ignore_transactions;
        api.upsertAccountIgnore(acct);
        this.$store.commit("updateAccountIgnore", acct);
      } catch (err) {
        acct.ignore_transactions = !acct.ignore_transactions;
        console.error(err);
      }
    },
    //Triggers display event when account is clicked in accounts pane
    showTransactionsForAccount: function(event, id) {
      this.$nextTick(() => {
        let y = id;
        let x = this.model;
      if (this.model.every(element => element.length === 0)) {
        this.$store.commit("updateFilteredTrans", this.transactions);
        return;
      }
      let accountsToFilter = [];
      for (let i in this.model) {
        for (let j in this.model[i]) {
          accountsToFilter.push(
            this.printedAccounts[i][this.model[i][j]].account_id
          );
        }
      }
      let filteredTrans = this.transactions.filter(x => {
        return accountsToFilter.indexOf(x.account_id) !== -1;
      });
        this.$store.commit("updateFilteredTrans", filteredTrans);
        return;
      });
    },
    //Formats balance to display correctly in the accounts pane
    formatBalance: function(bal, code) {
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: code
      }).format(bal);
    },
    //General function to filter, sort, and sum accounts for the accounts pane
    AmountMethod(accounts, typeToCheck, evalType) {
      var filtered = accounts.filter(item => {
        switch (typeToCheck) {
          case "Cash":
            return this.cashAccountTypes.includes(item.type);
            break;
          case "Credit":
            return this.creditAccountTypes.includes(item.type);
            break;
          case "Investment":
            return this.investmentAccountTypes.includes(item.type);
        }
      });
      if (evalType === "account") {
        filtered.sort(
          (a, b) =>
            a.institution.localeCompare(b.institution) ||
            a.name.localeCompare(b.name)
        );
        return filtered;
      }
      if (evalType === "sum") {
        let currencyIndex = [];
        let currencies = [...new Set(filtered.map(item => item.currency))];
        for (let currency of currencies) {
          let accountsToSum = filtered.filter(x => x.currency === currency);
          let currencyToPush = {};
          let initialAmount = Currency(0);
          let numAmount = accountsToSum
            .filter(({ ignore_transactions }) => !ignore_transactions)
            .reduce((a, b) => Currency(a).add(b.balance), initialAmount);
          currencyToPush.amount = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: currency
          }).format(numAmount);
          currencyToPush.currencyName = currency;
          currencyIndex.push(currencyToPush);
        }
        currencyIndex.sort((a, b) =>
          a.currencyName.localeCompare(b.currencyName)
        );
        return currencyIndex;
      }
    },
    //Gets time in a nice date for the last refresh in accounts pane
    timeSince(date, itemID) {
      let matchItem = this.itemTokens.filter(x => x.item_id == itemID);

      if (matchItem[0].interactive) {
        date = Date.parse(matchItem[0].last_refresh);
      } else {
        date = Date.parse(date);
      }
      var seconds = Math.floor((new Date() - date) / 1000);

      var interval = Math.floor(seconds / 31536000);

      if (interval > 1) {
        return interval + " years ago";
      }
      interval = Math.floor(seconds / 2592000);
      if (interval > 1) {
        return interval + " months ago";
      }
      interval = Math.floor(seconds / 86400);
      if (interval > 1) {
        return interval + " days ago";
      }
      interval = Math.floor(seconds / 3600);
      if (interval > 1) {
        return interval + " hours ago";
      }
      interval = Math.floor(seconds / 60);
      if (interval > 1) {
        return interval + " minutes ago";
      }
      return Math.floor(seconds) + " seconds ago";
    }
  }
};
</script>