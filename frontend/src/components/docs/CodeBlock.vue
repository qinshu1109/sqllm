<template>
  <div class="relative group my-4">
    <div class="absolute right-2 top-2 opacity-0 group-hover:opacity-100 transition-opacity">
      <button
        @click="copyCode"
        class="px-2 py-1 text-xs bg-slate-700 hover:bg-slate-600 text-slate-300 rounded transition-colors"
      >
        {{ copied ? '已复制' : '复制' }}
      </button>
    </div>
    <div v-if="title" class="px-4 py-2 bg-slate-800 border-b border-slate-700 rounded-t-lg">
      <span class="text-xs text-slate-400 font-mono">{{ title }}</span>
    </div>
    <pre
      :class="[
        'overflow-x-auto p-4 text-sm font-mono leading-relaxed',
        title ? 'rounded-b-lg' : 'rounded-lg',
        'bg-[#0d1117] text-slate-300'
      ]"
    ><code v-html="highlightedCode"></code></pre>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  code: string
  language?: string
  title?: string
}>()

const copied = ref(false)

function copyCode() {
  navigator.clipboard.writeText(props.code)
  copied.value = true
  setTimeout(() => {
    copied.value = false
  }, 2000)
}

// Simple syntax highlighting
const highlightedCode = computed(() => {
  let code = props.code
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  const lang = props.language || 'text'

  if (lang === 'bash' || lang === 'shell' || lang === 'curl') {
    // Highlight curl commands
    code = code
      .replace(/^(curl|wget|npm|pnpm|yarn)/gm, '<span class="text-purple-400">$1</span>')
      .replace(/(-[A-Za-z]+|--[a-z-]+)/g, '<span class="text-yellow-300">$1</span>')
      .replace(/(https?:\/\/[^\s"']+)/g, '<span class="text-green-400">$1</span>')
      .replace(/"([^"]+)"/g, '"<span class="text-blue-300">$1</span>"')
      .replace(/'([^']+)'/g, '\'<span class="text-orange-300">$1</span>\'')
  } else if (lang === 'python') {
    code = code
      .replace(/\b(import|from|as|def|class|return|if|else|elif|for|while|try|except|with|async|await)\b/g, '<span class="text-purple-400">$1</span>')
      .replace(/\b(True|False|None)\b/g, '<span class="text-orange-400">$1</span>')
      .replace(/(#.*)$/gm, '<span class="text-slate-500">$1</span>')
      .replace(/"([^"]+)"/g, '"<span class="text-green-400">$1</span>"')
      .replace(/'([^']+)'/g, '\'<span class="text-green-400">$1</span>\'')
      .replace(/\b(\d+)\b/g, '<span class="text-orange-300">$1</span>')
  } else if (lang === 'javascript' || lang === 'typescript' || lang === 'js' || lang === 'ts') {
    code = code
      .replace(/\b(import|export|from|const|let|var|function|return|if|else|for|while|try|catch|async|await|new|class|extends)\b/g, '<span class="text-purple-400">$1</span>')
      .replace(/\b(true|false|null|undefined)\b/g, '<span class="text-orange-400">$1</span>')
      .replace(/(\/\/.*)$/gm, '<span class="text-slate-500">$1</span>')
      .replace(/"([^"]+)"/g, '"<span class="text-green-400">$1</span>"')
      .replace(/'([^']+)'/g, '\'<span class="text-green-400">$1</span>\'')
      .replace(/`([^`]+)`/g, '`<span class="text-green-400">$1</span>`')
      .replace(/\b(\d+)\b/g, '<span class="text-orange-300">$1</span>')
  } else if (lang === 'json') {
    code = code
      .replace(/"([^"]+)":/g, '"<span class="text-blue-300">$1</span>":')
      .replace(/:\s*"([^"]+)"/g, ': "<span class="text-green-400">$1</span>"')
      .replace(/:\s*(\d+)/g, ': <span class="text-orange-300">$1</span>')
      .replace(/:\s*(true|false|null)/g, ': <span class="text-orange-400">$1</span>')
  }

  return code
})
</script>
