import { FormattedMessage } from "react-intl";
import Logo from "../../components/logo";
import { useAuth } from "../../contexts/authContext";
import styles from "./index.module.css";

function parseUrlParams() {
  const query = document.location.search.substr(1);
  const result: Record<string, string> = {};
  query.split("&").forEach(function (part) {
    const item = part.split("=");
    result[item[0]] = decodeURIComponent(item[1]);
  });

  return result;
}

const LoginCheck = () => {
  const auth = useAuth();
  const urlParams = parseUrlParams()

  return (<div className={styles.container}>
    <div className={styles.logoContainer}>
      <div className={styles.logo}>
        <Logo />
      </div>
    </div>
    <div className={styles.headline}>
      <FormattedMessage id='login_page.headline' />
    </div>
    <div className={styles.authForm}>
      <button className={styles.button} onClick={() => auth.login(urlParams["redirectAfterLogin"])}>
        <FormattedMessage id='login_page.buttons.submit' />
      </button>
    </div>
  </div>);
};

export default LoginCheck;
