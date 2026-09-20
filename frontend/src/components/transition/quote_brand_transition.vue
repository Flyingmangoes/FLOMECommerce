<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { gsap } from '@/lib/gsap'

// Same SF Compact -> Inter Extrabold substitution as BrandSponsorTransition.
const items = ['FLOMM', 'MOTTO_1', 'FLOMM', 'MOTTO_1']

const trackRef = ref<HTMLElement | null>(null)
let tween: ReturnType<typeof gsap.to> | undefined

onMounted(() => {
  if (!trackRef.value) return
  const loopWidth = trackRef.value.scrollWidth / 2
  tween = gsap.to(trackRef.value, {
    x: -loopWidth,
    duration: 18,
    ease: 'none',
    repeat: -1,
  })
})

onUnmounted(() => {
  tween?.kill()
})
</script>

<template>
  <div class="w-full h-[100px] bg-[#EAEAEA] overflow-hidden flex items-center">
    <div ref="trackRef" class="flex items-center whitespace-nowrap will-change-transform">
      <span
        v-for="(item, i) in [...items, ...items]"
        :key="i"
        class="font-sans font-extrabold text-[40px] text-[#1A1815] px-[65px]"
      >
        {{ item }}
      </span>
    </div>
  </div>
</template>
