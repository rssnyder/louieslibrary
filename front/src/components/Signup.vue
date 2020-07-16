<template>
  <div id="login">
    <label>Username:</label>
    <input type='text' v-model="username" name='username' value=''><br>
    <label>Email:</label>
    <input type='text' v-model="email" name='email' value=''><br>
    <label>Invite Code:</label>
    <input type='text' v-model="code" name='code' value=''><br>
    <label>Password:</label>
    <input type='password' v-model="password" name='password'><br>
    <input type='submit' value='Singup' v-on:click="signup">
  </div>
</template>

<script>
import Auth from '@/service/Auth.js';

export default {
  name: 'Login',
  data() {
    return {
      username: '',
      email: '',
      code: '',
      password: ''
    }
  },
  methods: {
    signup() {
      let user = {
        username: this.username,
        email: this.email,
        code: this.code,
        password: this.password
      }

      if (Auth.signup(user)) {
        window.location = "/login";
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
