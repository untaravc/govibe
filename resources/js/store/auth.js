import { defineStore } from "pinia";

import { apiErrorMessage, apiFieldErrors } from "../utils/apiError.js";
import api from "../api.js";

const STORAGE_KEY = "govibe.auth";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: "",
    refreshToken: "",
    user: null,
    loading: false,
    error: "",
    fieldErrors: {}
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token)
  },
  actions: {
    hydrate() {
      try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (!raw) return;
        const parsed = JSON.parse(raw);
        this.token = typeof parsed?.token === "string" ? parsed.token : "";
        this.refreshToken = typeof parsed?.refreshToken === "string" ? parsed.refreshToken : "";
        this.user = parsed?.user ?? null;
      } catch {
        // ignore invalid storage
      }
    },
    persist() {
      try {
        localStorage.setItem(
          STORAGE_KEY,
          JSON.stringify({
            token: this.token,
            refreshToken: this.refreshToken,
            user: this.user
          })
        );
      } catch {
        // ignore storage failures
      }
    },
    clear() {
      this.token = "";
      this.refreshToken = "";
      this.user = null;
      this.error = "";
      try {
        localStorage.removeItem(STORAGE_KEY);
      } catch {
        // ignore
      }
    },
    async login(email, password) {
      this.loading = true;
      this.error = "";
      this.fieldErrors = {};

      try {
        const { res, json } = await api.post("/api/login", { email, password }, { auth: false });

        if (!res.ok || json?.success === false) {
          this.fieldErrors = apiFieldErrors(json);
          this.error = apiErrorMessage(json, `Login failed (${res.status})`);
          return false;
        }

        const result = json?.result ?? {};
        this.token = typeof result?.access_token === "string" ? result.access_token : "";
        this.refreshToken = typeof result?.refresh_token === "string" ? result.refresh_token : "";
        this.user = result?.user ?? null;

        if (!this.token) {
          this.error = "Login failed: missing access token";
          return false;
        }

        this.persist();
        return true;
      } catch (err) {
        this.error = String(err);
        return false;
      } finally {
        this.loading = false;
      }
    },
    logout() {
      this.clear();
    }
  }
});
