<template>
  <footer class="border-t border-[#293a35] bg-[#17211f] px-4 py-10 text-white sm:px-6 lg:px-8">
    <div class="mx-auto grid max-w-7xl gap-8 md:grid-cols-[1.2fr_0.8fr_0.8fr]">
      <div>
        <p class="text-lg font-black tracking-tight">GoVibe Logistics</p>
        <p class="mt-3 max-w-md text-sm leading-6 text-white/65">
          Shipment provider for air way bill tracking, transit visibility, and destination status updates.
        </p>
      </div>
      <div>
        <p class="text-xs font-black uppercase tracking-[0.24em] text-[#f2c14e]">Service</p>
        <div class="mt-4 space-y-2 text-sm text-white/65">
          <p>Door to door delivery</p>
          <p>Transit monitoring</p>
          <p>Shipment tracking</p>
        </div>
      </div>
      <div>
        <p class="text-xs font-black uppercase tracking-[0.24em] text-[#f2c14e]">Contact</p>
        <div class="mt-4 space-y-2 text-sm text-white/65">
          <p>support@govibe.local</p>
          <p>Jakarta, Indonesia</p>
        </div>
      </div>
    </div>
    <div class="mx-auto mt-8 flex max-w-7xl flex-wrap items-center justify-between gap-3 border-t border-white/10 pt-5 text-xs font-semibold text-white/45">
      <p>© {{ year }} GoVibe Logistics</p>
      <p>Built for reliable shipment operations.</p>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>
