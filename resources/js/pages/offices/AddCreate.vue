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
            <select
              v-model="type"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              :disabled="officeTypesLoading"
              required
            >
              <option value="" disabled>Select type</option>
              <option v-if="type && !officeTypeValues.has(type)" :value="type">
                {{ type }}
              </option>
              <option v-for="t in officeTypes" :key="t.value" :value="t.value">
                {{ t.label }}
              </option>
            </select>
            <p v-if="officeTypesError" class="mt-2 text-sm text-danger">{{ officeTypesError }}</p>
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
            <label class="text-sm font-medium text-slate-700">Province</label>
            <select
              v-model="provinceId"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              :disabled="provincesLoading"
            >
              <option value="">(optional)</option>
              <option v-for="p in provinces" :key="p.id" :value="String(p.id)">
                {{ p.name }}
              </option>
            </select>
            <p v-if="provincesError" class="mt-2 text-sm text-danger">{{ provincesError }}</p>
            <p v-if="fieldErrors.province_id" class="mt-2 text-sm text-danger">{{ fieldErrors.province_id }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-slate-700">City</label>
            <select
              v-model="cityId"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
              :disabled="citiesLoading || !provinceId"
            >
              <option value="">{{ provinceId ? "(optional)" : "Select province first" }}</option>
              <option v-for="c in cities" :key="c.id" :value="String(c.id)">
                {{ c.name }}
              </option>
            </select>
            <p v-if="citiesError" class="mt-2 text-sm text-danger">{{ citiesError }}</p>
            <p v-if="fieldErrors.city_id" class="mt-2 text-sm text-danger">{{ fieldErrors.city_id }}</p>
          </div>

          <div class="md:col-span-2">
            <label class="text-sm font-medium text-slate-700">Image</label>
            <input
              type="file"
              accept="image/jpeg,image/png,image/webp,image/gif"
              class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4"
              :disabled="imageUploading"
              @change="onImageFileChange"
            />
            <p v-if="imageUploading" class="mt-2 text-sm text-slate-500">Uploading image...</p>
            <p v-if="imageUploadError" class="mt-2 text-sm text-danger">{{ imageUploadError }}</p>
            <p v-if="fieldErrors.image_url" class="mt-2 text-sm text-danger">{{ fieldErrors.image_url }}</p>
            <div v-if="imageUrl" class="mt-3 flex items-center gap-3">
              <img
                :src="imageUrl"
                alt=""
                class="h-16 w-24 rounded-lg border border-slate-200 object-cover"
              />
              <div class="min-w-0">
                <p class="truncate text-sm text-slate-700">{{ imageFileName || "Uploaded image" }}</p>
                <button
                  type="button"
                  class="mt-1 text-sm font-medium text-danger hover:underline"
                  :disabled="imageUploading"
                  @click="clearImage"
                >
                  Remove
                </button>
              </div>
            </div>
          </div>
        </div>

        <p v-if="message" class="text-sm" :class="messageToneClass">{{ message }}</p>

        <div class="flex items-center justify-end gap-2">
          <button
            type="submit"
            class="rounded-xl bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading || imageUploading"
          >
            {{ imageUploading ? "Uploading..." : loading ? "Saving..." : isEdit ? "Save changes" : "Create office" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
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
const bootstrapping = ref(true);
const pendingCityId = ref("");
const officeTypes = ref([]); // [{ value, label }]
const officeTypesLoading = ref(false);
const officeTypesError = ref("");
const provinces = ref([]);
const cities = ref([]);
const provincesLoading = ref(false);
const citiesLoading = ref(false);
const provincesError = ref("");
const citiesError = ref("");
const imageUrl = ref("");
const imageUploading = ref(false);
const imageUploadError = ref("");
const imageFileName = ref("");

const loading = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error
const fieldErrors = ref({});

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

const officeTypeValues = computed(() => new Set(officeTypes.value.map((t) => t.value)));

function parseOptionalUint(input) {
  const trimmed = String(input || "").trim();
  if (!trimmed) return undefined;
  const n = Number(trimmed);
  if (!Number.isFinite(n) || n <= 0) return undefined;
  return Math.floor(n);
}

function clearImage() {
  imageUrl.value = "";
  imageFileName.value = "";
  imageUploadError.value = "";
}

function isAllowedImageFile(file) {
  return ["image/jpeg", "image/png", "image/webp", "image/gif"].includes(file?.type || "");
}

async function loadOfficeTypes() {
  officeTypesLoading.value = true;
  officeTypesError.value = "";
  try {
    const { res, json } = await api.get("/api/categories?section=office", { auth: true });
    if (!res.ok || json?.success === false) {
      officeTypes.value = [];
      officeTypesError.value = apiErrorMessage(json, `Failed to load office types (${res.status})`);
      return;
    }

    const categories = Array.isArray(json?.result?.categories) ? json.result.categories : [];
    officeTypes.value = categories.map((c) => ({
      value: typeof c?.slug === "string" && c.slug.trim() ? c.slug.trim() : String(c?.name || "").trim(),
      label: typeof c?.name === "string" && c.name.trim() ? c.name.trim() : String(c?.slug || "").trim()
    }));
  } catch (err) {
    officeTypes.value = [];
    officeTypesError.value = String(err);
  } finally {
    officeTypesLoading.value = false;
  }
}

async function loadProvinces() {
  provincesLoading.value = true;
  provincesError.value = "";
  try {
    const { res, json } = await api.get("/api/provinces", { auth: true });
    if (!res.ok || json?.success === false) {
      provinces.value = [];
      provincesError.value = apiErrorMessage(json, `Failed to load provinces (${res.status})`);
      return;
    }
    provinces.value = Array.isArray(json?.result?.provinces) ? json.result.provinces : [];
  } catch (err) {
    provinces.value = [];
    provincesError.value = String(err);
  } finally {
    provincesLoading.value = false;
  }
}

async function loadCities(provinceIdValue) {
  const pid = String(provinceIdValue || "").trim();
  if (!pid) {
    cities.value = [];
    citiesError.value = "";
    return;
  }

  citiesLoading.value = true;
  citiesError.value = "";
  try {
    const { res, json } = await api.get(`/api/cities?province_id=${encodeURIComponent(pid)}`, { auth: true });
    if (!res.ok || json?.success === false) {
      cities.value = [];
      citiesError.value = apiErrorMessage(json, `Failed to load cities (${res.status})`);
      return;
    }
    cities.value = Array.isArray(json?.result?.cities) ? json.result.cities : [];
  } catch (err) {
    cities.value = [];
    citiesError.value = String(err);
  } finally {
    citiesLoading.value = false;
  }
}

async function onImageFileChange(event) {
  const input = event.target;
  const file = input?.files?.[0];
  imageUploadError.value = "";

  if (!file) return;

  if (!isAllowedImageFile(file)) {
    imageUploadError.value = "Only JPG, PNG, WEBP, or GIF images are allowed.";
    input.value = "";
    return;
  }

  imageUploading.value = true;
  try {
    const form = new FormData();
    form.append("file", file);
    form.append("folder", "offices");

    const { res, json } = await api.post("/api/upload", form, { auth: true });
    if (!res.ok || json?.success === false) {
      imageUploadError.value = apiErrorMessage(json, `Upload failed (${res.status})`);
      return;
    }

    const uploadedUrl = typeof json?.result?.download_url === "string" ? json.result.download_url : "";
    if (!uploadedUrl) {
      imageUploadError.value = "Upload failed: missing image URL";
      return;
    }

    imageUrl.value = uploadedUrl;
    imageFileName.value = file.name;
  } catch (err) {
    imageUploadError.value = String(err);
  } finally {
    imageUploading.value = false;
    input.value = "";
  }
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
    pendingCityId.value = o?.city_id ? String(o.city_id) : "";
    cityId.value = "";
    imageUrl.value = typeof o?.image_url === "string" ? o.image_url : "";
    imageFileName.value = imageUrl.value ? "Current image" : "";
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

  if (imageUploading.value) {
    message.value = "Please wait for the image upload to finish.";
    messageTone.value = "error";
    return;
  }

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

watch(
  provinceId,
  async (next, prev) => {
    // Province changed: reload cities and clear city selection.
    if (bootstrapping.value) return;
    if (String(next || "") === String(prev || "")) return;
    cityId.value = "";
    await loadCities(next);
  },
  { flush: "post" }
);

onMounted(() => {
  // Load static options first, then office data (edit mode), then dependent cities list.
  (async () => {
    bootstrapping.value = true;
    await loadOfficeTypes();
    await loadProvinces();
    await loadOffice();
    if (provinceId.value) {
      await loadCities(provinceId.value);
      if (pendingCityId.value) cityId.value = pendingCityId.value;
    }
    bootstrapping.value = false;
  })();
});
</script>
