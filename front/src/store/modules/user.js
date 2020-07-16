import axios from "axios";
import { authHeader } from "../../service/Auth.js"

// Hold data for one user
const state = {
  user: {
    Username: 'green',
    Role: 'guy',
    Reviews: [
      {
        Username: "yoyo"
      }
    ]
  }
};

const getters = {
  user: (state) => state.user
};

const actions = {
  // Get data on one user
  async getUser(state, data) {
    const response = await axios.get('/user/' + data, {headers: authHeader()});
    state.commit('setUser', response.data);
  }
};

const mutations = {
  setUser: (state, user) => (state.user = user)
};

export default {
  state,
  getters,
  actions,
  mutations
}