import { Button, Callout, Intent } from "@blueprintjs/core";
import { AnimatePresence, motion } from "framer-motion";

type Props = {
  title: string;
  show: boolean;
  intent?: Intent;
  onDismiss?: (event: React.MouseEvent<HTMLElement>) => void;
  children: React.ReactNode;
}

export default ({ title, show, intent, onDismiss, children }: Props) => {
  return (
    <AnimatePresence>
      {show && (
        <motion.div exit={{ opacity: 0 }}>
          <Callout intent={intent}>
            <div className="flex justify-between">
              <h4 className="bp4-heading">{title}</h4>
              <Button
                  icon="cross"
                  className="self-start"
                  minimal={true}
                  small={true}
                  onClick={(e: React.MouseEvent<HTMLElement>) => onDismiss != null && onDismiss(e)}
              />
            </div>
            {children}
          </Callout>
        </motion.div>
      )}
    </AnimatePresence>
  );
};