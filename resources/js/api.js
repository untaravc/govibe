import router from "./route.js";
import { useUiStore } from "./store/ui.js";
import { withApiBaseUrl } from "./config/baseUrl.js";

function getStoredToken() {
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

function getStoredRefreshToken() {
  try {
    return String(localStorage.getItem("refresh_token") || "").trim();
  } catch {
    return "";
  }
}

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

function maybeNavigateForStatus(status) {
  const currentPath = router.currentRoute?.value?.path || "";
  const routes = {
    401: "/auth/unauthenticated",
    402: "/auth/unauthorized",
    403: "/auth/unauthorized",
    404: "/auth/not-found"
  };
  const target = routes[status];
  if (!target) return;
  if (currentPath === target) return;
  router.push(target).catch(() => {});
}

let refreshPromise = null;

async function refreshTokens() {
  if (refreshPromise) return refreshPromise;

  const refreshToken = getStoredRefreshToken();
  if (!refreshToken) return null;

  refreshPromise = (async () => {
    const res = await fetch(withApiBaseUrl("/api/refresh-token"), {
      method: "POST",
      headers: {
        Authorization: `Bearer ${refreshToken}`
      }
    });

    const json = await res.json().catch(() => null);
    const result = json?.result || {};

    const accessToken = typeof result?.access_token === "string" ? result.access_token : "";
    const newRefreshToken = typeof result?.refresh_token === "string" ? result.refresh_token : "";

    if (!res.ok || !json?.success || !accessToken || !newRefreshToken) return null;

    try {
      localStorage.setItem("access_token", accessToken);
      localStorage.setItem("refresh_token", newRefreshToken);
    } catch {
      // ignore storage failures
    }

    return { accessToken, refreshToken: newRefreshToken };
  })();

  try {
    return await refreshPromise;
  } finally {
    refreshPromise = null;
  }
}

async function request(url, options = {}) {
  const {
    method = "GET",
    auth = false,
    headers = {},
    body = undefined
  } = options;
  const retried = Boolean(options?._retried);

  const finalHeaders = { ...headers };

  if (auth) {
    const token = getStoredToken();
    if (token) finalHeaders.Authorization = `Bearer ${token}`;
  }

  let finalBody = body;
  const isFormData = typeof FormData !== "undefined" && body instanceof FormData;
  if (!isFormData && body && typeof body === "object" && !(body instanceof Blob)) {
    finalHeaders["Content-Type"] = finalHeaders["Content-Type"] || "application/json";
    finalBody = JSON.stringify(body);
  }

  const track = isApiUrl(url);
  const ui = track ? useUiStore() : null;

  if (track) ui.apiStart();
  try {
    const res = await fetch(typeof url === "string" ? withApiBaseUrl(url) : url, {
      method,
      headers: finalHeaders,
      body: method === "GET" || method === "HEAD" ? undefined : finalBody
    });

    const json = await res.json().catch(() => null);

    if (!res.ok) {
      if (
        res.status === 401 &&
        !retried &&
        auth &&
        typeof url === "string" &&
        url !== "/api/login" &&
        url !== "/api/refresh-token"
      ) {
        const refreshed = await refreshTokens();
        if (refreshed) {
          return request(url, { ...options, _retried: true });
        }
        try {
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
        } catch {
          // ignore
        }
      }

      maybeNavigateForStatus(res.status);
    }

    return { res, json };
  } finally {
    if (track) ui.apiStop();
  }
}

async function get(url, options = {}) {
  return request(url, { ...options, method: "GET" });
}

async function post(url, body, options = {}) {
  return request(url, { ...options, method: "POST", body });
}

async function put(url, body, options = {}) {
  return request(url, { ...options, method: "PUT", body });
}

async function del(url, options = {}) {
  return request(url, { ...options, method: "DELETE" });
}

export const api = {
  fetch: request,
  get,
  post,
  put,
  delete: del
};

export default api;
