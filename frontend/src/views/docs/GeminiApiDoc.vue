<template>
  <ApiDocLayout :toc-items="tocItems">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-slate-900 dark:text-white mb-2">Gemini API</h1>
      <p class="text-slate-600 dark:text-slate-400">
        FAC 支持 Google Gemini API，您可以使用原生 SDK 或流式通过 Gemini CLI 调用。
      </p>
    </div>

    <!-- API 端点 -->
    <section id="endpoints" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        API 端点
      </h2>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">端点</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">方法</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1beta/models</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 rounded text-xs font-medium">GET</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">获取可用模型列表</td>
            </tr>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1beta/models/{model}</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 rounded text-xs font-medium">GET</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">获取模型详情</td>
            </tr>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1beta/models/{model}:generateContent</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">POST</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">生成内容</td>
            </tr>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1beta/models/{model}:streamGenerateContent</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">POST</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">流式生成内容</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- 认证方式 -->
    <section id="auth" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        认证方式
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">在请求头或 URL 参数中添加 API Key：</p>
      <ul class="list-disc list-inside text-slate-600 dark:text-slate-400 space-y-2 mb-4">
        <li><strong>请求头：</strong><code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-sm">x-goog-api-key: YOUR_API_KEY</code></li>
        <li><strong>或者 URL 参数：</strong><code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-sm">?key=YOUR_API_KEY</code></li>
      </ul>
    </section>

    <!-- 基础调用示例 -->
    <section id="basic-example" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        基础调用示例
      </h2>

      <h3 id="curl" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">curl</h3>
      <CodeBlock :code="curlExample" language="bash" title="curl" />

      <h3 id="python" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python (使用官方 SDK)</h3>
      <CodeBlock code="pip install google-generativeai" language="bash" />
      <CodeBlock :code="pythonExample" language="python" title="Python" />

      <h3 id="requests" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">使用 requests</h3>
      <CodeBlock :code="requestsExample" language="python" title="Python (requests)" />
    </section>

    <!-- 流式响应 -->
    <section id="streaming" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        流式响应
      </h2>

      <h3 id="curl-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">curl 流式</h3>
      <CodeBlock :code="curlStreamExample" language="bash" title="curl" />

      <h3 id="python-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python 流式</h3>
      <CodeBlock :code="pythonStreamExample" language="python" title="Python" />
    </section>

    <!-- 多轮对话 -->
    <section id="multi-turn" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        多轮对话
      </h2>
      <CodeBlock :code="multiTurnExample" language="python" title="Python" />
    </section>

    <!-- 多模态 (图片理解) -->
    <section id="multimodal" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        多模态 (图片理解)
      </h2>
      <CodeBlock :code="multimodalExample" language="python" title="Python" />
    </section>

    <!-- Gemini CLI 兼容 -->
    <section id="gemini-cli" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        Gemini CLI 兼容
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">FAC 支持 Gemini CLI 工具。配置方法：</p>

      <h3 id="gemini-cli-env" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">设置环境变量</h3>
      <CodeBlock :code="geminiCliConfig" language="bash" title="环境变量" />

      <h3 id="gemini-cli-usage" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">使用 Gemini CLI</h3>
      <CodeBlock :code="geminiCliUsage" language="bash" title="终端" />
    </section>

    <!-- 可用模型 -->
    <section id="models" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        可用模型
      </h2>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">模型 ID</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="model in geminiModels" :key="model.name" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">{{ model.name }}</code></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ model.desc }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p class="text-slate-600 dark:text-slate-400 mt-4">获取完整模型列表：</p>
      <CodeBlock :code="getModelsExample" language="bash" />
    </section>

    <!-- 请求参数 -->
    <section id="params" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        请求参数
      </h2>

      <h3 id="generate-content-params" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">generateContent 参数</h3>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">参数</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">类型</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="param in generateContentParams" :key="param.name" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">{{ param.name }}</code></td>
              <td class="py-2 px-3 text-slate-500">{{ param.type }}</td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ param.desc }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h3 id="generation-config" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">generationConfig 配置</h3>
      <CodeBlock :code="generationConfigExample" language="json" title="JSON" />
      <div class="overflow-x-auto mt-4">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">参数</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">类型</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="param in generationConfigParams" :key="param.name" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">{{ param.name }}</code></td>
              <td class="py-2 px-3 text-slate-500">{{ param.type }}</td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ param.desc }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- 响应格式 -->
    <section id="response" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        响应格式
      </h2>
      <CodeBlock :code="responseExample" language="json" title="Response" />
    </section>

    <!-- 错误处理 -->
    <section id="error-handling" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        错误处理
      </h2>
      <CodeBlock :code="errorHandlingExample" language="python" title="Python" />
    </section>

    <!-- 常见问题 -->
    <section id="faq" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        常见问题
      </h2>

      <div class="space-y-4">
        <div class="p-4 bg-slate-50 dark:bg-slate-800 rounded-lg">
          <h4 class="font-medium text-slate-900 dark:text-white mb-2">与官方 Gemini API 有什么区别?</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400 mb-2">FAC 完全兼容原生的 Gemini API 格式。主要区别规格：</p>
          <ul class="text-sm text-slate-600 dark:text-slate-400 list-disc list-inside space-y-1">
            <li>使用不同的 API 端点</li>
            <li>使用 FAC 的 API Key</li>
          </ul>
        </div>
        <div class="p-4 bg-slate-50 dark:bg-slate-800 rounded-lg">
          <h4 class="font-medium text-slate-900 dark:text-white mb-2">支持哪些 Gemini 功能?</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400 mb-2">大部分功能均支持，包括：</p>
          <ul class="text-sm text-slate-600 dark:text-slate-400 list-disc list-inside space-y-1">
            <li>文本生成</li>
            <li>流式输出</li>
            <li>多轮对话</li>
            <li>多模态 (图片理解)</li>
            <li>函数调用</li>
          </ul>
        </div>
        <div class="p-4 bg-slate-50 dark:bg-slate-800 rounded-lg">
          <h4 class="font-medium text-slate-900 dark:text-white mb-2">费用如何计算?</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400">
            按实际消耗的 Token 数计费，费率取决于选择的具体模型。详见账户额度页面。
          </p>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="mt-12 pt-6 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between text-sm text-slate-500 dark:text-slate-400">
      <router-link to="/docs/api/openai" class="text-primary-500 hover:text-primary-600 dark:hover:text-primary-400 flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        上一页: OpenAI 兼容 API
      </router-link>
      <span>最后更新于: 2024/11/4 20:56</span>
      <span></span>
    </footer>
  </ApiDocLayout>
