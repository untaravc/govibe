<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Users</h3>
          <p class="mt-1 text-sm text-slate-600">Manage application users.</p>
        </div>

        <button
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          @click="router.push('/admin/users/new')"
        >
          Add user
        </button>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[680px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">ID</th>
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Email</th>
              <th class="px-4 py-3">Created</th>
              <th class="px-4 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="u in users" :key="u.id">
              <td class="px-4 py-3 text-slate-700">{{ u.id }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ u.name }}</td>
              <td class="px-4 py-3 text-slate-700">{{ u.email }}</td>
              <td class="px-4 py-3 text-slate-700">{{ formatDate(u.created_at) }}</td>
              <td class="px-4 py-3 text-right">
                <button
                  type="button"
                  class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                  @click="openDetail(u)"
                >
                  Detail
                </button>
                <button
                  type="button"
                  class="ml-2 rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                  @click="router.push(`/admin/users/${u.id}/edit`)"
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="ml-2 rounded-lg px-3 py-2 text-sm font-medium text-danger hover:bg-rose-50"
                  @click="onDelete(u)"
                >
                  Delete
                </button>
              </td>
            </tr>

            <tr v-if="users.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="5">No users found.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="mt-4">
        <PaginationNav :meta="meta" @prev="goPrev" @next="goNext" />
      </div>
    </div>
  </div>

  <Modal :open="detailOpen" title="User detail" @close="detailOpen = false">
    <dl v-if="selectedUser" class="space-y-3">
      <div class="grid grid-cols-3 gap-3">
        <dt class="text-sm font-medium text-slate-600">ID</dt>
        <dd class="col-span-2 text-sm text-slate-900">{{ selectedUser.id }}</dd>
      </div>
      <div class="grid grid-cols-3 gap-3">
        <dt class="text-sm font-medium text-slate-600">Name</dt>
        <dd class="col-span-2 text-sm text-slate-900">{{ selectedUser.name }}</dd>
      </div>
      <div class="grid grid-cols-3 gap-3">
        <dt class="text-sm font-medium text-slate-600">Email</dt>
        <dd class="col-span-2 text-sm text-slate-900">{{ selectedUser.email }}</dd>
      </div>
      <div class="grid grid-cols-3 gap-3">
        <dt class="text-sm font-medium text-slate-600">Created</dt>
        <dd class="col-span-2 text-sm text-slate-900">{{ formatDate(selectedUser.created_at) }}</dd>
      </div>
      <div class="grid grid-cols-3 gap-3">
        <dt class="text-sm font-medium text-slate-600">Updated</dt>
        <dd class="col-span-2 text-sm text-slate-900">{{ formatDate(selectedUser.updated_at) }}</dd>
      </div>
    </dl>

    <template #footer>
      <div class="flex items-center justify-end gap-2">
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="detailOpen = false"
        >
          Close
        </button>
        <button
          v-if="selectedUser"
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          @click="
            detailOpen = false;
            router.push(`/admin/users/${selectedUser.id}/edit`);
          "
        >
          Edit
        </button>
      </div>
    </template>
  </Modal>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";

import Modal from "../../components/Modal.vue";
import PaginationNav from "../../components/PaginationNav.vue";
import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const router = useRouter();

const users = ref([]);
const page = ref(1);
const perPage = ref(10);
const meta = ref({
  page: 1,
  per_page: 10,
  total: 0,
  total_pages: 0,
  has_prev: false,
  has_next: false
});
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const detailOpen = ref(false);
const selectedUser = ref(null);

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

function formatDate(value) {
  if (!value) return "—";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);
  return d.toLocaleString();
}

async function loadUsers() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.get(`/api/users?per_page=${perPage.value}&page=${page.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      users.value = [];
      meta.value = { ...meta.value, total: 0, total_pages: 0, has_prev: false, has_next: false };
      return;
    }
    users.value = Array.isArray(json?.result?.users) ? json.result.users : [];
    meta.value = typeof json?.result?.meta === "object" && json.result.meta ? { ...meta.value, ...json.result.meta } : meta.value;
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    users.value = [];
    meta.value = { ...meta.value, total: 0, total_pages: 0, has_prev: false, has_next: false };
  }
}

async function onDelete(u) {
  const ok = confirm(`Delete user "${u.name}"?`);
  if (!ok) return;

  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.delete(`/api/users/${u.id}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Delete failed (${res.status})`);
      messageTone.value = "error";
      return;
    }
    message.value = "User deleted.";
    messageTone.value = "success";
    await loadUsers();
    if (users.value.length === 0 && page.value > 1 && meta.value.total_pages > 0 && page.value > meta.value.total_pages) {
      page.value = meta.value.total_pages;
      await loadUsers();
    }
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  }
}

function openDetail(u) {
  selectedUser.value = u;
  detailOpen.value = true;
}

function goPrev() {
  if (!meta.value.has_prev) return;
  page.value = Math.max(1, page.value - 1);
  loadUsers();
}

function goNext() {
  if (!meta.value.has_next) return;
  page.value = page.value + 1;
  loadUsers();
}

onMounted(() => {
  loadUsers();
});
</script>
