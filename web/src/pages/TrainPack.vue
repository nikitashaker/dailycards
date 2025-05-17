<template>
  <div v-if="loaded && cards.length" class="max-w-xl mx-auto flex flex-col gap-6">
    <!-- Заголовок -->
    <h2 class="text-lg font-bold text-center">Повторение карточек</h2>

    <!-- Карточка -->
    <div class="border rounded-box p-6 shadow-md bg-base-100">
      <!-- Вопрос -->
      <p class="text-lg font-semibold mb-2">{{ currentCard.question }}</p>

      <!-- 3.4. Пометка о прошлой ошибке -->
      <p v-if="currentCard.last_wrong" class="text-error text-sm mb-2">
        * В прошлый раз здесь была ошибка
      </p>

      <!-- Сложность -->
      <p class="text-sm opacity-60 mb-4">Сложность: {{ currentCard.rating }}</p>

      <!-- Поле для ответа -->
      <label class="form-control mb-4">
        <span class="label-text">Ваш ответ</span>
        <input
          v-model="userAnswer"
          class="input input-bordered w-full"
          @keyup.enter="submitAnswer"
        />
      </label>

      <!-- Кнопка Ответить -->
      <button class="btn btn-primary w-full mt-4" @click="submitAnswer">
        Ответить
      </button>
    </div>
  </div>

  <!-- Если нет карточек -->
  <p v-else-if="loaded" class="text-center mt-10 text-error">
    В этом паке ещё нет карточек
  </p>

  <!-- Модал результата одного ответа -->
  <dialog ref="resultDialog" class="modal">
    <div class="modal-box">
      <h3
        class="text-lg font-bold mb-2"
        :class="correct ? 'text-success' : 'text-error'"
      >
        {{ correct ? 'Вы ответили правильно!' : 'Вы ответили неправильно!' }}
      </h3>
      <p v-if="!correct" class="mb-2">
        Правильный ответ: <strong>{{ currentCard.answer }}</strong>
      </p>
      <p class="mb-4">Текущий рейтинг: <strong>{{ score }}</strong></p>
      <button class="btn btn-primary w-full" @click="nextCard">
        {{ isLast ? 'Завершить' : 'Следующий вопрос' }}
      </button>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>

  <!-- 2) Финальный модал -->
  <dialog ref="finalDialog" class="modal">
    <div class="modal-box">
      <!-- Если не было ни одной ошибки -->
      <h3 class="text-xl font-bold mb-4 text-success" v-if="incorrectCount === 0">
        🎉 Вы изучили все карточки без ошибок!
      </h3>
      <div v-else>
        <h3 class="text-xl font-bold mb-4">Статистика повторения</h3>
        <p class="mb-2">Правильных ответов: <strong>{{ correctCount }}</strong></p>
        <p class="mb-2">Неправильных ответов: <strong>{{ incorrectCount }}</strong></p>
      </div>
      <p class="mb-4">Итоговый рейтинг: <strong>{{ score }}</strong></p>
      <button class="btn btn-primary w-full" @click="closeAndHome">Закрыть</button>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route  = useRoute()
const router = useRouter()
const packId = route.params.id

const cards  = ref([])
const loaded = ref(false)

const index          = ref(0)
const userAnswer     = ref('')
const correct        = ref(false)
const score          = ref(0)
const correctCount   = ref(0)
const incorrectCount = ref(0)

const sessionStats = []

const resultDialog = ref(null)
const finalDialog  = ref(null)

const currentCard = computed(() => cards.value[index.value] || {})
const isLast      = computed(() => index.value >= cards.value.length - 1)

async function loadCards() {
  try {
    const res  = await fetch(`/api/packs/${packId}/repeat`, { credentials: 'include' })
    const data = await res.json()
    cards.value = data.map(c => ({
    id:         c.id,
    question:   c.question,
    answer:     c.answer,
    rating:     c.rating ?? 0,
    last_wrong: c.last_wrong ?? false
  }))
  } catch (e) {
    console.error('Ошибка загрузки карточек:', e)
  } finally {
    loaded.value = true
  }
}

function submitAnswer() {
  if (!userAnswer.value.trim()) return

  const isCorrect = userAnswer.value.trim().toLowerCase() === currentCard.value.answer.trim().toLowerCase()
  correct.value = isCorrect

  sessionStats.push({
    card_id: currentCard.value.id,
    correct: isCorrect
  })

  if (isCorrect) {
    correctCount.value++
    score.value++
  } else {
    incorrectCount.value++
    score.value--
  }

  userAnswer.value = ''
  resultDialog.value.showModal()
}

async function nextCard() {
  resultDialog.value.close()

  if (isLast.value) {
    await fetch(`/api/packs/${packId}/finish`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type':'application/json' },
      body: JSON.stringify({ stats: sessionStats })
    })
    finalDialog.value.showModal()
  } else {
    index.value++
  }
}

function closeAndHome() {
  finalDialog.value.close()
  router.push('/')
}

onMounted(loadCards)
</script>
