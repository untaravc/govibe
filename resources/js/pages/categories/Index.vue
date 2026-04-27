<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Categories</h3>
          <p class="mt-1 text-sm text-slate-600">Manage post categories.</p>
        </div>

        <button
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          @click="router.push('/admin/categories/new')"
        >
          Add category
        </button>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="mt-4 grid grid-cols-1 gap-3 md:grid-cols-3">
        <div>
          <label class="text-sm font-medium text-slate-700">Section</label>
          <select
            v-model="filterSection"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            :disabled="sectionsLoading"
          >
            <option value="">All</option>
            <option v-for="item in sections" :key="item" :value="item">
              {{ item }}
            </option>
          </select>
          <p v-if="sectionsError" class="mt-2 text-sm text-danger">{{ sectionsError }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Name</label>
          <input
            v-model.trim="filterName"
            type="text"
            placeholder="Search name"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
          />
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Status</label>
          <select
            v-model="filterStatus"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
          >
            <option value="">All</option>
            <option value="1">Active</option>
            <option value="0">Inactive</option>
          </select>
        </div>
      </div>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[860px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">No</th>
              <th class="px-4 py-3">Section</th>
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Slug</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="(c, idx) in categories" :key="c.id">
              <td class="px-4 py-3 text-slate-700">{{ idx + 1 }}</td>
              <td class="px-4 py-3 text-slate-700">{{ c.section }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ c.name }}</td>
              <td class="px-4 py-3 text-slate-700">{{ c.slug }}</td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium"
                  :class="c.status ? 'bg-emerald-50 text-emerald-700' : 'bg-rose-50 text-rose-700'"
                >
                  {{ c.status ? "Active" : "Inactive" }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <div :ref="(el) => setActionRoot(c.id, el)" class="relative inline-block text-left">
                  <button
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                    :aria-expanded="actionsOpenFor === c.id ? 'true' : 'false'"
                    aria-haspopup="menu"
                    @click="toggleActions(c.id)"
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
                    v-if="actionsOpenFor === c.id"
                    role="menu"
                    class="absolute right-0 z-10 mt-2 w-40 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-lg"
                  >
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-slate-50"
                      @click="
                        closeActions();
                        router.push(`/admin/categories/${c.id}/edit`);
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
                        onDelete(c);
                      "
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </td>
            </tr>

            <tr v-if="categories.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="6">No categories found.</td>
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

const categories = ref([]);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const sections = ref([]);
const sectionsLoading = ref(false);
const sectionsError = ref("");
const filterSection = ref("");
const filterName = ref("");
const filterStatus = ref("");
const actionsOpenFor = ref(null);

const actionRoots = new Map();

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

async function loadCategories() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const params = new URLSearchParams();
    if (filterSection.value) params.set("section", filterSection.value);
    if (filterName.value) params.set("name", filterName.value);
    if (filterStatus.value !== "") params.set("status", filterStatus.value);

    const url = params.toString() ? `/api/categories?${params.toString()}` : "/api/categories";
    const { res, json } = await api.get(url, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      categories.value = [];
      return;
    }
    categories.value = Array.isArray(json?.result?.categories) ? json.result.categories : [];
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    categories.value = [];
  }
}

async function loadSections() {
  sectionsLoading.value = true;
  sectionsError.value = "";
  try {
    const { res, json } = await api.get("/api/sections", { auth: true });
    if (!res.ok || json?.success === false) {
      sections.value = [];
      sectionsError.value = apiErrorMessage(json, `Failed to load sections (${res.status})`);
      return;
    }
    sections.value = Array.isArray(json?.result?.sections) ? json.result.sections : [];
  } catch (err) {
    sections.value = [];
    sectionsError.value = String(err);
  } finally {
    sectionsLoading.value = false;
  }
}

let filterTimer = null;
watch([filterSection, filterName, filterStatus], () => {
  if (filterTimer) clearTimeout(filterTimer);
  filterTimer = setTimeout(() => {
    loadCategories();
    filterTimer = null;
  }, 250);
});

async function onDelete(c) {
  const ok = confirm(`Delete category "${c.name}"?`);
  if (!ok) return;

  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.delete(`/api/categories/${c.id}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Delete failed (${res.status})`);
      messageTone.value = "error";
      return;
    }
    message.value = "Category deleted.";
    messageTone.value = "success";
    await loadCategories();
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  }
}

onMounted(() => {
  window.addEventListener("click", onWindowClick, true);
  window.addEventListener("keydown", onWindowKeydown);
  loadSections();
  loadCategories();
});

onBeforeUnmount(() => {
  window.removeEventListener("click", onWindowClick, true);
  window.removeEventListener("keydown", onWindowKeydown);
  if (filterTimer) clearTimeout(filterTimer);
});
</script>
