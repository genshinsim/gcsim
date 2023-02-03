export {};

declare global {
  const JWT_SECRET: string;
  const DISCORD_ID: string;
  const DISCORD_SECRET: string;
  const POSTGREST_ENDPOINT: string; //secret url for POSTGREST backend
  const PREVIEW_ENDPOINT: string; //secret url for embed server
  const ASSETS_ENDPOINT: string; //secret url for assets server
  const USER_TOKENS: KVNamespace;
  const GCSIM_ASSETS: R2Bucket; //bucket
}
