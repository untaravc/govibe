<template>
  <footer class="border-t border-slate-200/60 py-6 text-center text-xs text-slate-500">
    © {{ year }} GoVibe
  </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>

