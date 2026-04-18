<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">{{ isEdit ? "Edit role" : "Add role" }}</h3>
          <p class="mt-1 text-sm text-slate-600">
            {{ isEdit ? "Update role properties." : "Create a new role." }}
          </p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="router.push('/admin/roles')"
        >
          Back
        </button>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="text-sm font-medium text-slate-700">Role</label>
          <input
            v-model.trim="role"
            type="text"
            autocomplete="off"
            placeholder="e.g. admin"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.role" class="mt-2 text-sm text-danger">{{ fieldErrors.role }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Name</label>
          <input
            v-model.trim="name"
            type="text"
            autocomplete="off"
            placeholder="Display name"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.name" class="mt-2 text-sm text-danger">{{ fieldErrors.name }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Status</label>
          <select
            v-model.number="status"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
          >
            <option :value="1">Active</option>
            <option :value="0">Inactive</option>
          </select>
          <p v-if="fieldErrors.status" class="mt-2 text-sm text-danger">{{ fieldErrors.status }}</p>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : isEdit ? "Save changes" : "Create role" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const router = useRouter();
const route = useRoute();

const id = computed(() => String(route.params.id || "").trim());
const isEdit = computed(() => id.value.length > 0);

const role = ref("");
const name = ref("");
const status = ref(1);
const loading = ref(false);

const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

async function loadRole() {
  if (!isEdit.value) return;

  loading.value = true;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};

  try {
    const { res, json } = await api.get(`/api/roles/${id.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const r = json?.result?.role;
    role.value = typeof r?.role === "string" ? r.role : "";
    name.value = typeof r?.name === "string" ? r.name : "";
    status.value = typeof r?.status === "number" ? r.status : 1;
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
      role: role.value,
      name: name.value,
      status: Number(status.value)
    };

    const { res, json } = isEdit.value
      ? await api.put(`/api/roles/${id.value}`, body, { auth: true })
      : await api.post("/api/roles", body, { auth: true });
    if (!res.ok) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    message.value = isEdit.value ? "Role updated." : "Role created.";
    messageTone.value = "success";
    setTimeout(() => router.push("/admin/roles"), 300);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadRole();
});
</script>
