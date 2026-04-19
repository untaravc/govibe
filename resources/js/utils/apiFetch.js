import { useUiStore } from "../store/ui.js";
import { withApiBaseUrl } from "../config/baseUrl.js";

function isApiUrl(input) {
  try {
    if (typeof input === "string") return input.startsWith("/api");
    if (input instanceof URL) return input.pathname.startsWith("/api");
    if (typeof input?.url === "string") return input.url.startsWith("/api");
  } catch {
    // ignore
  }
  return false;
}

function getStoredAccessToken() {
  try {
    return (
      (localStorage.getItem("access_token") ||
        sessionStorage.getItem("access_token") ||
        "").trim()
    );
  } catch {
    return "";
  }
}

export async function apiFetch(input, init) {
  const track = isApiUrl(input);
  const ui = track ? useUiStore() : null;

  if (track) ui.apiStart();
  try {
    const finalInit = init ? { ...init } : {};
    const headers = new Headers(finalInit.headers || {});

    if (!headers.has("Authorization")) {
      const accessToken = getStoredAccessToken();
      if (accessToken) headers.set("Authorization", `Bearer ${accessToken}`);
    }

    finalInit.headers = headers;
    const finalInput = typeof input === "string" ? withApiBaseUrl(input) : input;
    return await fetch(finalInput, finalInit);
  } finally {
    if (track) ui.apiStop();
  }
}
