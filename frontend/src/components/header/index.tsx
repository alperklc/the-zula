
import Logo from "../logo";
import styles from "./index.module.css";

const Header = (props: { className?: string; onMenuIconClicked: () => void; }) => (
  <header className={`${styles.container} ${props.className}`}>
    <span className={styles.menuAndLogo}>
      <button aria-label="menu toggle" onClick={props.onMenuIconClicked}>
        ☰
      </button>

      <a href={'/'} className={styles.logoContainer}>
        <Logo />
      </a>
    </span>
  </header>
);

export default Header;
