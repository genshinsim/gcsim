import { createFileRoute } from "@tanstack/react-router";
import { DiscordCallback } from "../../Pages";

export const Route = createFileRoute("/auth/discord")({
	component: DiscordCallback,
});
