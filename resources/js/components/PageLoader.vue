<template>
  <div
    v-if="ui.isApiLoading"
    class="fixed inset-0 z-50 flex items-center justify-center bg-white/70 backdrop-blur"
    aria-live="polite"
    aria-busy="true"
  >
    <div class="flex items-center gap-3 rounded-2xl border border-slate-200 bg-white px-4 py-3 shadow-sm">
      <div class="h-5 w-5 animate-spin rounded-full border-2 border-slate-200 border-t-slate-900"></div>
      <p class="text-sm font-medium text-slate-700">Loading…</p>
    </div>
  </div>
</template>

<script setup>
import { useUiStore } from "../store/ui.js";

const ui = useUiStore();
</script>

