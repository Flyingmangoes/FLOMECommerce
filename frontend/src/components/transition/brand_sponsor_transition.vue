<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { gsap } from '@/lib/gsap'

// Figma specifies font-family "SF Compact" at weight 790 for this text.
// SF Compact is an Apple system font with no public web-font license, so
// it can't be loaded on the web — substituting Inter Extrabold (already
// used site-wide) as the closest available match. Flagging this rather
// than silently picking something.
const items = ['FLOMM  SPONSOR_1', 'FLOMM  SPONSOR_2', 'FLOMM  SPONSOR_3', 'FLOMM  SPONSOR_4']

const trackRef = ref<HTMLElement | null>(null)
let tween: ReturnType<typeof gsap.to> | undefined

onMounted(() => {
  if (!trackRef.value) return
  const loopWidth = trackRef.value.scrollWidth / 2
  tween = gsap.to(trackRef.value, {
    x: -loopWidth,
    duration: 22,
    ease: 'none',
    repeat: -1,
  })
})

onUnmounted(() => {
  tween?.kill()
})
</script>

<template>
  <div class="w-full h-[100px] bg-[#242323] overflow-hidden flex items-center">
    <div ref="trackRef" class="flex items-center whitespace-nowrap will-change-transform">
      <span
        v-for="(item, i) in [...items, ...items]"
        :key="i"
        class="font-sans font-extrabold text-[40px] text-white px-[55px]"
      >
        {{ item }}
      </span>
    </div>
  </div>
</template>
