import { Button, ButtonGroup, Checkbox, Label } from "@gcsim/primitives";
import axios from "axios";
import { LogOut, Save } from "lucide-react";
import { Viewport } from "../../Components";
import { user as userStore, useStore } from "../../Stores";
import { authProvider, Login } from "./Login";

function saveSettings() {
	axios
		.post("/api/user/save", userStore.get().data)
		.then(() => console.log("save ok"))
		.catch(() => console.log("save failed"));
}

function logout() {
	authProvider
		.logout()
		.catch((err) => console.warn("Error occured logging out: ", err))
		.finally(userStore.reset);
}

export function PageUserAccount() {
	const user = useStore(userStore);

	if (user.uid === "") {
		return <Login />;
	}

	return (
		<Viewport>
			<div className="flex flex-col ">
				<div className="flex flex-col gap-2">
					<div className="flex items-center gap-2">
						<Checkbox
							id="show-tips"
							checked={user.data.settings.showTips}
							onCheckedChange={() => {
								userStore.setSettings({
									...user.data.settings,
									showTips: !user.data.settings.showTips,
								});
							}}
						/>
						<Label htmlFor="show-tips">Show tips</Label>
					</div>
					<div className="flex items-center gap-2">
						<Checkbox
							id="show-builder"
							checked={user.data.settings.showBuilder}
							onCheckedChange={() => {
								userStore.setSettings({
									...user.data.settings,
									showBuilder: !user.data.settings.showBuilder,
								});
							}}
						/>
						<Label htmlFor="show-builder">Show builder</Label>
					</div>
				</div>
				<div className="flex flex-row place-content-center mt-2">
					<ButtonGroup>
						<Button variant="secondary" size="lg" onClick={saveSettings}>
							<Save />
							Save Settings
						</Button>
						<Button variant="destructive" size="lg" onClick={logout}>
							<LogOut />
							Logout
						</Button>
					</ButtonGroup>
				</div>
			</div>
		</Viewport>
	);
}
