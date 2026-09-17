import { Button, ButtonGroup, Checkbox, Label } from "@gcsim/primitives";
import { LogOut, Save } from "lucide-react";
import { Viewport } from "../../Components";
import {
	type AppThunk,
	useAppDispatch,
	useAppSelector,
} from "../../Stores/store";
import {
	initialState,
	saveUserSettings,
	userActions,
} from "../../Stores/userSlice";
import { authProvider, Login } from "./Login";

//thunks
function logout(): AppThunk {
	return (dispatch) => {
		authProvider
			.logout()
			.then(() => dispatch(userActions.setUser(initialState)))
			.catch((err) => {
				//log out the user
				console.warn("Error occured logging out: ", err);
				dispatch(userActions.setUser(initialState));
			});
	};
}

export function PageUserAccount() {
	const user = useAppSelector((state) => state.user);
	const dispatch = useAppDispatch();

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
								dispatch(
									userActions.setUserSettings({
										showTips: !user.data.settings.showTips,
										showBuilder: user.data.settings.showBuilder,
										showNameSearch: user.data.settings.showNameSearch,
									}),
								);
							}}
						/>
						<Label htmlFor="show-tips">Show tips</Label>
					</div>
					<div className="flex items-center gap-2">
						<Checkbox
							id="show-builder"
							checked={user.data.settings.showBuilder}
							onCheckedChange={() => {
								dispatch(
									userActions.setUserSettings({
										showTips: user.data.settings.showTips,
										showBuilder: !user.data.settings.showBuilder,
										showNameSearch: user.data.settings.showNameSearch,
									}),
								);
							}}
						/>
						<Label htmlFor="show-builder">Show builder</Label>
					</div>
				</div>
				<div className="flex flex-row place-content-center mt-2">
					<ButtonGroup>
						<Button
							variant="secondary"
							size="lg"
							onClick={() => dispatch(saveUserSettings())}
						>
							<Save />
							Save Settings
						</Button>
						<Button
							variant="destructive"
							size="lg"
							onClick={() => dispatch(logout())}
						>
							<LogOut />
							Logout
						</Button>
					</ButtonGroup>
				</div>
			</div>
		</Viewport>
	);
}
