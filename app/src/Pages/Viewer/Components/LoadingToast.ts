import { Intent, Toaster } from "@blueprintjs/core";
import React from "react";
import { useLocation } from "wouter";

export const LoadingToaster = Toaster.create({
  maxToasts: 1,
  autoFocus: true,
  canEscapeKeyClear: false,
});

export default function useLoadingToast(data: any) {
  const [loaded, setLoaded] = React.useState(data != null);
  const [_, setLocation] = useLocation();

  // abort loading if we leave the page
  React.useEffect(() => {
    return () => {
      if (data == null) {
        LoadingToaster.clear();
        setLoaded(false);
      }
    };
  }, []);
 
  React.useEffect(() => {
    if (data === null) {
      LoadingToaster.show({
        message: 'Loading...',
        icon: "refresh",
        intent: Intent.PRIMARY,
        timeout: 0,
        isCloseButtonShown: false,
        action: { onClick: () => setLocation("/viewer"), text: 'Abort' }
      });
    } else {
      LoadingToaster.clear();
      setLoaded(true);
    }
  }, [data]);

  return loaded;
}