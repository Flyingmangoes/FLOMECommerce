<script setup lang="ts">
import { computed } from 'vue'
import { Shirt, Layers, Footprints, ShoppingBag } from 'lucide-vue-next'
import type { Product } from '@/lib/models/Product'

const props = defineProps<{
  product: Product
}>()

const emit = defineEmits<{
  select: [product: Product]
}>()

// Placeholder icon per category, used until productImage is populated
// with real photography.
const categoryIcon = computed(() => {
  const key = (props.product.productCategory ?? '').toLowerCase()
  if (key.includes('foot')) return Footprints
  if (key.includes('access')) return ShoppingBag
  if (key.includes('top') || key.includes('outer')) return Shirt
  return Layers
})

const isOutOfStock = computed(() => props.product.productTotalStock <= 0)

const availableColorCount = computed(() => {
  const inStock = props.product.productExtra?.filter((e) => e.stock > 0) ?? []
  return new Set(inStock.map((e) => e.color)).size
})

const availableSizeCount = computed(() => {
  const inStock = props.product.productExtra?.filter((e) => e.stock > 0) ?? []
  return new Set(inStock.map((e) => e.size)).size
})

function handleClick() {
  if (isOutOfStock.value) return
  emit('select', props.product)
}
</script>

<template>
  <article
    class="w-full max-w-[350px] group"
    :class="isOutOfStock ? 'opacity-60 cursor-not-allowed' : 'cursor-pointer'"
    role="button"
    :tabindex="isOutOfStock ? -1 : 0"
    @click="handleClick"
    @keydown.enter="handleClick"
  >
    <!-- Picture -->
    <div
      class="relative w-full aspect-[350/340] rounded-2xl bg-neutral-100 flex items-center justify-center overflow-hidden"
    >
      <img
        v-if="product.productImage"
        :src="product.productImage"
        :alt="product.productName"
        class="w-full h-full object-cover transition-transform duration-300"
        :class="!isOutOfStock && 'group-hover:scale-[1.03]'"
      />
      <component
        :is="categoryIcon"
        v-else
        :size="40"
        :stroke-width="1.2"
        class="text-neutral-400"
      />

      <span
        v-if="isOutOfStock"
        class="absolute top-3 left-3 rounded-full bg-black/80 px-3 py-1 text-[11px] font-semibold text-white"
      >
        Out of Stock
      </span>
    </div>

    <!-- Description -->
    <div class="w-full bg-white px-3.5 pt-1.5 pb-5">
      <div class="flex items-start justify-between">
        <span class="text-[11px] font-semibold tracking-wide text-neutral-500 uppercase">
          {{ product.productCategory }}
        </span>
        <span class="text-lg font-extrabold text-black">${{ product.productPrice.toFixed(0) }}</span>
      </div>
      <h3 class="text-base font-bold leading-snug text-black mt-1 mb-1">
        {{ product.productName }}
      </h3>
      <p
        v-if="availableColorCount > 0 || availableSizeCount > 0"
        class="text-xs text-neutral-500"
      >
        <template v-if="availableColorCount > 0">
          {{ availableColorCount }} {{ availableColorCount === 1 ? 'Color' : 'Colors' }}
        </template>
        <template v-if="availableColorCount > 0 && availableSizeCount > 0"> &middot; </template>
        <template v-if="availableSizeCount > 0">
          {{ availableSizeCount }} {{ availableSizeCount === 1 ? 'Size' : 'Sizes' }}
        </template>
      </p>
    </div>
  </article>
</template>