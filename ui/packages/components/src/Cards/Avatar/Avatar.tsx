import type { ReactNode } from "react";
import { charBG } from "../../lib/helper";
import { cn } from "../../lib/utils";

type AvatarProps = {
	name: string;
	element?: string;
	background?: boolean;
	overlay?: boolean;
	onImageLoaded?: () => void;

	className?: string;
	imageWrapClassName?: string;
	imageClassName?: string;

	children?: ReactNode;
};

export const Avatar = ({
	name,
	element = "",
	background = true,
	overlay = false,
	onImageLoaded,
	className = "",
	imageWrapClassName = "flex justify-center",
	imageClassName = "relative object-contain h-24",
	children,
}: AvatarProps) => {
	return (
		<div className={cn("relative", background && charBG(element), className)}>
			{overlay ? (
				<div
					className="absolute top-0 left-0 right-0 bottom-0 !bg-cover !bg-center mix-blend-luminosity"
					style={{ background: `url(/api/assets/misc/overlay.jpg)` }}
				></div>
			) : null}

			<div className={imageWrapClassName}>
				<img
					className={imageClassName}
					key={name}
					alt={name}
					src={`/api/assets/avatar/${name}.png`}
					onError={(e) => {
						(e.target as HTMLImageElement).src = "/api/assets/misc/default.png";
						onImageLoaded?.();
					}}
					onLoad={onImageLoaded}
				/>
			</div>

			{children}
		</div>
	);
};
