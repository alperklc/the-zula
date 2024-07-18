import React from "react";
import { useAuth } from "../../contexts/authContext";

function parseUrlParams() {
  const query = document.location.search.substr(1);
  const result: Record<string, string> = {};
  query.split("&").forEach(function(part) {
    const item = part.split("=");
    result[item[0]] = decodeURIComponent(item[1]);
  });

  return result;
}

const LoginCheck = () => {
  const auth = useAuth();
  const urlParams = parseUrlParams()

  React.useEffect(() => {
    auth.login(urlParams["redirectAfterLogin"])
  }, [auth, urlParams])

  return <>Redirecting...</>;
};

export default LoginCheck;