import { Alert, AlertDescription, Button } from "@gcsim/primitives";
import { X } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";

type AlertVariant = "default" | "destructive" | "warning" | "success";

function variantForIntent(intent?: string): AlertVariant {
	switch (intent) {
		case "danger":
			return "destructive";
		case "warning":
			return "warning";
		case "success":
			return "success";
		default:
			return "default";
	}
}

type Props = {
	title: string;
	show: boolean;
	intent?: string;
	onDismiss?: (event: React.MouseEvent<HTMLElement>) => void;
	children: React.ReactNode;
};

export default ({ title, show, intent, onDismiss, children }: Props) => {
	return (
		<AnimatePresence>
			{show && (
				<motion.div exit={{ opacity: 0 }}>
					<Alert variant={variantForIntent(intent)}>
						<AlertDescription className="block text-current">
							<div className="flex justify-between">
								<h4 className="font-medium">{title}</h4>
								<Button
									variant="ghost"
									size="icon-sm"
									className="self-start"
									onClick={(e) => onDismiss?.(e)}
								>
									<X />
								</Button>
							</div>
							{children}
						</AlertDescription>
					</Alert>
				</motion.div>
			)}
		</AnimatePresence>
	);
};
