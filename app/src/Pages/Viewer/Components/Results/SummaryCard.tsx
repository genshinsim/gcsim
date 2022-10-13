import { Card, Colors, Dialog, Drawer, DrawerSize, Icon, Position, Classes } from "@blueprintjs/core";
import { Tooltip2 } from "@blueprintjs/popover2";
import classNames from "classnames";
import React from "react";
import { useTranslation } from "react-i18next";

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

// TODO: theme and size the overlay
export default ({
    title, color, value, label, auxStats, tooltip, children, drawerTitle }: CardProps) => {
  const [isOpen, setOpen] = React.useState(false);
  const interactable = children != undefined;

  return (
    <div className="flex basis-1/4 flex-auto pl-1 min-w-fit" style={{ background: color }}>
      <Card
          className="flex flex-auto flex-row items-stretch justify-between"
          interactive={interactable}
          onClick={() => value != undefined && setOpen(true)}>
        <div className="flex flex-col justify-start">
          <CardTitle title={title} tooltip={tooltip} />
          <CardValue value={value} label={label} />
          <CardAux aux={auxStats} />
        </div>
        <CardChevron interactable={interactable} />
      </Card>
      <CardDrawer title={drawerTitle} children={children} openState={[isOpen, setOpen]} />
    </div>
  );
};

const CardTitle = ({ title, tooltip }: { title: string, tooltip?: string | JSX.Element }) => {
  const helpIcon = tooltip == undefined ? null : <Icon icon="help" color={Colors.GRAY1} />;
  const out = (
    <div className="flex flex-row text-lg text-gray-400 items-center gap-3 outline-0">
      {title}
      {helpIcon}
    </div>
  );

  if (tooltip != undefined) {
    return (
      <div onClick={(e) => e.stopPropagation()}>
        <Tooltip2 content={tooltip}>{out}</Tooltip2>
      </div>
    );
  }
  return out;
}

const CardValue = ({ value, label }: { value?: number | string | null, label?: string }) => {
  const { i18n } = useTranslation();

  const out = value == undefined ? 1234 : value;
  const valueClass = classNames(
    "text-5xl font-bold tabular-nums",
    { "bp4-skeleton": value == undefined }
  );

  let lbl;
  if (label != undefined) {
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
}

const CardChevron = ({ interactable }: { interactable: boolean }) => {
  if (!interactable) {
    return null;
  }
  return (
    <div className="flex flex-grow justify-end self-stretch justify-self-end">
      <Icon icon="chevron-right" size={36} color={Colors.GRAY1} className="self-center" />
    </div>
  );
}

const CardAux = ({ aux }: { aux?: Array<AuxStat> }) => {
  if (aux === undefined) {
    return null;
  }

  return (
    <div className="grid grid-cols-3 gap-x-5 pt-1 justify-start text-sm font-mono text-gray-400 min-w-fit">
      {aux.map(e => <AuxItem key={e.title} stat={e} />)}
    </div>
  );
};

const AuxItem = ({ stat }: { stat: AuxStat }) => {
  const { i18n } = useTranslation();

  const cls = classNames(
    "font-black text-current text-sm !text-bp4-light-gray-500",
    { "bp4-skeleton": stat.value == undefined }
  );
  const val = stat.value == undefined ? 123.45 : stat.value;

  return (
    <div className="flex flex-row items-start gap-3">
      <div>{stat.title}</div>
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

  if (children === undefined) {
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
}