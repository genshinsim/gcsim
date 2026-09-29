import { Toaster } from "@gcsim/primitives";
import { Footer } from "@gcsim/ui/src/features/layout";
import Nav from "./Nav";

export default function Layout({ children }: { children: React.ReactNode }) {
	return (
		<>
			<Toaster position="top-right" theme="dark" />
			<Nav />
			{children}
			<Footer />
		</>
	);
}
