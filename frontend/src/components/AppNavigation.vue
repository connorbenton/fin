<template>
    <span>
        <v-navigation-drawer
            app
            v-model="drawer"
            dark
            disable-resize-watcher
        >
            <v-list>
                <template v-for="(item, index) in items">
                    <v-list-item :key="index" :to="item.url">
                        <v-list-item-content>
                            {{ item.title }}
                        </v-list-item-content>
                    </v-list-item>
                    <v-divider :key="`divider-${index}`"></v-divider>
                </template>
            </v-list>
        </v-navigation-drawer>
        <v-app-bar dense dark>
            <v-app-bar-nav-icon
                class="hidden-md-and-up"
                @click="drawer = !drawer"
            ></v-app-bar-nav-icon>
            <v-spacer class="hidden-md-and-up"></v-spacer>
            <!-- <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/"
                data-cy="summaryBtn"
                >Summary</v-btn
            > -->
            <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/transactions"
                data-cy="transactionsBtn"
                >Transactions</v-btn
            >
            <!-- <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/bills"
                data-cy="billsBtn"
                >Bills</v-btn
            >
            <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/invesments"
                data-cy="investmentsBtn"
                >Investments</v-btn
            > -->
            <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/analysis"
                data-cy="analysisBtn"
                >Analysis</v-btn
            >
            <!-- <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/categories"
                data-cy="categoriesBtn"
                >Categories</v-btn
            > -->
            <!-- <v-spacer class="hidden-sm-and-down"></v-spacer> -->
            <v-btn
                text
                class="hidden-sm-and-down nav-menu"
                to="/accounts"
                data-cy="accountsBtn"
                >Accounts</v-btn
            >
            <v-btn
                v-if="!isProduction"
                text
                class="hidden-sm-and-down nav-menu"
                to="/database"
                data-cy="dbBtn"
                >Currency DB Editor</v-btn >
            <v-btn
                v-if="!isProduction"
                text
                class="hidden-sm-and-down nav-menu"
                to="/databasego"
                data-cy="dbgoBtn"
                >DB Editor</v-btn >
            

        </v-app-bar>
    </span>
</template>

<script>
export default {
    name: 'AppNavigation',
    data() {
        return {
            appTitle: 'Fintrack',
            drawer: false,
            itemsDev: [
                // { title: 'Summary', url: '/' },
                { title: 'Transactions', url: '/transactions' },
                // { title: 'Bills', url: '/bills' },
                // { title: 'Investments', url: '/investments' },
                { title: 'Analysis', url: '/analysis' },
                // { title: 'Categories', url: '/categories' },
                { title: 'Accounts', url: '/accounts' },
                { title: 'Currency DB Editor', url: '/database' },
                { title: 'DB Editor', url: '/databasego' },
            ],
            itemsProd: [
                { title: 'Transactions', url: '/transactions' },
                // { title: 'Bills', url: '/bills' },
                // { title: 'Investments', url: '/investments' },
                { title: 'Analysis', url: '/analysis' },
                // { title: 'Categories', url: '/categories' },
                { title: 'Accounts', url: '/accounts' },
            ],
            isProduction: false,
            items: [], 
        };
    },
    created() {
        this.isProduction = process.env.NODE_ENV === 'production';
        this.items = this.isProduction ? this.itemsProd : this.itemsDev;
    }
};
</script>

<style scoped>
a {
    color: white;
    text-decoration: none;
}
/* html {overflow-y: auto} */
</style>
