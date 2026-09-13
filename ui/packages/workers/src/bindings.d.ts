export interface Env {
	API_ENDPOINT: string;
	ASSETS_ENDPOINT: string;
	ASSETS: Fetcher; //static-assets binding serving the SPA (../web/dist)
	GCSIM_WASM: R2Bucket; //bucket
}
