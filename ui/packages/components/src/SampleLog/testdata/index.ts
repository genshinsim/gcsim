import type { Sample } from "@gcsim/types";
import fixture from "./sample.json";

// The first 300 frames of a real sample; proto JSON omits zero values such as
// `cons: 0`, which the hand-written Sample type doesn't allow for.
export const sampleFixture = fixture as unknown as Sample;
