<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight">Roles</h3>
          <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Dummy roles list.</p>
        </div>
        <button
          type="button"
          class="rounded-xl bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
        >
          Add role
        </button>
      </div>

      <div class="mt-6 grid grid-cols-1 gap-3 sm:grid-cols-2">
        <div
          v-for="r in roles"
          :key="r"
          class="rounded-2xl border border-slate-200/60 bg-white p-5 shadow-sm dark:border-white/10 dark:bg-slate-950"
        >
          <p class="text-sm font-semibold text-slate-900 dark:text-slate-50">{{ r }}</p>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">Permissions: —</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const roles = ["Administrator", "Manager", "Staff"];
</script>

