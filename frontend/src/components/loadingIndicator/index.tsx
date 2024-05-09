import styles from "./index.module.css";

const LoadingIndicator = (props: { className?: string }) => (
  <div className={props.className}>
    <span className={styles.spinner} />
  </div>
);

export default LoadingIndicator;
