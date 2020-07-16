import Vue from 'vue'
import Router from 'vue-router'
import Login from '../components/Login.vue'
import Signup from '../components/Signup.vue'
import { authGuard } from '../service/Auth.js'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      // component: Home,
      beforeEnter: authGuard
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/signup',
      name: 'signup',
      component: Signup
    },
    {
      path: '/user/:username',
      name: 'user',
      component: () => import('../components/User.vue'),
      beforeEnter: authGuard
    }
  ]
})