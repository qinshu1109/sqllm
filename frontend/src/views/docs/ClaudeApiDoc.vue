<template>
  <ApiDocLayout :toc-items="tocItems">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-slate-900 dark:text-white mb-2">Claude API</h1>
      <p class="text-slate-600 dark:text-slate-400">
        FAC 完全兼容 Anthropic Claude API，您可以使用官方 SDK 或直接调用 REST API。
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
              <td class="py-2 px-3"><code class="text-primary-500">/v1/messages</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">POST</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">发送消息 (支持流式)</td>
            </tr>
            <tr class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">/v1/messages/count_tokens</code></td>
              <td class="py-2 px-3"><span class="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded text-xs font-medium">POST</span></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">计算 Token 数量</td>
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
        :code="authHeaders"
        language="text"
      />
    </section>

    <!-- 基础调用示例 -->
    <section id="basic-example" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        基础调用示例
      </h2>

      <h3 id="curl" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">curl</h3>
      <CodeBlock
        :code="curlExample"
        language="bash"
        title="curl"
      />

      <h3 id="python" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python</h3>
      <p class="text-slate-600 dark:text-slate-400 mb-3">
        需要安装 <code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-sm">anthropic</code> SDK：
      </p>
      <CodeBlock code="pip install anthropic" language="bash" />
      <CodeBlock :code="pythonExample" language="python" title="Python" />

      <h3 id="nodejs" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Node.js</h3>
      <p class="text-slate-600 dark:text-slate-400 mb-3">
        需要安装 <code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-sm">@anthropic-ai/sdk</code>：
      </p>
      <CodeBlock code="npm install @anthropic-ai/sdk" language="bash" />
      <CodeBlock :code="nodejsExample" language="javascript" title="Node.js" />
    </section>

    <!-- 流式响应 -->
    <section id="streaming" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        流式响应
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">
        流式响应可以让用户更快看到输出，适合聊天场景。
      </p>

      <h3 id="python-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Python 流式</h3>
      <CodeBlock :code="pythonStreamExample" language="python" title="Python" />

      <h3 id="nodejs-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">Node.js 流式</h3>
      <CodeBlock :code="nodejsStreamExample" language="javascript" title="Node.js" />

      <h3 id="curl-stream" class="text-lg font-semibold text-slate-800 dark:text-slate-200 mt-6 mb-3">curl 流式</h3>
      <CodeBlock :code="curlStreamExample" language="bash" title="curl" />
    </section>

    <!-- 多轮对话 -->
    <section id="multi-turn" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        多轮对话
      </h2>
      <CodeBlock :code="multiTurnExample" language="python" title="Python" />
    </section>

    <!-- System Prompt -->
    <section id="system-prompt" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        System Prompt
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">
        使用 <code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-sm">system</code> 参数设置角色状态：
      </p>
      <CodeBlock :code="systemPromptExample" language="python" title="Python" />
    </section>

    <!-- 可用模型 -->
    <section id="models" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        可用模型
      </h2>
      <p class="text-slate-600 dark:text-slate-400 mb-4">通过 API 获取可用模型列表：</p>
      <CodeBlock :code="getModelsExample" language="bash" />
      <p class="text-slate-600 dark:text-slate-400 mt-4 mb-2 font-medium">常用模型：</p>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">模型 ID</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="model in claudeModels" :key="model.name" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-primary-500">{{ model.name }}</code></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ model.desc }}</td>
            </tr>
          </tbody>
        </table>
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

    <!-- 常见错误说明 -->
    <section id="error-codes" class="mb-10">
      <h2 class="text-xl font-bold text-slate-900 dark:text-white mb-4 pb-2 border-b border-slate-200 dark:border-slate-700">
        常见错误说明
      </h2>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700">
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">状态码</th>
              <th class="text-left py-2 px-3 font-medium text-slate-700 dark:text-slate-300">说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="error in errorCodes" :key="error.code" class="border-b border-slate-100 dark:border-slate-800">
              <td class="py-2 px-3"><code class="text-red-500">{{ error.code }}</code></td>
              <td class="py-2 px-3 text-slate-600 dark:text-slate-400">{{ error.desc }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Footer -->
    <footer class="mt-12 pt-6 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between text-sm text-slate-500 dark:text-slate-400">
      <span>最后更新于: 2024/11/4 20:56</span>
      <router-link to="/docs/api/openai" class="text-primary-500 hover:text-primary-600 dark:hover:text-primary-400 flex items-center gap-1">
        下一页: OpenAI 兼容 API
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
  { id: 'basic-example', title: '基础调用示例', level: 2 },
  { id: 'curl', title: 'curl', level: 3 },
  { id: 'python', title: 'Python', level: 3 },
  { id: 'nodejs', title: 'Node.js', level: 3 },
  { id: 'streaming', title: '流式响应', level: 2 },
  { id: 'python-stream', title: 'Python 流式', level: 3 },
  { id: 'nodejs-stream', title: 'Node.js 流式', level: 3 },
  { id: 'curl-stream', title: 'curl 流式', level: 3 },
  { id: 'multi-turn', title: '多轮对话', level: 2 },
  { id: 'system-prompt', title: 'System Prompt', level: 2 },
  { id: 'models', title: '可用模型', level: 2 },
  { id: 'params', title: '请求参数', level: 2 },
  { id: 'response', title: '响应格式', level: 2 },
  { id: 'error-handling', title: '错误处理', level: 2 },
  { id: 'error-codes', title: '常见错误说明', level: 2 }
]

