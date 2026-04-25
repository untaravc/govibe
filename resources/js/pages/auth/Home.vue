<template>
  <div class="-mx-4 -my-6 overflow-hidden bg-[#f7f1e7] text-[#17211f]">
    <section class="relative isolate min-h-[720px] px-4 py-16 sm:px-6 lg:px-8">
      <div class="absolute inset-0 -z-10">
        <div
          v-for="(slide, idx) in slides"
          :key="slide.title"
          class="absolute inset-0 transition-opacity duration-700"
          :class="idx === activeSlide ? 'opacity-100' : 'opacity-0'"
        >
          <div class="absolute inset-0" :class="slide.background"></div>
          <div class="absolute left-[8%] top-[16%] h-48 w-48 rounded-full bg-white/30 blur-3xl"></div>
          <div class="absolute bottom-[10%] right-[10%] h-64 w-64 rounded-full bg-black/10 blur-3xl"></div>
        </div>
      </div>

      <div class="mx-auto grid min-h-[620px] max-w-7xl items-center gap-12 lg:grid-cols-[1.05fr_0.95fr]">
        <div>
          <p class="inline-flex rounded-full border border-[#17211f]/15 bg-white/55 px-4 py-2 text-xs font-bold uppercase tracking-[0.28em] text-[#234039] shadow-sm backdrop-blur">
            Same-day shipment network
          </p>
          <h1 class="mt-7 max-w-4xl text-5xl font-black leading-[0.95] tracking-[-0.06em] text-[#14201d] sm:text-7xl lg:text-8xl">
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
            <RouterLink
              to="/auth/login"
              class="rounded-full border border-[#14201d]/15 bg-white/70 px-6 py-3 text-sm font-bold text-[#14201d] backdrop-blur transition hover:-translate-y-0.5 hover:bg-white"
            >
              Admin login
            </RouterLink>
          </div>

          <div class="mt-12 grid max-w-2xl grid-cols-3 gap-3">
            <div v-for="metric in metrics" :key="metric.label" class="rounded-3xl border border-white/60 bg-white/45 p-4 shadow-sm backdrop-blur">
              <p class="text-2xl font-black tracking-tight text-[#14201d]">{{ metric.value }}</p>
              <p class="mt-1 text-xs font-semibold uppercase tracking-[0.18em] text-[#5a716b]">{{ metric.label }}</p>
            </div>
          </div>
        </div>

        <div class="relative">
          <div class="absolute -left-8 top-10 hidden h-28 w-28 rounded-[2rem] bg-[#f2c14e] shadow-2xl shadow-[#ad7c1d]/20 lg:block"></div>
          <div class="absolute -right-4 bottom-12 hidden h-36 w-36 rounded-full bg-[#ef6f4d] shadow-2xl shadow-[#8b331d]/20 lg:block"></div>

          <div class="relative overflow-hidden rounded-[2.25rem] border border-white/60 bg-[#14201d] p-5 text-white shadow-2xl shadow-[#14201d]/25">
            <div class="rounded-[1.75rem] bg-[#21342f] p-5">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-xs font-bold uppercase tracking-[0.24em] text-[#f2c14e]">Live board</p>
                  <p class="mt-1 text-2xl font-black tracking-tight">{{ currentSlide.title }}</p>
                </div>
                <div class="rounded-2xl bg-white/10 px-4 py-2 text-sm font-bold">{{ currentSlide.tag }}</div>
              </div>

              <div class="mt-8 space-y-4">
                <div v-for="route in currentSlide.routes" :key="route.code" class="rounded-3xl bg-white p-4 text-[#14201d] shadow-xl shadow-black/10">
                  <div class="flex items-start justify-between gap-4">
                    <div>
                      <p class="text-sm font-black">{{ route.code }}</p>
                      <p class="mt-1 text-xs font-semibold uppercase tracking-[0.18em] text-[#6a7c77]">{{ route.path }}</p>
                    </div>
                    <span class="rounded-full bg-[#dff4e8] px-3 py-1 text-xs font-black text-[#17613d]">{{ route.status }}</span>
                  </div>
                  <div class="mt-4 h-2 overflow-hidden rounded-full bg-[#edf1ec]">
                    <div class="h-full rounded-full bg-[#f2c14e]" :style="{ width: route.progress }"></div>
                  </div>
                </div>
              </div>
            </div>

            <div class="mt-5 flex items-center justify-between px-2">
              <button
                type="button"
                class="rounded-full border border-white/15 px-4 py-2 text-sm font-bold text-white/85 transition hover:bg-white/10"
                @click="previousSlide"
              >
                Prev
              </button>
              <div class="flex items-center gap-2">
                <button
                  v-for="(_, idx) in slides"
                  :key="idx"
                  type="button"
                  class="h-2.5 rounded-full transition-all"
                  :class="idx === activeSlide ? 'w-8 bg-[#f2c14e]' : 'w-2.5 bg-white/35'"
                  @click="activeSlide = idx"
                >
                  <span class="sr-only">Show slide {{ idx + 1 }}</span>
                </button>
              </div>
              <button
                type="button"
                class="rounded-full border border-white/15 px-4 py-2 text-sm font-bold text-white/85 transition hover:bg-white/10"
                @click="nextSlide"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="track" class="relative bg-[#fffaf1] px-4 py-20 sm:px-6 lg:px-8">
      <div class="mx-auto grid max-w-6xl gap-8 lg:grid-cols-[0.8fr_1.2fr]">
        <div>
          <p class="text-sm font-black uppercase tracking-[0.26em] text-[#d0603d]">Track air way bill</p>
          <h2 class="mt-4 text-4xl font-black tracking-[-0.04em] text-[#14201d] sm:text-5xl">
            Enter your shipment code.
          </h2>
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
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const slides = [
  {
    title: "Express lanes",
    tag: "Fast route",
    background: "bg-[radial-gradient(circle_at_20%_20%,#ffe8a3_0,#f7c55c_28%,#e5794f_62%,#21342f_100%)]",
    routes: [
      { code: "SHP-1048", path: "Jakarta to Bandung", status: "Perjalanan", progress: "58%" },
      { code: "SHP-1052", path: "Medan to Pekanbaru", status: "Dibuat", progress: "18%" }
    ]
  },
  {
    title: "Transit visibility",
    tag: "Hub scan",
    background: "bg-[radial-gradient(circle_at_75%_20%,#b8f2e6_0,#6fbfb2_32%,#315d56_64%,#101b19_100%)]",
    routes: [
      { code: "SHP-2110", path: "Surabaya Hub", status: "Tiba di Transit", progress: "72%" },
      { code: "SHP-2117", path: "Semarang to Solo", status: "Berangkat", progress: "83%" }
    ]
  },
  {
    title: "Destination proof",
    tag: "Arrival",
    background: "bg-[radial-gradient(circle_at_25%_75%,#f8d7c5_0,#e49a72_35%,#5b3f36_68%,#17211f_100%)]",
    routes: [
      { code: "SHP-3201", path: "Denpasar Destination", status: "Sampai", progress: "100%" },
      { code: "SHP-3208", path: "Makassar to Gowa", status: "Perjalanan", progress: "64%" }
    ]
  }
];

const metrics = [
  { value: "24/7", label: "Tracking" },
  { value: "5", label: "Statuses" },
  { value: "1 Code", label: "Lookup" }
];

const activeSlide = ref(0);
const trackingCode = ref("");
const trackingLoading = ref(false);
const trackingError = ref("");
const trackingResult = ref(null);
let slideTimer = null;

const currentSlide = computed(() => slides[activeSlide.value] || slides[0]);

function nextSlide() {
  activeSlide.value = (activeSlide.value + 1) % slides.length;
}

function previousSlide() {
  activeSlide.value = (activeSlide.value - 1 + slides.length) % slides.length;
}

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

onMounted(() => {
  slideTimer = window.setInterval(nextSlide, 5200);
});

onBeforeUnmount(() => {
  if (slideTimer) window.clearInterval(slideTimer);
});
</script>
