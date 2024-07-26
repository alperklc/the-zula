import { Link } from "react-router-dom";
import Button from "../form/button";
import Icons from "../icons";
import Logo from "../logo";
import styles from "./index.module.css";

const Header = (props: { className?: string; onMenuIconClicked: () => void; }) => (
  <header className={`${styles.container} ${props.className}`}>
    <span className={styles.menuAndLogo}>
      <Button aria-label="menu toggle" onClick={props.onMenuIconClicked}>
        <Icons.Menu />
      </Button>

      <Link to={'/'} className={styles.logoContainer}>
        <Logo />
      </Link>
    </span>
  </header>
);

export default Header;
