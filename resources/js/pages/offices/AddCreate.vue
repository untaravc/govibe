<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">{{ isEdit ? "Edit office" : "Add office" }}</h3>
          <p class="mt-1 text-sm text-slate-600">
            {{ isEdit ? "Update office details." : "Create a new office entry." }}
          </p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="router.push('/admin/offices')"
        >
          Back
        </button>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">Type</label>
            <input
              v-model.trim="type"
              type="text"
              placeholder="hq / branch / etc"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.type" class="mt-2 text-sm text-danger">{{ fieldErrors.type }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Code</label>
            <input
              v-model.trim="code"
              type="text"
              placeholder="OFFICE-001"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.code" class="mt-2 text-sm text-danger">{{ fieldErrors.code }}</p>
          </div>

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Name</label>
            <input
              v-model.trim="name"
              type="text"
              placeholder="Office name"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.name" class="mt-2 text-sm text-danger">{{ fieldErrors.name }}</p>
          </div>

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Address</label>
            <textarea
              v-model.trim="address"
              rows="3"
              placeholder="Address (optional)"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.address" class="mt-2 text-sm text-danger">{{ fieldErrors.address }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Phone</label>
            <input
              v-model.trim="phone"
              type="text"
              placeholder="+62..."
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.phone" class="mt-2 text-sm text-danger">{{ fieldErrors.phone }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Status</label>
            <select
              v-model.number="status"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            >
              <option :value="1">Active</option>
              <option :value="0">Inactive</option>
            </select>
            <p v-if="fieldErrors.status" class="mt-2 text-sm text-danger">{{ fieldErrors.status }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Province ID</label>
            <input
              v-model.trim="provinceId"
              type="number"
              min="0"
              placeholder="(optional)"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.province_id" class="mt-2 text-sm text-danger">{{ fieldErrors.province_id }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">City ID</label>
            <input
              v-model.trim="cityId"
              type="number"
              min="0"
              placeholder="(optional)"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.city_id" class="mt-2 text-sm text-danger">{{ fieldErrors.city_id }}</p>
          </div>

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Image URL</label>
            <input
              v-model.trim="imageUrl"
              type="text"
              placeholder="https://..."
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.image_url" class="mt-2 text-sm text-danger">{{ fieldErrors.image_url }}</p>
          </div>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : isEdit ? "Save changes" : "Create office" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import Swal from "sweetalert2";

import api from "../../api.js";
import { apiErrorMessage, apiFieldErrors } from "../../utils/apiError.js";

const router = useRouter();
const route = useRoute();

const id = computed(() => String(route.params.id || "").trim());
const isEdit = computed(() => id.value.length > 0);

const type = ref("");
const name = ref("");
const code = ref("");
const address = ref("");
const phone = ref("");
const status = ref(1);
const provinceId = ref("");
const cityId = ref("");
const imageUrl = ref("");

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

function parseOptionalUint(input) {
  const trimmed = String(input || "").trim();
  if (!trimmed) return undefined;
  const n = Number(trimmed);
  if (!Number.isFinite(n) || n <= 0) return undefined;
  return Math.floor(n);
}

async function loadOffice() {
  if (!isEdit.value) return;

  loading.value = true;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};

  try {
    const { res, json } = await api.get(`/api/offices/${id.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const o = json?.result?.office;
    type.value = typeof o?.type === "string" ? o.type : "";
    name.value = typeof o?.name === "string" ? o.name : "";
    code.value = typeof o?.code === "string" ? o.code : "";
    address.value = typeof o?.address === "string" ? o.address : "";
    phone.value = typeof o?.phone === "string" ? o.phone : "";
    status.value = typeof o?.status === "number" ? o.status : 1;
    provinceId.value = o?.province_id ? String(o.province_id) : "";
    cityId.value = o?.city_id ? String(o.city_id) : "";
    imageUrl.value = typeof o?.image_url === "string" ? o.image_url : "";
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

async function onSubmit() {
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};
  loading.value = true;

  try {
    const body = {
      type: type.value,
      name: name.value,
      code: code.value,
      status: Number(status.value)
    };

    const addr = address.value.trim();
    const ph = phone.value.trim();
    const img = imageUrl.value.trim();
    if (addr) body.address = addr;
    if (ph) body.phone = ph;
    if (img) body.image_url = img;

    const p = parseOptionalUint(provinceId.value);
    const city = parseOptionalUint(cityId.value);
    if (p) body.province_id = p;
    if (city) body.city_id = city;

    const { res, json } = isEdit.value
      ? await api.put(`/api/offices/${id.value}`, body, { auth: true })
      : await api.post("/api/offices", body, { auth: true });
    if (!res.ok) {
      fieldErrors.value = apiFieldErrors(json);
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    Swal.fire({
      toast: true,
      position: "top-end",
      icon: "success",
      title: isEdit.value ? "Office updated" : "Office created",
      showConfirmButton: false,
      timer: 1600,
      timerProgressBar: true
    });

    setTimeout(() => router.push("/admin/offices"), 300);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadOffice();
});
</script>

