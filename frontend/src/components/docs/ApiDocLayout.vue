<template>
  <div class="min-h-screen bg-white dark:bg-dark-900 transition-colors">
    <!-- Top Navbar -->
    <nav class="fixed top-0 left-0 right-0 z-50 h-14 border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-dark-900/80 backdrop-blur-md">
      <div class="h-full px-4 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <router-link
            to="/docs"
            class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-primary-500 to-cyan-400"
          >
            SQAI
          </router-link>
          <span class="text-sm text-slate-500 dark:text-slate-400">API 文档</span>
        </div>
        <div class="flex items-center gap-4">
          <!-- Search -->
          <div class="hidden md:flex items-center relative">
            <svg class="absolute left-3 w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              type="text"
              placeholder="搜索文档..."
              class="pl-10 pr-4 py-1.5 text-sm bg-slate-100 dark:bg-slate-800 border-none rounded-lg focus:ring-2 focus:ring-primary-500 w-48 text-slate-600 dark:text-slate-300 placeholder-slate-400 outline-none"
            />
          </div>
          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 transition-colors"
          >
            <svg v-if="isDark" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </button>
          <!-- Back to Home -->
          <router-link
            to="/docs"
            class="text-sm text-slate-600 dark:text-slate-400 hover:text-primary-500 transition-colors"
          >
            返回首页
          </router-link>
        </div>
      </div>
    </nav>

    <div class="pt-14 flex">
      <!-- Left Sidebar -->
      <aside class="fixed left-0 top-14 bottom-0 w-64 border-r border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-dark-800 overflow-y-auto hidden lg:block">
        <div class="p-4">
          <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">API 指南</h3>
          <nav class="space-y-1">
            <router-link
              v-for="item in sidebarItems"
              :key="item.path"
              :to="item.path"
              class="flex items-center gap-2 px-3 py-2 text-sm rounded-lg transition-colors"
              :class="[
                isActive(item.path)
                  ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
                  : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'
              ]"
            >
              <component :is="item.icon" class="w-4 h-4" />
              {{ item.title }}
            </router-link>
          </nav>
        </div>
      </aside>

      <!-- Main Content -->
      <main class="flex-1 lg:ml-64 lg:mr-56">
        <div class="max-w-3xl mx-auto px-6 py-8">
          <slot />
        </div>
      </main>

      <!-- Right TOC -->
      <aside class="fixed right-0 top-14 bottom-0 w-56 border-l border-slate-200 dark:border-slate-700 bg-white dark:bg-dark-900 overflow-y-auto hidden xl:block">
        <div class="p-4">
          <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">本页目录</h3>
          <nav class="space-y-1">
            <a
              v-for="item in tocItems"
              :key="item.id"
              :href="`#${item.id}`"
              class="block px-3 py-1.5 text-sm text-slate-600 dark:text-slate-400 hover:text-primary-500 dark:hover:text-primary-400 transition-colors"
              :class="{ 'pl-6': item.level === 3 }"
            >
              {{ item.title }}
            </a>
          </nav>
        </div>
      </aside>
    </div>

    <!-- Mobile Sidebar Toggle -->
    <button
      @click="showMobileSidebar = !showMobileSidebar"
      class="lg:hidden fixed bottom-4 left-4 z-50 p-3 bg-primary-500 text-white rounded-full shadow-lg"
    >
      <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>

    <!-- Mobile Sidebar Overlay -->
    <div
      v-if="showMobileSidebar"
      class="lg:hidden fixed inset-0 z-40 bg-black/50"
      @click="showMobileSidebar = false"
    >
      <aside class="absolute left-0 top-0 bottom-0 w-64 bg-white dark:bg-dark-800 overflow-y-auto" @click.stop>
        <div class="p-4 pt-6">
          <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">API 指南</h3>
          <nav class="space-y-1">
            <router-link
              v-for="item in sidebarItems"
              :key="item.path"
              :to="item.path"
              class="flex items-center gap-2 px-3 py-2 text-sm rounded-lg transition-colors"
              :class="[
                isActive(item.path)
                  ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
                  : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'
              ]"
              @click="showMobileSidebar = false"
            >
              <component :is="item.icon" class="w-4 h-4" />
              {{ item.title }}
            </router-link>
          </nav>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { useRoute } from 'vue-router'

