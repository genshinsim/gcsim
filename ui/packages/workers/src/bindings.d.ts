export interface Env {
	API_ENDPOINT: string;
	ASSETS: Fetcher; //static-assets binding serving the SPA (../web/dist)
	GCSIM_WASM: R2Bucket; //bucket
	GCSIM_ASSETS: R2Bucket; //game-asset cache + static files
	GCSIM_SHARES: R2Bucket; //gzipped share results
	SHARE_KEYS: string; //JSON: { [keyid]: { class: "prod" | "dev", key: hex } }
	ASSET_SOURCE_HOSTS?: string; //JSON: per-type ordered mihoyo source hosts
	OG_PREVIEW_SCALE?: string; //OG card supersample factor, clamped [1,3], default 1.5
	OG_PREVIEW_CACHE_REV?: string; //OG edge-cache revision (integer); bump to force re-render, default 1
}
