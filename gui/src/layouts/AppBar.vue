<template>
  <v-app-bar class="bg-grey-lighten-1">
    <v-app-bar-title class="text-h5 font-weight-bold"
      >Auto Analyser</v-app-bar-title
    >
    <template v-slot:append>
      <v-dialog v-model="dialog"
      width="800">
        <v-card>
        <v-card-text>
          <v-file-input class="ma-3" show-size v-model="dumpfile" label="File to analyse"></v-file-input>
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" block @click="selectdump()" 
        :disabled="dumpfile == null">Select dump</v-btn>
        </v-card-actions>
        </v-card>
      </v-dialog>

      <v-btn
        prepend-icon="mdi-plus"
        class="text-none me-2"
        rounded
        @click="dialog = true"
        variant="outlined">
        New
        
      </v-btn>
      
      <!-- currently not working -->
      <v-btn
        prepend-icon="mdi-send"
        class="text-none"
        rounded
        disabled
        variant="outlined"
      >
        Run</v-btn
      >
    </template>
  </v-app-bar>
</template>

<script setup lang="ts">

import { ref } from 'vue';
import { useModuleStore } from "@/stores/ModulesStore";

const store = useModuleStore();
const dialog = ref(false);
const dumpfile:any = ref(null);

const selectdump =() => {
  dialog.value = false;
  const file: File = dumpfile.value[0]
  console.log(file.name);
  store.setFile(file);
} 

</script>



<style scoped></style>
