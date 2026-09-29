import { useEffect, useState } from "react";

const MOBILE_QUERY = "(max-width: 47.999rem)";

export function useIsMobile(): boolean {
	const [isMobile, setIsMobile] = useState(
		() =>
			typeof window !== "undefined" && window.matchMedia(MOBILE_QUERY).matches,
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
