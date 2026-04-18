import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./route.js";
import "../css/app.css";
import { Icon } from "@iconify/vue";
import pageLoader from "./plugins/pageLoader.js";

const el = document.getElementById("app");
if (!el) {
  throw new Error("Missing #app element");
}

const title = el.dataset.title || "GoVibe (Fiber)";
const app = createApp(App, { title });

app.use(createPinia());
app.use(router);
app.use(pageLoader);
app.component("Icon", Icon);

app.mount(el);
