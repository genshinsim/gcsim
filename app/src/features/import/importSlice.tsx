import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ArtifactParam {
  name: string;
  value: number;
}

export interface Artifact {
  position: string;
  setName: string;
  mainTag: ArtifactParam;
  normalTags: ArtifactParam[];
  comments?: string;
}

export interface DataSet {
  [key: string]: Artifact[];
}

interface ImportState {
  data: DataSet | null;
}

const initialState: ImportState = {
  data: null,
};

export const importSlice = createSlice({
  name: "import",
  initialState,
  reducers: {
    setData: (state, action: PayloadAction<DataSet>) => {
      state.data = action.payload;
    },
  },
});

export const { setData } = importSlice.actions;
export default importSlice.reducer;
