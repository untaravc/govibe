<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Shipments</h3>
          <p class="mt-1 text-sm text-slate-600">Manage shipment records.</p>
        </div>

        <button
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          @click="router.push('/admin/shipments/new')"
        >
          Add shipment
        </button>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[980px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">ID</th>
              <th class="px-4 py-3">Code</th>
              <th class="px-4 py-3">Customer</th>
              <th class="px-4 py-3">Phone</th>
              <th class="px-4 py-3">Price</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="s in shipments" :key="s.id">
              <td class="px-4 py-3 text-slate-700">{{ s.id }}</td>
              <td class="px-4 py-3 font-medium text-slate-900">{{ s.code }}</td>
              <td class="px-4 py-3 text-slate-700">{{ s.customer_name }}</td>
              <td class="px-4 py-3 text-slate-700">{{ s.customer_phone }}</td>
              <td class="px-4 py-3 text-slate-700">{{ formatMoney(s.price) }}</td>
              <td class="px-4 py-3 text-slate-700">{{ s.status }}</td>
              <td class="px-4 py-3 text-right">
                <button
                  type="button"
                  class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
                  @click="router.push(`/admin/shipments/${s.id}/edit`)"
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="ml-2 rounded-lg px-3 py-2 text-sm font-medium text-danger hover:bg-rose-50"
                  @click="onDelete(s)"
                >
                  Delete
                </button>
              </td>
            </tr>

            <tr v-if="shipments.length === 0 && !message">
              <td class="px-4 py-8 text-center text-slate-500" colspan="7">No shipments found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const router = useRouter();

const shipments = ref([]);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error

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

async function loadShipments() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.get("/api/shipments", { auth: true });
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
  loadShipments();
});
</script>