</template>

<script setup lang="ts">
import ApiDocLayout from '@/components/docs/ApiDocLayout.vue'
import CodeBlock from '@/components/docs/CodeBlock.vue'

const tocItems = [
  { id: 'endpoints', title: 'API 端点', level: 2 },
  { id: 'auth', title: '认证方式', level: 2 },
  { id: 'basic-example', title: '基础调用示例', level: 2 },
  { id: 'curl', title: 'curl', level: 3 },
  { id: 'python', title: 'Python (使用官方 SDK)', level: 3 },
  { id: 'requests', title: '使用 requests', level: 3 },
  { id: 'streaming', title: '流式响应', level: 2 },
  { id: 'curl-stream', title: 'curl 流式', level: 3 },
  { id: 'python-stream', title: 'Python 流式', level: 3 },
  { id: 'multi-turn', title: '多轮对话', level: 2 },
  { id: 'multimodal', title: '多模态 (图片理解)', level: 2 },
  { id: 'gemini-cli', title: 'Gemini CLI 兼容', level: 2 },
  { id: 'models', title: '可用模型', level: 2 },
  { id: 'params', title: '请求参数', level: 2 },
  { id: 'response', title: '响应格式', level: 2 },
  { id: 'error-handling', title: '错误处理', level: 2 },
  { id: 'faq', title: '常见问题', level: 2 }
]

const curlExample = `curl "https://new.fastaiode.top/v1beta/models/gemini-2.0-flash:generateContent" \\
  -H "Content-Type: application/json" \\
  -d '{
    "contents": [
      {
        "parts": [
          {"text": "用简短的话解释量子计算"}
        ]
      }
    ]
  }'`

const pythonExample = `import google.generativeai as genai

genai.configure(
    api_key="YOUR_API_KEY",
    transport="rest",
    client_options={"api_endpoint": "new.fastaiode.top"}
)

model = genai.GenerativeModel("gemini-2.0-flash")
response = model.generate_content("用简短的话解释量子计算")

print(response.text)`

const requestsExample = `import requests

url = "https://new.fastaiode.top/v1beta/models/gemini-2.0-flash:generateContent"
headers = {
    "Content-Type": "application/json",
    "x-goog-api-key": "YOUR_API_KEY"
}

data = {
    "contents": [
        {
            "parts": [
                {"text": "用简短的话解释量子计算"}
            ]
        }
    ]
}

response = requests.post(url, headers=headers, json=data)
result = response.json()
print(result["candidates"][0]["content"]["parts"][0]["text"])`

