import { apiRequest } from "./api";

export async function trackEvent(eventName, options = {}) {
  if (!eventName) {
    return;
  }
  const { taskId, properties = {}, token } = options;
  // 前端埋点依赖用户 JWT。无 token 场景静默跳过，避免污染控制台。
  if (!token && typeof window !== "undefined" && !window.localStorage.getItem("jwt_token")) {
    return;
  }
  try {
    await apiRequest("/analytics/events", {
      method: "POST",
      token,
      body: {
        event_name: eventName,
        task_id: taskId || undefined,
        properties
      }
    });
  } catch (_) {
    // Ignore analytics failures.
  }
}
