export interface UserInfo {
  uid: string; //discord id
  name: string; //discord tag
  role: number; // role number
  permalinks: string[]; // list of permas
  data: UserData;
}

export interface UserData {
  settings: UserSettings;
}

export interface UserSettings {
  showTips: boolean;
  showBuilder: boolean;
  showNameSearch: boolean;
}
