<!-- File: src/pages/Stats.vue -->
<template>
  <div class="max-w-2xl p-6">
    <h1 class="text-3xl font-bold mb-6 text-left">Статистика пользователя</h1>

    <div v-if="error" class="text-error mb-4">
      {{ error }}
    </div>

    <div v-else-if="!loaded" class="text-center text-gray-500">
      Загрузка…
    </div>

    <div v-else class="space-y-4 text-lg text-left">
      <p>
        <span class="font-semibold">Общий рейтинг:</span>
        {{ stats.rating }}
      </p>
      <p>
        <span class="font-semibold">Количество паков:</span>
        {{ stats.packs_created }}
      </p>
      <p>
        <span class="font-semibold">Паков изучено:</span>
        {{ stats.packs_mastered }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const stats = ref({
  rating: 0,
  packs_created: 0,
  packs_mastered: 0
})
const loaded = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    const res = await fetch('/api/stats', {
      credentials: 'include'
    })
    if (res.status === 401) {
      throw new Error('Пользователь не авторизован')
    }
    if (!res.ok) {
      throw new Error(`Ошибка при загрузке: ${res.status}`)
    }
    const data = await res.json()
    // ожидаем, что бэкенд отдаёт поле exactly равно GetUserStatsRow:
    // { rating, packs_created, packs_mastered }
    stats.value = {
      rating: data.rating,
      packs_created: data.packs_created,
      packs_mastered: data.packs_mastered
    }
  } catch (e) {
    console.error(e)
    error.value = e.message || 'Не удалось загрузить статистику'
  } finally {
    loaded.value = true
  }
})
</script>

<style scoped>
/* опционально, например, чтобы убрать padding слева */
div {
  margin-left: 0;
}
</style>
