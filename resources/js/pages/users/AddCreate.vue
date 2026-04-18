<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">{{ isEdit ? "Edit user" : "Add user" }}</h3>
          <p class="mt-1 text-sm text-slate-600">
            {{ isEdit ? "Update user details." : "Create a new user account." }}
          </p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="router.push('/admin/users')"
        >
          Back
        </button>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="text-sm font-medium text-slate-700">Name</label>
          <input
            v-model.trim="name"
            type="text"
            autocomplete="name"
            placeholder="Full name"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.name" class="mt-2 text-sm text-danger">{{ fieldErrors.name }}</p>
        </div>

        <div>
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

        <div>
          <label class="text-sm font-medium text-slate-700">
            Password
            <span v-if="isEdit" class="font-normal text-slate-500">(leave empty to keep)</span>
          </label>
          <input
            v-model="password"
            type="password"
            :autocomplete="isEdit ? 'new-password' : 'new-password'"
            placeholder="••••••••"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            :required="!isEdit"
          />
          <p v-if="fieldErrors.password" class="mt-2 text-sm text-danger">{{ fieldErrors.password }}</p>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : isEdit ? "Save changes" : "Create user" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import Swal from "sweetalert2";

import api from "../../api.js";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const router = useRouter();
const route = useRoute();

const id = computed(() => String(route.params.id || "").trim());
const isEdit = computed(() => id.value.length > 0);

const name = ref("");
const email = ref("");
const password = ref("");
const loading = ref(false);

const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

async function loadUser() {
  if (!isEdit.value) return;

  loading.value = true;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};

  try {
    const { res, json } = await api.get(`/api/users/${id.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const u = json?.result?.user;
    name.value = typeof u?.name === "string" ? u.name : "";
    email.value = typeof u?.email === "string" ? u.email : "";
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

async function onSubmit() {
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};
  loading.value = true;

  try {
    const body = {
      name: name.value,
      email: email.value
    };

    if (!isEdit.value || password.value.trim().length > 0) {
      body.password = password.value;
    }

    const { res, json } = isEdit.value
      ? await api.put(`/api/users/${id.value}`, body, { auth: true })
      : await api.post("/api/users", body, { auth: true });
    if (!res.ok) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    message.value = isEdit.value ? "User updated." : "User created.";
    messageTone.value = "success";
    Swal.fire({
      toast: true,
      position: "top-end",
      icon: "success",
      title: isEdit.value ? "User updated" : "User created",
      showConfirmButton: false,
      timer: 1600,
      timerProgressBar: true
    });
    password.value = "";
    setTimeout(() => router.push("/admin/users"), 300);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadUser();
});
</script>
