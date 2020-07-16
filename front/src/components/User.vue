<template>
  <div id="user">
    <h1>{{ user.user.Username }}</h1>
    <h4>{{ user.user.Role }}</h4>
    <ul id="reviews">
      <li v-for="review in user.reviews" :key="review.ID">
        {{ review.Username }} - {{ review.Rating }} - {{ review.Review }}
      </li>
    </ul>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'Login',
  methods: {
    ...mapActions(['getUser'])
  },
  computed: mapGetters(['user']),
  created() {
    this.getUser(this.$route.params.username);
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
