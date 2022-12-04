import { Card, Colors, Dialog, Icon, Classes } from "@blueprintjs/core";
import classNames from "classnames";
import { ReactNode, useState } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle } from "../../Util";

type AuxStat = {
  title: string;
  value?: string;
}

type CardProps = {
  title: string;
  color: string;
  value?: string;
  label?: string;
  auxStats?: Array<AuxStat>;
  tooltip?: string | JSX.Element;
  children?: React.ReactNode;
  drawerTitle?: string;
};

// TODO: hash link
export const RollupCard = ({
    title, color, value, label, auxStats, tooltip, children, drawerTitle }: CardProps) => {
  const [isOpen, setOpen] = useState(false);
  const interactable = children != undefined;

  return (
    <div className="flex basis-1/4 flex-auto pl-1 min-w-fit" style={{ background: color }}>
      <Card
          className="flex flex-auto flex-row items-stretch justify-between"
          interactive={interactable}
          onClick={() => interactable && value != undefined && setOpen(true)}>
        <div className="flex flex-col justify-start">
          <CardTitle title={title} tooltip={tooltip} />
          <CardValue value={value} label={label} />
          <CardAux aux={auxStats} />
        </div>
        <CardChevron interactable={interactable} />
      </Card>
      <CardDrawer title={drawerTitle} openState={[isOpen, setOpen]}>
        {children}
      </CardDrawer>
    </div>
  );
};

const CardValue = ({ value, label }: { value?: number | string | null, label?: string }) => {
  const { i18n } = useTranslation();

  const out = value == null ? 1234 : value;
  const valueClass = classNames(
    "text-5xl font-bold tabular-nums",
    { "bp4-skeleton": value == null }
  );

  let lbl: ReactNode;
  if (label != null) {
    lbl = <div className="flex items-start text-base text-gray-400">{label}</div>;
  }

  return (
    <div className="flex flex-row py-2 gap-1 justify-start">
      <div className={valueClass}>
        {out.toLocaleString(i18n.language)}
      </div>
      {lbl}
    </div>
  );
};

const CardChevron = ({ interactable }: { interactable: boolean }) => {
  if (!interactable) {
    return null;
  }
  return (
    <div className="flex flex-grow justify-end self-stretch justify-self-end">
      <Icon icon="chevron-right" size={36} color={Colors.GRAY1} className="self-center" />
    </div>
  );
};

const CardAux = ({ aux }: { aux?: Array<AuxStat> }) => {
  if (aux == null) {
    return null;
  }

  return (
    <div className="grid grid-cols-3 gap-x-5 pt-1 justify-start text-sm font-mono min-w-fit">
      {aux.map(e => <AuxItem key={e.title} stat={e} />)}
    </div>
  );
};

const AuxItem = ({ stat }: { stat: AuxStat }) => {
  const { i18n } = useTranslation();

  const cls = classNames(
    "font-black text-current text-sm text-bp4-light-gray-500",
    { "bp4-skeleton": stat.value == null }
  );
  const val = stat.value == null ? 123.45 : stat.value;

  return (
    <div className="flex flex-row items-start gap-3">
      <div className="text-gray-400">{stat.title}</div>
      <div className={cls}>
        {val.toLocaleString(i18n.language)}
      </div>
    </div>
  );
};

const CardDrawer = ({ title, children, openState }: {
      title?: string,
      children?: React.ReactNode,
      openState: [boolean, React.Dispatch<React.SetStateAction<boolean>>]
    }) => {
  const [isOpen, setOpen] = openState;

  if (children == null) {
    return null;
  }

  return (
    <Dialog
        isOpen={isOpen}
        onClose={() => setOpen(false)}
        title={title}
        icon="list-detail-view"
        autoFocus={true}
        canEscapeKeyClose={true}
        canOutsideClickClose={true}
        enforceFocus={true}
        hasBackdrop={true}
        usePortal={true}
        style={{ width: "720px" }}>
      <div className={Classes.DIALOG_BODY}>
        {children}
      </div>
    </Dialog>
  );
};