<template>
  <div class="flex min-h-screen flex-col bg-slate-50 text-slate-900">
    <PublicTopbar />

    <main class="w-full flex-1 px-4 py-6">
      <RouterView />
    </main>

    <PublicFooter />
  </div>
</template>

<script setup>
import { RouterView } from "vue-router";

import PublicFooter from "../components/PublicFooter.vue";
import PublicTopbar from "../components/PublicTopbar.vue";
</script>

