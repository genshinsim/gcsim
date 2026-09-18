import { Toaster as Sonner, type ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
	return (
		<Sonner
			className="toaster group"
			style={
				{
					"--normal-bg": "hsl(var(--deprecated-popover))",
					"--normal-text": "hsl(var(--deprecated-popover-foreground))",
					"--normal-border": "hsl(var(--deprecated-border))",
				} as React.CSSProperties
			}
			{...props}
		/>
	);
};

export { toast } from "sonner";
export { Toaster };
