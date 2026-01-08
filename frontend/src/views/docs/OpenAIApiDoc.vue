<template>
  <ApiDocLayout :toc-items="tocItems">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-slate-900 dark:text-white mb-2">OpenAI 兼容 API</h1>
      <p class="text-slate-600 dark:text-slate-400">
        FAC 提供了兼容 OpenAI API 接口的接口。您可以使用 OpenAI SDK 调用 Claude、Gemini 等模型。
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
              <td class="py-2 px-3"><code class="text-primary-500">/v1/chat/completions</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">POST</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">Chat Completions API</td>
            </tr>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1/responses</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">POST</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">Responses API</td>
            </tr>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1/models</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 rounded text-xs font-medium">GET</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">获取可用模型列表</td>
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
      <p class="text-slate-600 dark:text-slate-400 mb-4">在请求头中添加 API Key：</p>
      <CodeBlock
        code='Authorization: Bearer YOUR_API_KEY'
        language="text"
      />
    </section>

    <!-- Chat Completions API -->
    <section id="chat-completions" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        Chat Completions API
      </h2>

      <h3 id="python" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python</h3>
      <p class="text-slate-600 dark:text-slate-400 mb-3">
        使用官方 <code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-sm">openai</code> SDK：
      </p>
      <CodeBlock code="pip install openai" language="bash" />
      <CodeBlock :code="pythonExample" language="python" title="Python" />

      <h3 id="nodejs" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Node.js</h3>
      <CodeBlock :code="nodejsExample" language="javascript" title="Node.js" />

      <h3 id="curl" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">curl</h3>
      <CodeBlock :code="curlExample" language="bash" title="curl" />
    </section>

    <!-- 流式响应 -->
    <section id="streaming" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        流式响应
      </h2>

      <h3 id="python-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python 流式</h3>
      <CodeBlock :code="pythonStreamExample" language="python" title="Python" />

      <h3 id="nodejs-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Node.js 流式</h3>
      <CodeBlock :code="nodejsStreamExample" language="javascript" title="Node.js" />
    </section>

    <!-- Responses API -->
    <section id="responses-api" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        Responses API
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">
        Responses API 是 OpenAI 推出的新一代 API，提供更丰富的功能。
      </p>

      <h3 id="responses-basic" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">基本使用</h3>
      <CodeBlock :code="responsesApiExample" language="python" title="Python" />

      <h3 id="responses-tools" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">使用工具 (Function Calling)</h3>
      <CodeBlock :code="responsesApiToolsExample" language="python" title="Python" />
    </section>

    <!-- 模型映射 -->
    <section id="model-mapping" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        模型映射
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">
        使用 OpenAI 兼容 API 时，模型名称会自动映射：
      </p>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">请求的模型</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">实际使用</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="mapping in modelMappings" :key="mapping.from" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">{{ mapping.from }}</code></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ mapping.to }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="mt-4 p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
        <p class="text-sm text-amber-800 dark:text-amber-200">
          <strong>提示：</strong>建议直接使用 Claude 模型名称，以获得更准确的控制。
        </p>
      </div>
    </section>

    <!-- 请求参数 -->
    <section id="params" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        请求参数
      </h2>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">参数</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">类型</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">必填</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="param in requestParams" :key="param.name" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">{{ param.name }}</code></td>
              <td class="py-2 px-3 text-slate-500">{{ param.type }}</td>
              <td class="py-2 px-3">
                <span v-if="param.required" class="text-red-500">是</span>
                <span v-else class="text-slate-400">否</span>
              </td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ param.desc }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- 与现有项目集成 -->
    <section id="integrations" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        与现有项目集成
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">
        无需重写代码逻辑，只需修改配置即可切换到 FAC：
      </p>

      <h3 id="integration-python" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python</h3>
      <CodeBlock :code="integrationPythonExample" language="python" title="Python" />

      <h3 id="integration-env" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">环境变量方式</h3>
      <CodeBlock :code="integrationEnvExample" language="bash" title="环境变量" />
      <p class="text-slate-600 dark:text-slate-400 mt-3 mb-2">然后在代码中直接使用：</p>
      <CodeBlock :code="integrationEnvUsageExample" language="python" title="Python" />
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
          <h4 class="font-medium text-slate-900 dark:text-white mb-2">为什么选用 OpenAI 兼容 API?</h4>
          <ul class="text-sm text-slate-600 dark:text-slate-400 list-disc list-inside space-y-1">
            <li>无需修改代码逻辑，只需更改 API 地址</li>
            <li>可以使用熟悉的 OpenAI SDK</li>
            <li>便于在不同服务商之间切换</li>
          </ul>
        </div>
        <div class="p-4 bg-slate-50 dark:bg-slate-800 rounded-lg">
          <h4 class="font-medium text-slate-900 dark:text-white mb-2">所有 OpenAI 功能都支持吗?</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400 mb-2">大部分常用功能都支持，包括：</p>
          <ul class="text-sm text-slate-600 dark:text-slate-400 list-disc list-inside space-y-1">
            <li>Chat Completions</li>
            <li>流式响应</li>
            <li>函数调用 / 工具使用</li>
            <li>Responses API</li>
          </ul>
          <p class="text-sm text-slate-600 dark:text-slate-400 mt-2">部分高级功能可能存在差异，建议以实际测试为准。</p>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="mt-12 pt-6 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between text-sm text-slate-500 dark:text-slate-400">
      <router-link to="/docs/api/claude" class="text-primary-500 hover:text-primary-600 dark:hover:text-primary-400 flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        上一页: Claude API
      </router-link>
      <span>最后更新于: 2024/11/4 20:56</span>
      <router-link to="/docs/api/gemini" class="text-primary-500 hover:text-primary-600 dark:hover:text-primary-400 flex items-center gap-1">
        下一页: Gemini API
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </router-link>
    </footer>
  </ApiDocLayout>
