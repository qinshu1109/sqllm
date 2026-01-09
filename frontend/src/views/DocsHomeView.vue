<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark text-slate-800 dark:text-slate-100 transition-colors duration-300">
    <!-- Background Effects -->
    <div class="fixed inset-0 z-0 bg-grid-pattern opacity-60 pointer-events-none"></div>
    <div class="fixed top-0 -left-4 w-96 h-96 bg-primary-500/20 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob dark:mix-blend-normal dark:bg-primary-500/10 pointer-events-none"></div>
    <div class="fixed top-0 -right-4 w-96 h-96 bg-cyan-500/20 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob animation-delay-2000 dark:mix-blend-normal dark:bg-cyan-500/10 pointer-events-none"></div>

    <!-- Navbar -->
    <nav class="fixed top-0 w-full z-50 border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-dark-900/80 backdrop-blur-md">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center gap-8">
            <router-link
              to="/home"
              class="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-primary-500 to-cyan-400"
            >
              SQAI
            </router-link>
            <div class="hidden md:flex items-center relative group">
              <svg class="absolute left-3 w-4 h-4 text-slate-400 group-focus-within:text-primary-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input
                type="text"
                :placeholder="t('docs.searchPlaceholder')"
                class="pl-10 pr-12 py-1.5 text-sm bg-slate-100 dark:bg-slate-800 border-none rounded-lg focus:ring-2 focus:ring-primary-500 w-64 transition-all text-slate-600 dark:text-slate-300 placeholder-slate-400 outline-none"
              />
              <div class="absolute right-3 px-1.5 py-0.5 text-xs text-slate-400 border border-slate-300 dark:border-slate-600 rounded bg-white dark:bg-slate-700 shadow-sm">
                Ctrl K
              </div>
            </div>
          </div>
          <div class="flex items-center gap-6">
            <div class="hidden lg:flex items-center gap-6 text-sm font-medium text-slate-600 dark:text-slate-300">
              <router-link
                to="/docs"
                class="text-primary-500 hover:text-primary-600 transition-colors"
              >
                {{ t('docs.nav.home') }}
              </router-link>
              <a
                href="#"
                class="hover:text-slate-900 dark:hover:text-white transition-colors flex items-center gap-1"
              >
                {{ t('docs.nav.quickStart') }}
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </a>
              <router-link
                to="/docs/api/claude"
                class="hover:text-slate-900 dark:hover:text-white transition-colors"
              >
                {{ t('docs.nav.apiGuide') }}
              </router-link>
              <a
                href="#"
                class="hover:text-slate-900 dark:hover:text-white transition-colors"
              >
                {{ t('docs.nav.faq') }}
              </a>
              <router-link
                to="/home"
                class="flex items-center gap-1 text-primary-500 hover:text-primary-600 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
                <span>{{ t('docs.nav.visitSite') }}</span>
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M7 17L17 7M17 7H7M17 7V17" />
                </svg>
              </router-link>
            </div>
            <div class="flex items-center gap-3 border-l border-slate-200 dark:border-slate-700 pl-6">
              <!-- Language Switcher -->
              <LocaleSwitcher />
              <!-- Theme Toggle -->
              <button
                @click="toggleTheme"
                aria-label="Toggle Dark Mode"
                class="p-2 rounded-full hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 transition-colors"
              >
                <svg v-if="isDark" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
                <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="relative pt-32 pb-20 overflow-hidden">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <!-- Hero Section -->
        <div class="grid lg:grid-cols-2 gap-12 items-center mb-24">
          <div class="text-center lg:text-left">
            <h2 class="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">
              <span class="bg-clip-text text-transparent bg-gradient-to-r from-primary-500 to-cyan-400">
                SQAI
              </span>
              <br />
              <span class="text-slate-900 dark:text-white mt-2 block">
                AI API {{ t('docs.hero.subtitle') }}
              </span>
            </h2>
            <div class="text-lg md:text-xl text-slate-600 dark:text-slate-400 mb-8 h-16 lg:h-auto">
              <p>
                {{ t('docs.hero.description') }}
                <span class="font-semibold text-slate-800 dark:text-slate-200">Claude</span>、<span class="font-semibold text-slate-800 dark:text-slate-200">OpenAI</span>、<span class="font-semibold text-slate-800 dark:text-slate-200">Gemini</span>
                {{ t('docs.hero.descriptionSuffix') }}
                <span class="inline-block w-0.5 h-5 bg-primary-500 align-middle ml-1 animate-cursor"></span>
              </p>
            </div>
            <div class="flex flex-wrap justify-center lg:justify-start gap-4 mb-10">
              <router-link
                :to="isAuthenticated ? '/dashboard' : '/login'"
                class="px-6 py-3 bg-primary-500 hover:bg-primary-600 text-white rounded-xl font-medium shadow-lg shadow-primary-500/25 transition-all hover:scale-105 active:scale-95 flex items-center gap-2"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.631 8.41m5.96 5.96a14.926 14.926 0 01-5.841 2.58m-.119-8.54a6 6 0 00-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 00-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 01-2.448-2.448 14.9 14.9 0 01.06-.312m-2.24 2.39a4.493 4.493 0 00-1.757 4.306 4.493 4.493 0 004.306-1.758M16.5 9a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z" />
                </svg>
                {{ t('docs.hero.quickStart') }}
              </router-link>
              <router-link
                to="/docs/api/claude"
                class="px-6 py-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-200 rounded-xl font-medium hover:bg-slate-50 dark:hover:bg-slate-700 transition-all hover:scale-105 active:scale-95 flex items-center gap-2"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
                {{ t('docs.hero.apiDocs') }}
              </router-link>
              <router-link
                to="/home"
                class="px-6 py-3 bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 rounded-xl font-medium hover:bg-slate-200 dark:hover:bg-slate-700 transition-all hover:scale-105 active:scale-95 flex items-center gap-2"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
                {{ t('docs.hero.visitSite') }}
              </router-link>
            </div>
            <!-- Supported Models -->
            <div class="flex flex-wrap justify-center lg:justify-start items-center gap-6 opacity-60 grayscale hover:grayscale-0 transition-all duration-500">
              <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider mr-2">
                Supported Models:
              </span>
              <div class="flex items-center gap-4">
                <!-- OpenAI Logo -->
                <svg class="h-6 w-auto dark:fill-white fill-slate-800" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M22.2819 9.8211a5.9847 5.9847 0 0 0-.5157-4.9108 6.0462 6.0462 0 0 0-6.5098-2.9A6.0651 6.0651 0 0 0 4.9807 4.1818a5.9847 5.9847 0 0 0-3.9977 2.9 6.0462 6.0462 0 0 0 .7427 7.0966 5.98 5.98 0 0 0 .511 4.9107 6.051 6.051 0 0 0 6.5146 2.9001A5.9847 5.9847 0 0 0 13.2599 24a6.0557 6.0557 0 0 0 5.7718-4.2058 5.9894 5.9894 0 0 0 3.9977-2.9001 6.0557 6.0557 0 0 0-.7475-7.0729zm-9.022 12.6081a4.4755 4.4755 0 0 1-2.8764-1.0408l.1419-.0804 4.7783-2.7582a.7948.7948 0 0 0 .3927-.6813v-6.7369l2.02 1.1686a.071.071 0 0 1 .038.052v5.5826a4.504 4.504 0 0 1-4.4945 4.4944zm-9.6607-4.1254a4.4708 4.4708 0 0 1-.5346-3.0137l.142.0852 4.783 2.7582a.7712.7712 0 0 0 .7806 0l5.8428-3.3685v2.3324a.0804.0804 0 0 1-.0332.0615L9.74 19.9502a4.4992 4.4992 0 0 1-6.1408-1.6464zM2.3408 7.8956a4.485 4.485 0 0 1 2.3655-1.9728V11.6a.7664.7664 0 0 0 .3879.6765l5.8144 3.3543-2.0201 1.1685a.0757.0757 0 0 1-.071 0l-4.8303-2.7865A4.504 4.504 0 0 1 2.3408 7.872zm16.5963 3.8558L13.1038 8.364 15.1195 7.2a.0757.0757 0 0 1 .071 0l4.8303 2.7913a4.4944 4.4944 0 0 1-.6765 8.1042v-5.6772a.79.79 0 0 0-.4069-.6666zm2.0107-3.0231l-.142-.0852-4.7735-2.7818a.7759.7759 0 0 0-.7854 0L9.409 9.2297V6.8974a.0662.0662 0 0 1 .0284-.0615l4.8303-2.7866a4.4992 4.4992 0 0 1 6.6802 4.6608zM8.3065 12.863l-2.02-1.1638a.0804.0804 0 0 1-.038-.0567V6.0742a4.4992 4.4992 0 0 1 7.3757-3.4537l-.142.0805L8.704 5.459a.7948.7948 0 0 0-.3927.6813zm1.09-1.1093l3.8654-2.2329 3.8654 2.2329-3.8654 2.2329-3.8654-2.2329z" />
                </svg>
                <!-- Gemini Logo -->
                <svg class="h-6 w-auto dark:fill-white fill-slate-800" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2L9.5 9.5 2 12l7.5 2.5L12 22l2.5-7.5L22 12l-7.5-2.5z" />
                </svg>
                <!-- Claude Logo -->
                <svg class="h-6 w-auto dark:fill-white fill-slate-800" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2L2 22h20L12 2zm0 4.5l6.5 13H5.5L12 6.5z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Code Example - Animated Terminal -->
          <div class="relative hidden lg:block">
            <AnimatedTerminal title="curl example" :lines="terminalLines" />
          </div>
        </div>

        <!-- Features Section -->
        <div id="features" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="(feature, index) in features"
            :key="index"
            class="group p-8 rounded-3xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
          >
            <div
              :class="[feature.colorBg, 'w-12 h-12 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform']"
            >
              <component :is="feature.icon" class="w-6 h-6" :class="feature.iconColor" />
            </div>
            <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-3">
              {{ feature.title }}
            </h3>
            <p class="text-sm text-slate-500 dark:text-slate-400 leading-relaxed">
              {{ feature.description }}
            </p>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="border-t border-slate-200 dark:border-slate-800 bg-white dark:bg-dark-900 py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col md:flex-row justify-between items-center gap-4">
        <div class="text-slate-500 dark:text-slate-400 text-sm">
          © 2024 SQAI AI API Gateway. All rights reserved.
        </div>
        <div class="flex gap-6 text-sm font-medium text-slate-600 dark:text-slate-400">
          <a href="#" class="hover:text-primary-500 transition-colors">
            Privacy Policy
          </a>
          <a href="#" class="hover:text-primary-500 transition-colors">
            Terms of Service
          </a>
          <a href="#" class="hover:text-primary-500 transition-colors">
            Contact
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, h, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import AnimatedTerminal, { type TerminalLine } from '@/components/common/AnimatedTerminal.vue'

