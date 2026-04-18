<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight">Users</h3>
          <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Dummy table layout.</p>
        </div>
        <button
          type="button"
          class="rounded-xl bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
        >
          Add user
        </button>
      </div>

      <div class="mt-6 overflow-auto rounded-xl border border-slate-200 dark:border-white/10">
        <table class="w-full min-w-[520px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500 dark:bg-white/5 dark:text-slate-400">
            <tr>
              <th class="px-4 py-3">Name</th>
              <th class="px-4 py-3">Email</th>
              <th class="px-4 py-3">Role</th>
              <th class="px-4 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200 dark:divide-white/10">
            <tr v-for="u in rows" :key="u.email">
              <td class="px-4 py-3 font-medium text-slate-900 dark:text-slate-50">{{ u.name }}</td>
              <td class="px-4 py-3 text-slate-600 dark:text-slate-300">{{ u.email }}</td>
              <td class="px-4 py-3">
                <span class="rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-700 dark:bg-white/10 dark:text-slate-200">
                  {{ u.role }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <button class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-white/10">
                  Edit
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
const rows = [
  { name: "Admin", email: "admin@govibe.local", role: "Administrator" },
  { name: "Staff", email: "staff@govibe.local", role: "Staff" }
];
</script>

