<template>
  <v-content v-if="apiStateLoaded" :key="reloadedData">
    <v-col cols="12" align="center">
      <v-col cols="10" style="maxWidth:1600px"> 
        <v-card height="800">
        <!-- <v-card> -->
          <v-row justify="center" align="start" style="maxWidth: 1200px">
            <!-- <v-col cols="12" sm="6" md="4"> -->
            <!-- <v-col cols="1" align="start"> -->
            <v-col cols="1" class="mr-4">
                <v-checkbox 
                class="mt-0" 
                v-model="showInvestment" 
                label="Show Investment"
                @click="loadTxData(select, showInvestment, true)"></v-checkbox>
            </v-col>
            <!-- <v-col cols="2" align="start"> -->
            <v-col class="flex-grow-0 flex-shrink-0">
              <v-btn
              class="ma-2"
              small
              @click="customFilter()"
              :disabled="isNaN(daysNum) || select !== 'Custom'"
              > {{prefix}} {{daysNum}} days
              </v-btn>
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
            <!-- <v-spacer></v-spacer> -->
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
                <!-- v-if="select === 'Custom'"  -->
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
              @change="loadTxData(select, showInvestment)"></v-select>
            </v-col>
          </v-row>

          <div class="treemap"
          
          >
            <!-- The SVG structure is explicitly defined in the template with attributes derived from component data -->
              <!-- :viewBox="[0.5, -30.5, width, height + 30]" -->
              <!-- :height="height"
              style="1argin-left: 0px; margin-right: 0px"
              :width="width"
              xmlns="http://www.w3.org/2000/svg" -->
            <svg
              :height="height + margin.bottom + margin.top"
              style="margin-left: 0px; margin-right: 0px"
              :width="width + margin.left + margin.right"
              xmlns="http://www.w3.org/2000/svg"
            >
              <g 
              style="shape-rendering: crispEdges;" 
              transform="translate(0,40)">
                <!-- We can use Vue transitions too! -->
                <transition-group name="list" tag="g" class="depth">
                  <!-- Generate each of the visible squares at a given zoom level (the current selected node) -->
                  <g
                    class="children"
                    v-for="(children, index) in selectedNode._children"
                    :key="'c_' + children.id"
                    v-if="selectedNode"
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

                    <!-- 
              The visible square rect element.
              You can attribute directly an event, that fires a method that changes the current node,
              restructuring the data tree, that reactivly gets reflected in the template.
                    -->
                    <rect
                      class="parent"
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
                      <title>{{ children.data.name }} | 
                        {{ formatBalance(children.value, 'USD') }} |
                        {{children.data.percent}} | 
                        {{children.data.count}} Tx</title>
                    </rect>

                    <!-- v-if="children.id.width > 170 && children.id.height > 93" -->
                    <!-- :x="x(children.x0) - positionPopupX(x(children.x1 - children.x0 + children.parent.x0), y(children.y1 - children.y0 + children.parent.y0))"  -->
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

                    <!-- <foreignObject 
            v-if="x(children.x1 - children.x0 + children.parent.x0) <= 170 ||
              y(children.y1 - children.y0 + children.parent.y0) <= 93" 
            id="popup"
            class="popup"
              :x="x(children.x0) - 170" 
              :y="y(children.y0) + 10" 
              width="180"
              height="30"
              visibility="hidden"
            >
              <div xmlns="http://www.w3.org/1999/xhtml"> 
                <v-col justify="start" cols=12 class="pa-0">
                <v-btn justify="start" small> See Transactions </v-btn>
                </v-col>
              </div>
                    </foreignObject>-->

                    <!-- The visible square text element with the title and value of the child node -->
                    <!-- :load="log(x(children.x1 - children.x0 + children.parent.x0))" -->
                    <!-- v-if="x(children.x1 - children.x0 + children.parent.x0) > 170"  -->
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
                    >{{ formatBalance(children.data.value, 'USD') }} - {{children.data.percent}}</text>

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
                    <!-- <foreignObject 
              :id="children.id + 'selector'"
              :key="children.id + 'selector'"
              :x="x(children.x0 + 6)" 
              :y="y(children.y0 + 58)" 
              width="180"
              height="30"
            >
              <div xmlns="http://www.w3.org/1999/xhtml"> 
                <v-col justify="start" cols=12 class="pa-0">
                <v-btn justify="start" small> See Transactions </v-btn>
                </v-col>
              </div>
                    </foreignObject>-->

                    <!-- <rect 
              :id="children.id + 'selector'"
              :key="children.id + 'selector'"
              :x="x(children.x0 + 6)" 
              :y="y(children.y0 + 58)" 
              rx="2"
              width="100px"
              height="20px"
              :style="{ fill: color(FAFA) }"
                    small>See Transactions</rect>-->
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
                    v-if="selectedNode.data != undefined"
                  >{{ selectedNode.id }} - {{selectedNode.data.count}} Transactions 
                  ({{selectedNode.data.trueCount}} Total) - {{formatBalance(selectedNode.data.value, 'USD')}}</text>
                    <!-- v-if="selectedNode.id != 'Transactions by Category'" -->
                  <foreignObject
                    id="popup"
                    class="popup"
                    x="580"
                    y="-34"
                    width="180"
                    height="30"
                    visibility="hidden"
                  >
                    <!-- <set attributeName="visibility" from="hidden" to="visible" begin="children.id.mouseover" end="children.id.mouseout"/> -->
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
        </v-card>
      </v-col>
    </v-col>
    <v-dialog v-model="dialog">
      <v-card>
        <TransactionsTable v-bind:transactionsToDisplay="displayedTransactions" @changed="hideDialog()"></TransactionsTable>
      </v-card>
    </v-dialog>
  </v-content>