</template>

<script setup lang="ts">
import ApiDocLayout from '@/components/docs/ApiDocLayout.vue'
import CodeBlock from '@/components/docs/CodeBlock.vue'

const tocItems = [
  { id: 'endpoints', title: 'API 端点', level: 2 },
  { id: 'auth', title: '认证方式', level: 2 },
  { id: 'chat-completions', title: 'Chat Completions API', level: 2 },
  { id: 'python', title: 'Python', level: 3 },
  { id: 'nodejs', title: 'Node.js', level: 3 },
  { id: 'curl', title: 'curl', level: 3 },
  { id: 'streaming', title: '流式响应', level: 2 },
  { id: 'python-stream', title: 'Python 流式', level: 3 },
  { id: 'nodejs-stream', title: 'Node.js 流式', level: 3 },
  { id: 'responses-api', title: 'Responses API', level: 2 },
  { id: 'model-mapping', title: '模型映射', level: 2 },
  { id: 'params', title: '请求参数', level: 2 },
  { id: 'integrations', title: '与现有项目集成', level: 2 },
  { id: 'response', title: '响应格式', level: 2 },
  { id: 'error-handling', title: '错误处理', level: 2 },
  { id: 'faq', title: '常见问题', level: 2 }
]

const pythonExample = `from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://new.fastaiode.top/v1"
)

response = client.chat.completions.create(
    model="claude-3-sonnet-20240229",
    messages=[
        {"role": "system", "content": "你是一个有帮助的助手。"},
        {"role": "user", "content": "什么是机器学习?"}
    ]
)

print(response.choices[0].message.content)`

const nodejsExample = `import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'YOUR_API_KEY',
  baseURL: 'https://new.fastaiode.top/v1'
});

async function main() {
  const response = await client.chat.completions.create({
    model: "claude-3-sonnet-20240229",
    messages: [
      { role: "system", content: "你是一个有帮助的助手。" },
      { role: "user", content: "什么是机器学习?" }
    ],
  });

  console.log(response.choices[0].message.content);
}

main();`

