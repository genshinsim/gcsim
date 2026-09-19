type HeroArtProps = {
	src: string;
	alt?: string;
	className?: string;
};

export function HeroArt({ src, alt, className }: HeroArtProps) {
	return (
		<div className={`g-hero-art${className ? ` ${className}` : ""}`}>
			<img src={src} alt={alt ?? ""} />
		</div>
	);
}