const { t } = useI18n()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)

// Terminal lines for animation
const terminalLines: TerminalLine[] = [
  {
    type: 'prompt',
    segments: [
      { text: 'curl', class: 'code-cmd' },
      { text: '-X POST', class: 'code-flag' },
      { text: '/v1/chat/completions', class: 'code-url' }
    ],
    delay: '0.3s'
  },
  {
    type: 'comment',
    text: '# Sending request to API...',
    delay: '1s'
  },
  {
    type: 'response',
    status: '200 OK',
    text: '{ "choices": [...] }',
    delay: '1.8s'
  },
  {
    type: 'cursor',
    delay: '2.5s'
  }
]

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
})

// Feature icons as functional components
const KeyIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z' })
])

const BotIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' })
])

const ChartIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' })
])

const BanknoteIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z' })
])

const ZapIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M13 10V3L4 14h7v7l9-11h-7z' })
])

const LockIcon = () => h('svg', { class: 'w-6 h-6', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z' })
])

const features = computed(() => [
  {
    icon: KeyIcon,
    title: t('docs.features.easyToUse.title'),
    description: t('docs.features.easyToUse.description'),
    colorBg: 'bg-blue-100 dark:bg-blue-900/30',
    iconColor: 'text-blue-600 dark:text-blue-400'
  },
  {
    icon: BotIcon,
    title: t('docs.features.multiModel.title'),
    description: t('docs.features.multiModel.description'),
    colorBg: 'bg-purple-100 dark:bg-purple-900/30',
    iconColor: 'text-purple-600 dark:text-purple-400'
  },
  {
    icon: ChartIcon,
    title: t('docs.features.transparent.title'),
    description: t('docs.features.transparent.description'),
    colorBg: 'bg-green-100 dark:bg-green-900/30',
    iconColor: 'text-green-600 dark:text-green-400'
  },
  {
    icon: BanknoteIcon,
    title: t('docs.features.flexible.title'),
    description: t('docs.features.flexible.description'),
    colorBg: 'bg-yellow-100 dark:bg-yellow-900/30',
    iconColor: 'text-yellow-600 dark:text-yellow-400'
  },
  {
    icon: ZapIcon,
    title: t('docs.features.reliable.title'),
    description: t('docs.features.reliable.description'),
    colorBg: 'bg-orange-100 dark:bg-orange-900/30',
    iconColor: 'text-orange-600 dark:text-orange-400'
  },
  {
    icon: LockIcon,
    title: t('docs.features.secure.title'),
    description: t('docs.features.secure.description'),
    colorBg: 'bg-red-100 dark:bg-red-900/30',
    iconColor: 'text-red-600 dark:text-red-400'
  }
])
</script>

<style scoped>
.bg-background-light {
  background-color: #f8fafc;
}
.bg-background-dark {
  background-color: #0B1120;
}

.bg-grid-pattern {
  background-size: 40px 40px;
  background-image: linear-gradient(
      to right,
      rgba(148, 163, 184, 0.1) 1px,
      transparent 1px
    ),
    linear-gradient(
      to bottom,
      rgba(148, 163, 184, 0.1) 1px,
      transparent 1px
    );
}

:deep(.dark) .bg-grid-pattern {
  background-image: linear-gradient(
      to right,
      rgba(148, 163, 184, 0.05) 1px,
      transparent 1px
    ),
    linear-gradient(
      to bottom,
      rgba(148, 163, 184, 0.05) 1px,
      transparent 1px
    );
}

.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #1e293b;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #475569;
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #64748b;
}

@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}

@keyframes cursor {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animate-cursor {
  animation: cursor 0.75s step-end infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.perspective-1000 {
  perspective: 1000px;
}
</style>
