declare module "*.png";
declare module "*.jpg";
declare module "*.css";
declare module "@fontsource/fira-mono";

interface ImportMetaEnv {
	readonly VITE_HIDE_KQM_DB?: string;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
