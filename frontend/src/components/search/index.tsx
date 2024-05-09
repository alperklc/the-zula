import { useState, useEffect } from "react";
import useDebounce from "../../utils/useDebounce";

const SearchInput = (props: {
  autoFocus?: boolean;
  className?: string;
  query: string;
  onQueryUpdate: (_: string) => void;
  inputref?: React.RefObject<HTMLInputElement>
}) => {
  const [initialized, setInitialized] = useState<boolean>(false);
  const [keyword, setKeyword] = useState<string>(props.query);

  useEffect(() => {
    setInitialized(true);
  }, []);

  const debouncedInput = useDebounce(keyword, 500);

  useEffect(() => {
    if (initialized) {
      props.onQueryUpdate(keyword);
    }
  }, [debouncedInput]);

  return (
    <span className={props.className}>
      <input
        ref={props?.inputref}
        autoFocus={props.query ? props.autoFocus : false}
        onChange={(e) => setKeyword(e.target.value)}
        placeholder={"Type here to search"}
        type="text"
        value={keyword}
      />
    </span>
  );
};

export default SearchInput;
