<!-- File: src/App.vue -->
<template>
  <Navbar>
    <!-- Левый слот: логотип/меню -->
    <template #start>
      <details class="dropdown">
        <summary class="text-xl btn btn-ghost font-bold">DayilyCards</summary>
        <ul class="menu dropdown-content bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm">
          <li><a>Учетная запись</a></li>
          <li><a>История</a></li>
        </ul>
      </details>
    </template>

    <!-- Правый слот: кнопки или имя пользователя -->
    <template #end>
      <template v-if="!loggedInUser">
        <button class="btn btn-neutral" @click="openLoginModal">Login</button>
        <button class="btn btn-neutral" @click="openRegisterModal">Register</button>
      </template>
      <template v-else>
        <span class="font-semibold text-primary text-lg tracking-wide">
          {{ loggedInUser }}
        </span>
      </template>
    </template>
  </Navbar>

  <!-- Разделитель -->
  <div class="flex w-full flex-col">
    <div class="pt-15 divider"></div>
  </div>

  <!-- Здесь будет отображаться либо Home.vue, либо EditPack.vue -->
  <router-view />

  <!-- Модалки регистрации и логина (живут в App, чтобы быть доступны на любом роуте) -->
  <dialog ref="regDialog" class="modal">
    <div class="modal-box">
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" @click="regDialog.close()">✕</button>
      <h3 class="text-lg font-bold mb-4">Регистрация</h3>
      <form @submit.prevent="handleRegister" class="flex flex-col gap-3">
        <label class="form-control">
          <span class="label-text mb-1">Имя пользователя</span>
          <input v-model="username" class="input input-bordered w-full" autocomplete="username" />
        </label>
        <label class="form-control">
          <span class="label-text mb-1">Пароль</span>
          <input v-model="password" type="password" class="input input-bordered w-full" autocomplete="new-password" />
        </label>
        <p v-if="regError" class="text-error text-sm">{{ regError }}</p>
        <button type="submit" class="btn btn-primary w-full mt-2">Зарегистрироваться</button>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>

  <dialog ref="loginDialog" class="modal">
    <div class="modal-box">
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" @click="loginDialog.close()">✕</button>
      <h3 class="text-lg font-bold mb-4">Вход</h3>
      <form @submit.prevent="handleLogin" class="flex flex-col gap-3">
        <label class="form-control">
          <span class="label-text mb-1">Имя пользователя</span>
          <input v-model="loginUsername" class="input input-bordered w-full" autocomplete="username" />
        </label>
        <label class="form-control">
          <span class="label-text mb-1">Пароль</span>
          <input v-model="loginPassword" type="password" class="input input-bordered w-full" autocomplete="current-password" />
        </label>
        <p v-if="loginError" class="text-error text-sm">{{ loginError }}</p>
        <button type="submit" class="btn btn-primary w-full mt-2">Войти</button>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/Navbar.vue'

const loggedInUser = ref(null)

// регистрация
const regDialog = ref(null)
const username  = ref('')
const password  = ref('')
const regError  = ref('')

function openRegisterModal() {
  regError.value = ''
  username.value = ''
  password.value = ''
  regDialog.value.showModal()
}

async function handleRegister() {
  regError.value = ''
  if (!username.value || !password.value) {
    regError.value = 'Заполните оба поля'
    return
  }
  try {
    const res = await fetch('/api/users', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      credentials: 'include',
      body: JSON.stringify({ username: username.value, password: password.value })
    })
    if (!res.ok) {
      const { error } = await res.json().catch(() => ({}))
      regError.value = error || 'Ошибка регистрации'
      return
    }
    regDialog.value.close()
    // после успешного флоу вы можете открывать логин-модалку
    loginDialog.value.showModal()
  } catch {
    regError.value = 'Сервер недоступен'
  }
}

// логин
const loginDialog   = ref(null)
const loginUsername = ref('')
const loginPassword = ref('')
const loginError    = ref('')

function openLoginModal() {
  loginError.value    = ''
  loginUsername.value = ''
  loginPassword.value = ''
  loginDialog.value.showModal()
}

async function handleLogin() {
  loginError.value = ''
  if (!loginUsername.value || !loginPassword.value) {
    loginError.value = 'Заполните оба поля'
    return
  }
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      credentials: 'include',
      body: JSON.stringify({ username: loginUsername.value, password: loginPassword.value })
    })
    if (!res.ok) throw res
    loggedInUser.value = loginUsername.value
    loginDialog.value.close()
  } catch (err) {
    try {
      const { error } = await err.json()
      loginError.value = error || 'Ошибка входа'
    } catch {
      loginError.value = 'Неверные учётные данные'
    }
  }
}

// при монтировании проверяем сессию
onMounted(async () => {
  try {
    const res = await fetch('/api/me', { credentials: 'include' })
    if (res.ok) {
      const { username } = await res.json()
      loggedInUser.value = username
    }
  // eslint-disable-next-line no-empty
  } catch {}
})
</script>
