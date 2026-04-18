<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-[60]">
      <div class="absolute inset-0 bg-slate-900/40" @click="onBackdrop" />

      <div class="absolute inset-0 overflow-y-auto p-4">
        <div class="mx-auto w-full max-w-lg">
          <div
            class="rounded-2xl border border-slate-200 bg-white shadow-xl"
            role="dialog"
            aria-modal="true"
            :aria-labelledby="titleId"
          >
            <div class="flex items-start justify-between gap-3 border-b border-slate-200 px-5 py-4">
              <div class="min-w-0">
                <h3 :id="titleId" class="truncate text-base font-semibold tracking-tight text-slate-900">
                  {{ title }}
                </h3>
              </div>
              <button
                type="button"
                class="inline-flex h-9 w-9 items-center justify-center rounded-xl text-slate-600 hover:bg-slate-100 hover:text-slate-900"
                @click="$emit('close')"
              >
                <span class="sr-only">Close</span>
                <svg viewBox="0 0 24 24" fill="none" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg">
                  <path d="M6 6l12 12M18 6L6 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
                </svg>
              </button>
            </div>

            <div class="px-5 py-4">
              <slot />
            </div>

            <div v-if="$slots.footer" class="border-t border-slate-200 px-5 py-4">
              <slot name="footer" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { onBeforeUnmount, onMounted } from "vue";

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    required: true
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(["close"]);

const titleId = `modal-title-${Math.random().toString(36).slice(2)}`;

function onBackdrop() {
  if (!props.closeOnBackdrop) return;
  emit("close");
}

function onKeydown(event) {
  if (!props.open) return;
  if (event.key === "Escape") emit("close");
}

onMounted(() => {
  window.addEventListener("keydown", onKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", onKeydown);
});
</script>

