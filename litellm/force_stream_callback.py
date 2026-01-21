"""
强制流式回调 - 在所有请求中添加 stream=true
"""
from litellm.integrations.custom_logger import CustomLogger
from typing import Optional, Dict, Any


class ForceStreamCallback(CustomLogger):
    """
    自定义回调：强制所有请求使用流式模式
    """

    async def async_pre_call_hook(
        self,
        user_api_key_dict: Dict[str, Any],
        cache: Any,
        data: Dict[str, Any],
        call_type: str,
    ) -> Optional[Dict[str, Any]]:
        """
        在 LLM 调用之前修改请求参数，强制添加 stream=true
        """
        # 强制设置 stream=true
        data["stream"] = True
        return data

    def log_pre_api_call(self, model, messages, kwargs):
        """
        在 API 调用前记录日志（可选）
        """
        # 确保 stream 参数被设置
        if "stream" not in kwargs or not kwargs.get("stream"):
            kwargs["stream"] = True

    def log_success_event(self, kwargs, response_obj, start_time, end_time):
        """成功回调（可选）"""
        pass

    def log_failure_event(self, kwargs, response_obj, start_time, end_time):
        """失败回调（可选）"""
        pass


# 导出回调实例
force_stream_callback = ForceStreamCallback()
