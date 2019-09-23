<template>

<v-content>
 <v-card>
     <v-text-field
        v-model="search"
        prepend-inner-icon="search"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
  <v-data-table
    :headers="headers"
    :items="transactions"
    :search="search"
    multi-sort
    class="elevation-1"
  >
      <template v-slot:item.category="{ item }">
      <v-card-text class="pa-0" @click="openEditMenu($event, item)">{{item.category}}</v-card-text>
      </template> 
  </v-data-table>

      </v-card>

      <v-menu v-model="editMenu" :position-x="x" :position-y="y" absolute offset-y>
          <v-list class="pa-0">
        <v-menu offset-x open-on-hover
              v-for="(cat, index) in filterTopCategory(categories)"
              :key="index"
        >
          <template v-slot:activator="{ on }">
            <v-hover v-slot:default="{ hover }">
            <v-list-item
              @click="editCategory(cat.topCategory)"
              v-on="on"
            :class="`${hover? 'class1': 'class2'}`"
            >
            <v-list-item-title 
            >
              {{cat.topCategory}}
            </v-list-item-title>
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
              @click="editCategory(subcat.subCategory)"
              :class="`${hover? 'class1': 'class2'}`"
            >
            <v-list-item-title
            >
            {{subcat.subCategory}}
            </v-list-item-title>
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
  background-color:#BDBDBD;
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
      currentTopCategory: null,
      offset: true,
      headers: [
        {text: 'Date', value: 'date' },
        {text: 'Description', value: 'description' },
        {text: 'Category', value: 'category' },
        {text: 'Amount', value: 'amount' },
        {text: 'Account', value: 'accountName' },
      ],
      loading: false,
      transactions: [],
      categories: [],
      topCategories: [],
      subCategories: [],
    }
  },
  async created () {
    this.refreshTransactions()
    this.refreshCategories()
  },
  methods: {
    async refreshCategories () {
      this.loading = true
      this.categories = await api.getCategories()
      //console.log(JSON.stringify(this.categories))
      var topArr = []
      this.categories.forEach(function (obj) {topArr.push(obj.topCategory)})
      this.topCategories = [...new Set(topArr)]
      //console.log(JSON.stringify(this.topCategories))
      this.loading = false
    },
    async refreshTransactions () {
      this.loading = true
      this.transactions = await api.getTransactions()
      this.loading = false
    },
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
    editCategory(cat) {
      var catToSave = cat
       //console.log("Before - " + this.transactions[this.editedIndex].category)
       this.transactions[this.editedIndex].category = catToSave
       api.updateTransaction(this.transactions[this.editedIndex].id, this.transactions[this.editedIndex])
       //console.log("After - " + this.transactions[this.editedIndex].category)
       this.editMenu = false
    },
    filterSubCategory(topCat, categories) {
      var filterTop = topCat 
      //console.log(this.currentTopCategory)
      var filtered=categories.filter(function(item){
        return item.topCategory==filterTop
      })
      //console.log(JSON.stringify(filtered))
      return filtered
    },
    filterTopCategory(categories) {
      var filtered=categories.filter(function(item){
        return item.topCategory==item.subCategory
      })
      //console.log(JSON.stringify(filtered))
      return filtered
    },
  }
}
</script>