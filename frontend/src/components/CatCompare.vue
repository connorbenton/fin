<template>
  <v-row justify="center">
    <v-dialog v-model="dialog" persistent width="90%" max-width="1200">
      <v-card>
        <v-card-title class="justify-center">Manually Assign Matching Categories</v-card-title>
        <v-card-text>Unassigned matches will be left as Uncategorized</v-card-text>
        <v-data-table
          :headers="headers"
          :items="compareCats"
          :footer-props="{
            'items-per-page-options': [10, 20, 40, -1]
          }"
          :items-per-page="20"
          :key="redrawTable"
        >
          <template v-slot:item.assignedCat="{ item }">
            <v-card-text
              class="pa-0"
              @click="openEditMenu($event, item)"
            >{{catName(item.assignedCat, categories)}}</v-card-text>
          </template>
        </v-data-table>
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

        <br />
        <div class="text-center">
          <v-btn class="ma-2" color="success" @click="emitDone">Done</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<style>
.class1 {
  background-color: #bdbdbd;
}

.class2 {
}
</style>

<script>
import api from "@/api";
export default {
  data: () => ({
    reactive: true,
    headers: [
      { text: "Category to Assign", value: "category" },
      { text: "Assigned Category", value: "assignedCat" }
    ],
    y: 0,
    x: 0,
    compareCats: [],
    editMenu: false,
    editedIndex: -1,
    categories: [],
    topCategories: [],
    subCategories: [],
    dialog: false,
    redrawTable: '',
    console
  }),
  created() {
    //console.log(this.PLAID_PUBLIC_KEY)
    this.$store.subscribe((mutation, state) => {
      if (mutation.type === "newAssignCats") {
        this.compareCats = mutation.payload.compareCats;
        this.categories = mutation.payload.dbCats;
        // this.refreshCategories();
        this.dialog = true;
      }
    });
  },
  beforeDestroy() {
    // this.unsubscribe();
  },

  methods: {
    async refreshCategories() {
      this.categories = await api.getCategories();
      var topArr = [];
      this.categories.forEach(function(obj) {
        topArr.push(obj.topCategory);
      });
      this.topCategories = [...new Set(topArr)];
      this.catload = false;
    },
    catName(id, categories) {
      if (!id || 0 === id.length) { return "Uncategorized"; }
      var cat = categories.find(x => x.id === id).subCategory;
      return cat;
    },
    openEditMenu: function(event, item) {
      event.preventDefault();
      this.editMenu = false;
      this.x = event.clientX;
      this.y = event.clientY;
      this.editedIndex = this.compareCats.indexOf(item);
      // this.editMenu = true;
      this.$nextTick(() => {
        this.editMenu = true;
      });
    },
    editCategory(cat, categories) {
      var catToSave = cat;
      this.compareCats[this.editedIndex].assignedCat = categories.find(
        x => x.subCategory === catToSave
      ).id;
      this.compareCats[this.editedIndex].assignedCatName = categories.find(
        x => x.subCategory === catToSave
      ).subCategory;
      //  api.updateTransaction(this.transactions[this.editedIndex].id, this.transactions[this.editedIndex])
      this.editMenu = false;
      this.redrawTable = catToSave;
      // this.vm.$forceUpdate();
    },
    filterSubCategory(topCat, categories) {
      var filterTop = topCat;
      var filtered = categories.filter(function(item) {
        return item.topCategory == filterTop;
      });
      return filtered;
    },
    filterTopCategory(categories) {
      var filtered = categories.filter(function(item) {
        return item.topCategory == item.subCategory;
      });
      return filtered;
    },
    emitDone() {
      this.$store.commit("assignDone", this.compareCats);
      this.dialog = false;
    }
    // emitYes() {
    // this.$store.commit("answerGiven", "yes");
    // }
  },
  filters: {
    // pretty: function(value) {
    // return JSON.stringify(value, null, 2);
    // }
  },
  computed: {
    // getTrans1: {
    //   get() {
    //     return this.$store.getters.getTrans1;
    //   }
    // },
    // getTrans2: {
    //   get() {
    //     return this.$store.getters.getTrans2;
    //   }
    // },
    // dialog: {
    //   get() {
    //     return this.$store.state.compareMatch;
    //   }
    // }
  }
};
</script>