const curlExample = `curl https://new.fastaiode.top/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -d '{
    "model": "claude-3-sonnet-20240229",
    "messages": [
      {"role": "system", "content": "你是一个有帮助的助手。"},
      {"role": "user", "content": "什么是机器学习?"}
    ]
  }'`

const pythonStreamExample = `from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://new.fastaiode.top/v1"
)

stream = client.chat.completions.create(
    model="claude-3-sonnet-20240229",
    messages=[
        {"role": "user", "content": "讲一个有趣的故事"}
    ],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)`

const nodejsStreamExample = `import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'YOUR_API_KEY',
  baseURL: 'https://new.fastaiode.top/v1'
});

async function main() {
  const stream = await client.chat.completions.create({
    model: "claude-3-sonnet-20240229",
    messages: [
      { role: "user", content: "讲一个有趣的故事" }
    ],
    stream: true,
  });

  for await (const chunk of stream) {
    const content = chunk.choices[0].delta.content;
    if (content) {
      process.stdout.write(content);
    }
  }
}

main();`

const responsesApiExample = `from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://new.fastaiode.top/v1"
)

response = client.responses.create(
    model="claude-3-sonnet-20240229",
    input="解释什么是相对论"
)
print(response.output.text)`

const responsesApiToolsExample = `response = client.responses.create(
    model="claude-3-sonnet-20240229",
    input="北京现在的天气怎么样?",
    tools=[
        {
            "type": "function",
            "function": {
                "name": "get_weather",
                "description": "获取指定地区的天气信息",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "city": {
                            "type": "string",
                            "description": "城市名称"
                        }
                    },
                    "required": ["city"]
                }
            }
        }
    ]
)`

const integrationPythonExample = `# 原本的代码
# client = OpenAI()

# 修改为：
client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://new.fastaiode.top/v1"
)`

const integrationEnvExample = `export OPENAI_API_KEY=YOUR_API_KEY
export OPENAI_BASE_URL=https://new.fastaiode.top/v1`

const integrationEnvUsageExample = `from openai import OpenAI
client = OpenAI() # 会自动读取环境变量`

const responseExample = `{
  "id": "chatcmpl-333",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "claude-3-sonnet-20240229",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "机器学习是一门...."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 20,
    "completion_tokens": 150,
    "total_tokens": 170
  }
}`

const errorHandlingExample = `from openai import OpenAI, APIError, RateLimitError

try:
    # 调用 API 的代码
except RateLimitError:
    print("请求过于频繁，请稍后再试")
except APIError as e:
    print(f"API 错误: {e}")`

const modelMappings = [
  { from: 'gpt-4o', to: '可复用的默认模型' },
  { from: 'gpt-4-turbo', to: '可复用的默认模型' },
  { from: 'claude-3-sonnet-20240229', to: 'Claude Sonnet 4' },
  { from: 'claude-opus-20240229', to: 'Claude Opus 4' }
]

const requestParams = [
  { name: 'model', type: 'string', required: true, desc: '模型 ID' },
  { name: 'messages', type: 'array', required: true, desc: '对话消息数组' },
  { name: 'max_tokens', type: 'integer', required: false, desc: '最大输出 token 数' },
  { name: 'temperature', type: 'number', required: false, desc: '温度，0-2，默认 1' },
  { name: 'top_p', type: 'number', required: false, desc: 'Top-P 采样' },
  { name: 'stream', type: 'boolean', required: false, desc: '是否流式输出' },
  { name: 'stop', type: 'array', required: false, desc: '停止序列' },
  { name: 'presence_penalty', type: 'number', required: false, desc: '存在惩罚' },
  { name: 'frequency_penalty', type: 'number', required: false, desc: '频率惩罚' }
]
</script>
