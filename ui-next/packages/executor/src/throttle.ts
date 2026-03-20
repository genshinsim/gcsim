/**
 * Simple throttle implementation that limits function execution to at most
 * once per `wait` milliseconds. Calls on the leading and trailing edge.
 */
export function throttle<T extends (...args: unknown[]) => void>(fn: T, wait: number): T {
  let lastCallTime = 0;
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  const throttled = (...args: unknown[]) => {
    const now = Date.now();
    const remaining = wait - (now - lastCallTime);

    if (remaining <= 0) {
      // Leading edge: enough time has passed, execute immediately
      if (timeoutId !== null) {
        clearTimeout(timeoutId);
        timeoutId = null;
      }
      lastCallTime = now;
      fn(...args);
    } else if (timeoutId === null) {
      // Trailing edge: schedule execution for when the wait period ends
      timeoutId = setTimeout(() => {
        lastCallTime = Date.now();
        timeoutId = null;
        fn(...args);
      }, remaining);
    }
  };

  return throttled as T;
}
