<template>
  <header class="sticky top-0 z-20 border-b border-[#eadfcd] bg-[#fffaf1]/90 backdrop-blur">
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-3 px-4 py-3 sm:px-6 lg:px-8">
      <RouterLink to="/" class="flex min-w-0 items-center gap-3">
        <span class="flex h-10 w-10 items-center justify-center rounded-2xl bg-[#17211f] text-sm font-black text-[#f2c14e]">
          GV
        </span>
        <span class="min-w-0">
          <span class="block text-sm font-black tracking-tight text-[#17211f]">GoVibe Logistics</span>
          <span class="block text-xs font-semibold text-[#687b75]">Shipment provider</span>
        </span>
      </RouterLink>

      <nav class="flex items-center gap-2">
        <RouterLink
          to="/"
          class="hidden rounded-full px-4 py-2 text-sm font-bold text-[#344943] hover:bg-[#efe5d4] sm:inline-flex"
          active-class="bg-[#efe5d4] text-[#17211f]"
          exact-active-class="bg-[#efe5d4] text-[#17211f]"
        >
          Home
        </RouterLink>
        <RouterLink
          :to="{ path: '/', hash: '#track' }"
          class="hidden rounded-full px-4 py-2 text-sm font-bold text-[#344943] hover:bg-[#efe5d4] sm:inline-flex"
        >
          Track
        </RouterLink>
        <RouterLink
          v-if="!hasToken"
          to="/auth/register"
          class="hidden rounded-full px-4 py-2 text-sm font-bold text-[#344943] hover:bg-[#efe5d4] md:inline-flex"
          active-class="bg-[#efe5d4] text-[#17211f]"
        >
          Register
        </RouterLink>
        <RouterLink
          v-if="!hasToken"
          to="/auth/login"
          class="rounded-full bg-[#17211f] px-4 py-2 text-sm font-bold text-white hover:bg-[#263a35]"
        >
          Login
        </RouterLink>
        <RouterLink
          v-if="hasToken"
          to="/admin"
          class="rounded-full border border-[#d8cbb8] bg-white px-4 py-2 text-sm font-bold text-[#17211f] hover:bg-[#fff4de]"
        >
          Admin site
        </RouterLink>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

const token = ref("");

const hasToken = computed(() => token.value.trim().length > 0);

function syncTokenFromStorage() {
  try {
    token.value = localStorage.getItem("access_token") || "";
  } catch {
    token.value = "";
  }
}

function onStorage(event) {
  if (event.key === "access_token") syncTokenFromStorage();
}

onMounted(() => {
  syncTokenFromStorage();
  window.addEventListener("storage", onStorage);
});

onBeforeUnmount(() => {
  window.removeEventListener("storage", onStorage);
});
</script>
