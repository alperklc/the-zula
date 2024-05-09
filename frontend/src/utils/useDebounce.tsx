import React from "react";

function useDebounce<T>(value: T, timeout: number) {
  const [state, setState] = React.useState<T>(value);

  React.useEffect(() => {
    const handler = setTimeout(() => setState(value), timeout);

    return () => clearTimeout(handler);
  }, [value, timeout]);

  return state;
}

export default useDebounce;
