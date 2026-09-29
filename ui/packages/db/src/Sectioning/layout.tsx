import Nav from "./Nav";

export default function Layout({ children }: { children: React.ReactNode }) {
	return (
		<>
			<Nav />
			{children}
		</>
	);
}
