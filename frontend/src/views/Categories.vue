<template>

<v-content>
  <v-row>
  <v-col>
 <v-card class="mx-auto" width="300">
          <v-list class="overflow-y-auto" style = "max-height:90vh">
            <v-list-item v-for="s_cat in saltedge_top" :key = s_cat>
            <v-list-item-title v-text="s_cat"></v-list-item-title>
            </v-list-item>
          </v-list>
 </v-card>
  </v-col>
  <v-col>
 <v-card class="mx-auto" width="300">
          <v-list class="pa-0 overflow-y-auto" style = "max-height:90vh">
        <v-menu offset-x open-on-hover
              v-for="(cat, index) in filtertop_category(categories)"
              :key="index"
        >
          <template v-slot:activator="{ on }">
            <v-hover v-slot:default="{ hover }">
            <v-list-item
              @click="editCategory(cat.top_category)"
              v-on="on"
            :class="`${hover? 'class1': 'class2'}`"
            >
            <v-list-item-title 
            >
              {{cat.top_category}}
            </v-list-item-title>
            <v-list-item-action>
              <v-btn icon>
              <v-icon color="grey lighten-1">delete</v-icon>
              </v-btn>
            </v-list-item-action>
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
              @click="editCategory(subcat.sub_category)"
              :class="`${hover? 'class1': 'class2'}`"
            >
            <v-list-item-title
            >
            {{subcat.sub_category}}
            </v-list-item-title>
            <v-list-item-action>
              <v-btn icon>
              <v-icon v-if="cat.exclude_from_analysis == 1" color="red lighten-1">visibility_off</v-icon>
              <v-icon v-if="cat.exclude_from_analysis == 0" color="green lighten-1">visibility</v-icon>
              </v-btn>
            </v-list-item-action>
            <v-list-item-action>
              <v-btn icon>
              <v-icon color="grey lighten-1">delete</v-icon>
              </v-btn>
            </v-list-item-action>
            </v-list-item>
            </v-hover>
            </v-list>
            </v-menu>
          </v-list>
      </v-card>
      </v-col>
      </v-row>
    </v-content>

</template>

<style>
.class1 {
  background-color:#BDBDBD;
}
html {
  overflow-y: auto
}
.class2 {
}

</style>

<script>
import api from '@/api' 
export default { 
  data () {
     return {
      search: '',
      editMenu: false,
      editedIndex: -1,
      x: 0,
      y: 0,
      currentItem: {},
      currenttop_category: null,
      offset: true,
      headers: [
        {text: 'Date', value: 'date' , dataType: 'Date'},
        {text: 'Description', value: 'description' },
        {text: 'Category', value: 'category' },
        {text: 'Amount', value: 'amount' },
        {text: 'Account', value: 'accountName' },
      ],
      loading: false,
      transactions: [],
      categories: [],
      plaid_categories: [],
      // plaid_display: [],
      saltedge_categories: [],
      saltedge_top: ['personal', 'business'],
      // saltedge_display: [],
      // topCategories: [],
      // subCategories: [],
    }
  },
  async created () {
    // this.refreshTransactions()
    this.refreshCategories()
  },
  methods: {
    async refreshCategories () {
      // console.log("refresh")
      this.loading = true
      this.categories = await api.getCategories()
      this.plaid_categories = await api.getPlaidCategories()
      this.saltedge_categories = await api.getSaltEdgeCategories()
      // var pdis = []
      // var sdis = []
      // this.plaid_categories.forEach(function (obj) {pdis.push(obj.hierarchy)})
      // this.saltedge_categories.forEach(function (obj) {sdis.push(obj.sub_category)})
      // this.plaid_display = [...new Set(pdis)]
      // this.saltedge_display = [...new Set(sdis)]
      // console.log(JSON.stringify(this.categories))
      // console.log(JSON.stringify(this.saltedge_categories))
      // var topArr = []
      // this.categories.forEach(function (obj) {topArr.push(obj.top_category)})
      // this.topCategories = [...new Set(topArr)]
      //console.log(JSON.stringify(this.topCategories))
      this.loading = false
    },
    // async refreshTransactions () {
    //   this.loading = true
    //   this.transactions = await api.getTransactions()
    //   this.loading = false
    // },
    openEditMenu: function(event, item) {
      event.preventDefault()
      this.editMenu = false
      this.x = event.clientX
      this.y = event.clientY
      this.editedIndex = this.transactions.indexOf(item)
      this.$nextTick(() => {
        this.editMenu = true
      })
    },
    // customSort(items, index, isDesc) { //don't need this since date is sorting OK
    //   items.sort((a, b) => {
    //     if (index === "date") {
    //       if (!isDesc) {
    //         return dateHelp.compare(a.date, b.date);
    //       } else {
    //         return dateHelp.compare(b.date, a.date);
    //       }
    //     } else {
    //       if (!isDesc) {
    //         return a[index] < b[index] ? -1 : 1;
    //       } else {
    //         return b[index] < a[index] ? -1 : 1;
    //       }
    //     }
    //   });
    //   return items;
    // },
    editCategory(cat) {
      var catToSave = cat
       //console.log("Before - " + this.transactions[this.editedIndex].category)
       this.transactions[this.editedIndex].category = catToSave
       api.updateTransaction(this.transactions[this.editedIndex].id, this.transactions[this.editedIndex])
       //console.log("After - " + this.transactions[this.editedIndex].category)
       this.editMenu = false
    },
    filtersub_category(topCat, categories) {
      var filterTop = topCat 
      //console.log(this.currenttop_category)
      var filtered=categories.filter(function(item){
        return item.top_category==filterTop
      })
      //console.log(JSON.stringify(filtered))
      return filtered
    },
    filterSaltSub(input){
      //console.log(JSON.stringify(filtered))
      var filtered = [...new Set (input.reduce((r, a) => r.concat(a.top_category), []))]
      console.log(JSON.stringify(filtered))
      return filtered
    },
    filtertop_category(categories) {
      var filtered=categories.filter(function(item){
        return item.top_category==item.sub_category
      })
      // console.log(JSON.stringify(filtered))
      return filtered
    },
  }
}
</script>