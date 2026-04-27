<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Users</h3>
          <p class="mt-1 text-sm text-slate-600">Manage application users.</p>
        </div>

        <div class="flex flex-wrap items-center justify-end gap-2">
          <div class="flex items-center gap-2">
            <input
              v-model="nameFilter"
              type="text"
              class="h-10 w-56 rounded-xl border border-slate-200 bg-white px-3 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-4 focus:ring-slate-900/10"
              placeholder="Filter by name..."
              autocomplete="off"
            />
            <button
              v-if="nameFilter.trim()"
              type="button"
              class="h-10 rounded-xl border border-slate-200 bg-white px-3 text-sm font-medium text-slate-700 hover:bg-slate-50"
              @click="nameFilter = ''"
            >
              Clear
            </button>
          </div>

          <button
            type="button"
            class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            @click="router.push('/admin/users/new')"
          >
            Add user
          </button>
        </div>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[680px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">No</th>
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Email</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="(u, idx) in users" :key="u.id">
              <td class="px-4 py-3 text-slate-700">{{ rowNumber(idx) }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ u.name }}</td>
              <td class="px-4 py-3 text-slate-700">{{ u.email }}</td>
              <td class="px-4 py-3 text-right">
                <div :ref="(el) => setActionRoot(u.id, el)" class="relative inline-block text-left">
                  <button
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                    :aria-expanded="actionsOpenFor === u.id ? 'true' : 'false'"
                    aria-haspopup="menu"
                    @click="toggleActions(u.id)"
                  >
                    Actions
                    <svg viewBox="0 0 24 24" fill="none" class="h-4 w-4" xmlns="http://www.w3.org/2000/svg">
                      <path
                        d="M6 9l6 6 6-6"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </button>

                  <div
                    v-if="actionsOpenFor === u.id"
                    role="menu"
                    class="absolute right-0 z-10 mt-2 w-44 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-lg"
                  >
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-slate-50"
                      @click="
                        closeActions();
                        openDetail(u);
                      "
                    >
                      Detail
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-slate-50"
                      @click="
                        closeActions();
                        router.push(`/admin/users/${u.id}/edit`);
                      "
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-rose-600 hover:bg-rose-50"
                      @click="
                        closeActions();
                        onDelete(u);
                      "
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </td>
            </tr>

            <tr v-if="users.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="4">No users found.</td>
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
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";

import Modal from "../../components/Modal.vue";
import PaginationNav from "../../components/PaginationNav.vue";
import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const router = useRouter();

const users = ref([]);
const page = ref(1);
const perPage = ref(10);
const nameFilter = ref("");
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
const actionsOpenFor = ref(null);

const actionRoots = new Map();
let nameFilterTimer = null;

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

const rowNumber = (idx) => (page.value - 1) * perPage.value + idx + 1;

function formatDate(value) {
  if (!value) return "—";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);
  return d.toLocaleString();
}

function setActionRoot(id, el) {
  if (!id) return;
  if (!el) {
    actionRoots.delete(id);
    return;
  }
  actionRoots.set(id, el);
}

function closeActions() {
  actionsOpenFor.value = null;
}

function toggleActions(id) {
  actionsOpenFor.value = actionsOpenFor.value === id ? null : id;
}

function onWindowClick(e) {
  if (!actionsOpenFor.value) return;
  const root = actionRoots.get(actionsOpenFor.value);
  if (!root) return closeActions();
  if (root === e.target || root.contains(e.target)) return;
  closeActions();
}

function onWindowKeydown(e) {
  if (!actionsOpenFor.value) return;
  if (e.key === "Escape") closeActions();
}

async function loadUsers() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const params = new URLSearchParams();
    params.set("per_page", String(perPage.value));
    params.set("page", String(page.value));
    if (nameFilter.value.trim()) params.set("name", nameFilter.value.trim());

    const { res, json } = await api.get(`/api/users?${params.toString()}`, { auth: true });
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
  window.addEventListener("click", onWindowClick, true);
  window.addEventListener("keydown", onWindowKeydown);
  loadUsers();
});

onBeforeUnmount(() => {
  window.removeEventListener("click", onWindowClick, true);
  window.removeEventListener("keydown", onWindowKeydown);
  if (nameFilterTimer) clearTimeout(nameFilterTimer);
});

watch(nameFilter, () => {
  if (nameFilterTimer) clearTimeout(nameFilterTimer);
  nameFilterTimer = setTimeout(() => {
    page.value = 1;
    loadUsers();
  }, 350);
});
</script>
