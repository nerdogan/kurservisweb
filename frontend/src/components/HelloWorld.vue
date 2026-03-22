<script setup lang="ts">
defineProps<{
  msg: string;
}>();
</script>

<template>
  <div class="greetings">
    <h1 class="green">{{ msg }}</h1>
    <h3>
      Altın ve Gümüş USD Hesaplama <br>
      <a href="https://kilitcisakir.com/" target="_blank" rel="noopener">Kilitçi Şakir Web Sitesi</a> <br>
      <a href="https://kilitcisakir.com/katalog/" target="_blank" rel="noopener">Kilitçi Şakir Katalog</a><br>
      ...
    </h3>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h3 {
    text-align: left;
  }
}
</style>
