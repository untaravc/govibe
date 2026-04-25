<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">
            {{ isEdit ? "Edit category" : "Add category" }}
          </h3>
          <p class="mt-1 text-sm text-slate-600">
            {{ isEdit ? "Update category details." : "Create a new category." }}
          </p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="router.push('/admin/categories')"
        >
          Back
        </button>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">Section</label>
            <select
              v-model="section"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              :disabled="sectionsLoading"
              required
            >
              <option value="" disabled>Select section</option>
              <option v-for="item in sections" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
            <p v-if="sectionsLoading" class="mt-2 text-sm text-slate-500">Loading sections...</p>
            <p v-if="sectionsError" class="mt-2 text-sm text-danger">{{ sectionsError }}</p>
            <p v-if="fieldErrors.section" class="mt-2 text-sm text-danger">{{ fieldErrors.section }}</p>
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

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Name</label>
            <input
              v-model.trim="name"
              type="text"
              placeholder="Category name"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.name" class="mt-2 text-sm text-danger">{{ fieldErrors.name }}</p>
          </div>

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Slug</label>
            <input
              v-model.trim="slug"
              type="text"
              placeholder="category-slug"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.slug" class="mt-2 text-sm text-danger">{{ fieldErrors.slug }}</p>
          </div>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : isEdit ? "Save changes" : "Create category" }}
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

const section = ref("");
const name = ref("");
const slug = ref("");
const status = ref(1);
const sections = ref([]);
const sectionsLoading = ref(false);
const sectionsError = ref("");

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

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
    if (!section.value && sections.value.length > 0) {
      section.value = sections.value[0];
    }
  } catch (err) {
    sections.value = [];
    sectionsError.value = String(err);
  } finally {
    sectionsLoading.value = false;
  }
}

async function loadCategory() {
  if (!isEdit.value) return;

  loading.value = true;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};

  try {
    const { res, json } = await api.get(`/api/categories/${id.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const c = json?.result?.category;
    section.value = typeof c?.section === "string" ? c.section : "";
    name.value = typeof c?.name === "string" ? c.name : "";
    slug.value = typeof c?.slug === "string" ? c.slug : "";
    status.value = typeof c?.status === "number" ? c.status : 1;
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
      section: section.value,
      name: name.value,
      slug: slug.value,
      status: Number(status.value)
    };

    const { res, json } = isEdit.value
      ? await api.put(`/api/categories/${id.value}`, body, { auth: true })
      : await api.post("/api/categories", body, { auth: true });
    if (!res.ok) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    Swal.fire({
      toast: true,
      position: "top-end",
      icon: "success",
      title: isEdit.value ? "Category updated" : "Category created",
      showConfirmButton: false,
      timer: 1600,
      timerProgressBar: true
    });

    setTimeout(() => router.push("/admin/categories"), 300);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadSections();
  loadCategory();
});
</script>
