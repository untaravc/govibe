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
              <th class="px-4 py-3">ID</th>
              <th class="px-4 py-3">Section</th>
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Slug</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3">Created</th>
              <th class="px-4 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="c in categories" :key="c.id">
              <td class="px-4 py-3 text-slate-700">{{ c.id }}</td>
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
              <td class="px-4 py-3 text-slate-700">{{ formatDate(c.created_at) }}</td>
              <td class="px-4 py-3 text-right">
                <button
                  type="button"
                  class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                  @click="router.push(`/admin/categories/${c.id}/edit`)"
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="ml-2 rounded-lg px-3 py-2 text-sm font-medium text-danger hover:bg-rose-50"
                  @click="onDelete(c)"
                >
                  Delete
                </button>
              </td>
            </tr>

            <tr v-if="categories.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="7">No categories found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
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
  loadSections();
  loadCategories();
});
</script>
