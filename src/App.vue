<template>
    <v-app id="app">
        <app-navigation></app-navigation>

        <v-content transition="slide-x-transition">
            <router-view></router-view>
        </v-content>
        <transaction-compare></transaction-compare>
        <cat-compare></cat-compare>
    </v-app>
</template>

<script>
import AppNavigation from '@/components/AppNavigation';
import CatCompare from '@/components/CatCompare';
import TransactionCompare from '@/components/TransactionCompare';

export default {
    name: 'App',
    components: {
        AppNavigation,
        CatCompare,
        TransactionCompare
    }
};
</script>

<style></style>