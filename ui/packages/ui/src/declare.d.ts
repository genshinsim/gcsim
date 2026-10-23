// Ambient declarations for non-code imports and `import.meta.env`. (This
// package has no `vite` dependency, so it cannot reference `vite/client`;
// declare the shapes it actually uses directly.)
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
