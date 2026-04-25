<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">{{ isEdit ? "Edit shipment" : "Add shipment" }}</h3>
          <p class="mt-1 text-sm text-slate-600">
            {{ isEdit ? "Update shipment details." : "Create a new shipment." }}
          </p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="router.push('/admin/shipments')"
        >
          Back
        </button>
      </div>

      <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="text-sm font-medium text-slate-700">Code</label>
            <input
              v-model.trim="code"
              type="text"
              placeholder="SHP-0001"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.code" class="mt-2 text-sm text-danger">{{ fieldErrors.code }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Status</label>
            <input
              v-model.number="status"
              type="number"
              min="0"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.status" class="mt-2 text-sm text-danger">{{ fieldErrors.status }}</p>
          </div>

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Customer Name</label>
            <input
              v-model.trim="customerName"
              type="text"
              placeholder="Customer name"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.customer_name" class="mt-2 text-sm text-danger">{{ fieldErrors.customer_name }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Customer Phone</label>
            <input
              v-model.trim="customerPhone"
              type="text"
              placeholder="+62..."
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              required
            />
            <p v-if="fieldErrors.customer_phone" class="mt-2 text-sm text-danger">{{ fieldErrors.customer_phone }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Customer Email</label>
            <input
              v-model.trim="customerEmail"
              type="email"
              placeholder="(optional)"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.customer_email" class="mt-2 text-sm text-danger">{{ fieldErrors.customer_email }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Origin Office</label>
            <select
              v-model="officeOriginId"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              :disabled="officesLoading"
              required
            >
              <option value="" disabled>Select office</option>
              <option v-for="o in offices" :key="o.id" :value="String(o.id)">
                {{ o.name }}
              </option>
            </select>
            <p v-if="officesError" class="mt-2 text-sm text-danger">{{ officesError }}</p>
            <p v-if="fieldErrors.office_origin_id" class="mt-2 text-sm text-danger">{{ fieldErrors.office_origin_id }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Destination Office</label>
            <select
              v-model="officeDestinationId"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              :disabled="officesLoading"
              required
            >
              <option value="" disabled>Select office</option>
              <option v-for="o in offices" :key="o.id" :value="String(o.id)">
                {{ o.name }}
              </option>
            </select>
            <p v-if="officesError" class="mt-2 text-sm text-danger">{{ officesError }}</p>
            <p v-if="fieldErrors.office_destination_id" class="mt-2 text-sm text-danger">{{ fieldErrors.office_destination_id }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Price</label>
            <input
              v-model.number="price"
              type="number"
              min="0"
              step="0.01"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
            <p v-if="fieldErrors.price" class="mt-2 text-sm text-danger">{{ fieldErrors.price }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">Price Type</label>
            <select
              v-model="priceType"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              required
            >
              <option value="weight">weight</option>
              <option value="dimension">dimension</option>
            </select>
            <p v-if="fieldErrors.price_type" class="mt-2 text-sm text-danger">{{ fieldErrors.price_type }}</p>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 md:grid-cols-5">
          <div>
            <label class="text-sm font-medium text-slate-700">Wight</label>
            <input
              v-model.number="wight"
              type="number"
              min="0"
              step="0.001"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
          </div>
          <div>
            <label class="text-sm font-medium text-slate-700">Length</label>
            <input
              v-model.number="length"
              type="number"
              min="0"
              step="0.001"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
          </div>
          <div>
            <label class="text-sm font-medium text-slate-700">Width</label>
            <input
              v-model.number="width"
              type="number"
              min="0"
              step="0.001"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
          </div>
          <div>
            <label class="text-sm font-medium text-slate-700">Height</label>
            <input
              v-model.number="height"
              type="number"
              min="0"
              step="0.001"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
            />
          </div>
          <div>
            <label class="text-sm font-medium text-slate-700">Items</label>
            <button
              type="button"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
              @click="addDetail"
            >
              Add item
            </button>
          </div>
        </div>

        <div class="overflow-auto rounded-xl border border-slate-200">
          <table class="w-full min-w-[840px] text-left text-sm">
            <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
              <tr>
                <th class="px-4 py-3">Item</th>
                <th class="px-4 py-3">Price</th>
                <th class="px-4 py-3">Category</th>
                <th class="px-4 py-3 text-right">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200">
              <tr v-for="(d, idx) in details" :key="idx">
                <td class="px-4 py-3">
                  <input
                    v-model.trim="d.item_name"
                    type="text"
                    class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
                    placeholder="Item name"
                    required
                  />
                </td>
                <td class="px-4 py-3">
                  <input
                    v-model.number="d.item_price"
                    type="number"
                    min="0"
                    step="0.01"
                    class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
                  />
                </td>
                <td class="px-4 py-3">
                  <select
                    v-model="d.category_id"
                    class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
                    :disabled="categoriesLoading"
                  >
                    <option value="">No category</option>
                    <option v-for="c in categories" :key="c.id" :value="String(c.id)">
                      {{ c.name }}
                    </option>
                  </select>
                </td>
                <td class="px-4 py-3 text-right">
                  <button
                    type="button"
                    class="rounded-lg px-3 py-2 text-sm font-medium text-danger hover:bg-rose-50"
                    @click="removeDetail(idx)"
                  >
                    Remove
                  </button>
                </td>
              </tr>
              <tr v-if="details.length === 0">
                <td class="px-4 py-8 text-center text-slate-500" colspan="4">No items.</td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
          >
            {{ loading ? "Saving..." : isEdit ? "Save changes" : "Create shipment" }}
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

const code = ref("");
const customerName = ref("");
const customerPhone = ref("");
const customerEmail = ref("");
const officeOriginId = ref("");
const officeDestinationId = ref("");
const price = ref(0);
const status = ref(0);
const priceType = ref("weight");
const wight = ref(0);
const length = ref(0);
const width = ref(0);
const height = ref(0);

const details = ref([]);

const offices = ref([]);
const officesLoading = ref(false);
const officesError = ref("");

const categories = ref([]);
const categoriesLoading = ref(false);
const categoriesError = ref("");

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

function addDetail() {
  details.value.push({ item_name: "", item_price: 0, category_id: "" });
}

function removeDetail(idx) {
  details.value.splice(idx, 1);
}

async function loadOffices() {
  officesLoading.value = true;
  officesError.value = "";
  try {
    const { res, json } = await api.get("/api/offices", { auth: true });
    if (!res.ok || json?.success === false) {
      offices.value = [];
      officesError.value = apiErrorMessage(json, `Failed to load offices (${res.status})`);
      return;
    }
    offices.value = Array.isArray(json?.result?.offices) ? json.result.offices : [];
  } catch (err) {
    offices.value = [];
    officesError.value = String(err);
  } finally {
    officesLoading.value = false;
  }
}

async function loadCategories() {
  categoriesLoading.value = true;
  categoriesError.value = "";
  try {
    const { res, json } = await api.get("/api/categories", { auth: true });
    if (!res.ok || json?.success === false) {
      categories.value = [];
      categoriesError.value = apiErrorMessage(json, `Failed to load categories (${res.status})`);
      return;
    }
    categories.value = Array.isArray(json?.result?.categories) ? json.result.categories : [];
  } catch (err) {
    categories.value = [];
    categoriesError.value = String(err);
  } finally {
    categoriesLoading.value = false;
  }
}

async function loadShipment() {
  if (!isEdit.value) return;

  loading.value = true;
  message.value = "";
  messageTone.value = "neutral";
  fieldErrors.value = {};

  try {
    const { res, json } = await api.get(`/api/shipments/${id.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    const s = json?.result?.shipment;
    code.value = typeof s?.code === "string" ? s.code : "";
    customerName.value = typeof s?.customer_name === "string" ? s.customer_name : "";
    customerPhone.value = typeof s?.customer_phone === "string" ? s.customer_phone : "";
    customerEmail.value = typeof s?.customer_email === "string" ? s.customer_email : "";
    officeOriginId.value = s?.office_origin_id ? String(s.office_origin_id) : "";
    officeDestinationId.value = s?.office_destination_id ? String(s.office_destination_id) : "";
    price.value = Number(s?.price || 0);
    status.value = Number(s?.status || 0);
    priceType.value = typeof s?.price_type === "string" ? s.price_type : "weight";
    wight.value = Number(s?.wight || 0);
    length.value = Number(s?.length || 0);
    width.value = Number(s?.width || 0);
    height.value = Number(s?.height || 0);

    const incoming = Array.isArray(s?.details) ? s.details : [];
    details.value = incoming.map((d) => ({
      item_name: typeof d?.item_name === "string" ? d.item_name : "",
      item_price: Number(d?.item_price || 0),
      category_id: d?.category_id ? String(d.category_id) : ""
    }));
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
      code: code.value,
      customer_name: customerName.value,
      customer_phone: customerPhone.value,
      customer_email: customerEmail.value.trim() ? customerEmail.value : null,
      office_origin_id: officeOriginId.value ? Number(officeOriginId.value) : null,
      office_destination_id: officeDestinationId.value ? Number(officeDestinationId.value) : null,
      price: Number(price.value || 0),
      status: Number(status.value || 0),
      price_type: priceType.value,
      wight: Number(wight.value || 0),
      length: Number(length.value || 0),
      width: Number(width.value || 0),
      height: Number(height.value || 0),
      details: details.value.map((d) => ({
        item_name: d.item_name,
        item_price: Number(d.item_price || 0),
        category_id: d.category_id ? Number(d.category_id) : null
      }))
    };

    const { res, json } = isEdit.value
      ? await api.put(`/api/shipments/${id.value}`, body, { auth: true })
      : await api.post("/api/shipments", body, { auth: true });

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
      title: isEdit.value ? "Shipment updated" : "Shipment created",
      showConfirmButton: false,
      timer: 1600,
      timerProgressBar: true
    });

    setTimeout(() => router.push("/admin/shipments"), 250);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadOffices();
  loadCategories();
  loadShipment();
  if (!isEdit.value && details.value.length === 0) addDetail();
});
</script>

