import React from "react";

export const debounce = (fn: (..._: unknown[]) => void, time: number) => {
  let timeoutId: NodeJS.Timeout | number | null

  return wrapper

  function wrapper(...args: unknown[]) {
    if (timeoutId) {
      clearTimeout(timeoutId as number)
    }

    timeoutId = setTimeout(() => {
      timeoutId = null

      fn(...args)
    }, time)
  }
}

function useDebounce<T>(value: T, timeout: number) {
  const [state, setState] = React.useState<T>(value);

  React.useEffect(() => {
    const handler = setTimeout(() => setState(value), timeout);

    return () => clearTimeout(handler);
  }, [value, timeout]);

  return state;
}

export default useDebounce;
