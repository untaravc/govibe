<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">{{ isEdit ? "Edit post" : "Add post" }}</h3>
          <p class="mt-1 text-sm text-slate-600">
            {{ isEdit ? "Update post content." : "Create a new post." }}
          </p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="router.push('/admin/posts')"
        >
          Back
        </button>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="text-sm font-medium text-slate-700">Title</label>
          <input
            v-model.trim="title"
            type="text"
            autocomplete="off"
            placeholder="Post title"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            required
          />
          <p v-if="fieldErrors.title" class="mt-2 text-sm text-danger">{{ fieldErrors.title }}</p>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Subtitle</label>
          <input
            v-model.trim="subtitle"
            type="text"
            autocomplete="off"
            placeholder="Optional subtitle"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
          />
          <p v-if="fieldErrors.subtitle" class="mt-2 text-sm text-danger">{{ fieldErrors.subtitle }}</p>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">Slug</label>
            <input
              v-model.trim="slug"
              type="text"
              autocomplete="off"
              placeholder="post-slug"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.slug" class="mt-2 text-sm text-danger">{{ fieldErrors.slug }}</p>
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
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">Image URL</label>
            <input
              v-model.trim="imageUrl"
              type="url"
              autocomplete="off"
              placeholder="https://..."
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.image_url" class="mt-2 text-sm text-danger">{{ fieldErrors.image_url }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Release At</label>
            <input
              v-model="releaseAt"
              type="datetime-local"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.release_at" class="mt-2 text-sm text-danger">{{ fieldErrors.release_at }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">User ID</label>
            <input
              v-model.number="userId"
              type="number"
              min="1"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.user_id" class="mt-2 text-sm text-danger">{{ fieldErrors.user_id }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Category ID</label>
            <input
              v-model.number="categoryId"
              type="number"
              min="1"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.category_id" class="mt-2 text-sm text-danger">{{ fieldErrors.category_id }}</p>
          </div>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Content</label>
          <div class="mt-2">
            <RichTextEditor v-model="content" />
          </div>
          <p v-if="fieldErrors.content" class="mt-2 text-sm text-danger">{{ fieldErrors.content }}</p>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : isEdit ? "Save changes" : "Create post" }}
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
import RichTextEditor from "../../components/RichTextEditor.vue";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const router = useRouter();
const route = useRoute();

const id = computed(() => String(route.params.id || "").trim());
const isEdit = computed(() => id.value.length > 0);

const title = ref("");
const subtitle = ref("");
const slug = ref("");
const content = ref("");
const status = ref(1);
const imageUrl = ref("");
const userId = ref(1);
const categoryId = ref(null);
const releaseAt = ref("");

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

function toDatetimeLocal(iso) {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function toISOFromDatetimeLocal(value) {
  if (!value) return null;
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return null;
  return d.toISOString();
}

async function loadPost() {
  if (!isEdit.value) return;

  loading.value = true;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};

  try {
    const { res, json } = await api.get(`/api/posts/${id.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const p = json?.result?.post;
    title.value = typeof p?.title === "string" ? p.title : "";
    subtitle.value = typeof p?.subtitle === "string" ? p.subtitle : "";
    slug.value = typeof p?.slug === "string" ? p.slug : "";
    content.value = typeof p?.content === "string" ? p.content : "";
    status.value = typeof p?.status === "number" ? p.status : 1;
    imageUrl.value = typeof p?.image_url === "string" ? p.image_url : "";
    userId.value = typeof p?.user_id === "number" ? p.user_id : 1;
    categoryId.value = typeof p?.category_id === "number" ? p.category_id : null;
    releaseAt.value = toDatetimeLocal(p?.release_at);
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
      title: title.value,
      subtitle: subtitle.value.trim() ? subtitle.value : null,
      slug: slug.value,
      content: content.value,
      status: Number(status.value),
      image_url: imageUrl.value.trim() ? imageUrl.value : null,
      user_id: Number(userId.value),
      release_at: toISOFromDatetimeLocal(releaseAt.value),
      category_id: categoryId.value ? Number(categoryId.value) : null
    };

    const { res, json } = isEdit.value
      ? await api.put(`/api/posts/${id.value}`, body, { auth: true })
      : await api.post("/api/posts", body, { auth: true });

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
      title: isEdit.value ? "Post updated" : "Post created",
      showConfirmButton: false,
      timer: 1600,
      timerProgressBar: true
    });

    setTimeout(() => router.push("/admin/posts"), 250);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadPost();
});
</script>

