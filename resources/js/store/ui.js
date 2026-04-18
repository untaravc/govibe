import { defineStore } from "pinia";

export const useUiStore = defineStore("ui", {
  state: () => ({
    apiPendingCount: 0
  }),
  getters: {
    isApiLoading: (state) => state.apiPendingCount > 0
  },
  actions: {
    apiStart() {
      this.apiPendingCount += 1;
    },
    apiStop() {
      this.apiPendingCount = Math.max(0, this.apiPendingCount - 1);
    }
  }
});

