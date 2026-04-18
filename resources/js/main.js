import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./route.js";
import "../css/app.css";

const el = document.getElementById("app");
if (!el) {
  throw new Error("Missing #app element");
}

const title = el.dataset.title || "GoVibe (Fiber)";
const app = createApp(App, { title });

app.use(createPinia());
app.use(router);

app.mount(el);
