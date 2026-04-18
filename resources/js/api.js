import router from "./route.js";
import { useUiStore } from "./store/ui.js";

function getStoredToken() {
  try {
    return (localStorage.getItem("token") || sessionStorage.getItem("token") || "").trim();
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

async function request(url, options = {}) {
  const {
    method = "GET",
    auth = false,
    headers = {},
    body = undefined
  } = options;

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
    const res = await fetch(url, {
      method,
      headers: finalHeaders,
      body: method === "GET" || method === "HEAD" ? undefined : finalBody
    });

    const json = await res.json().catch(() => null);

    if (!res.ok) {
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

