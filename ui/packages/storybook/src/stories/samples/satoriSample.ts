import type { model } from "@gcsim/types";
import { cloneDeep } from "lodash-es";
// Import the JSON directly (not via ./index, which uses Vite's ?raw suffix) so
// this module also loads under Node/tsx in the fixture generator script.
import sampleResult from "./sampleResult.json";

// Scale the baked portrait fixture is composited at. The production OG render
// uses 3x; the fixture uses 2x purely to keep the committed data smaller — the
// slice/outline geometry is identical, only the pixel density differs.
export const SATORI_SAMPLE_SCALE = 2;

// A sample tuned to exercise the portrait compositor's artifact cases in one
// card, so the "Composited" story is a visual check of the outline + slice:
//   - bennett:   2pc + 2pc  (left half of set A joined to the right half of B)
//   - xiangling: lone 2pc   (left-half slice)
//   - xingqiu:   single 4pc (full square flower)
//   - zhongli:   single 5pc (full square flower)
export const satoriSample: model.SimulationResult = (() => {
	const s = cloneDeep(sampleResult) as unknown as model.SimulationResult;
	const chars = s.character_details ?? [];
	if (chars[0]) {
		chars[0].sets = { crimsonwitchofflames: 2, emblemofseveredfate: 2 };
	}
	if (chars[1]) {
		chars[1].sets = { noblesseoblige: 2 };
	}
	return s;
})();
