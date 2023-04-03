<template>
  <v-expansion-panels variant="accordion">
    <v-expansion-panel v-for="element in defaultModules" :key="element.name">
      <v-expansion-panel-title>{{ element.name }}</v-expansion-panel-title>
      <v-expansion-panel-text>
        <v-chip
          v-for="comp in element.components"
          class="ma-2"
          color="success"
          variant="outlined"
          rounded
          @click="add_module(comp)"
        >
          {{ comp.name }}
          <template v-slot:append>
            <v-badge inline :content="number_of_elements(comp.name)" />
          </template>
        </v-chip>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script setup lang="ts">
import { IAnalyserElement } from "@/models/AnalyserElement";
import { useModuleStore } from "@/stores/ModulesStore";
import { storeToRefs } from 'pinia';


const store = useModuleStore();
const { add_to_selected, number_of_elements } =store
const { defaultModules } =storeToRefs(store);

// add module element to list for tree build
const add_module = (element: IAnalyserElement) => {
  add_to_selected(element);
};
</script>

<style scoped></style>
