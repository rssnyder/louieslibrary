<template>
  <div id="login">
    <label>Username:</label>
    <input type='text' v-model="username" name='username' value=''><br>
    <label>Password:</label>
    <input type='password' v-model="password" name='password'><br>
    <input type='submit' value='Login' v-on:click="login">
  </div>
</template>

<script>
import Auth from '@/service/Auth.js';

export default {
  name: 'Login',
  data() {
    return {
      username: '',
      password: '',
      token: ''
    }
  },
  methods: {
    login() {
      let user = {
        username: this.username,
        password: this.password
      }

      Auth.login(user)
    },

    validate() {
      if (Auth.valid()) {
        console.log("USER IS VALID!!!")
      } else {
        console.log("__user not valid___")
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
