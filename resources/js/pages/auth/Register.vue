<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <h2 class="text-2xl font-semibold tracking-tight text-slate-900 dark:text-slate-50">Create account</h2>
      <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Register a new user.</p>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Name</label>
          <input
            v-model.trim="name"
            type="text"
            autocomplete="name"
            placeholder="Your name"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            required
          />
          <p v-if="fieldErrors.name" class="mt-2 text-sm text-rose-600 dark:text-rose-400">
            {{ fieldErrors.name }}
          </p>
        </div>

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
            autocomplete="new-password"
            placeholder="••••••••"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            required
          />
          <p v-if="fieldErrors.password" class="mt-2 text-sm text-rose-600 dark:text-rose-400">
            {{ fieldErrors.password }}
          </p>
        </div>

        <button
          type="submit"
          class="inline-flex w-full items-center justify-center rounded-xl bg-slate-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-slate-800 focus:outline-none focus:ring-4 focus:ring-slate-900/20 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
          :disabled="loading"
        >
          {{ loading ? "Creating..." : "Create account" }}
        </button>

        <p v-if="message" class="text-sm" :class="messageToneClass">
          {{ message }}
        </p>

        <div class="flex items-center justify-between text-sm">
          <RouterLink to="/auth/login" class="font-medium text-slate-900 hover:underline dark:text-slate-50">
            Already have an account?
          </RouterLink>
          <RouterLink to="/" class="text-slate-600 hover:underline dark:text-slate-300">Home</RouterLink>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";

import { apiFetch } from "../../utils/apiFetch.js";

const router = useRouter();

const name = ref("");
const email = ref("");
const password = ref("");
const loading = ref(false);

const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error

const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-emerald-700 dark:text-emerald-300";
  if (messageTone.value === "error") return "text-rose-700 dark:text-rose-300";
  return "text-slate-700 dark:text-slate-200";
});

async function onSubmit() {
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};
  loading.value = true;

  try {
    const res = await apiFetch("/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: name.value,
        email: email.value,
        password: password.value
      })
    });

    const json = await res.json().catch(() => null);

    if (res.status === 201) {
      message.value = "Account created. Redirecting to login...";
      messageTone.value = "success";
      setTimeout(() => router.push("/auth/login"), 600);
      return;
    }

    if (res.status === 422 && json && typeof json === "object") {
      fieldErrors.value = json?.result?.errors || {};
      message.value = json.message || "Validation error";
      messageTone.value = "error";
      return;
    }

    message.value = (json && json.message) || `Request failed (${res.status})`;
    messageTone.value = "error";
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}
</script>
