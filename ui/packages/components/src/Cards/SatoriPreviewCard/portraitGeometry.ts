// Portrait-internal layout, in the card's LOGICAL px. Single source of truth for
// BOTH the layered Portrait render (Portraits.tsx, browser/Storybook) and the
// Photon compositor (portraitCompositor.ts, edge Worker) — the two must place
// every layer identically or the composited card drifts from the live preview.
// Values mirror the live AvatarPortrait.

export const PORTRAIT_W = 127;
export const PORTRAIT_H = 106;

// Avatar: square, horizontally centred, top-aligned under the card's pt-2.
export const AVATAR_SIZE = 96;
export const AVATAR_MARGIN_TOP = 8;

// Weapon: square, bottom-right. `right:-4` overhangs the right edge by 4px (the
// portrait's overflow:hidden clips it) — matching the CSS, so a negative value
// means "further right".
export const WEAPON_SIZE = 55;
export const WEAPON_BOTTOM = 4;
export const WEAPON_RIGHT = -4;

// Artifact icon: square (or a 17.5-wide half for the two-set/lone-2pc slice),
// bottom-left.
export const ARTIFACT_SIZE = 35;
export const ARTIFACT_HALF = 17.5;
export const ARTIFACT_BOTTOM = 1;
export const ARTIFACT_LEFT = 1;

// Weapon and artifact wrappers are opacity-85; the empty-slot placeholder is
// opacity-50.
export const ICON_OPACITY = 0.85;
export const PLACEHOLDER_OPACITY = 0.5;

// Incomplete-build "WIP" bar, vertically at the top third of the portrait.
export const WIP_TOP = "33%";
