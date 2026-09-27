import { ResultsView, type ResultsViewProps } from "@gcsim/components";
import { useEffect, useRef } from "react";
import { useLocation } from "react-router";

export default (props: ResultsViewProps) => {
	useScrollToLocation();
	return <ResultsView {...props} />;
};

function useScrollToLocation() {
	const scrolled = useRef(false);
	const { key, hash } = useLocation();
	const prevKey = useRef(key);

	useEffect(() => {
		if (hash == null) {
			return;
		}

		if (prevKey.current !== key) {
			prevKey.current = key;
			scrolled.current = false;
		}

		if (scrolled.current) {
			return;
		}
		const id = hash.replace("#", "");
		if (!id) {
			return;
		}
		const element = document.getElementById(id);
		if (element) {
			element.scrollIntoView({ behavior: "smooth" });
			scrolled.current = true;
		}
	});
}
