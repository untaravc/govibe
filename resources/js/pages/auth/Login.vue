<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <h2 class="text-2xl font-semibold tracking-tight text-slate-900 dark:text-slate-50">Login</h2>
      <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Sign in to continue.</p>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Email</label>
          <input
            v-model.trim="email"
            type="email"
            autocomplete="email"
            placeholder="you@example.com"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            required
          />
          <p v-if="fieldErrors.email" class="mt-2 text-sm text-rose-600 dark:text-rose-400">
            {{ fieldErrors.email }}
          </p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Password</label>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            placeholder="••••••••"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            required
          />
          <p v-if="fieldErrors.password" class="mt-2 text-sm text-rose-600 dark:text-rose-400">
            {{ fieldErrors.password }}
          </p>
        </div>

        <div class="flex items-center justify-between">
          <RouterLink to="/auth/forgot-password" class="text-sm font-medium text-slate-900 hover:underline dark:text-slate-50">
            Forgot password?
          </RouterLink>
        </div>

        <button
          type="submit"
          class="inline-flex w-full items-center justify-center rounded-xl bg-slate-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-slate-800 focus:outline-none focus:ring-4 focus:ring-slate-900/20 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
          :disabled="auth.loading"
        >
          {{ auth.loading ? "Signing in..." : "Sign in" }}
        </button>

        <p class="text-center text-sm text-slate-600 dark:text-slate-300">
          <RouterLink to="/auth/register" class="font-medium text-slate-900 hover:underline dark:text-slate-50">
            Create an account
          </RouterLink>
          <span class="px-2 text-slate-400">·</span>
          <RouterLink to="/" class="font-medium text-slate-900 hover:underline dark:text-slate-50">Home</RouterLink>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";

import { useAuthStore } from "../../store/auth.js";

const router = useRouter();
const auth = useAuthStore();

const email = ref("");
const password = ref("");

const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-emerald-700 dark:text-emerald-300";
  if (messageTone.value === "error") return "text-rose-700 dark:text-rose-300";
  return "text-slate-700 dark:text-slate-200";
});

const fieldErrors = computed(() => auth.fieldErrors || {});

async function onSubmit() {
  message.value = "";
  messageTone.value = "neutral";

  const ok = await auth.login(email.value, password.value);
  if (!ok) {
    message.value = auth.error || "Login failed";
    messageTone.value = "error";
    return;
  }

  try {
    localStorage.setItem("access_token", auth.token);
    localStorage.setItem("refresh_token", auth.refreshToken || "");
  } catch {
    // ignore storage failures
  }

  message.value = "Signed in.";
  messageTone.value = "success";
  router.push("/admin");
}
</script>
