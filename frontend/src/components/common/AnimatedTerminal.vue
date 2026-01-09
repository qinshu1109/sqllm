<template>
  <div class="terminal-container">
    <div class="terminal-window">
      <!-- Window header -->
      <div class="terminal-header">
        <div class="terminal-buttons">
          <span class="btn-close"></span>
          <span class="btn-minimize"></span>
          <span class="btn-maximize"></span>
        </div>
        <span class="terminal-title">{{ title }}</span>
      </div>
      <!-- Terminal content -->
      <div class="terminal-body">
        <div
          v-for="(line, index) in lines"
          :key="index"
          class="code-line"
          :style="{ animationDelay: line.delay || `${0.3 + index * 0.7}s` }"
        >
          <!-- Prompt line -->
          <template v-if="line.type === 'prompt'">
            <span class="code-prompt">$</span>
            <span
              v-for="(segment, sIndex) in line.segments"
              :key="sIndex"
              :class="segment.class"
            >{{ segment.text }}</span>
          </template>

          <!-- Comment line -->
          <template v-else-if="line.type === 'comment'">
            <span class="code-comment">{{ line.text }}</span>
          </template>

          <!-- Response line -->
          <template v-else-if="line.type === 'response'">
            <span v-if="line.status" class="code-success">{{ line.status }}</span>
            <span v-if="line.text" class="code-response">{{ line.text }}</span>
          </template>

          <!-- Cursor line -->
          <template v-else-if="line.type === 'cursor'">
            <span class="code-prompt">$</span>
            <span class="cursor"></span>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface TerminalSegment {
  text: string
  class: string
}

export interface TerminalLine {
  type: 'prompt' | 'comment' | 'response' | 'cursor'
  segments?: TerminalSegment[]
  text?: string
  status?: string
  delay?: string
}

withDefaults(
  defineProps<{
    title?: string
    lines: TerminalLine[]
  }>(),
  {
    title: 'terminal'
  }
)
</script>

<style scoped>
/* Terminal Container */
.terminal-container {
  position: relative;
  display: inline-block;
}

/* Terminal Window */
.terminal-window {
  width: 420px;
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

/* Terminal Header */
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}

/* Terminal Body */
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #38bdf8;
}
.code-flag {
  color: #a78bfa;
}
.code-url {
  color: #14b8a6;
}
.code-comment {
  color: #64748b;
  font-style: italic;
}
.code-success {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #fbbf24;
}

/* Blinking Cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Dark mode adjustments */
:deep(.dark) .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(20, 184, 166, 0.2),
    0 0 40px rgba(20, 184, 166, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}
</style>
