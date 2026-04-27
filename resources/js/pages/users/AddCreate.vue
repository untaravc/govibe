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

        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">Role</label>
            <select
              v-model="roleId"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            >
              <option value="">No role</option>
              <option v-for="r in roles" :key="r.id" :value="String(r.id)">{{ r.name }}</option>
            </select>
            <p v-if="fieldErrors.role_id" class="mt-2 text-sm text-danger">{{ fieldErrors.role_id }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Phone</label>
            <input
              v-model="phone"
              type="tel"
              autocomplete="tel"
              placeholder="Phone number"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.phone" class="mt-2 text-sm text-danger">{{ fieldErrors.phone }}</p>
          </div>
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

    <div v-if="isEdit" class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">User offices</h3>
          <p class="mt-1 text-sm text-slate-600">Assign this user to offices.</p>
        </div>
      </div>

      <form class="mt-6 grid gap-4 md:grid-cols-12" @submit.prevent="addUserOffice">
        <div class="md:col-span-7">
          <label class="text-sm font-medium text-slate-700">Office</label>
          <select
            v-model="newOfficeId"
            class="mt-2 h-11 w-full rounded-xl border border-slate-200 bg-white px-4 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            required
          >
            <option value="" disabled>Select office...</option>
            <option v-for="o in offices" :key="o.id" :value="String(o.id)">{{ o.name }} ({{ o.code }})</option>
          </select>
        </div>

        <div class="md:col-span-3">
          <label class="text-sm font-medium text-slate-700">Status</label>
          <select
            v-model="newOfficeStatus"
            class="mt-2 h-11 w-full rounded-xl border border-slate-200 bg-white px-4 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
          >
            <option :value="1">Active</option>
            <option :value="0">Inactive</option>
          </select>
        </div>

        <div class="flex items-end md:col-span-2">
          <button
            type="submit"
            class="h-11 w-full rounded-xl bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="userOfficesLoading"
          >
            Add
          </button>
        </div>
      </form>

      <p v-if="userOfficesMessage" class="mt-4 text-sm" :class="userOfficesMessageToneClass">{{ userOfficesMessage }}</p>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <table class="w-full min-w-[520px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">Office</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="uo in userOffices" :key="uo.id">
              <td class="px-4 py-3 text-slate-900">
                <div class="font-medium">{{ uo.office?.name || `Office #${uo.office_id}` }}</div>
                <div v-if="uo.office?.code" class="text-xs text-slate-500">{{ uo.office.code }}</div>
              </td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2 py-1 text-xs font-semibold"
                  :class="uo.status === 1 ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-700'"
                >
                  {{ uo.status === 1 ? "Active" : "Inactive" }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <button
                  type="button"
                  class="rounded-lg px-3 py-2 text-sm font-medium text-danger hover:bg-rose-50 disabled:cursor-not-allowed disabled:opacity-60"
                  :disabled="userOfficesLoading"
                  @click="removeUserOffice(uo)"
                >
                  Remove
                </button>
              </td>
            </tr>

            <tr v-if="userOffices.length === 0">
              <td class="px-4 py-8 text-center text-slate-500" colspan="3">No offices assigned.</td>
            </tr>
          </tbody>
        </table>
      </div>
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
const phone = ref("");
const roles = ref([]);
const roleId = ref("");
const offices = ref([]);
const userOffices = ref([]);
const newOfficeId = ref("");
const newOfficeStatus = ref(1);
const userOfficesLoading = ref(false);
const userOfficesMessage = ref("");
const userOfficesMessageTone = ref("neutral"); // neutral | success | error
const loading = ref(false);

const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

const userOfficesMessageToneClass = computed(() => {
  if (userOfficesMessageTone.value === "success") return "text-success";
  if (userOfficesMessageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

async function loadRoles() {
  try {
    const { res, json } = await api.get("/api/roles", { auth: true });
    if (!res.ok) {
      roles.value = [];
      return;
    }
    roles.value = Array.isArray(json?.result?.roles) ? json.result.roles : [];
  } catch {
    roles.value = [];
  }
}

async function loadOffices() {
  try {
    const { res, json } = await api.get("/api/offices", { auth: true });
    if (!res.ok) {
      offices.value = [];
      return;
    }
    offices.value = Array.isArray(json?.result?.offices) ? json.result.offices : [];
  } catch {
    offices.value = [];
  }
}

async function loadUserOffices() {
  if (!isEdit.value) return;

  userOfficesMessage.value = "";
  userOfficesMessageTone.value = "neutral";
  userOfficesLoading.value = true;
  try {
    const { res, json } = await api.get(`/api/users/${id.value}/offices`, { auth: true });
    if (!res.ok) {
      userOfficesMessage.value = apiErrorMessage(json, `Request failed (${res.status})`);
      userOfficesMessageTone.value = "error";
      userOffices.value = [];
      return;
    }
    userOffices.value = Array.isArray(json?.result?.user_offices) ? json.result.user_offices : [];
  } catch (err) {
    userOfficesMessage.value = String(err);
    userOfficesMessageTone.value = "error";
    userOffices.value = [];
  } finally {
    userOfficesLoading.value = false;
  }
}

async function addUserOffice() {
  if (!isEdit.value) return;
  const officeIdNum = Number.parseInt(String(newOfficeId.value), 10);
  if (!Number.isFinite(officeIdNum) || officeIdNum <= 0) return;

  userOfficesMessage.value = "";
  userOfficesMessageTone.value = "neutral";
  userOfficesLoading.value = true;
  try {
    const body = {
      office_id: officeIdNum,
      status: newOfficeStatus.value
    };
    const { res, json } = await api.post(`/api/users/${id.value}/offices`, body, { auth: true });
    if (!res.ok) {
      userOfficesMessage.value = apiErrorMessage(json, `Request failed (${res.status})`);
      userOfficesMessageTone.value = "error";
      return;
    }

    userOfficesMessage.value = "Office assigned.";
    userOfficesMessageTone.value = "success";
    newOfficeId.value = "";
    newOfficeStatus.value = 1;
    await loadUserOffices();
  } catch (err) {
    userOfficesMessage.value = String(err);
    userOfficesMessageTone.value = "error";
  } finally {
    userOfficesLoading.value = false;
  }
}

async function removeUserOffice(uo) {
  if (!isEdit.value) return;
  const ok = confirm(`Remove office assignment${uo?.office?.name ? ` "${uo.office.name}"` : ""}?`);
  if (!ok) return;

  userOfficesMessage.value = "";
  userOfficesMessageTone.value = "neutral";
  userOfficesLoading.value = true;
  try {
    const { res, json } = await api.delete(`/api/users/${id.value}/offices/${uo.office_id}`, { auth: true });
    if (!res.ok) {
      userOfficesMessage.value = apiErrorMessage(json, `Request failed (${res.status})`);
      userOfficesMessageTone.value = "error";
      return;
    }
    userOfficesMessage.value = "Office removed.";
    userOfficesMessageTone.value = "success";
    await loadUserOffices();
  } catch (err) {
    userOfficesMessage.value = String(err);
    userOfficesMessageTone.value = "error";
  } finally {
    userOfficesLoading.value = false;
  }
}

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
    phone.value = typeof u?.phone === "string" ? u.phone : "";
    roleId.value = u?.role_id != null ? String(u.role_id) : "";
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
      email: email.value,
      phone: phone.value
    };

    if (roleId.value) {
      const parsed = Number.parseInt(String(roleId.value), 10);
      if (Number.isFinite(parsed) && parsed > 0) body.role_id = parsed;
    }

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
  loadRoles();
  loadUser();
  loadOffices();
  loadUserOffices();
});
</script>
