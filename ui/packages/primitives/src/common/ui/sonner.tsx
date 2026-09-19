import { Toaster as Sonner, type ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
	return (
		<Sonner
			className="toaster group"
			style={
				{
					"--normal-bg": "var(--g-surface)",
					"--normal-text": "var(--g-text)",
					"--normal-border": "var(--g-border)",
				} as React.CSSProperties
			}
			{...props}
		/>
	);
};

export { toast } from "sonner";
export { Toaster };
