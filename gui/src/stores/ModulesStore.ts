import { IAnalyserModule } from "@/models/AnalyserModule";
import { IAnalyserElement } from "@/models/AnalyserElement";
import { defineStore } from "pinia";
import availablemod from "@/data/AvailableModules.json";
import { computed, Ref, ref } from "vue";


export const useModuleStore = defineStore("module", () => {
  const dumpfile = ref(new File([""], "empty"));
  const defaultModules: Ref<IAnalyserModule[]> = ref(availablemod);
  const selectedModules: Ref<IAnalyserElement[]> = ref([]);

  const can_edit = computed(() => dumpfile.value.name != "empty");
  
  //  function to get nthe number of selected elements to display next to element name
  const number_of_elements = (name:string) => { 
    const lenghtof =selectedModules.value.filter(p => p.name == name).length;
    return lenghtof 
};
// has any selected elements
  const has_selected = computed(() => selectedModules.value.length > 0);

  // add a new element to the selection 
  const add_to_selected = (module: IAnalyserElement) => {
    selectedModules.value.push(module);
  };

  // clear all selected module element
  const clear_selected = () => {
    selectedModules.value = [];
  };

  // set entry file
  const setFile = (file: File) => {
    dumpfile.value = file
  }

  return {
    dumpfile,
    defaultModules,
    selectedModules,
    can_edit,
    has_selected,
    add_to_selected,
    clear_selected,
    setFile,
    number_of_elements
  };
});
