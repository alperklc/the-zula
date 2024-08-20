import { BaseListItem } from "./BaseListItem";
import { NoteChange } from "../../types/Api";
import TimeDisplay from "../timeDisplay";
import styles from "./index.module.css";

export const NoteChangeListItem = ({ item }: { item: NoteChange }) => {
    return <BaseListItem
        href={item.shortId!}
        title={item.updatedAt ?? ""}
        description={item.updatedBy}
        sideInfo={
            <a className={styles.updatedAt} href={`/notes/${item.shortId}/changes/${item.shortId}`}>
                <TimeDisplay isoDate={item.updatedAt ?? ""} />
            </a>
        } />
}
