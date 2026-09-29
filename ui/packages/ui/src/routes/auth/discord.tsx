import { createFileRoute } from "@tanstack/react-router";
import { DiscordCallback } from "../../features/user/DiscordCallback";

export const Route = createFileRoute("/auth/discord")({
	component: DiscordCallback,
});
