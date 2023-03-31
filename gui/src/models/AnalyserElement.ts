export interface IAnalyserElement {
  name: string;
  description?: string;
  has_result: boolean;
  share_input: boolean;
  type: string;
  output: string[];
  input: string[];
}
