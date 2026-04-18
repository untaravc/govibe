<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <h2 class="text-2xl font-semibold tracking-tight text-slate-900">Forgot password</h2>
      <p class="mt-1 text-sm text-slate-600">Request a reset token using email or phone.</p>

      <div class="mt-6">
        <div class="flex items-center gap-2 rounded-xl border border-slate-200 bg-slate-50 p-1">
          <button
            type="button"
            class="flex-1 rounded-lg px-3 py-2 text-sm font-medium"
            :class="tab === 'email' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'"
            @click="setTab('email')"
          >
            Email
          </button>
          <button
            type="button"
            class="flex-1 rounded-lg px-3 py-2 text-sm font-medium"
            :class="tab === 'phone' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'"
            @click="setTab('phone')"
          >
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

