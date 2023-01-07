export {};

declare global {
  const API_ENDPOINT: string;
  const PREVIEW_ENDPOINT: string;
  const GCSIM_ASSETS: R2Bucket; //bucket
  const GCSIM_WASM: R2Bucket; //bucket
}
