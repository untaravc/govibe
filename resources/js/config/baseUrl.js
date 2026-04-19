function normalizeBaseUrl(value) {
  const url = String(value || "").trim();
  if (!url) return "";
  return url.endsWith("/") ? url.slice(0, -1) : url;
}

// Default for local development. Override via Vite env:
// - VITE_BASE_URL
// - VITE_API_BASE_URL
export const BASE_URL = normalizeBaseUrl(import.meta.env.VITE_BASE_URL || "http://127.0.0.1:3000");
export const API_BASE_URL = normalizeBaseUrl(import.meta.env.VITE_API_BASE_URL || BASE_URL);

export function withApiBaseUrl(input) {
  if (typeof input !== "string") return input;
  const url = input.trim();
  if (!url) return url;
  if (/^https?:\/\//i.test(url) || url.startsWith("//")) return url;
  if (!url.startsWith("/")) return url;
  return `${API_BASE_URL}${url}`;
}

