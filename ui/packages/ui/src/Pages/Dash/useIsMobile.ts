import { useEffect, useState } from "react";

// The Dash ships two distinct layouts (see the design handoff): a desktop hero
// and a dedicated mobile home, not one responsive block. This hook picks which
// one mounts so only the active layout fetches data. Breakpoint matches the
// Tailwind `md` (48rem) phone/tablet cutoff.
const MOBILE_QUERY = "(max-width: 47.999rem)";

export function useIsMobile(): boolean {
	const [isMobile, setIsMobile] = useState(
		() => typeof window !== "undefined" && window.matchMedia(MOBILE_QUERY).matches,
	);

	useEffect(() => {
		const mql = window.matchMedia(MOBILE_QUERY);
		const onChange = () => setIsMobile(mql.matches);
		onChange();
		mql.addEventListener("change", onChange);
		return () => mql.removeEventListener("change", onChange);
	}, []);

	return isMobile;
}
