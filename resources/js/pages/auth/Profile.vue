<template>
  <div class="space-y-4">
    <div
      class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900"
    >
      <h3 class="text-lg font-semibold tracking-tight text-slate-900 dark:text-slate-50">Profile</h3>
      <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Update your account details.</p>

      <form class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2" @submit.prevent="onSubmit">
        <div class="sm:col-span-2">
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Name</label>
          <input
            v-model.trim="form.name"
            type="text"
            autocomplete="name"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            placeholder="Your name"
            required
          />
          <p v-if="fieldErrors.name" class="mt-2 text-sm text-danger">{{ fieldErrors.name }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Email</label>
          <input
            v-model.trim="form.email"
            type="email"
            autocomplete="email"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            placeholder="you@example.com"
            required
          />
          <p v-if="fieldErrors.email" class="mt-2 text-sm text-danger">{{ fieldErrors.email }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Phone</label>
          <input
            v-model.trim="form.phone"
            type="tel"
            autocomplete="tel"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            placeholder="e.g. +628123456789"
          />
          <p v-if="fieldErrors.phone" class="mt-2 text-sm text-danger">{{ fieldErrors.phone }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">New password</label>
          <input
            v-model.trim="form.password"
            type="password"
            autocomplete="new-password"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            placeholder="••••••••"
          />
          <p v-if="fieldErrors.password" class="mt-2 text-sm text-danger">{{ fieldErrors.password }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Confirm password</label>
          <input
            v-model.trim="form.passwordConfirm"
            type="password"
            autocomplete="new-password"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            placeholder="••••••••"
          />
          <p v-if="confirmError" class="mt-2 text-sm text-danger">{{ confirmError }}</p>
        </div>

        <div class="sm:col-span-2 space-y-3">
          <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>
          <button
            type="submit"
            class="inline-flex items-center justify-center rounded-xl bg-slate-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : "Save changes" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";

import api from "../../api.js";
import { useAuthStore } from "../../store/auth.js";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const auth = useAuthStore();

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});
const confirmError = ref("");

const form = ref({
  name: "",
  email: "",
  phone: "",
  password: "",
  passwordConfirm: ""
});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700 dark:text-slate-200";
});

function setFromUser(user) {
  form.value.name = String(user?.name || "").trim();
  form.value.email = String(user?.email || "").trim();
  form.value.phone = String(user?.phone || "").trim();
}

async function loadProfile() {
  setFromUser(auth.user);

  try {
    const { res, json } = await api.get("/api/profile", { auth: true, navigate: false });
    if (!res.ok) return;
    const user = json?.result?.user;
    if (!user) return;
    setFromUser(user);

    auth.user = { ...(auth.user || {}), ...(user || {}) };
    auth.persist();
  } catch {
    // ignore
  }
}

async function onSubmit() {
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};
  confirmError.value = "";

  const password = String(form.value.password || "").trim();
  const passwordConfirm = String(form.value.passwordConfirm || "").trim();
  if (password || passwordConfirm) {
    if (password !== passwordConfirm) {
      confirmError.value = "Passwords do not match";
      return;
    }
  }

  const payload = {
    name: String(form.value.name || "").trim(),
    email: String(form.value.email || "").trim(),
    phone: String(form.value.phone || "").trim()
  };
  if (password) payload.password = password;

  loading.value = true;
  try {
    const { res, json } = await api.fetch("/api/profile", {
      method: "PATCH",
      auth: true,
      navigate: false,
      body: payload
    });

    if (res.status === 422) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, "Validation error");
      messageTone.value = "error";
      return;
    }

    if (res.status === 409) {
      message.value = apiErrorMessage(json, "Email already registered");
      messageTone.value = "error";
      return;
    }

    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const user = json?.result?.user || null;
    if (user) {
      auth.user = { ...(auth.user || {}), ...(user || {}) };
      auth.persist();
    }

    form.value.password = "";
    form.value.passwordConfirm = "";

    message.value = json?.message || "Profile updated";
    messageTone.value = "success";
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadProfile();
});
</script>

