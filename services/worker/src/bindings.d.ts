export {};

declare global {
  const JWT_SECRET: string;
  const DISCORD_ID: string;
  const DISCORD_SECRET: string;
  const POSTGREST_ENDPOINT: string;
  const PREVIEW_ENDPOINT: string;
  const USER_TOKENS: KVNamespace;
}
