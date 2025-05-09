<template>
  <main>
    <button class="btn btn-square" :class="{'btn-primary': count > 5}" @click="increment"><Arrow /></button>
    <p>{{ count }}</p>
    <p v-if="count <= 5">Lower or equal</p>
    <p v-else>Higher</p>
  </main>
  <div>
    
  </div>
</template>

<script setup lang="js">
import Arrow from "../components/icons/Arrow.vue"
import { ref } from "vue"

const count = ref(0)

function increment() {
  count.value = count.value + 1
}
</script>