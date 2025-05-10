<template>
  <ul class="list bg-base-100 rounded-box shadow-md">
    <!-- заголовок -->
    <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">
      <slot name="header">{{ header }}</slot>
    </li>

    <!-- строки списка -->
    <li
      v-for="item in items"
      :key="item.id"
      class="list-row flex items-center justify-between p-4"
    >
      <div class="flex flex-col gap-0.5">
        <slot name="title" :item="item">
          {{ item.title }}
        </slot>
        <slot name="subtitle" :item="item">
          <span
            v-if="item.subtitle"
            class="text-xs uppercase font-semibold opacity-60"
          >
            {{ item.subtitle }}
          </span>
        </slot>
      </div>

      <div class="flex gap-2">
        <!-- 1. Play (если не отключён) -->
        <button
          v-if="!item.noPlay"
          class="btn btn-square btn-ghost"
          @click="onPlay(item)"
          aria-label="Воспроизвести"
        >
          <Play class="size-[1.2em]" />
        </button>

        <!-- 2. Edit (теперь скрывается, если item.noEdit) -->
        <button
          v-if="!item.noEdit"
          class="btn btn-square btn-ghost"
          @click="onEdit(item)"
          aria-label="Редактировать"
        >
          <Pencil class="size-[1.2em]" />
        </button>

        <!-- 3. Delete всегда показываем -->
        <button
          class="btn btn-square btn-ghost"
          @click="onDelete(item)"
          aria-label="Удалить"
        >
          <Bin class="size-[1.2em]" />
        </button>
      </div>
    </li>
  </ul>
</template>

<script setup>
import Play   from './icons/Play.vue'
import Pencil from './icons/Pencil.vue'
import Bin    from './icons/Bin.vue'

const { items, header } = defineProps({
  items: {
    type: Array,
    required: true,
    default: () => []
  },
  header: {
    type: String,
    default: 'Ваши паки'
  }
})

const emit = defineEmits(['play', 'edit', 'delete'])

function onPlay(item) {
  emit('play', item)
}
function onEdit(item) {
  emit('edit', item)
}
function onDelete(item) {
  emit('delete', item)
}
</script>
