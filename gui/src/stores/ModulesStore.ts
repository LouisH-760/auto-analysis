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
  
  const number_of_elements = (name:string) => { 
    const lenghtof =selectedModules.value.filter(p => p.name == name).length;
    return lenghtof 
};
  const has_selected = computed(() => selectedModules.value.length > 0);

  const add_to_selected = (module: IAnalyserElement) => {
    selectedModules.value.push(module);
  };

  const clear_selected = () => {
    selectedModules.value = [];
  };

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
