// Ambient module declarations for non-code, side-effect-only imports that ship
// no type declarations. (This package has no `vite` dependency, so it cannot
// reference `vite/client`; declare the shapes it actually uses directly.)
declare module "*.png";
declare module "*.jpg";
declare module "*.css";
declare module "@fontsource/fira-mono";
