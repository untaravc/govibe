<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Posts</h3>
          <p class="mt-1 text-sm text-slate-600">Manage published content.</p>
        </div>

        <button
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          @click="router.push('/admin/posts/new')"
        >
          Add post
        </button>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="mt-6 overflow-auto rounded-xl border border-slate-200">
        <table class="w-full min-w-[860px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">ID</th>
              <th class="px-4 py-3">Title</th>
              <th class="px-4 py-3">Slug</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3">Release</th>
              <th class="px-4 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="p in posts" :key="p.id">
              <td class="px-4 py-3 text-slate-700">{{ p.id }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ p.title }}</td>
              <td class="px-4 py-3 text-slate-700">{{ p.slug }}</td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium"
                  :class="p.status ? 'bg-emerald-50 text-emerald-700' : 'bg-rose-50 text-rose-700'"
                >
                  {{ p.status ? "Active" : "Inactive" }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-700">{{ formatDate(p.release_at) }}</td>
              <td class="px-4 py-3 text-right">
                <button
                  type="button"
                  class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                  @click="router.push(`/admin/posts/${p.id}/edit`)"
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="ml-2 rounded-lg px-3 py-2 text-sm font-medium text-danger hover:bg-rose-50"
                  @click="onDelete(p)"
                >
                  Delete
                </button>
              </td>
            </tr>

            <tr v-if="posts.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="6">No posts found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const router = useRouter();

const posts = ref([]);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error

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

async function loadPosts() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.get("/api/posts", { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      posts.value = [];
      return;
    }
    posts.value = Array.isArray(json?.result?.posts) ? json.result.posts : [];
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    posts.value = [];
  }
}

async function onDelete(p) {
  const ok = confirm(`Delete post "${p.title}"?`);
  if (!ok) return;

  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.delete(`/api/posts/${p.id}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Delete failed (${res.status})`);
      messageTone.value = "error";
      return;
    }
    message.value = "Post deleted.";
    messageTone.value = "success";
    await loadPosts();
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  }
}

onMounted(() => {
  loadPosts();
});
</script>

