<template>
  <!-- <v-content :key="reloadedData"> -->
  <v-content>
    <v-col class="flex-grow-1 flex-shrink-0 mx-auto">
        <!-- v-on:keyup.enter="$event.target.blur()" -->
      <v-text-field
        v-model="searchDisplay"
        prepend-inner-icon="search"
        label="Search"
        single-line
        hide-details
        v-on:blur="selected=[];search = searchDisplay"
        v-on:keyup.enter.prevent="selected=[];search = searchDisplay"
      ></v-text-field>
      <v-data-table
        :headers="headers"
        :search="search"
        :items="transactions"
        :custom-sort="customSort"
        :sort-by.sync="sortBy"
        :sort-desc.sync="sortDesc"
        :footer-props="{
            'items-per-page-options': [25, 50, 100, -1]
          }"
        :items-per-page="25"
        :single-select="singleSelectStatus"
        @current-items="getCurrentItems"
        v-model="selected"
        multi-sort
        class="elevation-1"
      >
        <template v-slot:header.catName="">Category 
          <v-btn 
            class="ml-4"
            v-if="selected.length > 0" 
            color="success" 
            small
            v-on:click.stop="openEditMenu($event)">Change Category</v-btn>
        </template>
        <template v-slot:item="{ item, isSelected, select }">
          <tr
            :class="isSelected?'grey':''"
            v-on:click.left.exact="toggle(isSelected, select, $event)"
            v-on:click.shift.left.exact="shiftToggle(item)"
            v-on:click.ctrl.left.exact="ctrlToggle(isSelected, select, $event)"
            v-on:mousedown.shift.exact.prevent
          >
            <td>{{item.date}}</td>
            <td>{{item.description}}</td>
              <!-- @click="openEditMenu($event, item)" -->
            <td
              :class="{ 'font-weight-bold': item.catName == 'Uncategorized' }"
            >{{item.catName}}</td>
            <td>{{formatBalance(item.amount, item.currency_code)}}</td>
            <td>{{item.accName}}</td>
          </tr>
        </template>
        <!-- <template v-slot:item.account="{ item }">{{item.accName}}</template>
        <template v-slot:item.amount="{ item }">{{formatBalance(item.amount, item.currency_code)}}</template>
        <template v-slot:item.catName="{ item }">
          <v-card-text
            class="pa-0"
            @click="openEditMenu($event, item)"
            :class="{ 'font-weight-bold': item.catName == 'Uncategorized' }"
          >{{item.catName}}</v-card-text>
        </template>-->
      </v-data-table>
    </v-col>

    <v-menu v-model="editMenu" :position-x="x" :position-y="y" absolute offset-y>
      <v-list class="pa-0">
        <v-menu
          offset-x
          open-on-hover
          v-for="(cat, index) in filterTopCategory(categories)"
          :key="index"
        >
          <template v-slot:activator="{ on }">
            <v-hover v-slot:default="{ hover }">
              <v-list-item
                @click="editCategory(cat.topCategory, categories)"
                v-on="on"
                :class="`${hover? 'class1': 'class2'}`"
              >
                <v-list-item-title>{{cat.topCategory}}</v-list-item-title>
              </v-list-item>
            </v-hover>
          </template>
          <v-list
            v-for="(subcat, index) in filterSubCategory(cat.topCategory,categories)"
            :key="index"
            class="pa-0"
          >
            <v-hover v-slot:default="{ hover }">
              <v-list-item
                @click="editCategory(subcat.subCategory, categories)"
                :class="`${hover? 'class1': 'class2'}`"
              >
                <v-list-item-title>{{subcat.subCategory}}</v-list-item-title>
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
      singleSelectStatus: true,
      currentItems: [],
      searchDisplay: "",
      search: "",
      selected: [],
      editMenu: false,
      // editedIndex: -1,
      x: 0,
      y: 0,
      headers: [
        { text: "Date", value: "date", dataType: "Date", width: "110" },
        { text: "Description", value: "description", width: "40%" },
        { text: "Category", value: "catName" },
        // //Eliminated currency since now balance is formatted instead
        // {text: 'Currency', value: 'currency_code', align: 'end' },
        // {text: "Currency", value: "currency_code" },
        { text: "Amount", value: "amount" },
        { text: "Account", value: "accName" },
        { text: "accID", value: "account_id", align: " d-none" }
      ],
      sortDesc: true,
      sortBy: "date",
      transactions: this.transactionsToDisplay.map(t => Object.assign({}, t)),
      accounts: [],
      categories: []
      // console
    };
  },
  watch: {
    //Loads new transactions from the parent when accounts are set to filter there
    transactionsToDisplay(newVal, oldVal) {
      this.transactions = newVal.map(t => Object.assign({}, t));
    }
  },
  created() {
    //Initialize categories and accounts from store
    this.categories = this.$store.getters.getAllCategories;
    this.accounts = this.$store.getters.getAllAccounts;
  },
  methods: {
    //Multi-selects current page items with shift
    shiftToggle(item) {
      if (this.selected.length > 2) {
        this.selected = this.currentItems.filter(x => x.id === item.id);
      return;
      } else {
      this.singleSelectStatus = false;
        let a = this.currentItems.findIndex(x => x.id === this.selected[0].id);
        let b = this.currentItems.findIndex(x => x.id === item.id);

        this.selected = (a >= b) ? this.currentItems.slice(b, a + 1) : this.currentItems.slice(a, b + 1); 

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

    //Updates first backend, then store with new category
    //Now using selected
    async editCategory(cat, categories) {
      let catToSave = cat;
      let foundCat = categories.find(x => x.subCategory === catToSave);
      const editArray = this.selected.map(item => {

      // for (let i in this.selected) {
        // console.log(this.selected[i].description);
        // let a = this.transactions.find(x => x.id === this.selected[i].id);
        let a = this.transactions.find(x => x.id === item.id);
        // console.log(a.description);
      // return;
      // this.showConfirm = true;
      // this.transactions[this.editedIndex].category = foundCat.id;
      // this.transactions[this.editedIndex].catName = foundCat.subCategory;
      a.category = foundCat.id;
      a.catName = foundCat.subCategory;

      api.updateTransaction(a.id, a);

      this.$store.commit( "updateTransaction", a);
      return 'update done';
      });
      
      await Promise.all(editArray).then(() => {

      // }

      this.selected = [];
      this.editMenu = false;
      // this.$forceUpdate;
      // this.reloadedData += 1;
      this.$emit("changed");
      this.$store.dispatch("reanalyze");
      });
      // .then(() => this.$store.dispatch("reanalyze")); 
      // this.dispatchReanalyze();
    },
    // async dispatchReanalyze() {
      // this.$store.dispatch("reanalyze");
    // },
    //Populate array of sub categories for category dropdown
    filterSubCategory(topCat, categories) {
      let filterTop = topCat;
      let filtered = categories.filter(function(item) {
        return item.topCategory == filterTop;
      });
      return filtered;
    },
    //Populate array of top categories for category dropdown
    filterTopCategory(categories) {
      let filtered = categories.filter(function(item) {
        return item.topCategory == item.subCategory;
      });
      return filtered;
    }
  }
};
</script>