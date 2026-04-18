export function apiErrorMessage(payload, fallbackMessage = "Request failed") {
  if (!payload || typeof payload !== "object") return fallbackMessage;

  const message =
    typeof payload.message === "string" && payload.message.trim().length > 0
      ? payload.message.trim()
      : fallbackMessage;

  const errors = payload?.result?.errors;
  if (!errors || typeof errors !== "object") return message;

  const parts = [];
  for (const value of Object.values(errors)) {
    if (typeof value === "string" && value.trim().length > 0) parts.push(value.trim());
  }

  if (parts.length === 0) return message;
  return `${message}: ${parts.join(", ")}`;
}

export function apiFieldErrors(payload) {
  const errors = payload?.result?.errors;
  if (!errors || typeof errors !== "object") return {};

  const out = {};
  for (const [key, value] of Object.entries(errors)) {
    if (typeof value === "string" && value.trim().length > 0) out[key] = value.trim();
  }
  return out;
}