const authHeaders = `x-api-key: YOUR_API_KEY
anthropic-version: 2023-06-01`

const curlExample = `curl https://sqllm.dpdns.org/v1/messages \\
  -H "Content-Type: application/json" \\
  -H "x-api-key: YOUR_API_KEY" \\
  -H "anthropic-version: 2023-06-01" \\
  -d '{
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "什么是人工智能?"}
    ]
  }'`

const pythonExample = `import anthropic

client = anthropic.Anthropic(
    api_key="YOUR_API_KEY",
    base_url="https://sqllm.dpdns.org"
)

message = client.messages.create(
    model="claude-3-sonnet-20240229",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "什么是人工智能?"}
    ]
)

print(message.content[0].text)`

const nodejsExample = `import Anthropic from '@anthropic-ai/sdk';

const client = new Anthropic({
  apiKey: 'YOUR_API_KEY',
  baseURL: 'https://sqllm.dpdns.org'
});

async function main() {
  const message = await client.messages.create({
    model: "claude-3-sonnet-20240229",
    max_tokens: 1024,
    messages: [
      { role: "user", content: "什么是人工智能?" }
    ],
  });

  console.log(message.content[0].text);
}

main();`

const pythonStreamExample = `import anthropic

client = anthropic.Anthropic(
    api_key="YOUR_API_KEY",
    base_url="https://sqllm.dpdns.org"
)

with client.messages.stream(
    model="claude-3-sonnet-20240229",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "写一篇关于季节的诗"}
    ]
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)`

const nodejsStreamExample = `import Anthropic from '@anthropic-ai/sdk';

const client = new Anthropic({
  apiKey: 'YOUR_API_KEY',
  baseURL: 'https://sqllm.dpdns.org'
});

async function main() {
  const stream = await client.messages.stream({
    model: "claude-3-sonnet-20240229",
    max_tokens: 1024,
    messages: [
      { role: "user", content: "写一篇关于季节的诗" }
    ],
  });

  for await (const event of stream) {
    if (event.type === 'content_block_delta' &&
        event.delta.type === 'text_delta') {
      process.stdout.write(event.delta.text);
    }
  }
}

main();`

const curlStreamExample = `curl https://sqllm.dpdns.org/v1/messages \\
  -H "Content-Type: application/json" \\
  -H "x-api-key: YOUR_API_KEY" \\
  -H "anthropic-version: 2023-06-01" \\
  -d '{
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 1024,
    "stream": true,
    "messages": [
      {"role": "user", "content": "写一篇关于季节的诗"}
    ]
  }'`

