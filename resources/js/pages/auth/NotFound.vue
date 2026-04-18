<template>
  <div class="mx-auto max-w-xl text-center">
    <p class="text-sm font-medium uppercase tracking-wider text-slate-500">404</p>
    <h2 class="mt-2 text-3xl font-semibold tracking-tight text-slate-900">Page not found</h2>
    <p class="mt-3 text-slate-600">The page you’re looking for doesn’t exist or was moved.</p>

    <div class="mt-6 flex flex-wrap items-center justify-center gap-2">
      <RouterLink
        to="/"
        class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        Go home
      </RouterLink>
      <RouterLink
        to="/auth/login"
        class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
      >
        Login
      </RouterLink>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";
</script>

