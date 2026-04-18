<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-slate-900">
      <h3 class="text-lg font-semibold tracking-tight">Settings</h3>
      <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Dummy settings form</p>

      <form class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div class="sm:col-span-2">
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Site name</label>
          <input
            type="text"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            placeholder="GoVibe"
          />
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Timezone</label>
          <select
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50"
          >
            <option>Asia/Jakarta</option>
            <option>UTC</option>
            <option>America/Los_Angeles</option>
          </select>
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700 dark:text-slate-200">Theme</label>
          <select
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50"
          >
            <option>System</option>
            <option>Light</option>
            <option>Dark</option>
          </select>
        </div>

        <div class="sm:col-span-2">
          <button
            type="button"
            class="rounded-xl bg-slate-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-slate-800 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
          >
            Save changes
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

