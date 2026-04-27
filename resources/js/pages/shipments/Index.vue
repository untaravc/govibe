<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">{{ pageTitle }}</h3>
          <p class="mt-1 text-sm text-slate-600">{{ pageDescription }}</p>
        </div>

        <button
          v-if="canAddShipment"
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          @click="router.push('/admin/shipments/new')"
        >
          Add shipment
        </button>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="mt-4 grid grid-cols-1 gap-3 md:grid-cols-3">
        <div>
          <label class="text-sm font-medium text-slate-700">Code</label>
          <input
            v-model="codeFilter"
            type="text"
            placeholder="Search code"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            autocomplete="off"
          />
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Customer</label>
          <input
            v-model="customerNameFilter"
            type="text"
            placeholder="Search customer name"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            autocomplete="off"
          />
        </div>

        <div>
          <label class="text-sm font-medium text-slate-700">Created date</label>
          <input
            v-model="createdDateFilter"
            type="date"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
          />
        </div>
      </div>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[980px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">No</th>
              <th class="px-4 py-3">Code</th>
              <th class="px-4 py-3">Customer</th>
              <th class="px-4 py-3">Price</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="(s, idx) in shipments" :key="s.id">
              <td class="px-4 py-3 text-slate-700">{{ idx + 1 }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ s.code }}</td>
              <td class="px-4 py-3 text-slate-700">{{ s.customer_name }}</td>
              <td class="px-4 py-3 text-slate-700">{{ formatMoney(s.price) }}</td>
              <td class="px-4 py-3 text-slate-700">{{ s.status }}</td>
              <td class="px-4 py-3 text-right">
                <div :ref="(el) => setActionRoot(s.id, el)" class="relative inline-block text-left">
                  <button
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                    :aria-expanded="actionsOpenFor === s.id ? 'true' : 'false'"
                    aria-haspopup="menu"
                    @click="toggleActions(s.id)"
                  >
                    Actions
                    <svg viewBox="0 0 24 24" fill="none" class="h-4 w-4" xmlns="http://www.w3.org/2000/svg">
                      <path
                        d="M6 9l6 6 6-6"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </button>

                  <div
                    v-if="actionsOpenFor === s.id"
                    role="menu"
                    class="absolute right-0 z-10 mt-2 w-40 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-lg"
                  >
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-slate-50"
                      @click="
                        closeActions();
                        router.push(`/admin/shipments/${s.id}/edit`);
                      "
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      class="block w-full px-4 py-2.5 text-left text-sm text-rose-600 hover:bg-rose-50"
                      @click="
                        closeActions();
                        onDelete(s);
                      "
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </td>
            </tr>

            <tr v-if="shipments.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="6">No shipments found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const router = useRouter();
const route = useRoute();

const shipments = ref([]);
const codeFilter = ref("");
const customerNameFilter = ref("");
const createdDateFilter = ref("");
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const actionsOpenFor = ref(null);

const actionRoots = new Map();
let filterTimer = null;

const shipmentTypeLabels = {
  departure: "Departure",
  transit: "Transit",
  arrive: "Arrive"
};

const shipmentType = computed(() => {
  const value = Array.isArray(route.query.type) ? route.query.type[0] : route.query.type;
  return typeof value === "string" ? value.trim().toLowerCase() : "";
});

const pageTitle = computed(() => {
  const label = shipmentTypeLabels[shipmentType.value];
  return label ? `${label} Shipments` : "Shipments";
});

const pageDescription = computed(() => {
  const label = shipmentTypeLabels[shipmentType.value];
  return label ? `Manage ${label.toLowerCase()} shipment records.` : "Manage shipment records.";
});

const canAddShipment = computed(() => shipmentType.value === "departure");

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

function formatMoney(value) {
  const n = Number(value);
  if (!Number.isFinite(n)) return String(value ?? "");
  return n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function setActionRoot(id, el) {
  if (!id) return;
  if (!el) {
    actionRoots.delete(id);
    return;
  }
  actionRoots.set(id, el);
}

function closeActions() {
  actionsOpenFor.value = null;
}

function toggleActions(id) {
  actionsOpenFor.value = actionsOpenFor.value === id ? null : id;
}

function onWindowClick(e) {
  if (!actionsOpenFor.value) return;
  const root = actionRoots.get(actionsOpenFor.value);
  if (!root) return closeActions();
  if (root === e.target || root.contains(e.target)) return;
  closeActions();
}

function onWindowKeydown(e) {
  if (!actionsOpenFor.value) return;
  if (e.key === "Escape") closeActions();
}

async function loadShipments() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const params = new URLSearchParams();
    if (shipmentType.value) params.set("type", shipmentType.value);
    if (codeFilter.value.trim()) params.set("code", codeFilter.value.trim());
    if (customerNameFilter.value.trim()) params.set("customer_name", customerNameFilter.value.trim());
    if (createdDateFilter.value) params.set("created_date", createdDateFilter.value);

    const url = params.toString() ? `/api/shipments?${params.toString()}` : "/api/shipments";
    const { res, json } = await api.get(url, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      shipments.value = [];
      return;
    }
    shipments.value = Array.isArray(json?.result?.shipments) ? json.result.shipments : [];
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    shipments.value = [];
  }
}

async function onDelete(s) {
  const ok = confirm(`Delete shipment "${s.code}"?`);
  if (!ok) return;

  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.delete(`/api/shipments/${s.id}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Delete failed (${res.status})`);
      messageTone.value = "error";
      return;
    }
    message.value = "Shipment deleted.";
    messageTone.value = "success";
    await loadShipments();
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  }
}

onMounted(() => {
  window.addEventListener("click", onWindowClick, true);
  window.addEventListener("keydown", onWindowKeydown);
  loadShipments();
});

onBeforeUnmount(() => {
  window.removeEventListener("click", onWindowClick, true);
  window.removeEventListener("keydown", onWindowKeydown);
  if (filterTimer) clearTimeout(filterTimer);
});

watch(
  () => route.query.type,
  () => {
    loadShipments();
  }
);

watch([codeFilter, customerNameFilter, createdDateFilter], () => {
  if (filterTimer) clearTimeout(filterTimer);
  filterTimer = setTimeout(() => {
    loadShipments();
    filterTimer = null;
  }, 350);
});
</script>
