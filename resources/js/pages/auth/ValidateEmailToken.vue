<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <h2 class="text-2xl font-semibold tracking-tight text-slate-900">Reset password</h2>
      <p class="mt-1 text-sm text-slate-600">Set a new password for your account.</p>

      <div class="mt-6">
        <p v-if="validating" class="text-sm text-slate-700">Validating token...</p>
        <p v-else-if="!token" class="text-sm text-danger">Missing token.</p>
        <p v-else-if="!valid" class="text-sm text-danger">
          {{ message || "Invalid or expired token." }}
        </p>
      </div>

      <form v-if="token && valid && !validating && !done" class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="text-sm font-medium text-slate-700">New password</label>
          <input
            v-model.trim="newPassword"
            type="password"
            autocomplete="new-password"
            placeholder="••••••••"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.new_password" class="mt-2 text-sm text-danger">{{ fieldErrors.new_password }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Confirm new password</label>
          <input
            v-model.trim="confirmPassword"
            type="password"
            autocomplete="new-password"
            placeholder="••••••••"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="confirmError" class="mt-2 text-sm text-danger">{{ confirmError }}</p>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <button
          type="submit"
          class="inline-flex w-full items-center justify-center rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="loading"
        >
          {{ loading ? "Updating..." : "Update password" }}
        </button>

        <p class="text-center text-sm text-slate-600">
          <RouterLink to="/auth/login" class="font-medium text-slate-900 hover:underline">Back to login</RouterLink>
        </p>
      </form>

      <div v-else-if="done" class="mt-6 space-y-4">
        <p class="text-sm text-success">{{ message || "Password updated. You can now log in." }}</p>
        <RouterLink
          to="/auth/login"
          class="inline-flex w-full items-center justify-center rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          Go to login
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRoute } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const route = useRoute();

const token = computed(() => String(route.query?.token || "").trim());

const validating = ref(false);
const valid = ref(false);
const done = ref(false);

const newPassword = ref("");
const confirmPassword = ref("");
const confirmError = ref("");

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

async function validateToken() {
  if (!token.value) {
    valid.value = false;
    return;
  }

  validating.value = true;
  message.value = "";
  messageTone.value = "neutral";

  try {
    const url = `/api/validate-email-token?token=${encodeURIComponent(token.value)}`;
    const { res, json } = await api.get(url, { auth: false, navigate: false });

    if (res.ok && json?.success && json?.result?.valid === true) {
      valid.value = true;
      return;
    }

    valid.value = false;
    message.value = apiErrorMessage(json, res.status === 404 ? "Invalid or expired token." : `Request failed (${res.status})`);
    messageTone.value = "error";
  } catch (err) {
    valid.value = false;
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    validating.value = false;
  }
}

async function onSubmit() {
  confirmError.value = "";
  fieldErrors.value = {};
  message.value = "";
  messageTone.value = "neutral";

  if (newPassword.value !== confirmPassword.value) {
    confirmError.value = "Passwords do not match";
    return;
  }

  loading.value = true;
  try {
    const { res, json } = await api.post(
      "/api/update-password-with-token",
      { email_token: token.value, new_password: newPassword.value },
      { auth: false, navigate: false }
    );

    if (res.status === 422) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, "Validation error");
      messageTone.value = "error";
      return;
    }

    if (res.status === 404) {
      valid.value = false;
      message.value = apiErrorMessage(json, "Invalid or expired token.");
      messageTone.value = "error";
      return;
    }

    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    message.value = json?.message || "Password updated. You can now log in.";
    messageTone.value = "success";
    done.value = true;
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  validateToken();
});
</script>