defineProps<{
  tocItems?: Array<{ id: string; title: string; level: number }>
}>()

const route = useRoute()
const showMobileSidebar = ref(false)
const isDark = ref(document.documentElement.classList.contains('dark'))

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

function isActive(path: string) {
  return route.path === path
}

onMounted(() => {
  initTheme()
})

// Sidebar icons
const ClaudeIcon = () => h('svg', { class: 'w-4 h-4', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M12 2L2 22h20L12 2zm0 4.5l6.5 13H5.5L12 6.5z' })
])

const GeminiIcon = () => h('svg', { class: 'w-4 h-4', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M12 2L9.5 9.5 2 12l7.5 2.5L12 22l2.5-7.5L22 12l-7.5-2.5z' })
])

const OpenAIIcon = () => h('svg', { class: 'w-4 h-4', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M22.2819 9.8211a5.9847 5.9847 0 0 0-.5157-4.9108 6.0462 6.0462 0 0 0-6.5098-2.9A6.0651 6.0651 0 0 0 4.9807 4.1818a5.9847 5.9847 0 0 0-3.9977 2.9 6.0462 6.0462 0 0 0 .7427 7.0966 5.98 5.98 0 0 0 .511 4.9107 6.051 6.051 0 0 0 6.5146 2.9001A5.9847 5.9847 0 0 0 13.2599 24a6.0557 6.0557 0 0 0 5.7718-4.2058 5.9894 5.9894 0 0 0 3.9977-2.9001 6.0557 6.0557 0 0 0-.7475-7.0729zm-9.022 12.6081a4.4755 4.4755 0 0 1-2.8764-1.0408l.1419-.0804 4.7783-2.7582a.7948.7948 0 0 0 .3927-.6813v-6.7369l2.02 1.1686a.071.071 0 0 1 .038.052v5.5826a4.504 4.504 0 0 1-4.4945 4.4944zm-9.6607-4.1254a4.4708 4.4708 0 0 1-.5346-3.0137l.142.0852 4.783 2.7582a.7712.7712 0 0 0 .7806 0l5.8428-3.3685v2.3324a.0804.0804 0 0 1-.0332.0615L9.74 19.9502a4.4992 4.4992 0 0 1-6.1408-1.6464zM2.3408 7.8956a4.485 4.485 0 0 1 2.3655-1.9728V11.6a.7664.7664 0 0 0 .3879.6765l5.8144 3.3543-2.0201 1.1685a.0757.0757 0 0 1-.071 0l-4.8303-2.7865A4.504 4.504 0 0 1 2.3408 7.872zm16.5963 3.8558L13.1038 8.364 15.1195 7.2a.0757.0757 0 0 1 .071 0l4.8303 2.7913a4.4944 4.4944 0 0 1-.6765 8.1042v-5.6772a.79.79 0 0 0-.4069-.6666zm2.0107-3.0231l-.142-.0852-4.7735-2.7818a.7759.7759 0 0 0-.7854 0L9.409 9.2297V6.8974a.0662.0662 0 0 1 .0284-.0615l4.8303-2.7866a4.4992 4.4992 0 0 1 6.6802 4.6608zM8.3065 12.863l-2.02-1.1638a.0804.0804 0 0 1-.038-.0567V6.0742a4.4992 4.4992 0 0 1 7.3757-3.4537l-.142.0805L8.704 5.459a.7948.7948 0 0 0-.3927.6813zm1.09-1.1093l3.8654-2.2329 3.8654 2.2329-3.8654 2.2329-3.8654-2.2329z' })
])

const sidebarItems = [
  { path: '/docs/api/claude', title: 'Claude API', icon: ClaudeIcon },
  { path: '/docs/api/gemini', title: 'Gemini API', icon: GeminiIcon },
  { path: '/docs/api/openai', title: 'OpenAI 兼容 API', icon: OpenAIIcon }
]
</script>

<style scoped>
/* Custom scrollbar */
aside::-webkit-scrollbar {
  width: 6px;
}
aside::-webkit-scrollbar-track {
  background: transparent;
}
aside::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}
.dark aside::-webkit-scrollbar-thumb {
  background: #475569;
}
</style>
