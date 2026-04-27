<template>
  <div class="-mx-4 -my-6 overflow-hidden bg-[#f7f1e7] text-[#17211f]">
    <section class="relative isolate min-h-[720px] px-4 py-16 sm:px-6 lg:px-8">
      <div class="absolute inset-0 -z-10">
        <div class="absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,#ffe8a3_0,#f7c55c_28%,#e5794f_62%,#21342f_100%)]"></div>
        <div class="absolute left-[8%] top-[16%] h-48 w-48 rounded-full bg-white/30 blur-3xl"></div>
        <div class="absolute bottom-[10%] right-[10%] h-64 w-64 rounded-full bg-black/10 blur-3xl"></div>
      </div>

      <div class="mx-auto flex min-h-[620px] max-w-5xl items-center">
        <div>
          <p
            class="inline-flex rounded-full border border-[#17211f]/15 bg-white/55 px-4 py-2 text-xs font-bold uppercase tracking-[0.28em] text-[#234039] shadow-sm backdrop-blur"
          >
            Same-day shipment network
          </p>
          <h1
            class="mt-7 max-w-4xl text-5xl font-black leading-[0.95] tracking-[-0.06em] text-[#14201d] sm:text-7xl lg:text-8xl"
          >
            Shipment updates without the waiting room.
          </h1>
          <p class="mt-7 max-w-2xl text-lg leading-8 text-[#314641] sm:text-xl">
            Track air way bills, monitor transit movement, and keep customers informed from pickup to arrival.
          </p>

          <div class="mt-9 flex flex-wrap gap-3">
            <a
              href="#track"
              class="rounded-full bg-[#14201d] px-6 py-3 text-sm font-bold text-white shadow-xl shadow-[#14201d]/20 transition hover:-translate-y-0.5 hover:bg-[#263a35]"
            >
              Track shipment
            </a>
          </div>
        </div>
      </div>
    </section>

    <section id="track" class="relative bg-[#fffaf1] px-4 py-20 sm:px-6 lg:px-8">
      <div class="mx-auto grid max-w-6xl gap-8 lg:grid-cols-[0.8fr_1.2fr]">
        <div>
          <p class="text-sm font-black uppercase tracking-[0.26em] text-[#d0603d]">Track air way bill</p>
          <h2 class="mt-4 text-4xl font-black tracking-[-0.04em] text-[#14201d] sm:text-5xl">Enter your shipment code.</h2>
          <p class="mt-5 text-base leading-7 text-[#536760]">
            Use the air way bill or shipment code printed on your receipt to check the latest delivery status.
          </p>
        </div>

        <div class="rounded-[2rem] border border-[#eadfcd] bg-white p-5 shadow-2xl shadow-[#6d5630]/10 sm:p-7">
          <form class="grid gap-3 sm:grid-cols-[1fr_auto]" @submit.prevent="trackShipment">
            <label class="sr-only" for="shipment-code">Air way bill or shipment code</label>
            <input
              id="shipment-code"
              v-model.trim="trackingCode"
              type="text"
              class="h-14 rounded-2xl border border-[#d8cbb8] bg-[#fffaf1] px-5 text-base font-bold uppercase tracking-[0.08em] text-[#14201d] outline-none ring-[#f2c14e]/25 placeholder:normal-case placeholder:tracking-normal placeholder:text-[#8b9a95] focus:border-[#f2c14e] focus:ring-4"
              placeholder="Example: SHP-0001"
            />
            <button
              type="submit"
              class="h-14 rounded-2xl bg-[#d0603d] px-7 text-sm font-black uppercase tracking-[0.14em] text-white shadow-lg shadow-[#d0603d]/20 transition hover:-translate-y-0.5 hover:bg-[#b84f30] disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="trackingLoading"
            >
              {{ trackingLoading ? "Checking" : "Track" }}
            </button>
          </form>

          <p v-if="trackingError" class="mt-4 rounded-2xl bg-[#fff0ed] px-4 py-3 text-sm font-semibold text-[#a23a24]">
            {{ trackingError }}
          </p>

          <div v-if="trackingResult" class="mt-5 rounded-[1.5rem] bg-[#17211f] p-5 text-white">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <p class="text-xs font-bold uppercase tracking-[0.24em] text-[#f2c14e]">Shipment found</p>
                <p class="mt-2 text-2xl font-black tracking-tight">{{ trackingResult.code }}</p>
              </div>
              <span class="rounded-full bg-white px-4 py-2 text-sm font-black text-[#17211f]">
                {{ trackingResult.status_name }}
              </span>
            </div>
            <dl class="mt-6 grid gap-3 sm:grid-cols-2">
              <div class="rounded-2xl bg-white/10 p-4">
                <dt class="text-xs font-bold uppercase tracking-[0.18em] text-white/55">Status</dt>
                <dd class="mt-2 text-lg font-black">{{ trackingResult.status }}</dd>
              </div>
              <div class="rounded-2xl bg-white/10 p-4">
                <dt class="text-xs font-bold uppercase tracking-[0.18em] text-white/55">Last update</dt>
                <dd class="mt-2 text-lg font-black">{{ formatDate(trackingResult.updated_at) }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref } from "vue";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const trackingCode = ref("");
const trackingLoading = ref(false);
const trackingError = ref("");
const trackingResult = ref(null);

function formatDate(value) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}

async function trackShipment() {
  const code = trackingCode.value.trim();
  trackingError.value = "";
  trackingResult.value = null;

  if (!code) {
    trackingError.value = "Please enter an air way bill or shipment code.";
    return;
  }

  trackingLoading.value = true;
  try {
    const { res, json } = await api.get(`/api/shipment-track?code=${encodeURIComponent(code)}`, {
      auth: false,
      navigate: false
    });
    if (!res.ok) {
      trackingError.value = apiErrorMessage(json, res.status === 404 ? "Shipment code was not found." : `Request failed (${res.status})`);
      return;
    }
    trackingResult.value = json?.result?.shipment || null;
    if (!trackingResult.value) trackingError.value = "Shipment code was not found.";
  } catch (err) {
    trackingError.value = String(err);
  } finally {
    trackingLoading.value = false;
  }
}
</script>

