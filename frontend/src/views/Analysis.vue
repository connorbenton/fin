<template>
  <v-content v-show="apiStateLoaded" :key="reloadedData">
    <v-col cols="12" align="center">
      <v-col cols="10" style="maxWidth:2400px" align="center" justify="center">
        <v-row justify="center" align="start" style="maxWidth: 1400px">
          <v-col cols="1" class="mr-4">
            <v-switch
              class="mt-0"
              v-model="showInvestment"
              label="Show Invest"
              @click="loadTxData(select, showInvestment, true)"
            ></v-switch>
          </v-col>
          <v-col class="flex-grow-0 flex-shrink-0">
            <v-btn
              class="ma-2"
              small
              @click="customFilter()"
              :disabled="isNaN(daysNum) || select !== 'Custom'"
            >{{prefix}} {{daysNum}} days</v-btn>
          </v-col>
          <v-col cols="2">
            <v-menu
              ref="menu0"
              v-model="menu0"
              :close-on-content-click="false"
              :nudge-right="40"
              transition="scale-transition"
              offset-y
              min-width="290px"
            >
              <template v-slot:activator="{ on }">
                <v-text-field
                  v-model="startDate"
                  outlined
                  dense
                  hint="Start date"
                  persistent-hint
                  prepend-icon="event"
                  :readonly="select !== 'Custom'"
                  :disabled="select !== 'Custom'"
                  @blur="startDate = setDateFromTextField(0)"
                  :rules="[rules.validDate]"
                  v-on="on"
                ></v-text-field>
              </template>
              <v-date-picker
                :disabled="select !== 'Custom'"
                v-model="startDatePicker"
                @input="setDateFromPicker(0)"
              ></v-date-picker>
            </v-menu>
          </v-col>
          <v-col cols="2">
            <v-menu
              ref="menu1"
              v-model="menu1"
              :close-on-content-click="false"
              :nudge-right="40"
              transition="scale-transition"
              offset-y
              min-width="290px"
            >
              <template v-slot:activator="{ on }">
                <v-text-field
                  v-model="endDate"
                  outlined
                  dense
                  hint="End date"
                  persistent-hint
                  prepend-icon="event"
                  :readonly="select !== 'Custom'"
                  :disabled="select !== 'Custom'"
                  @blur="endDate = setDateFromTextField(1)"
                  :rules="[rules.validDate]"
                  v-on="on"
                ></v-text-field>
              </template>
              <v-date-picker
                :disabled="select !== 'Custom'"
                v-model="endDatePicker"
                @input="setDateFromPicker(1)"
              ></v-date-picker>
            </v-menu>
          </v-col>
          <v-col cols="2">
            <v-select
              :items="dateRanges"
              outlined
              dense
              persistent-hint
              hint="Date Range"
              v-model="select"
              @change="loadTxData(select, showInvestment)"
            ></v-select>
          </v-col>
        </v-row>

        <v-row justify="center">
          <div class="treemap">
            <!-- The SVG structure is explicitly defined in the template with attributes derived from component data -->
            <svg
              :height="height + margin.bottom + margin.top"
              style="margin-left: 0px; margin-right: 0px"
              :width="width + margin.left + margin.right"
              xmlns="http://www.w3.org/2000/svg"
            >
              <g style="shape-rendering: crispEdges;" transform="translate(0,40)">
                <!-- We can use Vue transitions too! -->
                <transition-group name="list" tag="g" class="depth" v-if="selectedNode">
                  <!-- Generate each of the visible squares at a given zoom level (the current selected node) -->
                  <g
                    class="children"
                    v-for="(children, index) in selectedNode._children"
                    :key="'c_' + children.id"
                  >
                    <!-- Generate the children squares (only visible on hover of a square) -->
                    <rect
                      v-for="child in children._children"
                      class="child"
                      :id="child.id"
                      :key="child.id"
                      :height="y(child.y1) - y(child.y0)"
                      :width="x(child.x1) - x(child.x0)"
                      :x="x(child.x0)"
                      :y="y(child.y0)"
                    />

                    <rect
                      :class="'parent ' + (isDark?'darkRect':'lightRect')"
                      v-on:click="selectNode"
                      :id="children.id"
                      :key="children.id"
                      :x="x(children.x0)"
                      :y="y(children.y0)"
                      :width="x(children.x1 - children.x0 + children.parent.x0)"
                      :height="y(children.y1 - children.y0 + children.parent.y0)"
                      :style="{ fill: color(index) }"
                    >
                      <!-- The title attribute -->
                      <title>
                        {{ children.data.name }} |
                        {{ formatBalance(children.value) }} |
                        {{children.data.percent}} |
                        {{children.data.count}} Tx |
                        {{ formatBalance(children.data.per30)}} per 30
                      </title>
                    </rect>
                    <foreignObject
                      id="popup"
                      class="popup"
                      :x="x(children.x0) + positionPopupX(x(children.x1 - children.x0 + children.parent.x0), y(children.y1 - children.y0 + children.parent.y0))"
                      :y="y(children.y0) + positionPopupY(x(children.x1 - children.x0 + children.parent.x0), y(children.y1 - children.y0 + children.parent.y0))"
                      width="180"
                      height="30"
                      visibility="hidden"
                    >
                      <div xmlns="http://www.w3.org/1999/xhtml">
                        <v-col justify="start" cols="12" class="pa-0">
                          <v-btn
                            justify="start"
                            @click="ClickedCat(children)"
                            small
                          >See Transactions</v-btn>
                        </v-col>
                      </div>
                    </foreignObject>

                    <text
                      v-if="y(children.y1 - children.y0 + children.parent.y0) > 23"
                      dy="1em"
                      :key="'t_' + children.id + 'name'"
                      :x="x(children.x0) + 6"
                      :y="y(children.y0) + 6"
                      style="fill-opacity: 1;"
                    >{{ children.data.name }}</text>

                    <text
                      v-if="y(children.y1 - children.y0 + children.parent.y0) > 55"
                      dy="2.25em"
                      :key="'t_' + children.id + 'value'"
                      :x="x(children.x0) + 6"
                      :y="y(children.y0) + 6"
                      style="fill-opacity: 1"
                      font-weight="bold"
                    >{{ formatBalance(children.data.value) }} - {{children.data.percent}}</text>

                    <text
                      v-if="x(children.x1 - children.x0 + children.parent.x0) > 120 &&
              y(children.y1 - children.y0 + children.parent.y0) > 55"
                      dy="3.5em"
                      :key="'t_' + children.id + 'count'"
                      :x="x(children.x0) + 6"
                      :y="y(children.y0) + 6"
                      style="fill-opacity: 1;"
                    >{{ children.data.count}} Transactions</text>
                    <text
                      v-else-if="y(children.y1 - children.y0 + children.parent.y0) > 55"
                      dy="3.5em"
                      :key="'t_' + children.id + 'count_short'"
                      :x="x(children.x0) + 6"
                      :y="y(children.y0) + 6"
                      style="fill-opacity: 1;"
                    >{{ children.data.count}} Tx</text>
                    <text
                      v-if="y(children.y1 - children.y0 + children.parent.y0) > 78"
                      dy="4.75em"
                      :key="'t_' + children.id + 'per30'"
                      :x="x(children.x0) + 6"
                      :y="y(children.y0) + 6"
                      style="fill-opacity: 1;"
                    >{{ formatBalance(children.data.per30)}} per 30</text>
                  </g>
                </transition-group>

                <!-- The top most element, representing the previous node -->
                <g class="grandparent">
                  <rect
                    :height="margin.top"
                    :width="width"
                    :y="(margin.top * -1)"
                    v-on:click="selectNode"
                    :id="parentId"
                  />

                  <!-- The visible square text element with the id (basically a breadcumb, if you will) -->
                  <text
                    dy=".65em"
                    x="6"
                    y="-24"
                    v-if="selectedNode.data != undefined && selectedNode.id != 'Transactions by Category'"
                  >
                    {{ selectedNode.id.replace('Transactions by Category.','') }} - {{selectedNode.data.count}} Transactions
                    ({{selectedNode.data.trueCount}} Total) - {{formatBalance(selectedNode.data.value)}}
                    | {{formatBalance(selectedNode.data.per30)}} per 30 | {{formatBalance(rootNode.data.income_total)}} Income
                  </text>
                  <text
                    dy=".65em"
                    x="6"
                    y="-24"
                    v-if="selectedNode.id == 'Transactions by Category'"
                  >
                    Overall - {{selectedNode.data.count}} Transactions
                    ({{selectedNode.data.trueCount}} Total) - {{formatBalance(selectedNode.data.value)}}
                    | {{formatBalance(selectedNode.data.per30)}} per 30 | {{formatBalance(rootNode.data.income_total)}} Income
                  </text>
                  <foreignObject
                    id="popup"
                    class="popup"
                    x="600"
                    y="-34"
                    width="180"
                    height="30"
                    visibility="hidden"
                  >
                    <div xmlns="http://www.w3.org/1999/xhtml">
                      <v-col justify="start" cols="12" class="pa-0">
                        <v-btn
                          justify="start"
                          @click="ClickedCat(selectedNode)"
                          small
                        >See Transactions</v-btn>
                      </v-col>
                    </div>
                  </foreignObject>
                </g>
              </g>
            </svg>
          </div>
        </v-row>
        <!-- </v-card> -->
      </v-col>
    </v-col>
    <v-dialog v-model="dialog" @input="v => v || refreshTrees()">
      <v-card>
        <TransactionsTable v-bind:transactionsToDisplay="displayedTransactions"></TransactionsTable>
        <!-- @changed="hideDialog()" -->
      </v-card>
    </v-dialog>
  </v-content>
