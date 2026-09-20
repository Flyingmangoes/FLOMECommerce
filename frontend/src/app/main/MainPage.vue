<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Header from '@/components/header.vue'
import Footer from '@/components/footer.vue'
import ProductGrid from '@/components/product/product_grid.vue'
import BrandSponsorTransition from '@/components/transition/brand_sponsor_transition.vue'
import QuoteBrandTransition from '@/components/transition/quote_brand_transition.vue'
import type { Product } from '@/lib/models/Product'
import { gsap, ScrollTrigger } from '@/lib/gsap'
import { mockProducts } from '@/test/mock_product'

const productList: Product[] = mockProducts
function handleSelectProduct(product: Product) {
  // TODO: navigate to product detail page once vue-router + routes/ are set up
  console.log('Selected product:', product)
}

/* ------------------------------------------------------------------ *
 * Highlighted Product Info — interactive size selector + quantity
 * ------------------------------------------------------------------ */
const sizes = ['S', 'M', 'L', 'XL']
const selectedSize = ref('M')
const quantity = ref(1)

function incrementQty() {
  quantity.value++
}
function decrementQty() {
  if (quantity.value > 1) quantity.value--
}

/* ------------------------------------------------------------------ *
 * Filter bar — Category dropdown open/close (Figma showed this as a
 * static "open state (reference only)" layer; making it a real
 * toggle here since this is functioning code, not a design mockup)
 * ------------------------------------------------------------------ */
const isCategoryOpen = ref(false)
const categorySearch = ref('')
const categories = ['Outerwear', 'Tops', 'Bottoms', 'Footwear', 'Accessories']
const filteredCategories = computed(() =>
  categories.filter((c) => c.toLowerCase().includes(categorySearch.value.toLowerCase()))
)
const selectedCategory = ref('Outerwear')

function pickCategory(cat: string) {
  selectedCategory.value = cat
  isCategoryOpen.value = false
}

/* ------------------------------------------------------------------ *
 * Pagination
 * ------------------------------------------------------------------ */
const currentPage = ref(1)
const totalPages = 10

/* ------------------------------------------------------------------ *
 * Scroll-reveal animations (gsap + ScrollTrigger)
 * ------------------------------------------------------------------ */
const heroRef = ref<HTMLElement | null>(null)
const highlightedRef = ref<HTMLElement | null>(null)
const catalogHeaderRef = ref<HTMLElement | null>(null)

const triggers: ScrollTrigger[] = []

function reveal(el: Element | null, opts: Record<string, unknown> = {}) {
  if (!el) return
  const tween = gsap.from(el, {
    y: 40,
    opacity: 0,
    duration: 0.8,
    ease: 'power3.out',
    scrollTrigger: { trigger: el, start: 'top 85%' },
    ...opts,
  })
  if (tween.scrollTrigger) triggers.push(tween.scrollTrigger)
}

onMounted(() => {
  // Hero animates in on load rather than on scroll, since it's visible
  // immediately above the fold.
  if (heroRef.value) {
    gsap.from(heroRef.value.querySelectorAll('.hero-animate'), {
      y: 30,
      opacity: 0,
      duration: 0.8,
      stagger: 0.12,
      ease: 'power3.out',
      delay: 0.1,
    })
  }
  reveal(highlightedRef.value)
  reveal(catalogHeaderRef.value)
})

onUnmounted(() => {
  triggers.forEach((t) => t.kill())
})
</script>

