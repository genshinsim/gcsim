import { DashDesktop } from "./DashDesktop";
import { DashMobile } from "./DashMobile";
import { useIsMobile } from "./useIsMobile";

export function Dash() {
	return useIsMobile() ? <DashMobile /> : <DashDesktop />;
}
