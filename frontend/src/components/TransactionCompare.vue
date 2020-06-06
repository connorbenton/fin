<template>
  <v-row justify="center">
    <v-dialog 
    v-model="dialog" 
    persistent 
    width="90%"
    max-width="1200">
      <v-card>
        <v-card-title class="justify-center">Are these transactions duplicate?</v-card-title>
        <v-row class="d-flex justify-space-around ma-2">
        <v-card class="pa-6"> 
        <pre>{{ getTrans1 | pretty }}</pre>
        </v-card>
        <v-card class="pa-6"> 
        <pre>{{ getTrans2 | pretty }}</pre> </v-col>
        </v-card>
        </v-row>
        <br>
        <div class="text-center">
        <v-btn
        class="ma-2"
        color="error"
        @click="emitNo"
        >
        No
        </v-btn>
        <v-btn
        class="ma-2"
        color="success"
        @click="emitYes"
        >
        Yes
        </v-btn>
        </div>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
  export default {
    data: () => ({
      // sheet: false,
    }),
    methods: {
      emitNo() {
      this.$store.commit('answerGiven', false);
      },
      emitYes() {
      this.$store.commit('answerGiven', true);
      }
    },
    filters: {
      pretty: function(value) {
        return JSON.stringify(value, null, 2);
      }
    },
    computed: {
      getTrans1: {
        get () {
          return this.$store.getters.getTrans1
        }
      },
      getTrans2: {
        get () {
          return this.$store.getters.getTrans2
        }
      },
      dialog: {
        get () {
          return this.$store.state.compareMatch
        }
      }
    }
  }
</script>