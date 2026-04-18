<template>
  <div class="flex flex-wrap items-center justify-between gap-3">
    <p class="text-sm text-slate-600">
      Page <span class="font-medium text-slate-900">{{ meta.page }}</span> of
      <span class="font-medium text-slate-900">{{ meta.total_pages }}</span>
      · Total <span class="font-medium text-slate-900">{{ meta.total }}</span>
    </p>

    <div class="flex items-center gap-2">
      <button
        type="button"
        class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="!meta.has_prev"
        @click="$emit('prev')"
      >
        Prev
      </button>
      <button
        type="button"
        class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="!meta.has_next"
        @click="$emit('next')"
      >
        Next
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  meta: {
    type: Object,
    required: true
  }
});

defineEmits(["prev", "next"]);
</script>

