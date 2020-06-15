import Vue from 'vue';
import Vuetify, { colors } from 'vuetify/lib';

Vue.use(Vuetify);

const vuetify = new Vuetify({
  icons: {
    iconfont: 'md',
  },
  theme: {
    themes: {
      dark : {
        // primary: '#1e1e1e',
        // primary: '#b39670',
        // secondary: '#2d2d2d',
        // secondary: '#EEEEEE',
        accent: '#121212',
        // accent: '#dddddd',
      },
    },
  },
});

export default vuetify;
