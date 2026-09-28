import type { Sample } from "@gcsim/types";
import { saveAs } from "file-saver";
import Pako from "pako";

export function downloadSample(sample: Sample) {
	const out = Pako.deflate(JSON.stringify(sample));
	const blob = new Blob([out], { type: "application/base64" });
	saveAs(blob, "sample.gz");
}
