# sub2api_passthrough.py
# 将客户端的 API Key 透传给 Sub2API，实现动态密钥绑定

from litellm.integrations.custom_logger import CustomLogger
from litellm.proxy._types import UserAPIKeyAuth
from fastapi import Request

# 使用全局字典存储请求的原始 API Key（按请求 ID）
_request_api_keys = {}


async def user_api_key_auth(request: Request, api_key: str) -> UserAPIKeyAuth:
    """
    自定义认证：捕获原始的客户端 API Key
    使用请求的唯一标识符存储原始 key
    """
    # 跳过空 key（如 health 检查）
    if not api_key:
        return UserAPIKeyAuth(api_key="", user_id="")

    # 使用请求对象的 id 作为键
    request_id = id(request)
    _request_api_keys[request_id] = api_key

    # 也存储到 request.state 中，方便后续获取
    request.state.original_api_key = api_key

    print(f"[Sub2API Auth] Captured original API key: {api_key[:20]}... (request_id={request_id})")

    # 返回认证信息，将原始 key 放入 user_id 字段（hack）
    return UserAPIKeyAuth(
        api_key=api_key,
        user_id=f"sub2api:{api_key}"  # 将原始 key 编码到 user_id 中
    )


class Sub2APIPassthrough(CustomLogger):
    """
    将客户端的 API Key 透传给 Sub2API
    同时强制使用流式模式（Sub2API 非流式有 bug）
    """

    async def async_pre_call_hook(self, user_api_key_dict, cache, data, call_type):
        # 强制使用流式模式
        data["stream"] = True

        # 方式1: 从请求体获取 api_key（客户端通过 extra_body 传递）
        if "api_key" in data and data["api_key"]:
            print(f"[Sub2API Passthrough] Using api_key from request body: {data['api_key'][:20]}...")
            return data

        # 方式2: 从 user_id 字段获取（我们在认证时编码的）
        if hasattr(user_api_key_dict, 'user_id') and user_api_key_dict.user_id:
            user_id = user_api_key_dict.user_id
            if user_id.startswith("sub2api:"):
                original_key = user_id[8:]  # 移除 "sub2api:" 前缀
                data["api_key"] = original_key
                print(f"[Sub2API Passthrough] Using api_key from user_id: {original_key[:20]}...")
                return data

        # 方式3: 打印所有可用属性用于调试
        print(f"[Sub2API Passthrough] user_api_key_dict attributes: {dir(user_api_key_dict)}")
        if hasattr(user_api_key_dict, '__dict__'):
            for key, value in user_api_key_dict.__dict__.items():
                if value and isinstance(value, str) and len(value) > 0:
                    print(f"[Sub2API Passthrough]   {key}: {str(value)[:30]}...")

        print("[Sub2API Passthrough] No API key found!")
        return data


sub2api_passthrough = Sub2APIPassthrough()
print("[Sub2API Passthrough] Module loaded successfully")