<template>
  <div class="w-full bg-white font-sans">
    <Header />

    <!-- ================= Hero / Opening ================= -->
    <section ref="heroRef" class="relative w-full h-[768px] bg-white overflow-hidden px-10">
      <!-- Eyebrow -->
      <div class="hero-animate absolute top-10 right-10 flex items-center gap-2">
        <span class="h-2.5 w-2.5 rounded-full bg-[#FF5A1F]" />
        <span class="text-xs font-bold tracking-[0.15em] text-black">// NEW COLLECTION 2026</span>
      </div>

      <!-- Left headline -->
      <h1 class="hero-animate absolute left-10 top-20 font-extrabold leading-[0.95] text-black">
        <span class="block text-[80px]">WHERE</span>
        <span class="block text-[80px] ml-10">STYLE</span>
      </h1>

      <!-- Right headline -->
      <h1 class="hero-animate absolute right-[240px] top-[120px] font-extrabold leading-[0.95] text-right">
        <span class="block text-[80px] text-black">LIVES</span>
        <span class="block text-[80px] text-[#FF5A1F]">NOW</span>
      </h1>

      <!-- Center product image placeholder -->
      <div class="hero-animate absolute left-1/2 top-[15px] -translate-x-1/2 w-[400px] h-[700px] bg-neutral-100 rounded-2xl flex items-center justify-center">
        <span class="text-neutral-400 text-sm">Model photo</span>
      </div>

      <!-- Supporting paragraph -->
      <p class="hero-animate absolute left-10 top-[384px] w-[400px] text-base text-neutral-600 leading-relaxed">
        Discover curated streetwear drops and everyday essentials, all designed for how you actually move through the day.
      </p>

      <!-- Stat -->
      <div class="hero-animate absolute right-10 top-[590px] text-right">
        <div class="text-[40px] font-extrabold text-black leading-none">280K</div>
        <div class="text-xs tracking-[0.08em] text-neutral-500 mt-1">PEOPLE WE INSPIRE</div>
      </div>
    </section>

    <!-- ================= Highlighted Product Info ================= -->
    <section ref="highlightedRef" class="relative w-full min-h-[768px] bg-white px-10 py-6">
      <!-- Bestseller badge -->
      <div class="rounded-full bg-[#FF5A1F] w-[140px] h-7 flex items-center justify-center mb-1">
        <span class="text-white text-xs font-bold">BESTSELLER</span>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-[560px_1fr] gap-10">
        <!-- Left: product header + price/order -->
        <div>
          <!--
            NOTE: this text currently reads exactly as it does in Figma
            right now ("PRODUCT NAME" / "UNIQUENESS" / "CATEGORY") — see
            chat, this reverted from the real copy I'd set earlier.
          -->
          <h2 class="text-[45px] font-extrabold leading-tight text-black">PRODUCT NAME</h2>
          <p class="text-2xl font-medium text-black mt-2">UNIQUENESS</p>
          <p class="text-base text-neutral-600 mt-2">CATEGORY</p>

          <p class="text-base text-neutral-600 mt-6 max-w-[500px] leading-relaxed">
            A relaxed knit polo paired with tapered trousers, designed for effortless day-to-night styling.
          </p>
          <p class="text-base text-neutral-800 mt-4 leading-relaxed">
            • 100% breathable cotton-blend<br />
            • Relaxed, tailored fit<br />
            • Machine washable
          </p>

          <!-- Price & order panel -->
          <div class="mt-8 max-w-[482px]">
            <div class="flex items-center justify-between mb-4">
              <span class="text-xs font-semibold tracking-wide text-black">SIZE</span>
              <span class="text-xs font-semibold text-green-700">&#10003; In Stock</span>
            </div>
            <div class="flex gap-2 mb-6">
              <button
                v-for="size in sizes"
                :key="size"
                type="button"
                class="w-11 h-9 rounded-md text-sm font-semibold border transition-colors"
                :class="
                  selectedSize === size
                    ? 'bg-black text-white border-black'
                    : 'bg-white text-black border-black/20 hover:border-black'
                "
                @click="selectedSize = size"
              >
                {{ size }}
              </button>
            </div>

            <div class="flex items-center gap-6">
              <span class="text-2xl font-extrabold text-black">$89.00</span>
              <div class="flex items-center gap-3">
                <button
                  type="button"
                  aria-label="Decrease quantity"
                  class="w-8 h-8 rounded-full border border-black/20 flex items-center justify-center hover:border-black"
                  @click="decrementQty"
                >
                  −
                </button>
                <span class="w-6 text-center font-semibold">{{ quantity }}</span>
                <button
                  type="button"
                  aria-label="Increase quantity"
                  class="w-8 h-8 rounded-full border border-black/20 flex items-center justify-center hover:border-black"
                  @click="incrementQty"
                >
                  +
                </button>
              </div>
            </div>

            <button
              type="button"
              class="mt-5 w-full h-[58px] rounded-full bg-[#FF5A1F] text-white font-bold text-lg hover:bg-black transition-colors"
            >
              ADD TO CART
            </button>
          </div>
        </div>

        <!-- Right: product image with soft background shape + next product teaser -->
        <div class="relative flex justify-center">
          <!--
            Figma's background is a 13-piece rotated-rectangle collage
            with a soft gradient fill. Simplified to a single soft blob
            here — the 13-div version doesn't translate cleanly to a
            functional component; can revisit if pixel-parity matters.
          -->
          <div
            class="absolute w-[400px] h-[350px] top-[100px] rounded-[40%_60%_55%_45%] opacity-70"
            style="background: linear-gradient(135deg, #e5e5e5, #f8f8f8)"
          />
          <div class="relative w-[382px] h-[500px] bg-neutral-100 rounded-2xl flex items-center justify-center z-10">
            <span class="text-neutral-400 text-sm">Product photo</span>
          </div>

          <!-- Next Highlighted teaser -->
          <div class="absolute right-0 top-0 w-[260px]">
            <div class="flex items-center justify-end gap-1 text-sm font-light tracking-[0.15em] text-black mb-3">
              NEXT
              <span class="ml-1">→</span>
            </div>
            <div class="w-full h-[330px] bg-neutral-100 rounded-2xl flex items-center justify-center">
              <span class="text-neutral-400 text-xs">Next product</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <BrandSponsorTransition />

    <!-- ================= Product Catalog ================= -->
    <section class="w-full bg-[#FAF7F2] px-10 py-14">
      <div ref="catalogHeaderRef" class="flex items-end justify-between mb-10 flex-wrap gap-4">
        <div>
          <div class="flex items-center gap-2 mb-2">
            <span class="h-2 w-2 rounded-full bg-[#FF5A1F]" />
            <span class="text-xs font-bold tracking-[0.15em] text-black">// SHOP THE COLLECTION</span>
          </div>
          <h2 class="text-4xl font-extrabold leading-none tracking-tight text-black">PRODUCT CATALOG</h2>
        </div>
        <div class="flex items-center gap-6">
          <span class="text-base text-neutral-500">10 Results</span>
          <button class="flex items-center gap-2 rounded-full border border-black px-5 py-2.5 text-sm font-bold text-black">
            BEST SELLER
          </button>
        </div>
      </div>

      <!-- Filter bar -->
      <div class="relative mb-10">
        <div class="flex flex-wrap items-center gap-3 mb-4">
          <div class="relative">
            <button
              type="button"
              class="rounded-full border border-black bg-white px-5 h-11 text-sm font-semibold text-black"
              @click="isCategoryOpen = !isCategoryOpen"
            >
              Category ⌄
            </button>

            <div
              v-if="isCategoryOpen"
              class="absolute z-20 mt-2 w-[260px] rounded-xl border border-neutral-200 bg-white shadow-lg p-4"
            >
              <input
                v-model="categorySearch"
                type="text"
                placeholder="🔍 Search category..."
                class="w-full h-9 rounded-lg border border-neutral-200 px-3 text-sm mb-2 focus:outline-none focus:border-[#FF5A1F]"
              />
              <button
                v-for="cat in filteredCategories"
                :key="cat"
                type="button"
                class="w-full text-left px-3 py-2 rounded-lg text-sm hover:bg-neutral-100"
                :class="selectedCategory === cat && 'bg-neutral-100 font-semibold'"
                @click="pickCategory(cat)"
              >
                {{ cat }}
              </button>
            </div>
          </div>

          <button type="button" class="rounded-full border border-black bg-white px-5 h-11 text-sm font-semibold text-black">
            Price ⌄
          </button>
          <button type="button" class="rounded-full border border-black bg-white px-5 h-11 text-sm font-semibold text-black">
            Size ⌄
          </button>
        </div>

        <!-- Applied filter chips -->
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-sm text-neutral-500 mr-1">Applied Filters:</span>
          <span class="rounded-full bg-neutral-100 px-4 py-1.5 text-sm">{{ selectedCategory }} ×</span>
          <span class="rounded-full bg-neutral-100 px-4 py-1.5 text-sm">$40 – $99 ×</span>
          <span class="rounded-full bg-neutral-100 px-4 py-1.5 text-sm">Size {{ selectedSize }} ×</span>
        </div>
      </div>

      <ProductGrid :products=productList @select-product="handleSelectProduct" />

      <!-- Pagination -->
      <div class="flex items-center justify-between mt-16 pt-8 border-t border-black/10">
        <button
          type="button"
          class="text-lg text-neutral-500 disabled:opacity-40"
          :disabled="currentPage === 1"
          @click="currentPage = Math.max(1, currentPage - 1)"
        >
          Previous
        </button>

        <div class="flex items-center gap-4 text-lg tracking-widest">
          <button
            v-for="page in [1, 2, 3]"
            :key="page"
            type="button"
            class="w-9 h-9 rounded-full flex items-center justify-center transition-colors"
            :class="currentPage === page ? 'bg-[#FF5A1F] text-white font-semibold' : 'text-neutral-400 hover:text-black'"
            @click="currentPage = page"
          >
            {{ page }}
          </button>
          <span class="text-neutral-400">...</span>
          <button
            v-for="page in [8, 9, 10]"
            :key="page"
            type="button"
            class="w-9 h-9 rounded-full flex items-center justify-center transition-colors"
            :class="currentPage === page ? 'bg-[#FF5A1F] text-white font-semibold' : 'text-neutral-400 hover:text-black'"
            @click="currentPage = page"
          >
            {{ page }}
          </button>
        </div>

        <button
          type="button"
          class="text-lg text-black disabled:opacity-40"
          :disabled="currentPage === totalPages"
          @click="currentPage = Math.min(totalPages, currentPage + 1)"
        >
          Next
        </button>
      </div>
    </section>

    <QuoteBrandTransition />

    <Footer />
  </div>
</template>