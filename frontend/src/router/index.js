import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import HomeVieweur from "../views/HomeVieweur.vue";
import HomeViewtl from "../views/HomeViewtl.vue";
import HomeViewhas from "../views/HomeViewhas.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/eur",
      name: "home1",
      component: HomeVieweur,
    },
    {
      path: "/tl",
      name: "home2",
      component: HomeViewtl,
    },
    {
      path: "/has",
      name: "home3",
      component: HomeViewhas,
    },
   
  ],
});

export default router;
