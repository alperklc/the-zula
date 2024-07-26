import { useNavigate } from "react-router-dom";
import styles from "./index.module.css";

export const BaseListItem = ({ title, description, sideInfo, href }: { title: React.ReactNode, description: React.ReactNode, sideInfo: React.ReactNode, href: string }) => {
    const navigate = useNavigate()

    return <article className={styles.entry} onClick={() => navigate(href)}>
        <div className={styles.content}>
            <span className={styles.title}>{title}</span>
            <span className={styles.rightSide}>{sideInfo}</span>
        </div>
        <div className={styles.description}>{description}</div>
    </article>
}
