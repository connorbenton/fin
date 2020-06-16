<template>
  <!-- <v-content :key="reloadedData"> -->
  <v-content>
    <v-col class="flex-grow-1 flex-shrink-0 mx-auto">
      <!-- v-on:keyup.enter="$event.target.blur()" -->
      <v-text-field
        class="mt-0 pt-0"
        v-model="searchDisplay"
        prepend-inner-icon="search"
        label="Search"
        single-line
        hide-details
        v-on:blur="selected=[];search = searchDisplay"
        v-on:keyup.enter.prevent="selected=[];search = searchDisplay"
      ></v-text-field>
      <v-data-table
        :mobile-breakpoint="NaN"
        :headers="computedHeaders"
        :search="search"
        :items="transactions"
        :custom-sort="customSort"
        :sort-by.sync="sortBy"
        :sort-desc.sync="sortDesc"
        :footer-props="{
            'items-per-page-options': [20, 50, 100, -1],
            'items-per-page-text': $vuetify.breakpoint.smAndUp ? 'Rows per page:':'Rows',
          }"
        :items-per-page="25"
        :single-select="singleSelectStatus"
        @current-items="getCurrentItems"
        v-model="selected"
        multi-sort
        class="elevation-1"
      >
        <template v-slot:header.category_name>
          Category
          <v-btn icon v-on:click.stop="openFilterMenu($event)">
            <v-icon :color="catIsFiltered?'primary':''">filter_alt</v-icon>
          </v-btn>
          <v-btn icon v-show="catIsFiltered" v-on:click.stop="clearFilter($event)">
            <v-icon>not_interested</v-icon>
          </v-btn>
          <v-btn
            class="ml-4"
            v-if="selected.length > 0"
            color="success"
            small
            v-on:click.stop="openEditMenu($event)"
          >
            <v-icon>build</v-icon>
          </v-btn>
        </template>
        <template v-slot:item="{ item, isSelected, select }">
          <tr
            :class="isSelected?'grey':''"
            v-on:click.left.exact="toggle(isSelected, select, $event)"
            v-on:click.shift.left.exact="shiftToggle(item)"
            v-on:click.ctrl.left.exact="ctrlToggle(isSelected, select, $event)"
            v-on:mousedown.shift.exact.prevent
          >
            <td
              v-if="$vuetify.breakpoint.smAndUp"
            >{{item.date.split("T")[0]}}</td>
            <td
              v-else
            >{{item.date.split("T")[0].substr(item.date.split("T")[0].length - 5)}}</td>
            <!-- <td>{{item.date}}</td> -->
            <td>{{item.description}}</td>
            <td
              v-if="$vuetify.breakpoint.smAndUp"
              :class="{ 'font-weight-bold': item.category_name == 'Uncategorized' }"
            >{{item.category_name}}</td>
            <td>{{formatBalance(item.amount, item.currency_code)}}</td>
            <td>{{item.account_name}}</td>
          </tr>
        </template>
      </v-data-table>
    </v-col>

    <v-menu v-model="editMenu" :position-x="x" :position-y="y" absolute offset-y>
      <v-list class="pa-0">
        <v-menu
          offset-x
          open-on-hover
          v-for="(cat, index) in filtertop_category(categories)"
          :key="index"
        >
          <template v-slot:activator="{ on }">
            <v-hover v-slot:default="{ hover }">
              <v-list-item
                @click="editCategory(cat.top_category, categories)"
                v-on="on"
                :class="`${hover? 'class1': 'class2'}`"
              >
                <v-list-item-title>{{cat.top_category}}</v-list-item-title>
              </v-list-item>
            </v-hover>
          </template>
          <v-list
            v-for="(subcat, index) in filtersub_category(cat.top_category,categories)"
            :key="index"
            class="pa-0"
          >
            <v-hover v-slot:default="{ hover }">
              <v-list-item
                @click="editCategory(subcat.sub_category, categories)"
                :class="`${hover? 'class1': 'class2'}`"
              >
                <v-list-item-title>{{subcat.sub_category}}</v-list-item-title>
              </v-list-item>
            </v-hover>
          </v-list>
        </v-menu>
      </v-list>
    </v-menu>
    <v-menu v-model="filterMenu" :position-x="x" :position-y="y" absolute offset-y>
      <v-list class="pa-0">
        <v-menu
          offset-x
          open-on-hover
          v-for="(cat, index) in filtertop_category(categories)"
          :key="index"
        >
          <template v-slot:activator="{ on }">
            <v-hover v-slot:default="{ hover }">
              <v-list-item
                @click="filterCategory(cat.top_category, categories)"
                v-on="on"
                :class="`${hover? 'class1': 'class2'}`"
              >
                <v-list-item-title>{{cat.top_category}}</v-list-item-title>
              </v-list-item>
            </v-hover>
          </template>
          <v-list
            v-for="(subcat, index) in filtersub_category(cat.top_category,categories)"
            :key="index"
            class="pa-0"
          >
            <v-hover v-slot:default="{ hover }">
              <v-list-item
                @click="filterCategory(subcat.sub_category, categories)"
                :class="`${hover? 'class1': 'class2'}`"
              >
                <v-list-item-title>{{subcat.sub_category}}</v-list-item-title>
              </v-list-item>
            </v-hover>
          </v-list>
        </v-menu>
      </v-list>
    </v-menu>
  </v-content>
