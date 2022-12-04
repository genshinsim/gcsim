import { Colors } from "@blueprintjs/core";
import { scaleOrdinal } from "@visx/scale";

function safeGet(colors: string[], i: number) {
  return colors[i % colors.length];
}

type ElementColor = {
  label: string;
  value: string;
}

const elements: Map<string, ElementColor> = new Map([
  ["electro", { label: Colors.VIOLET4, value: Colors.VIOLET3 }],
  ["pyro", { label: Colors.VERMILION4, value: Colors.VERMILION3 }],
  ["cryo", { label: "#95CACB", value: "#4B8DAA" }],
  ["hydro", { label: Colors.CERULEAN4, value: Colors.CERULEAN3 }],
  ["dendro",{ label: Colors.FOREST4, value: Colors.FOREST3 }],
  ["anemo",{ label: Colors.TURQUOISE4, value: Colors.TURQUOISE3 }],
  ["geo", { label: Colors.GOLD4, value: Colors.GOLD3 }],
  ["physical",{ label: Colors.SEPIA4, value: Colors.SEPIA3 }],

  // not possible, but defined in attributes/element.go so here just in case
  ["frozen", { label: "#000", value: "#000" }],
  ["quicken", { label: "#FFF", value: "#FFF" }],
]);

const elementColor = scaleOrdinal<string, string>({
  domain: Array.from(elements.keys()),
  range: Array.from(elements.values()).map(e => e.value),
});

const elementLabelColor = scaleOrdinal<string, string>({
  domain: Array.from(elements.keys()),
  range: Array.from(elements.values()).map(e => e.label),
});


// Qualitative follows a specific order defined by bp4 to maximize distinction
const qualitative1 = [
  Colors.CERULEAN1,
  Colors.FOREST1,
  Colors.GOLD1,
  Colors.VERMILION1,
  Colors.VIOLET1,
  Colors.TURQUOISE1,
  Colors.ROSE1,
  Colors.LIME1,
  Colors.SEPIA1,
  Colors.INDIGO1,
];

const qualitative2 = [
  Colors.CERULEAN2,
  Colors.FOREST2,
  Colors.GOLD2,
  Colors.VERMILION2,
  Colors.VIOLET2,
  Colors.TURQUOISE2,
  Colors.ROSE2,
  Colors.LIME2,
  Colors.SEPIA2,
  Colors.INDIGO2,
];

const qualitative3 = [
  Colors.CERULEAN3,
  Colors.FOREST3,
  Colors.GOLD3,
  Colors.VERMILION3,
  Colors.VIOLET3,
  Colors.TURQUOISE3,
  Colors.ROSE3,
  Colors.LIME3,
  Colors.SEPIA3,
  Colors.INDIGO3,
];

const qualitative4 = [
  Colors.CERULEAN4,
  Colors.FOREST4,
  Colors.GOLD4,
  Colors.VERMILION4,
  Colors.VIOLET4,
  Colors.TURQUOISE4,
  Colors.ROSE4,
  Colors.LIME4,
  Colors.SEPIA4,
  Colors.INDIGO4,
];

const qualitative5 = [
  Colors.CERULEAN5,
  Colors.FOREST5,
  Colors.GOLD5,
  Colors.VERMILION5,
  Colors.VIOLET5,
  Colors.TURQUOISE5,
  Colors.ROSE5,
  Colors.LIME5,
  Colors.SEPIA5,
  Colors.INDIGO5,
];

export const DataColors = {
  element: elementColor,
  elementLabel: elementLabelColor,
  
  // TODO: better colors for characters?
  character: (i: number) => qualitative3[i],
  characterLabel: (i: number) => qualitative4[i],

  qualitative1: (i: number) => safeGet(qualitative1, i),
  qualitative2: (i: number) => safeGet(qualitative2, i),
  qualitative3: (i: number) => safeGet(qualitative3, i),
  qualitative4: (i: number) => safeGet(qualitative4, i),
  qualitative5: (i: number) => safeGet(qualitative5, i),
};