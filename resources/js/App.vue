<template>
  <PageLoader />
  <RouterView />
</template>

<script setup>
import { RouterView } from "vue-router";

defineProps({
  title: {
    type: String,
    required: true
  }
});
</script>
