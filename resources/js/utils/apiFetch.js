import { useUiStore } from "../store/ui.js";

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

export async function apiFetch(input, init) {
  const track = isApiUrl(input);
  const ui = track ? useUiStore() : null;

  if (track) ui.apiStart();
  try {
    return await fetch(input, init);
  } finally {
    if (track) ui.apiStop();
  }
}

