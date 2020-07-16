<template>
  <div id="header">
    <nav>
      <router-link to="/">Home</router-link>
      <router-link to="/book/all">Browse Books</router-link>
      <router-link to="/request/all">Browse Requests</router-link>
      <router-link to="/login" class="right">Login</router-link>
      <router-link to="/signup" class="right">Signup</router-link>
      <button v-on:click="logout">Logout</button>
    </nav>
    {{ flash }}
  </div>
</template>

<script>
import Auth from '@/service/Auth.js';

export default {
  name: 'Header',
  props: {
    flash: String,
    username: String
  },
  methods: {
    logout() {
      Auth.logout()
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
nav {
  padding: 2px calc((100% - 1000px) / 2) 0;
}

nav {
  border-bottom: 1px solid #E4E5E7;
  padding-top: 17px;
  padding-bottom: 15px;
  background: #F7F9FA;
  height: 15px;
  color: #6A6C6F;
}

nav a, nav form button {
  display: inline-block;
  float: left;
  margin-right: 1.5em;
  position: relative;
}

nav form button {
  position: relative;
  top: 3px;
}

.right {
  float: right;
}

nav a[href^="/about"], nav form button {
  float: right;
  margin-left: 1.5em;
  margin-right: 0;
}

nav a[href^="/messages/"], nav form button {
  float: right;
  margin-left: 1.5em;
  margin-right: 0;
}

nav a.live {
  color: #34495E;
  cursor: default;
}

nav a.live:hover {
  text-decoration: none;
}

nav a.live:after {
    content: '';
    display: block;
    position: relative;
    left: calc(50% - 7px);
    top: 9px;
    width: 14px;
    height: 14px;
    background: #F7F9FA;
    border-left:1px solid #E4E5E7;
    border-bottom:1px solid #E4E5E7;
    -moz-transform:rotate(45deg);
    -webkit-transform:rotate(-45deg);
}

</style>
