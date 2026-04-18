<template>
  <div class="space-y-4">
    <div class="rounded-2xl border border-slate-200/60 bg-white p-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold tracking-tight text-slate-900">Menu Role</h3>
          <p class="mt-1 text-sm text-slate-600">Assign menu access per role.</p>
        </div>

        <button
          type="button"
          class="rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="!selectedRoleId || saving"
          @click="onSave"
        >
          {{ saving ? "Saving..." : "Save" }}
        </button>
      </div>

      <div class="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <label class="text-sm font-medium text-slate-700">Role</label>
          <select
            v-model.number="selectedRoleId"
            class="mt-2 w-full rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-slate-900 shadow-sm outline-none ring-slate-900/10 focus:border-slate-300 focus:ring-4"
          >
            <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.role }} — {{ r.name }}</option>
          </select>
        </div>
      </div>

      <p v-if="message" class="mt-4 text-sm" :class="messageToneClass">{{ message }}</p>

      <div class="relative mt-6 overflow-auto rounded-xl border border-slate-200">
        <PageLoader :fullscreen="false" />
        <table class="w-full min-w-[860px] text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-4 py-3">Menu</th>
              <th v-for="m in uiMethods" :key="m.key" class="px-4 py-3 text-center">{{ m.label }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            <tr v-for="row in menuRows" :key="row.id" class="align-middle">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2" :style="{ paddingLeft: `${row.depth * 16}px` }">
                  <span v-if="row.depth > 0" class="text-slate-300">—</span>
                  <span class="font-medium text-slate-900">{{ row.name }}</span>
                  <span v-if="row.link" class="text-xs text-slate-500">({{ row.link }})</span>
                </div>
              </td>

              <td v-for="m in uiMethods" :key="m.key" class="px-4 py-3 text-center">
                <input
                  type="checkbox"
                  class="h-4 w-4 rounded border-slate-300 text-slate-900 focus:ring-slate-900"
                  :checked="isChecked(row.id, m.key)"
                  @change="toggleChecked(row.id, m.key, $event.target.checked)"
                />
              </td>
            </tr>

            <tr v-if="menuRows.length === 0 && selectedRoleId">
              <td class="px-4 py-8 text-center text-slate-500" :colspan="1 + uiMethods.length">No menus found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import Swal from "sweetalert2";

import api from "../../api.js";
import { apiErrorMessage } from "../../utils/apiError.js";

const roles = ref([]);
const selectedRoleId = ref(0);

const menus = ref([]); // tree [{..., children: []}]
const grants = ref([]); // [{menu_id, method}]
const checked = ref({}); // key => true

const saving = ref(false);
const message = ref("");
const messageTone = ref("neutral"); // neutral | success | error

const messageToneClass = computed(() => {
  if (messageTone.value === "success") return "text-success";
  if (messageTone.value === "error") return "text-danger";
  return "text-slate-700";
});

const uiMethods = [
  { key: "get", label: "GET", api: "get" },
  { key: "post", label: "POST", api: "create" },
  { key: "patch", label: "PATCH", api: "update" },
  { key: "delete", label: "DELETE", api: "delete" }
];

function grantKey(menuId, methodKey) {
  return `${menuId}:${methodKey}`;
}

function isChecked(menuId, methodKey) {
  return Boolean(checked.value[grantKey(menuId, methodKey)]);
}

function toggleChecked(menuId, methodKey, isOn) {
  const key = grantKey(menuId, methodKey);
  checked.value = { ...checked.value, [key]: isOn };
}

function flattenTree(tree, depth = 0, out = []) {
  for (const n of tree) {
    out.push({ id: n.id, name: n.name, link: n.link || "", depth });
    if (Array.isArray(n.children) && n.children.length > 0) flattenTree(n.children, depth + 1, out);
  }
  return out;
}

const menuRows = computed(() => {
  return flattenTree(Array.isArray(menus.value) ? menus.value : []);
});

function applyGrants(grantRows) {
  const next = {};

  for (const g of grantRows) {
    if (!g || typeof g !== "object") continue;
    const menuId = Number(g.menu_id);
    const method = String(g.method || "").toLowerCase();
    if (!menuId || !method) continue;

    // Map API methods to UI keys.
    if (method === "get") next[grantKey(menuId, "get")] = true;
    if (method === "create") next[grantKey(menuId, "post")] = true;
    if (method === "update") next[grantKey(menuId, "patch")] = true;
    if (method === "delete") next[grantKey(menuId, "delete")] = true;
  }

  checked.value = next;
}

async function loadRoles() {
  message.value = "";
  messageTone.value = "neutral";
  try {
    const { res, json } = await api.get("/api/roles", { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      roles.value = [];
      return;
    }
    roles.value = Array.isArray(json?.result?.roles) ? json.result.roles : [];
    if (!selectedRoleId.value && roles.value.length > 0) {
      selectedRoleId.value = roles.value[0].id;
    }
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    roles.value = [];
  }
}

async function loadMenus() {
  message.value = "";
  messageTone.value = "neutral";

  try {
    const { res, json } = await api.get("/api/menus", { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      menus.value = [];
      return;
    }

    menus.value = Array.isArray(json?.result?.menus) ? json.result.menus : [];
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    menus.value = [];
  }
}

async function loadGrantsForRole() {
  if (!selectedRoleId.value) return;

  message.value = "";
  messageTone.value = "neutral";
  checked.value = {};

  try {
    const { res, json } = await api.get(`/api/menu-roles?role_id=${selectedRoleId.value}`, { auth: true });
    if (!res.ok) {
      message.value = apiErrorMessage(json, `Request failed (${res.status})`);
      messageTone.value = "error";
      grants.value = [];
      return;
    }

    grants.value = Array.isArray(json?.result?.grants) ? json.result.grants : [];
    applyGrants(grants.value);
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
    grants.value = [];
  }
}

async function onSave() {
  if (!selectedRoleId.value) return;

  saving.value = true;
  message.value = "";
  messageTone.value = "neutral";
  try {
    const grantItems = [];

    for (const row of menuRows.value) {
      const menuId = row.id;
      for (const m of uiMethods) {
        if (!isChecked(menuId, m.key)) continue;
        grantItems.push({ menu_id: menuId, method: m.api });
      }
    }

    const { res, json } = await api.post(
      "/api/menu-roles",
      {
        role_id: Number(selectedRoleId.value),
        grants: grantItems
      },
      { auth: true }
    );

    if (!res.ok) {
      message.value = apiErrorMessage(json, `Save failed (${res.status})`);
      messageTone.value = "error";
      return;
    }

    Swal.fire({
      toast: true,
      position: "top-end",
      icon: "success",
      title: "Saved",
      showConfirmButton: false,
      timer: 1400,
      timerProgressBar: true
    });
  } catch (err) {
    message.value = String(err);
    messageTone.value = "error";
  } finally {
    saving.value = false;
  }
}

watch(
  () => selectedRoleId.value,
  () => {
    loadGrantsForRole();
  }
);

onMounted(() => {
  loadRoles();
  loadMenus();
});
</script>
