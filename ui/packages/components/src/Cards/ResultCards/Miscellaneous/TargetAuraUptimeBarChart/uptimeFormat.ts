export const uptimeFormat = (language: string) => (n?: number) => {
	if (n == null) {
		return undefined;
	}
	const num = n.toLocaleString(language, {
		minimumFractionDigits: 2,
		maximumFractionDigits: 2,
	});
	return `${num}%`;
};