</template>

<script>
import TransactionsTable from "../components/TransactionsTable.vue";
import * as d3 from "d3";
import api from "@/api";
import moment from "moment";
import vuetify from "../plugins/vuetify";
export default {
  name: "Treemap",
  // the component's data
  data: vm => ({
    isDark: null,
    showInvestment: false,
    reloadedData: 0,
    apiStateLoaded: "",
    dialog: false,
    jsonData: null,
    xMenu: 0,
    yMenu: 0,
    rootNode: {},
    analysisTrees: [],
    treeJSON: {},
    transactionsTree: {},
    transactionsTreeNoInvest: {},
    transactionsTreeDisplayed: {},
    categories: [],
    accounts: [],
    transactions: [],
    displayedTransactions: [],
    date: new Date().toISOString().substr(0, 10),
    baseCurrency:
      process.env.VUE_APP_BASE_CURRENCY || window._env_.BASE_CURRENCY,
    isCustomFilterApplied: false,
    select: "Last 30 Days",
    dateRanges: [
      "Custom",
      "Last 30 Days",
      "This Month",
      "Last Month",
      "Last 6 Months",
      "This Year",
      "Last Year",
      "From Beginning"
    ],
    currentDate: new Date().toISOString().substr(0, 10),
    startDate: new Date().toISOString().substr(0, 10),
    startDatePicker: new Date().toISOString().substr(0, 10),
    endDate: new Date().toISOString().substr(0, 10),
    endDatePicker: new Date().toISOString().substr(0, 10),
    margin: {
      top: 40,
      // top: 0,
      right: 0,
      bottom: 0,
      left: 0
    },
    width: 960,
    height: 930,
    selected: null,
    color: {},
    menu0: false,
    menu1: false,
    rules: {
      validDate: value =>
        moment(value, "YYYY-MM-DD", true).isValid() || "Invalid date"
    }
  }),
  // You can do whatever when the selected node changes
  watch: {
    "$store.state.isDark": function() {
      this.isDark = this.$store.state.isDark;
    }
    // selectedNode(newData, oldData) {
    // console.log("The selected node changed...");
    // }
  },
  //Checks if API has loaded into Vuex store, then loads data
  async created() {
    this.startDate = moment();
    this.startDate = this.startDate.subtract(29, "days");
    this.startDate = this.startDate.format("YYYY-MM-DD");
    this.startDatePicker = this.startDate;
    this.apiStateLoaded = this.$store.state.apiStateLoaded;
    if (this.apiStateLoaded) {
      this.initialImportData();
    }
    this.unsub = this.$store.subscribe(async (mutation, state) => {
      if (mutation.type === "doneLoading") {
        if (mutation.payload) {
          this.initialImportData();
        }
        this.apiStateLoaded = mutation.payload;
      }
      if (mutation.type === "setReloading") {
        if (mutation.payload) {
          this.initialImportData();

          this.loadTxData(this.select, this.showInvestment);
          this.reloadedData += 1;
        }
      }
    });
  },
  mounted() {
    window.addEventListener("resize", this.setTreemapSize);
    this.isDark = this.$vuetify.theme.dark;
    let string = this.isDark ? "darkRect" : "lightRect";
  },
  beforeDestroy() {
    this.unsub();
    window.removeEventListener("resize", this.setTreemapSize);
  },
  computed: {
    // The grandparent id
    daysNum() {
      let a = moment(this.startDate);
      let b = moment(this.endDate);
      let c = b.diff(a, "days") + 1;
      if (isNaN(c)) return "＿";
      return c;
    },
    prefix() {
      if (this.select === "Custom" && !this.isCustomFilterApplied)
        return "Custom Filter";
      return "Showing";
    },
    parentId() {
      if (
        this.selectedNode.parent === undefined ||
        this.selectedNode.parent === null
      ) {
        return this.selectedNode.id;
      } else {
        return this.selectedNode.parent.id;
      }
    },
    // Returns the x position within the current domain
    // Maybe it can be replaced by a vuejs method
    x() {
      return d3
        .scaleLinear()
        .domain([0, this.width])
        .range([0, this.width]);
    },
    // Returns the y position within the current domain
    // Maybe it can be replaced by a vuejs method
    y() {
      return d3
        .scaleLinear()
        .domain([0, this.height - this.margin.top - this.margin.bottom])
        .range([0, this.height - this.margin.top - this.margin.bottom]);
    },
    // The D3 function that recursively subdivides an area into rectangles
    treemap() {
      return d3
        .treemap()
        .size([this.width, this.height])
        .round(true)
        .paddingInner(0);
    },
    // The current selected node
    selectedNode() {
      let node = null;
      if (this.selected) {
        let nd = this.getNodeById(this.rootNode, this.selected, this);
        if (nd._children) {
          node = nd;
        } else {
          node = nd.parent;
        }
      } else {
        node = this.rootNode;
      }
      // Recalculates the y and x domains
      this.x.domain([node.x0, node.x0 + (node.x1 - node.x0)]);
      this.y.domain([node.y0, node.y0 + (node.y1 - node.y0)]);
      return node;
    }
  },
  methods: {
    setTreemapSize() {
      this.height = window.innerHeight * 0.7;
      this.width = window.innerWidth * 0.8;
      this.reloadData();
    },
    //Initializes data and then draws treemap
    async initialImportData() {
      this.categories = this.$store.getters.getAllCategories;
      this.accounts = this.$store.getters.getAllAccounts;
      this.analysisTrees = this.$store.getters.getAllTrees;
      this.color = d3.scaleOrdinal(d3.schemeCategory10);
      this.loadTxData(this.select);
    },
    reloadData() {
      this.rootNode = {};
      this.jsonData = this.transactionsTreeDisplayed;
      this.initialize();
      this.accumulate(this.rootNode, this);
      this.treemap(this.rootNode);
    },
    customFilter() {
      this.$store.dispatch("customFilter", {
        startFromPage: this.startDate,
        endFromPage: this.endDate
      });
    },
    refreshTrees() {
      this.$store.dispatch("getTrees");
    },

    log(item) {
      console.log(item);
    },
    hideDialog() {
      this.dialog = false;
    },
    loadTxData(range, showInvestment, toggleFinancial = 0) {
      if (toggleFinancial) this.showInvestment = !showInvestment;
      this.selected = this.rootNode.id;
      switch (range) {
        case "Custom":
          this.treeJSON = this.analysisTrees.find(x => x.name === "custom");
          break;
        case "Last 30 Days":
          this.treeJSON = this.analysisTrees.find(x => x.name === "last30");
          break;
        case "This Month":
          this.treeJSON = this.analysisTrees.find(x => x.name === "thisMonth");
          break;
        case "Last Month":
          this.treeJSON = this.analysisTrees.find(x => x.name === "lastMonth");
          break;
        case "Last 6 Months":
          this.treeJSON = this.analysisTrees.find(
            x => x.name === "last6Months"
          );
          break;
        case "This Year":
          this.treeJSON = this.analysisTrees.find(x => x.name === "thisYear");
          break;
        case "Last Year":
          this.treeJSON = this.analysisTrees.find(x => x.name === "lastYear");
          break;
        case "From Beginning":
          this.treeJSON = this.analysisTrees.find(
            x => x.name === "fromBeginning"
          );
      }

      this.height = window.innerHeight * 0.7;
      this.width = window.innerWidth * 0.8;

      this.startDate = this.treeJSON.first_date;
      this.endDate = this.treeJSON.last_date;

      this.startDatePicker = this.startDate;
      this.endDatePicker = this.endDate;

      this.transactionsTree = JSON.parse(this.treeJSON.data);
      this.transactionsTreeNoInvest = JSON.parse(this.treeJSON.data_no_invest);
      this.transactions = this.$store.state.transactions;
      this.displayedTransactions = this.transactions;
      this.transactionsTreeDisplayed = {};
      if (this.showInvestment)
        Object.assign(this.transactionsTreeDisplayed, this.transactionsTree);
      else
        Object.assign(
          this.transactionsTreeDisplayed,
          this.transactionsTreeNoInvest
        );
      this.reloadedData += 1;
      this.reloadData();
      // this.reloadData();
    },
    ClickedCat(item) {
      const start = moment(this.startDate, "YYYY-MM-DD", true);
      const end = moment(this.endDate, "YYYY-MM-DD", true);
      if (item.depth === 0) {
        this.displayedTransactions = this.transactions.filter(trans => {
          return moment(trans.date).isBetween(start, end, undefined, "[]");
        });
      } else {
        let CatsToFilter = [];
        if (item.children)
          item.children.forEach(child => CatsToFilter.push(child.data.dbID));
        CatsToFilter.push(item.data.dbID);
        this.displayedTransactions = this.transactions.filter(trans => {
          const inrange = moment(trans.date).isBetween(
            start,
            end,
            undefined,
            "[]"
          );
          return CatsToFilter.includes(trans.category) && inrange;
        });
      }
      this.dialog = true;
    },
    //Positions popups in X & Y to make sure they don't occlude the title/label
    positionPopupX(width, height) {
      if (width < 150 && height > 110) return -170;
      if (width < 150 && height < 110) return -170;
      if (width > 150 && height < 110) return 120;
      return -5;
    },
    positionPopupY(width, height) {
      if (width > 150 && height > 110) return 75;
      else if (width > 150 && height < 110) return 10;
      return 5;
    },
    //Format balance to display correctly as currency
    formatBalance(bal) {
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: this.baseCurrency
      }).format(bal);
    },
    setDateFromPicker(index) {
      if (index === 0) {
        this.startDate = this.startDatePicker;
        this.menu0 = false;
      }
      if (index === 1) {
        this.endDate = this.endDatePicker;
        this.menu1 = false;
      }
    },
    setDateFromTextField(index) {
      if (index === 0) {
        if (moment(this.startDate, "YYYY-MM-DD", true).isValid()) {
          this.startDatePicker = this.startDate;
        }
        return this.startDate;
      }
      if (index === 1) {
        if (moment(this.endDate, "YYYY-MM-DD", true).isValid()) {
          this.endDatePicker = this.endDate;
        }
        return this.endDate;
      }
    },
    formatDate(date) {
      if (!date) return null;

      const [year, month, day] = date.split("-");
      return `${month}/${day}/${year}`;
    },
    // Called once, to create the hierarchical data representation
    initialize() {
      let that = this;
      if (that.jsonData) {
        let counter = 0;
        that.rootNode = d3
          .hierarchy(that.jsonData)
          .eachBefore(function(d) {
            d.id = (d.parent ? d.parent.id + "." : "") + d.data.name;
          })
          .sum(d => d.value)
          .sort(function(a, b) {
            return b.height - a.height || b.value - a.value;
          });
        that.rootNode.x = that.rootNode.y = 0;
        that.rootNode.x1 = that.width;
        that.rootNode.y1 = that.height;
        that.rootNode.depth = 0;
      }
    },
    // Calculates the accumulated value (sum of children values) of a node - its weight,
    // represented afterwards in the area of its square
    accumulate(d, context) {
      d._children = d.children;
      if (d._children) {
        d.value = d._children.reduce(function(p, v) {
          return p + context.accumulate(v, context);
        }, 0);
        return d.value;
      } else {
        return d.value;
      }
    },
    // Helper method - gets a node by its id attribute
    getNodeById(node, id, context) {
      if (node.id === id) {
        return node;
      } else if (node._children) {
        for (var i = 0; i < node._children.length; i++) {
          var nd = context.getNodeById(node._children[i], id, context);
          if (nd) {
            return nd;
          }
        }
      }
    },
    // which fires the computed selectedNode, which in turn finds the Node by the id of the square clicked
    // and the template reflects the changes
    selectNode(event) {
      this.selected = event.target.id;
    }
  }
};
</script>

<style scoped>
text {
  pointer-events: none;
}
.grandparent text {
  font-weight: bold;
}
rect {
  fill: none;
  fill-opacity: 0.7;
}
rect.lightRect {
  stroke: #fff;
}
rect.darkRect {
  stroke: #000;
}

rect.parent,
.grandparent rect {
  stroke-width: 2px;
}
.grandparent rect {
  fill: orange;
}
.grandparent:hover rect {
  fill: #ee9700;
}
.grandparent:hover foreignObject.popup {
  visibility: visible;
}
.children rect.parent,
.grandparent rect {
  cursor: pointer;
}
.children rect.parent {
  fill: #bbb;
}
rect.child {
  fill: #bbb;
}
.children:hover foreignObject.popup {
  visibility: visible;
}
.children text {
  font-size: 0.8em;
}
.grandparent text {
  font-size: 0.9em;
}
.list-enter-active,
.list-leave-active {
  transition: all 1s;
}
.list-enter, .list-leave-to /* .list-leave-active for <2.1.8 */ {
  opacity: 0;
}
</style>