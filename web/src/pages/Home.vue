<template>
  <div>
    <Navbar />

    <!-- ====== СПИСОК ПАКОВ ====== -->
    <List
      :items="packs"
      header="Мои крутые паки"
      @edit="handleEdit"
      @play="handlePlay"
      @delete="handleDelete"
    >
      <template #header>
        <div class="flex items-center gap-2">
          <h2 class="text-lg font-bold">Пользовательские паки</h2>
          <button
            class="btn btn-square btn-ghost"
            @click="openCreatePackModal"
            aria-label="Создать пак"
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

    <!-- ====== МОДАЛКА СОЗДАНИЯ ПАКА ====== -->
    <dialog ref="createPackDialog" class="modal">
      <div class="modal-box">
        <button
          class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
          @click="createPackDialog.close()"
        >✕</button>
        <h3 class="text-lg font-bold mb-4">Создание пака</h3>

        <form @submit.prevent="handleCreatePack" class="flex flex-col gap-3">
          <label class="form-control">
            <span class="label-text mb-1">Имя пака</span>
            <input v-model="newPackName" class="input input-bordered w-full" />
          </label>
          <label class="form-control">
            <span class="label-text mb-1">Категория</span>
            <input v-model="newPackCategory" class="input input-bordered w-full" />
          </label>
          <p v-if="createPackError" class="text-error text-sm">{{ createPackError }}</p>
          <button type="submit" class="btn btn-primary w-full mt-2">Создать</button>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop"></form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/Navbar.vue'
import List   from '@/components/List.vue'
import Plus   from '@/components/icons/Plus.vue'

const router = useRouter()
const packs = ref([])

/* ---------- Навигация по кнопкам ---------- */
function handleEdit(item) {
  router.push(`/editpack/${ item.id }`)
}
function handlePlay(item) {
  router.push(`/train/${ item.id }`)
}

/* ---------- Загрузка паков ---------- */
async function loadPacks() {
  try {
    const res  = await fetch('/api/packs', { credentials: 'include' })
    if (!res.ok) throw new Error('Ошибка загрузки паков')
    const data = await res.json()
    packs.value = data.map(p => ({
      id:       p.ID,
      title:    p.Name,
      subtitle: p.Category
    }))
  } catch (err) {
    console.error(err)
  }
}

/* ---------- Создание пака ---------- */
const createPackDialog = ref(null)
const newPackName      = ref('')
const newPackCategory  = ref('')
const createPackError  = ref('')

function openCreatePackModal() {
  createPackError.value = ''
  newPackName.value     = ''
  newPackCategory.value = ''
  createPackDialog.value.showModal()
}

async function handleCreatePack() {
  createPackError.value = ''
  if (!newPackName.value || !newPackCategory.value) {
    createPackError.value = 'Заполните оба поля'
    return
  }
  try {
    const res = await fetch('/api/packs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        name:     newPackName.value,
        category: newPackCategory.value
      })
    })
    if (!res.ok) {
      const { error } = await res.json().catch(() => ({}))
      createPackError.value = error || 'Ошибка создания пака'
      return
    }
    await loadPacks()
    createPackDialog.value.close()
  } catch {
    createPackError.value = 'Сервер недоступен'
  }
}

/* ---------- Удаление пака ---------- */
async function handleDelete(item) {
  if (!confirm(`Удалить пак «${ item.title }»?`)) return

  try {
    const res = await fetch(`/api/packs/${ item.id }`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!res.ok) throw new Error(`Ошибка удаления: ${ res.status }`)
    await loadPacks()
  } catch (err) {
    console.error(err)
    alert(err.message || 'Не удалось удалить пак')
  }
}

onMounted(loadPacks)
</script>
