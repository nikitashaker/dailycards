import { createRouter, createWebHistory } from 'vue-router'
import Home     from '@/pages/Home.vue'
import EditPack from '@/pages/EditPack.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',              component: Home      },
    { path: '/editpack/:id',  component: EditPack, props: true },
    { path: '/train/:id',     name: 'Train',    component: () => import('@/pages/TrainPack.vue') }
  ],
})
