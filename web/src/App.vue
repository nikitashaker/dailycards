<template>
  <Navbar>
    <!-- Левый слот: логотип/меню -->
    <template #start>
      <RouterLink to="/" class="text-xl btn btn-ghost font-bold normal-case">
        DailyCards
      </RouterLink>
    </template>

    <!-- Правый слот: кнопки или имя пользователя -->
    <template #end>
      <template v-if="!loggedInUser">
        <button class="btn btn-neutral" @click="openLoginModal">Login</button>
        <button class="btn btn-neutral" @click="openRegisterModal">Register</button>
      </template>
      <template v-else>
        <button class="flex items-center gap-2 btn btn-ghost" @click="openUserModal">
          <Icon class="w-6 h-6 text-white" />
          <span class="font-semibold text-white text-lg tracking-wide">
            {{ loggedInUser }}
          </span>
        </button>
      </template>
    </template>
  </Navbar>

  <!-- Разделитель -->
  <div class="flex w-full flex-col">
    <div class="pt-15 divider"></div>
  </div>

  <!-- Основной контент страниц -->
  <router-view />

  <!-- Модал регистрации -->
  <dialog ref="regDialog" class="modal">
    <div class="modal-box">
      <button
        class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
        @click="regDialog.close()"
      >✕</button>
      <h3 class="text-lg font-bold mb-4">Регистрация</h3>
      <form @submit.prevent="handleRegister" class="flex flex-col gap-3">
        <label class="form-control">
          <span class="label-text mb-1">Имя пользователя</span>
          <input
            v-model="username"
            class="input input-bordered w-full"
            autocomplete="username"
          />
        </label>
        <label class="form-control">
          <span class="label-text mb-1">Пароль</span>
          <input
            v-model="password"
            type="password"
            class="input input-bordered w-full"
            autocomplete="new-password"
          />
        </label>
        <p v-if="regError" class="text-error text-sm">{{ regError }}</p>
        <button type="submit" class="btn btn-primary w-full mt-2">
          Зарегистрироваться
        </button>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>

  <!-- Модал входа -->
  <dialog ref="loginDialog" class="modal">
    <div class="modal-box">
      <button
        class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
        @click="loginDialog.close()"
      >✕</button>
      <h3 class="text-lg font-bold mb-4">Вход</h3>
      <form @submit.prevent="handleLogin" class="flex flex-col gap-3">
        <label class="form-control">
          <span class="label-text mb-1">Имя пользователя</span>
          <input
            v-model="loginUsername"
            class="input input-bordered w-full"
            autocomplete="username"
          />
        </label>
        <label class="form-control">
          <span class="label-text mb-1">Пароль</span>
          <input
            v-model="loginPassword"
            type="password"
            class="input input-bordered w-full"
            autocomplete="current-password"
          />
        </label>
        <p v-if="loginError" class="text-error text-sm">{{ loginError }}</p>
        <button type="submit" class="btn btn-primary w-full mt-2">
          Войти
        </button>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>

  <!-- Модал пользователя -->
  <dialog ref="userModal" class="modal">
    <div class="modal-box">
      <button
        class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
        @click="userModal.close()"
      >✕</button>
      <h3 class="text-lg font-bold mb-4">Меню пользователя</h3>
      <div class="flex flex-col gap-2">
        <RouterLink
          to="/stats"
          class="btn btn-ghost w-full text-left"
          @click="userModal.close()"
        >
          Общая статистика
        </RouterLink>
        <button
          class="btn btn-ghost w-full text-left"
          @click="logout"
        >
          Выйти
        </button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop"></form>
  </dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import Navbar from '@/components/Navbar.vue'
import Icon from '@/components/icons/Icon.vue'

const router = useRouter()
const loggedInUser = ref(null)

// refs для модалей
const regDialog = ref(null)
const loginDialog = ref(null)
const userModal = ref(null)

// открывающие функции
function openRegisterModal() {
  regError.value = ''
  username.value = ''
  password.value = ''
  regDialog.value.showModal()
}
function openLoginModal() {
  loginError.value = ''
  loginUsername.value = ''
  loginPassword.value = ''
  loginDialog.value.showModal()
}
function openUserModal() {
  userModal.value.showModal()
}

// регистрация
const username = ref('')
const password = ref('')
const regError = ref('')

async function handleRegister() {
  regError.value = ''
  if (!username.value || !password.value) {
    regError.value = 'Заполните оба поля'
    return
  }
  try {
    const res = await fetch('/api/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username: username.value, password: password.value }),
    })
    if (!res.ok) {
      const { error } = await res.json().catch(() => ({}))
      regError.value = error || 'Ошибка регистрации'
      return
    }
    regDialog.value.close()
    loginDialog.value.showModal()
  } catch {
    regError.value = 'Сервер недоступен'
  }
}

// вход
const loginUsername = ref('')
const loginPassword = ref('')
const loginError = ref('')

async function handleLogin() {
  loginError.value = ''
  if (!loginUsername.value || !loginPassword.value) {
    loginError.value = 'Заполните оба поля'
    return
  }
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username: loginUsername.value, password: loginPassword.value }),
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

// выход
function logout() {
  fetch('/api/logout', { method: 'POST', credentials: 'include' })
    .finally(() => {
      loggedInUser.value = null
      userModal.value.close()
      router.push('/login')
    })
}

// проверка сессии
onMounted(async () => {
  try {
    const res = await fetch('/api/me', { credentials: 'include' })
    if (res.ok) {
      const { username } = await res.json()
      loggedInUser.value = username
    }
  } catch {}
})

// уведомления раз в час
onMounted(() => {
  if (!('Notification' in window)) {
    console.warn('Browser does not support notifications')
    return
  }
  Notification.requestPermission().then(permission => {
    if (permission !== 'granted') {
      console.warn('User denied notifications')
      return
    }
    setInterval(() => {
      new Notification('DailyCards напоминание', {
        body: 'Пора повторить карточки!'
      })
    }, 5 * 60 * 1000)
  })
})
</script>