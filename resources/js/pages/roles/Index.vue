<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Roles</h3>
          <p class="mt-1 text-sm text-slate-600">Manage roles and access labels.</p>
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
            @click="router.push('/admin/roles/new')"
          >
            Add role
          </button>
        </div>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[720px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">No</th>
              <th class="px-4 py-3">Role</th>
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="(r, idx) in roles" :key="r.id">
              <td class="px-4 py-3 text-slate-700">{{ idx + 1 }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ r.role }}</td>
              <td class="px-4 py-3 text-slate-700">{{ r.name }}</td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium"
                  :class="r.status ? 'bg-emerald-50 text-emerald-700' : 'bg-rose-50 text-rose-700'"
                >
                  {{ r.status ? "Active" : "Inactive" }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <div :ref="(el) => setActionRoot(r.id, el)" class="relative inline-block text-left">
                  <button
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                    :aria-expanded="actionsOpenFor === r.id ? 'true' : 'false'"
                    aria-haspopup="menu"
                    @click="toggleActions(r.id)"
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
                    v-if="actionsOpenFor === r.id"
                    role="menu"
                    class="absolute right-0 z-10 mt-2 w-40 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-lg"
                  >
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-slate-50"
                      @click="
                        closeActions();
                        router.push(`/admin/roles/${r.id}/edit`);
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
                        onDelete(r);
                      "
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </td>
            </tr>

            <tr v-if="roles.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="5">No roles found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const router = useRouter();

const roles = ref([]);
const nameFilter = ref("");
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const actionsOpenFor = ref(null);

const actionRoots = new Map();
let nameFilterTimer = null;

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

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

async function loadRoles() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const params = new URLSearchParams();
    if (nameFilter.value.trim()) params.set("name", nameFilter.value.trim());
    const suffix = params.toString();

    const { res, json } = await api.get(`/api/roles${suffix ? `?${suffix}` : ""}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      roles.value = [];
      return;
    }
    roles.value = Array.isArray(json?.result?.roles) ? json.result.roles : [];
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    roles.value = [];
  }
}

async function onDelete(r) {
  const ok = confirm(`Delete role "${r.role}"?`);
  if (!ok) return;

  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.delete(`/api/roles/${r.id}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Delete failed (${res.status})`);
      messageTone.value = "error";
      return;
    }
    message.value = "Role deleted.";
    messageTone.value = "success";
    await loadRoles();
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  }
}

onMounted(() => {
  window.addEventListener("click", onWindowClick, true);
  window.addEventListener("keydown", onWindowKeydown);
  loadRoles();
});

onBeforeUnmount(() => {
  window.removeEventListener("click", onWindowClick, true);
  window.removeEventListener("keydown", onWindowKeydown);
  if (nameFilterTimer) clearTimeout(nameFilterTimer);
});

watch(nameFilter, () => {
  if (nameFilterTimer) clearTimeout(nameFilterTimer);
  nameFilterTimer = setTimeout(() => {
    loadRoles();
  }, 350);
});
</script>
