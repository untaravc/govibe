<template>
  <div class="mx-auto max-w-xl text-center">
    <p class="text-sm font-medium uppercase tracking-wider text-slate-500">401</p>
    <h2 class="mt-2 text-3xl font-semibold tracking-tight text-slate-900">You’re not signed in</h2>
    <p class="mt-3 text-slate-600">Please sign in to continue.</p>

    <div class="mt-6 flex flex-wrap items-center justify-center gap-2">
      <RouterLink
        to="/auth/login"
        class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        Go to login
      </RouterLink>
      <RouterLink
        to="/"
        class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
      >
        Back to home
      </RouterLink>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";
</script>