const multiTurnExample = `import anthropic

client = anthropic.Anthropic(
    api_key="YOUR_API_KEY",
    base_url="https://sqllm.dpdns.org"
)

# 维护对话历史
conversation = []

def chat(user_message):
    conversation.append({"role": "user", "content": user_message})
    response = client.messages.create(
        model="claude-3-sonnet-20240229",
        max_tokens=1024,
        messages=conversation
    )
    assistant_message = response.content[0].text
    conversation.append({"role": "assistant", "content": assistant_message})
    return assistant_message

# 开始对话
print(chat("你好，我想学习 Python"))
print(chat("有什么好的学习资源推荐吗?"))
print(chat("谢谢，优先从基础语法开始学起"))`

const systemPromptExample = `message = client.messages.create(
    model="claude-3-sonnet-20240229",
    max_tokens=1024,
    system="你是一个专业的 Python 编程导师，回答简明扼要。",
    messages=[
        {"role": "user", "content": "如何实现 递归 文件查询?"}
    ]
)`

const getModelsExample = `curl https://sqllm.dpdns.org/v1/models \\
  -H "x-api-key: YOUR_API_KEY"`

const responseExample = `{
  "id": "msg_01XFDJH29jABCvnyhtVHYEL",
  "type": "message",
  "role": "assistant",
  "content": [
    {
      "type": "text",
      "text": "人工智能 (AI) 是...."
    }
  ],
  "model": "claude-3-sonnet-20240229",
  "stop_reason": "end_turn",
  "usage": {
    "input_tokens": 12,
    "output_tokens": 155
  }
}`

const errorHandlingExample = `import anthropic

client = anthropic.Anthropic(
    api_key="YOUR_API_KEY",
    base_url="https://sqllm.dpdns.org"
)

try:
    message = client.messages.create(
        model="claude-3-sonnet-20240229",
        max_tokens=1024,
        messages=[{"role": "user", "content": "Hello"}]
    )
except anthropic.APIConnectionError:
    print("连接失败，请检查网络环境")
except anthropic.RateLimitError:
    print("请求过于频繁，请稍后再试")
except anthropic.APIStatusError as e:
    print(f"API 错误: {e.status_code} - {e.message}")`

const claudeModels = [
  { name: 'claude-3-sonnet-20240229', desc: 'Claude 3 Sonnet, 平衡性能和成本' },
  { name: 'claude-3-opus-20240229', desc: 'Claude 3 Opus, 最强能力' },
  { name: 'claude-3-5-haiku-20241022', desc: 'Claude 3.5 Haiku, 极速响应' },
  { name: 'claude-3-5-sonnet-20240620', desc: 'Claude 3.5 Sonnet' }
]

const requestParams = [
  { name: 'model', type: 'string', required: true, desc: '模型 ID' },
  { name: 'messages', type: 'array', required: true, desc: '对话消息数组' },
  { name: 'max_tokens', type: 'integer', required: true, desc: '最大输出 token 数' },
  { name: 'system', type: 'string', required: false, desc: '系统提示词' },
  { name: 'stream', type: 'boolean', required: false, desc: '是否流式输出' },
  { name: 'temperature', type: 'number', required: false, desc: '温度, 0-1, 默认 1' },
  { name: 'top_p', type: 'number', required: false, desc: 'Top-P 采样' },
  { name: 'top_k', type: 'integer', required: false, desc: 'Top-K 采样' },
  { name: 'stop_sequences', type: 'array', required: false, desc: '停止序列' }
]

const errorCodes = [
  { code: '400', desc: '请求参数错误' },
  { code: '401', desc: 'API Key 无效' },
  { code: '402', desc: '余额不足' },
  { code: '429', desc: '请求过于频繁' },
  { code: '500', desc: '服务器内部错误' },
  { code: '529', desc: '上游服务过载' }
]
</script>
