import Vue from 'vue';
import Vuex from 'vuex';
import user from './modules/user';

// Load vuex
Vue.use(Vuex);

// Create store
export default new Vuex.Store({
  modules: {
    user
  }
});