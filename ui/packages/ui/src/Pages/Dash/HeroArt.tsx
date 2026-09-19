import { cn } from "@gcsim/primitives";

type HeroArtProps = {
	src: string;
	alt?: string;
	className?: string;
};

export function HeroArt({ src, alt, className }: HeroArtProps) {
	return (
		<div className={cn("g-hero-art", className)}>
			<img src={src} alt={alt ?? ""} />
		</div>
	);
}
