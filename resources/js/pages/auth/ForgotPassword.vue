<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <h2 class="text-2xl font-semibold tracking-tight text-slate-900">Forgot password</h2>
      <p class="mt-1 text-sm text-slate-600">Request a reset token using email or phone.</p>

      <div class="mt-6">
        <div
          role="tablist"
          aria-label="Reset method"
          class="relative flex items-center rounded-2xl border border-slate-200 bg-slate-50 p-1"
        >
          <div
            aria-hidden="true"
            class="absolute inset-y-1 left-1 w-[calc(50%-0.25rem)] rounded-xl bg-white shadow-sm ring-1 ring-slate-900/5 transition-transform duration-200"
            :class="tab === 'phone' ? 'translate-x-full' : ''"
          ></div>

          <button
            role="tab"
            type="button"
            class="relative z-10 flex flex-1 items-center justify-center gap-2 rounded-xl px-3 py-2 text-sm font-medium transition-colors"
            :aria-selected="tab === 'email' ? 'true' : 'false'"
            :class="tab === 'email' ? 'bg-white text-slate-900' : 'text-slate-600 hover:text-slate-900'"
            @click="setTab('email')"
          >
            <svg viewBox="0 0 24 24" fill="none" class="h-4 w-4" xmlns="http://www.w3.org/2000/svg">
              <path
                d="M4 7.5A2.5 2.5 0 0 1 6.5 5h11A2.5 2.5 0 0 1 20 7.5v9A2.5 2.5 0 0 1 17.5 19h-11A2.5 2.5 0 0 1 4 16.5v-9Z"
                stroke="currentColor"
                stroke-width="2"
              />
              <path
                d="m5.5 7.5 6.1 4.2a1 1 0 0 0 1.2 0l6.1-4.2"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            Email
          </button>

          <button
            role="tab"
            type="button"
            class="relative z-10 flex flex-1 items-center justify-center gap-2 rounded-xl px-3 py-2 text-sm font-medium transition-colors"
            :aria-selected="tab === 'phone' ? 'true' : 'false'"
            :class="tab === 'phone' ? 'bg-white text-slate-900' : 'text-slate-600 hover:text-slate-900'"
            @click="setTab('phone')"
          >
            <svg viewBox="0 0 24 24" fill="none" class="h-4 w-4" xmlns="http://www.w3.org/2000/svg">
              <path
                d="M8 3h8a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2Z"
                stroke="currentColor"
                stroke-width="2"
                stroke-linejoin="round"
              />
              <path
                d="M10 18h4"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              />
            </svg>
            Phone
          </button>
        </div>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div v-if="tab === 'email'">
          <label class="text-sm font-medium text-slate-700">Email</label>
          <input
            v-model.trim="email"
            type="email"
            autocomplete="email"
            placeholder="you@example.com"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.email" class="mt-2 text-sm text-danger">{{ fieldErrors.email }}</p>
        </div>

        <div v-else>
          <label class="text-sm font-medium text-slate-700">Phone</label>
          <input
            v-model.trim="phone"
            type="tel"
            autocomplete="tel"
            placeholder="e.g. +628123456789"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.phone" class="mt-2 text-sm text-danger">{{ fieldErrors.phone }}</p>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <button
          type="submit"
          class="inline-flex w-full items-center justify-center rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="loading"
        >
          {{ loading ? "Sending..." : "Send reset token" }}
        </button>

        <p class="text-center text-sm text-slate-600">
          <RouterLink to="/auth/login" class="font-medium text-slate-900 hover:underline">Back to login</RouterLink>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { RouterLink } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const tab = ref("email"); // email | phone

const email = ref("");
const phone = ref("");
const loading = ref(false);

const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

function setTab(next) {
  tab.value = next;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};
}

async function onSubmit() {
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};
  loading.value = true;

  try {
    const payload = tab.value === "email" ? { email: email.value } : { phone: phone.value };
    const { res, json } = await api.post("/api/request-reset-password", payload, { auth: false });

    if (res.status === 422) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, "Validation error");
      messageTone.value = "error";
      return;
    }

    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    message.value = json?.message || "If the account exists, a reset token has been sent.";
    messageTone.value = "success";
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}
</script>
