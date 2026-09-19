<script setup lang="ts">
import { ref } from 'vue'
import { ArrowRight, Twitter, Instagram, Linkedin, Youtube } from 'lucide-vue-next'

const email = ref('')

function onSubscribe() {
  // TODO: wire up to your newsletter/subscribe API endpoint
  console.log('Subscribe:', email.value)
  email.value = ''
}

interface FooterColumn {
  title: string
  links: string[]
}

const columns: FooterColumn[] = [
  { title: 'Product', links: ['Background', 'Improvement', 'Features', 'Category', 'Prices'] },
  { title: 'Company', links: ['About Us', 'Careers', 'Press Kits', 'Partners', 'Contact'] },
  { title: 'Resources', links: ['Blogs', 'Help Center', 'Community'] },
  { title: 'Legal', links: ['Privacy Policy', 'Terms of Service', 'Compliance', 'Sitemap'] },
]

const socials = [
  { name: 'Twitter', href: '#', icon: Twitter },
  { name: 'Instagram', href: '#', icon: Instagram },
  { name: 'LinkedIn', href: '#', icon: Linkedin },
  { name: 'YouTube', href: '#', icon: Youtube },
]

const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="w-full bg-[#2B2620] text-white px-8 md:px-14 pt-14">
    <div class="grid grid-cols-1 lg:grid-cols-[320px_1fr_60px] gap-12 lg:gap-16 pb-14 border-b border-white/10">
      <div>
        <div class="font-sans font-extrabold text-2xl mb-6">FLOM</div>
        <p class="text-white/70 text-sm mb-8 max-w-[280px]">
          Imaginative minds for imaginative brands.
        </p>
        <div>
          <p class="text-xs uppercase tracking-wide text-white/60 mb-3">
            Subscribe to our newsletter
          </p>
          <form class="flex" @submit.prevent="onSubscribe">
            <input
              v-model="email"
              type="email"
              required
              placeholder="Your email address"
              class="flex-1 min-w-0 bg-white/5 border border-white/15 rounded-l-full px-5 h-12 text-sm text-white placeholder-white/40 focus:outline-none focus:border-[#FF591F] transition-colors"
            />
            <button
              type="submit"
              aria-label="Subscribe"
              class="w-12 h-12 shrink-0 rounded-r-full bg-[#FF591F] flex items-center justify-center hover:bg-white transition-colors group"
            >
              <ArrowRight :size="16" class="text-white group-hover:text-[#FF591F] transition-colors" />
            </button>
          </form>
        </div>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-4 gap-8">
        <div v-for="col in columns" :key="col.title">
          <p class="text-sm font-semibold mb-5">{{ col.title }}</p>
          <ul class="flex flex-col gap-3">
            <li v-for="link in col.links" :key="link">
              <a href="#" class="text-sm text-white/60 hover:text-white transition-colors">
                {{ link }}
              </a>
            </li>
          </ul>
        </div>
      </div>

      <div class="flex lg:flex-col gap-4 lg:gap-6 items-start lg:items-center">
        <a
          v-for="s in socials"
          :key="s.name"
          :href="s.href"
          :aria-label="s.name"
          class="w-10 h-10 rounded-full border border-white/15 flex items-center justify-center hover:border-[#FF591F] hover:text-[#FF591F] transition-colors"
        >
          <component :is="s.icon" :size="18" />
        </a>
      </div>
    </div>

    <div class="flex flex-col sm:flex-row gap-4 items-center justify-between py-8 text-xs text-white/50">
      <p>&copy; {{ currentYear }} FLOM. All rights reserved.</p>
      <div class="flex gap-6">
        <a href="/privacy" class="hover:text-white transition-colors">Privacy</a>
        <a href="/terms" class="hover:text-white transition-colors">Terms</a>
        <a href="/cookies" class="hover:text-white transition-colors">Cookies</a>
      </div>
    </div>
  </footer>
</template>