export interface UserInfo {
  user_id: number; //discord id
  user_name: string; //discord tag
  token?: string; //jwt token
  settings?: UserSettings;
}

export interface UserSettings {
  showTips: boolean;
  showBuilder: boolean;
}
