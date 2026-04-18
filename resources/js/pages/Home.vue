<template>
  <div class="space-y-4">
    <div class="rounded-xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <h2 class="text-xl font-semibold tracking-tight text-slate-900 dark:text-slate-50">Dashboard</h2>
      <p class="mt-1 text-slate-600 dark:text-slate-300">Fiber backend + Vue + Tailwind CSS (compiled by Vite).</p>
    </div>

    <div class="rounded-xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <button
        type="button"
        class="inline-flex items-center rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-700 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
        @click="callHealth"
      >
        Call /api/health
      </button>
      <pre class="mt-3 overflow-auto rounded-lg bg-slate-950 p-4 text-xs text-slate-100">{{ output }}</pre>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const output = ref("");

async function callHealth() {
  output.value = "Loading...";
  try {
    const res = await fetch("/api/health");
    const json = await res.json();
    output.value = JSON.stringify(json, null, 2);
  } catch (err) {
    output.value = String(err);
  }
}
</script>