</template>

<style>
.class1 {
  background-color: #bdbdbd;
}
</style>

<script>
import api from "@/api";
export default {
  props: ["transactionsToDisplay"],
  data() {
    return {
      // showConfirm: false,
      // reloadedData: 0,
      catIsFiltered: false,
      singleSelectStatus: true,
      currentItems: [],
      searchDisplay: "",
      search: "",
      selected: [],
      editMenu: false,
      filterMenu: false,
      // editedIndex: -1,
      x: 0,
      y: 0,
      headers: [
        { text: "Date", value: "date", dataType: "Date", width: "110" },
        // { text: "Date", value: "date", dataType: "Date"},
        // { text: "Description", value: "description", width: "40%" },
        { text: "Description", value: "description" },
        { text: "Category", value: "category_name" },
        { text: "Amount", value: "amount" },
        { text: "Account", value: "account_name" }
      ],
      headersSm: [
        { text: "Date", value: "date", dataType: "Date"},
        // { text: "Date", value: "date", dataType: "Date"},
        // { text: "Description", value: "description", width: "40%" },
        { text: "Description", value: "description"},
        { text: "Category", value: "category_name", align: " d-none" },
        { text: "Amount", value: "amount"},
        { text: "Account", value: "account_name"}
      ],
      sortDesc: true,
      sortBy: "date",
      masterTransactions: this.transactionsToDisplay,
      transactions: this.transactionsToDisplay,
      accounts: [],
      categories: []
      // console
    };
  },
  computed: {
    computedHeaders() {
      if (this.$vuetify.breakpoint.smAndUp) {
        return this.headers;
      } else {
        return this.headersSm;
      }
    }
  },
  watch: {
    //Loads new transactions from the parent when accounts are set to filter there
    transactionsToDisplay(newVal, oldVal) {
      // this.transactions = newVal.map(t => Object.assign({}, t));
      this.transactions = [];
      this.transactions = newVal;
      this.masterTransactions = newVal;
      Object.freeze(this.transactions);
      Object.freeze(this.masterTransactions);
      this.catIsFiltered = false;
    }
  },
  created() {
    //Initialize categories and accounts from store
    this.categories = this.$store.getters.getAllCategories;
    this.accounts = this.$store.getters.getAllAccounts;
  },
  methods: {
    openFilterMenu: function(event) {
      event.preventDefault();
      this.filterMenu = false;
      this.x = event.clientX;
      this.y = event.clientY;
      // this.editedIndex = this.transactions.indexOf(item);
      this.$nextTick(() => {
        this.filterMenu = true;
      });
    },
    clearFilter() {
      this.transactions = this.masterTransactions;
      this.catIsFiltered = false;
    },
    //Multi-selects current page items with shift
    shiftToggle(item) {
      if (this.selected.length > 2) {
        this.selected = this.currentItems.filter(x => x.id === item.id);
        return;
      } else {
        this.singleSelectStatus = false;
        let a = this.currentItems.findIndex(x => x.id === this.selected[0].id);
        let b = this.currentItems.findIndex(x => x.id === item.id);

        this.selected =
          a >= b
            ? this.currentItems.slice(b, a + 1)
            : this.currentItems.slice(a, b + 1);

        // this.selected = this.currentItems.slice(
        // this.currentItems.findIndex(x => x.id === this.selected[0].id),
        // this.currentItems.findIndex(x => x.id === this.selected[1].id + 1),
        // )

        // console.log(this.selected);
        // console.log(this.currentItems);
        this.singleSelectStatus = true;
      }
    },
    //Ctrl+selects current page item for category change
    ctrlToggle(isSelected, select, e) {
      this.singleSelectStatus = false;
      this.$nextTick(() => {
        select(!isSelected);
        this.singleSelectStatus = true;
      });
    },
    //Selects current page item for category change
    toggle(isSelected, select, e) {
      select(!isSelected);
    },
    getCurrentItems(e) {
      this.selected = [];
      this.currentItems = e;
    },
    //Sorting table to allow correct multi-sorting on the various columns
    customSort(items, indexes, isDesc) {
      let cats = this.categories;
      items.sort(function(a, b) {
        return indexes.reduce(function(bool, k) {
          let i = indexes.indexOf(k);
          if (k === "amount") {
            if (!isDesc[i]) {
              return bool || parseFloat(a[k]) - parseFloat(b[k]);
            } else {
              return bool || parseFloat(b[k]) - parseFloat(a[k]);
            }
          }
          if (!isDesc[i]) {
            return bool || a[k].localeCompare(b[k]);
          } else {
            return bool || b[k].localeCompare(a[k]);
          }
        }, 0);
      });
      return items;
    },
    //Format balance to display correctly in Amount column
    formatBalance: function(bal, code) {
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: code
      }).format(bal);
    },
    //Opens category dropdown at correct location
    //Doesn't use item directly anymore (now with multi-select change)
    // openEditMenu: function(event, item) {
    openEditMenu: function(event) {
      event.preventDefault();
      this.editMenu = false;
      this.x = event.clientX;
      this.y = event.clientY;
      // this.editedIndex = this.transactions.indexOf(item);
      this.$nextTick(() => {
        this.editMenu = true;
      });
    },

    filterCategory(cat, categories) {
      let catToSave = cat;
      // console.log(catToSave);
      let foundCat = categories.find(x => x.sub_category === catToSave);
      // console.log(foundCat)
      let newTransSet = [];
      if (foundCat.sub_category === foundCat.top_category) {
        let subCats = categories.filter(
          x => x.top_category === foundCat.top_category
        );
        newTransSet = this.transactions.filter(x => {
          return subCats.filter(e => e.id === x.category).length > 0;
        });
      } else {
        newTransSet = this.transactions.filter(x => x.category === foundCat.id);
      }
      // console.log(newTransSet);
      this.transactions = newTransSet;
      this.catIsFiltered = true;
      this.filterMenu = false;
    },
    //Updates first backend, then store with new category
    //Now using selected
    async editCategory(cat, categories) {
      let catToSave = cat;
      let foundCat = categories.find(x => x.sub_category === catToSave);

      // let editArray = this.selected.map(item => {
      let editArray = [];
      this.selected.forEach(item => {
        let a = this.transactions.find(x => x.id === item.id);

        a.category = foundCat.id;
        a.category_name = foundCat.sub_category;
        editArray.push(a);
      });

      await this.$store.commit("updateTransaction", editArray);
      await api.upsertTransaction(editArray).then(() => {
        this.selected = [];
        this.editMenu = false;

        this.$emit("changed");
      });
    },

    filtersub_category(topCat, categories) {
      let filterTop = topCat;
      let filtered = categories.filter(function(item) {
        return item.top_category == filterTop;
      });
      return filtered;
    },
    //Populate array of top categories for category dropdown
    filtertop_category(categories) {
      let filtered = categories.filter(function(item) {
        return item.top_category == item.sub_category;
      });
      return filtered;
    }
  }
};
</script>

<style scoped>

/* * {
  -webkit-tap-highlight-color: rgba(0, 0, 0, 0);
  -webkit-tap-higlight-color: transparent;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
} */

/* .v-data-table
  /deep/
  tbody
  /deep/
  tr:hover:not(.v-data-table__expanded__content) {
  background: rgba(255,255,255,0) !important;
} */

/* td {
  padding: 0px;
} */

</style>