import { IAnalyserElement } from "./AnalyserElement";

// basic definition of a groupe of modules based on the structure of the modules folder.
export interface IAnalyserModule {
  name: string;
  components: IAnalyserElement[];
}
