<template>
  <footer class="border-t border-slate-200/60 py-6 text-center text-xs text-slate-500 dark:border-white/10 dark:text-slate-400">
    © {{ year }} GoVibe · Administrator Dashboard
  </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>

