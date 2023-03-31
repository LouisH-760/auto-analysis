<template>
  <v-container class="bg-grey-lighten-2 h-100 elevation-5">
      <div v-if="can_edit" class="d-flex flex-column h-100">
      <!-- function to add modules -->
      <div class="ma-2 pa-2">
        <v-btn
          prepend-icon="mdi-plus"
          rounded
          class="mr-1"
          variant="outlined">
          add node
        </v-btn>
        <v-btn
          prepend-icon="mdi-backspace"
          rounded
          variant="outlined">
          clear
        </v-btn>
      </div>
      <div class="ma-2">
      </div>
    </div>
    <div v-else>
      
    </div>
  </v-container>
</template>

<script setup lang="ts">

import { useModuleStore } from "@/stores/ModulesStore";

import { storeToRefs } from 'pinia';


const store = useModuleStore();

const {can_edit} = storeToRefs(store);
</script>

<style scoped></style>
