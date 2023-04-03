<script setup lang="ts">
import { ref } from "vue";
import AppBar from "@/layouts/AppBar.vue";
import SideBar from "./layouts/SideBar.vue";
const author = ref<String>("Louis Hermier | Nico Gruschke");
</script>

<template>
<v-app>
    <v-layout>
      <AppBar />
      <v-main>
        <v-row no-gutters class="h-100">
          <v-col cols="2">
            <SideBar />
          </v-col>
          <v-col>
          </v-col>
        </v-row>
      </v-main>
    </v-layout>
    <v-footer class="flex-grow-0 bg-grey-lighten-1">
      <div class="px-4 py-2 text-center w-100">
        {{ new Date().getFullYear() }} - <strong>{{ author }}</strong>
      </div>
    </v-footer>
  </v-app>
</template>

<style scoped>

</style>