</template>

<script>
import TransactionsTable from "../components/TransactionsTable.vue";
// import {scaleLinear, scaleOrdinal} from 'd3-scale';
// import {schemeCategory10} from 'd3-scale-chromatic';
// import {json} from 'd3-request';
// import {hierarchy, treemap} from 'd3-hierarchy';
import * as d3 from "d3";
import api from "@/api";
import moment from 'moment';
// To be explicit about which methods are from D3 let's wrap them around an object
// Is there a better way to do this?
// let d3 = {
//   scaleLinear: scaleLinear,
//   scaleOrdinal: scaleOrdinal,
//   schemeCategory10: schemeCategory10,
//   json: json,
//   hierarchy: hierarchy,
//   treemap: treemap
// }
export default {
  name: "Treemap",
  // the component's data
  data: vm => ({
    showInvestment: false,
    reloadedData: 0,
    apiStateLoaded: "",
    dialog: false,
    jsonData: null,
    xMenu: 0,
    yMenu: 0,
    rootNode: {},
    transactionsTree: {},
    transactionsTreeNoInvest: {},
    transactionsTreeDisplayed: {},
    categories: [],
    accounts: [],
    transactions: [],
    displayedTransactions: [],
    date: new Date().toISOString().substr(0, 10),

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
    // startDate: vm.formatDate(new Date().toISOString().substr(0, 10)),
    // startDate: new Date(new Date().getFullYear(), new Date().getMonth(), 1),
    // endDate: vm.formatDate(new Date().toISOString().substr(0, 10)),
    // endDate: vm.formatDate(new Date().toISOString().substr(0, 10)),
    // startDateText: '',
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
    height: 630,
    selected: null,
    color: {},
    menu0: false,
    menu1: false,
    rules: {
      // required: value => !!value || 'Required',
      validDate: value => moment(value, "YYYY-MM-DD", true).isValid() || 'Invalid date'
    }
      // validDate: value => {

        // const pattern = /(\d{4})-(\d{2})-(\d{2})/
        // return pattern.test(value) || 'Invalid date'
    // }
  }),
  // You can do whatever when the selected node changes
  watch: {
    // selectedNode(newData, oldData) {
      // console.log("The selected node changed...");
    // }
  },
  //Checks if API has loaded into Vuex store, then loads data
  async created() {
    this.startDate = moment();
    this.startDate = this.startDate.subtract(29, 'days');
    this.startDate = this.startDate.format('YYYY-MM-DD');
    this.startDatePicker = this.startDate;
    // this.startDatePicker = new Date(beginMonth).toISOString().substr(0, 10);
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

          // this.reloadData();
          this.loadTxData(this.select, this.showInvestment);
          this.reloadedData += 1;
        }
        // this.apiStateLoaded = mutation.payload;
      }
    });
  },
  beforeDestroy() {
    this.unsub();
  },
  // The reactive computed variables that fire rerenders
  computed: {
    // startDate() {
    //   let monthStart = new Date(this.currentDate.getTime())
    //   monthStart.setDate(1);
    //   return this.formatDate(new Date(monthStart).toISOString().substr(0, 10));
    // },
    // endDate() {
    //   return this.formatDate(this.currentDate.toISOString().substr(0, 10));
    // },
    // The grandparent id
    daysNum() {
      let a = moment(this.startDate);
      let b = moment(this.endDate);
      let c =  b.diff(a, 'days') + 1;
      if (isNaN(c)) return '＿';
      return c;
    },
    prefix() {
      if (this.select === 'Custom' && !this.isCustomFilterApplied) return 'Custom Filter'
      return 'Showing'
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
        // .rangeRound([0, 580])
          // .domain([0, 580])
          // .range([0, 580])
        .domain([0, this.width])
        .range([0, this.width]);
    },
    // Returns the y position within the current domain
    // Maybe it can be replaced by a vuejs method
    y() {
      return d3
        .scaleLinear()
        // .rangeRound([0, 590])
          // .domain([0, 590])
          // .range([0, 590])
        .domain([0, this.height - this.margin.top - this.margin.bottom])
        .range([0, this.height - this.margin.top - this.margin.bottom]);
    },
    // The D3 function that recursively subdivides an area into rectangles
    treemap() {
      return d3
        .treemap()
        .size([this.width, this.height])
        // .size([580, 590])
        // .round(false)
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
    //Initializes data and then draws treemap
    async initialImportData() {
      this.categories = this.$store.getters.getAllCategories;
      this.accounts = this.$store.getters.getAllAccounts;
      this.color = d3.scaleOrdinal(d3.schemeCategory10);
      // this.transactions = this.$store.getters.getAllTransactions;
      // this.transactions = this.$store.state.txData.txSets.last30TxSet;
      // // let x = this.$store.state.txData.txSets.last30TxSet;
      // this.transactionsTree = this.$store.state.txData.txTrees.last30TxTree;
      // // let y = this.$store.state.txData.txTrees.last30TxTree;
      // // let z = this.$store.state.transactions;
      // // await this.refreshTransactions();
      // this.displayedTransactions = this.transactions;
      // this.transactionsTreeDisplayed = this.transactionsTree;
      this.loadTxData(this.select);
      // this.jsonData = this.transactionsTreeDisplayed;
      // this.initialize();
      // this.accumulate(this.rootNode, this);
      // this.treemap(this.rootNode);
    },
    reloadData() {
      // this.categories = this.$store.getters.getAllCategories;
      // this.accounts = this.$store.getters.getAllAccounts;
      // this.transactions = this.$store.getters.getAllTransactions;
      // this.transactions = this.$store.state.txData.txSets.last30TxSet;
      // let x = this.$store.state.txData.txSets.last30TxSet;
      // this.transactionsTree = this.$store.state.txData.txTrees.last30TxTree;
      // let y = this.$store.state.txData.txTrees.last30TxTree;
      // let z = this.$store.state.transactions;
      // await this.refreshTransactions();
      this.rootNode = {};
      // this.color = d3.scaleOrdinal(d3.schemeCategory10);
      this.jsonData = this.transactionsTreeDisplayed;
      this.initialize();
      this.accumulate(this.rootNode, this);
      this.treemap(this.rootNode);
    },
    customFilter() {
      // this.$store.commit('setCustomRange', {start:this.startDate, end:this.endDate});
      this.$store.dispatch('customFilter', {startFromPage:this.startDate, endFromPage:this.endDate});
    },
    // toggleFinancial(showInvestment) {
      // this.transactionsTreeDisplayed = {};
      // if (!showInvestment) Object.assign(this.transactionsTreeDisplayed, this.transactionsTree);
      // else Object.assign(this.transactionsTreeDisplayed, this.transactionsTreeNoInvest);
      // this.showInvestment = !showInvestment;
      // this.loadTxData(this.select);
      // this.reloadedData += 1;
      // this.reloadData();
    // },
    log(item) {
      console.log(item);
    },
    hideDialog() {
      // this.reloadedData += 1;
      this.dialog = false;
      // this.reloadedData += 1;
      // vm.$forceUpdate;
    },
    loadTxData(range, showInvestment, toggleFinancial = 0) {
      if (toggleFinancial) this.showInvestment = !showInvestment;
      this.selected = this.rootNode.id;
      this.endDate = moment().format('YYYY-MM-DD');
      switch(range) {
        case "Custom":
          this.startDate = this.$store.state.customStart;
          this.endDate = this.$store.state.customEnd;
          this.transactions = this.$store.state.txData.txSets.customTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.customTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.customTxTreeNoInvest;
          break;
        case "Last 30 Days":
          this.transactions = this.$store.state.txData.txSets.last30TxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.last30TxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.last30TxTreeNoInvest;
          break;
        case "This Month":
          this.startDate = moment(this.endDate).startOf('month').format('YYYY-MM-DD');
          this.transactions = this.$store.state.txData.txSets.thisMonthTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.thisMonthTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.thisMonthTxTreeNoInvest;
          break;
        case "Last Month":
          this.startDate = moment(this.endDate).subtract(1, 'month').startOf('month').format('YYYY-MM-DD');
          this.endDate = moment(this.endDate).subtract(1, 'month').endOf('month').format('YYYY-MM-DD');
          this.transactions = this.$store.state.txData.txSets.lastMonthTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.lastMonthTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.lastMonthTxTreeNoInvest;
          break;
        case "Last 6 Months":
          this.startDate = moment(this.endDate).subtract(6, 'month').format('YYYY-MM-DD');
          this.transactions = this.$store.state.txData.txSets.lastSixMonthsTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.lastSixMonthsTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.lastSixMonthsTxTreeNoInvest;
          break;
        case "This Year":
          this.startDate = moment(this.endDate).startOf('year').format('YYYY-MM-DD');
          this.transactions = this.$store.state.txData.txSets.thisYearTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.thisYearTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.thisYearTxTreeNoInvest;
          break;
        case "Last Year":
          this.startDate = moment(this.endDate).subtract(1, 'year').startOf('year').format('YYYY-MM-DD');
          this.endDate = moment(this.endDate).subtract(1, 'year').endOf('year').format('YYYY-MM-DD');
          this.transactions = this.$store.state.txData.txSets.lastYearTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.lastYearTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.lastYearTxTreeNoInvest;
          break;
        case "From Beginning":
          // this.startDate = '2000-01-01';
          this.startDate = this.$store.state.txData.oldestDate;
          this.transactions = this.$store.state.txData.txSets.fromBeginningTxSet;
          this.transactionsTree = this.$store.state.txData.txTrees.fromBeginningTxTree;
          this.transactionsTreeNoInvest = this.$store.state.txData.txTrees.fromBeginningTxTreeNoInvest;
      }
      this.displayedTransactions = this.transactions;
      this.transactionsTreeDisplayed = {};
      if (this.showInvestment) Object.assign(this.transactionsTreeDisplayed, this.transactionsTree);
      else Object.assign(this.transactionsTreeDisplayed, this.transactionsTreeNoInvest);
      this.reloadedData += 1;
      this.reloadData();
      // this.reloadData();
    },
    ClickedCat(item) {
      if (item.depth === 0) this.displayedTransactions = this.transactions; 
      else {
      let CatsToFilter = [];
      if (item.children) item.children.forEach(child => CatsToFilter.push(child.data.dbID));
      CatsToFilter.push(item.data.dbID);
      this.displayedTransactions = this.transactions.filter(trans => {
        // CatsToFilter.indexOf(trans.category) > -1;
        // CatsToFilter.indexOf(trans.category) < 0;
        return CatsToFilter.includes(trans.category);
        // let x = CatsToFilter.indexOf(trans.category);
        // return x;
      })
      }
      // let y = this.displayedTransactions;
      this.dialog = true;

      // console.log(item);
      // console.log(CatsToFilter);
    },
    //Positions popups in X & Y to make sure they don't occlude the title/label
    positionPopupX(width, height) {
      if (width < 170 && height > 93) return -170;
      if (width < 170 && height < 93) return -170;
      if (width > 170 && height < 93) return 120;
      return -5;
    },
    positionPopupY(width, height) {
      if (width > 170 && height > 93) return 58;
      else if (width > 170 && height < 93) return 10;
      return 5;
    },
    //Format balance to display correctly as currency
    formatBalance(bal, code) {
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: code
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
      // if (!date) return null;
      // if (!date) {
      //   switch (item) {
      //     case 0:
      //       this.menu = false;
      //       return;
      //       break;
      //     case 1:
      //       this.menu1 = false;
      //       return;
      //   }
      // }
      // try {
        // if (moment(date, "YYYY-MM-DD", true).isValid()) {
          // pickerDate = date; 
        // }
        // return date;
        // return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`
        // let mydate = new Date(parts[0], parts[1] - 1, parts[2]);
        // return mydate.toISOString().substr(0, 10);
        // console.log(mydate.getUTCDate())
        // return new Date(mydate).toISOString().substr(0, 10);
      // }
      // catch {
        // return null;
        // switch (item) {
        //   case 0:
        //     // this.menu0 = false;
        //     return null;
        //     break;
        //   case 1:
        //     // this.menu1 = false;
        //     return null;
        // }
      // }
    },
    // dateFilteredTransactions() {},
    // buttonFilteredTransactions() {},
    // async refreshTransactions() {
    //   // this.transload = true
    //   this.categories = this.$store.getters.getAllCategories;
    //   this.accounts = this.$store.getters.getAllAccounts;
    //   // this.categories = await api.getCategories();
    //   // this.accounts = await api.getAccounts();
    //   let cats = this.categories.map(c => Object.assign({}, c));
    //   let accs = this.accounts.map(a => Object.assign({}, a));

    //   // let loadcats = this.$store.getters.getAllCategories;
    //   // let loadaccs = this.$store.getters.getAllAccounts;
    //   // let cats = loadcats.slice(0);
    //   // let accs = loadaccs.slice(0);
    //   // let y = this.$store.getters.getAllCategories;
    //   // this.transactions = await api.getTransactions();
    //   this.transactions = [...this.$store.getters.getAllTransactions];
    //   for (let i in this.transactions) {
    //     this.transactions[i].accName = accs.find(
    //       x => x.account_id === this.transactions[i].account_id
    //     ).name;
    //     let matchCat = cats.find(x => x.id === this.transactions[i].category);
    //     this.transactions[i].catName = matchCat.subCategory;
    //     if (matchCat.count === undefined) matchCat.count = 0;
    //     matchCat.count = matchCat.count + 1;
    //     if (matchCat.total === undefined) matchCat.total = 0;
    //     matchCat.total =
    //       matchCat.total + parseFloat(this.transactions[i].normalized_amount);

    //     // matchCat.count = matchCat.count + 1;
    //     // matchCat.total = matchCat.total + parseFloat(this.transactions[i].amount);
    //   }

    //   this.transactionsTree = {};

    //   this.transactionsTree.name = "Transactions by Category";
    //   this.transactionsTree.children = [];
    //   this.transactionsTree.value = 0;
    //   for (let j in cats) {
    //     if (cats[j].excludeFromAnalysis || cats[j].topCategory === "Income")
    //       continue;
    //     if (cats[j].subCategory === cats[j].topCategory) {
    //       let newChild = {};
    //       newChild.name = cats[j].topCategory;
    //       newChild.children = [];
    //       newChild.value = 0;
    //       newChild.count = 0;
    //       // newChild.dbID = '';

    //       let children = cats.filter(
    //         x => x.topCategory === cats[j].topCategory
    //       );
    //       for (let k in children) {
    //         let subCatChildToPush = {};
    //         if (children[k].subCategory === children[k].topCategory) {
    //           subCatChildToPush.name = children[k].subCategory + ` (General)`;
    //           newChild.dbID = children[k].id;
    //         } else {
    //           subCatChildToPush.name = children[k].subCategory;
    //         }

    //         subCatChildToPush.dbID = children[k].id;
    //         let value = children[k].total;
    //         let count = children[k].count;
    //         subCatChildToPush.value =
    //           value === undefined ? 0 : -1 * value;
    //         subCatChildToPush.count = count === undefined ? 0 : count;
    //         subCatChildToPush.percent = "";
    //         newChild.value = newChild.value + subCatChildToPush.value;
    //         newChild.count = newChild.count + subCatChildToPush.count;
    //         //  subCatChildToPush.count = children[k].count;
    //         newChild.children.push(subCatChildToPush);
    //       }
    //       // newChild.children.map(obj => ({...obj, percent: (obj.value / newChild.value).toFixed(1)+"%"}));
    //       newChild.children.forEach(function(child) {
    //         // newChild.children[k].percent = '';
    //         child.percent =
    //           ((child.value / newChild.value) * 100).toFixed(1) + "%";
    //         let x = "y";
    //       });
    //       this.transactionsTree.children.push(newChild);
    //       this.transactionsTree.value =
    //         this.transactionsTree.value + newChild.value;
    //     }
    //   }
    //   let total = this.transactionsTree.value;
    //   this.transactionsTree.children.forEach(child => {
    //     child.percent = ((child.value / total) * 100).toFixed(1) + "%";
    //   });
    //   let a = this.transactionsTree;
    //   let b = "ha";
    // },
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
          // .eachAfter(function (d) { d.parent.count = (d.parent.count === undefined) ? d.count : d.parent.count + d.count})
          .sum(d => d.value)
          // .eachAfter(function (d) { if (d.parent != null) {
          //   d.count = Number.isInteger(d.data.count) ? d.data.count : 0;
          //   let total = Number.isInteger(d.parent.count) ? d.parent.count : 0;
          //   d.parent.count = total + d.count;
          //   d.parent.data.count = total + d.count;
          //   }})
          // .eachAfter(function (d) { if (d.parent != null) {
          //   d.percent = (d.value / d.parent.value * 100).toFixed(1)+"%";
          //   }})
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
        // d.count = d._children.reduce(function (p, v) { return p + context.accumulate(v, context) }, 0)
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
/* .treemap {
    width: 960px;
    height: 630px;
    overflow: scroll;
} */
text {
  pointer-events: none;
}
.grandparent text {
  font-weight: bold;
}
rect {
  fill: none;
  stroke: #fff;
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
  fill-opacity: 0.5;
}
rect.child {
  fill: #bbb;
}
.children:hover foreignObject.popup {
  visibility: visible;
}
/* .children:hover rect.child { 
  fill: #bbb;
} */
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