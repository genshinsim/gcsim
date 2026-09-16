import "i18next";
import type { resources } from "./index";

// react-i18next 17 no longer declares its own CustomTypeOptions; it reads
// i18next's TypeOptions, so the augmentation targets the i18next module.
declare module "i18next" {
	interface CustomTypeOptions {
		defaultNS: "translation";
		resources: (typeof resources)["en"];
	}
}