const curlStreamExample = `curl "https://new.fastaiode.top/v1beta/models/gemini-2.0-flash:streamGenerateContent" \\
  -H "Content-Type: application/json" \\
  -d '{
    "contents": [
      {
        "parts": [
          {"text": "写一篇关于人工智能的散文"}
        ]
      }
    ]
  }'`

const pythonStreamExample = `import google.generativeai as genai

genai.configure(
    api_key="YOUR_API_KEY",
    transport="rest",
    client_options={"api_endpoint": "new.fastaiode.top"}
)

model = genai.GenerativeModel("gemini-2.0-flash")
response = model.generate_content(
    "写一篇关于人工智能的散文",
    stream=True
)

for chunk in response:
    print(chunk.text, end="", flush=True)`

const multiTurnExample = `import google.generativeai as genai

genai.configure(api_key="YOUR_API_KEY", transport="rest", client_options={"api_endpoint": "new.fastaiode.top"})
model = genai.GenerativeModel("gemini-2.0-flash")
chat = model.start_chat(history=[])

# 第一轮对话
response = chat.send_message("你好，我想了解一下太阳系")
print(response.text)

# 第二轮对话
response = chat.send_message("最大的行星是哪个?")
print(response.text)

# 第三轮对话
response = chat.send_message("它的核心成分是什么?")
print(response.text)`

const multimodalExample = `import google.generativeai as genai
import PIL.Image

genai.configure(api_key="YOUR_API_KEY", transport="rest", client_options={"api_endpoint": "new.fastaiode.top"})
model = genai.GenerativeModel("gemini-2.0-flash")

# 加载图片
image = PIL.Image.open("image.jpg")

# 发送描述请求
response = model.generate_content([
    "请描述这张图片的内容",
    image
])
print(response.text)`

const geminiCliConfig = `export GEMINI_API_KEY=YOUR_API_KEY
export GEMINI_API_BASE=https://new.fastaiode.top`

const geminiCliUsage = `# 安装 Gemini CLI
npm install -g @anthropic-ai/gemini-cli

# 测试
gemini "你好，请介绍一下你自己"`

const getModelsExample = `curl "https://new.fastaiode.top/v1beta/models?key=YOUR_API_KEY"`

const generationConfigExample = `{
  "generationConfig": {
    "temperature": 0.7,
    "topP": 0.95,
    "topK": 40,
    "maxOutputTokens": 1024,
    "stopSequences": ["END"]
  }
}`

const responseExample = `{
  "candidates": [
    {
      "content": {
        "parts": [
          { "text": "量子计算是..." }
        ],
        "role": "model"
      },
      "finishReason": "STOP",
      "index": 0
    }
  ],
  "usageMetadata": {
    "promptTokenCount": 10,
    "candidatesTokenCount": 150,
    "totalTokenCount": 160
  }
}`

const errorHandlingExample = `import google.generativeai as genai
from google.api_core import exceptions

# ... 配置代码 ...

try:
    response = model.generate_content("hello")
except exceptions.InvalidArgument as e:
    print(f"参数错误: {e}")
except exceptions.ResourceExhausted as e:
    print(f"配额超限: {e}")
except exceptions.GoogleAPIError as e:
    print(f"API 错误: {e}")`

const geminiModels = [
  { name: 'gemini-2.0-flash', desc: 'Gemini 2.0 Flash，极速响应' },
  { name: 'gemini-2.0-flash-thinking', desc: 'Gemini 2.0 Flash 思考预览版' },
  { name: 'gemini-1.5-pro', desc: 'Gemini 1.5 Pro' },
  { name: 'gemini-1.5-flash', desc: 'Gemini 1.5 Flash' }
]

const generateContentParams = [
  { name: 'contents', type: 'array', desc: '对话内容数组' },
  { name: 'generationConfig', type: 'object', desc: '生成配置' },
  { name: 'safetySettings', type: 'array', desc: '安全设置' },
  { name: 'systemInstruction', type: 'object', desc: '系统指令' }
]

const generationConfigParams = [
  { name: 'temperature', type: 'number', desc: '温度，0-2' },
  { name: 'topK', type: 'integer', desc: 'Top-K 采样' },
  { name: 'topP', type: 'number', desc: 'Top-P 采样' },
  { name: 'maxOutputTokens', type: 'integer', desc: '最大输出 Token' },
  { name: 'stopSequences', type: 'array', desc: '停止序列' }
]
</script>
