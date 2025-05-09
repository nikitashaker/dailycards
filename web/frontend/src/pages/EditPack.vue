<template>
    <div>
      <Navbar />
  
      <!-- список карточек -->
      <List :items="cards" header="Создание карточек" @edit="handleEditCard">
        <template #header>
          <div class="flex items-center gap-2">
            <h2 class="text-lg font-bold">Создание карточек</h2>
            <button
              class="btn btn-square btn-ghost"
              @click="openCreateCard"
              aria-label="Создать карточку"
            >
              <Plus class="size-[1.2em]" />
            </button>
          </div>
        </template>
  
        <template #title="{ item }">
          <div class="font-semibold">{{ item.title }}</div>
        </template>
  
        <template #subtitle="{ item }">
          <span class="text-xs uppercase font-semibold opacity-60">
            {{ item.subtitle }}
          </span>
        </template>
      </List>
  
      <!-- модалка создания карточки -->
      <dialog ref="cardDialog" class="modal">
        <div class="modal-box">
          <button
            class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
            @click="cardDialog.close()"
          >
            ✕
          </button>
          <h3 class="text-lg font-bold mb-4">Создание карточки</h3>
          <form @submit.prevent="createCard" class="flex flex-col gap-3">
            <label class="form-control">
              <span class="label-text mb-1">Вопрос</span>
              <input v-model="newQ" class="input input-bordered w-full" />
            </label>
  
            <label class="form-control">
              <span class="label-text mb-1">Ответ</span>
              <input v-model="newA" class="input input-bordered w-full" />
            </label>
  
            <label class="form-control">
              <span class="label-text mb-1">Сложность (1–5)</span>
              <input
                type="number"
                min="1"
                max="5"
                v-model="newR"
                class="input input-bordered w-full"
              />
            </label>
  
            <p v-if="cardErr" class="text-error text-sm">{{ cardErr }}</p>
  
            <button type="submit" class="btn btn-primary w-full mt-2">Создать</button>
          </form>
        </div>
        <form method="dialog" class="modal-backdrop"></form>
      </dialog>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import Navbar from '@/components/Navbar.vue'
  import List from '@/components/List.vue'
  import Plus from '@/components/icons/Plus.vue'
  
  const route = useRoute()
  const packId = route.params.id
  
  const cards = ref([])
  const loggedInUser = ref(null)
  
  const cardDialog = ref(null)
  const newQ = ref('')
  const newA = ref('')
  const newR = ref(1)
  const cardErr = ref('')
  
  function openCreateCard() {
    cardErr.value = ''
    newQ.value = ''
    newA.value = ''
    newR.value = 1
    cardDialog.value.showModal()
  }
  
  async function loadCards() {
    try {
      const res = await fetch(`/api/packs/${packId}/cards`, {
        credentials: 'include'
      })
      if (!res.ok) throw new Error(`Ошибка загрузки карточек: ${res.status}`)
      const data = await res.json()
      cards.value = data.map(c => ({
        id:       c.id,
        title:    c.question, // <-- вот здесь
        subtitle: `Сложность: ${c.rating ?? 0}`, // <-- rating.Int32
        noPlay:   true
        }))
    } catch (err) {
      console.error(err)
    }
  }
  
  async function createCard() {
    cardErr.value = ''
    if (!newQ.value || !newA.value || !newR.value) {
      cardErr.value = 'Заполните все поля'
      return
    }
    try {
      const res = await fetch(`/api/packs/${packId}/cards`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          question: newQ.value,
          answer: newA.value,
          rating: Number(newR.value)
        })
      })
      if (!res.ok) {
        const { error } = await res.json().catch(() => ({}))
        cardErr.value = error || `Ошибка создания: ${res.status}`
        return
      }
      await loadCards()
      cardDialog.value.close()
    } catch (err) {
      console.error(err)
      cardErr.value = 'Сервер недоступен'
    }
  }
  
  onMounted(async () => {
    try {
      const me = await fetch('/api/me', { credentials: 'include' })
      if (me.ok) {
        const { username } = await me.json()
        loggedInUser.value = username
      }
    } catch (err) {
      console.error('Ошибка проверки сессии:', err)
    }
    await loadCards()
  })
  
  function handleEditCard(card) {
    console.log('Редактировать карточку', card)
  }
  </script>
  