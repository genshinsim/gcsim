export interface Env {
	API_ENDPOINT: string;
	ASSETS: Fetcher; //static-assets binding serving the SPA (../web/dist)
	GCSIM_WASM: R2Bucket; //bucket
	GCSIM_ASSETS: R2Bucket; //game-asset cache + static files
	ASSET_SOURCE_HOSTS?: string; //JSON: per-type ordered mihoyo source hosts
}
