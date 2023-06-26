import { Item } from "./Item";

type Props = {
  signKey?: string;
}

export const DevBuild = ({ signKey }: Props) => {
  if (signKey == null || signKey == "prod") {
    return null;
  }

  if (signKey != "dev") {
    return <Item value="unsigned" intent="danger" bright bold />;
  }
  return <Item value="dev build" intent="danger" bright bold />;
};