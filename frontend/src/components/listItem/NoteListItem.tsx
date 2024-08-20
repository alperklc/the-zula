import { BaseListItem } from "./BaseListItem";
import { Note } from "../../types/Api";
import TagsDisplay from "../tagsDisplay";
import TimeDisplay from "../timeDisplay";
import styles from "./index.module.css";

export const NoteListItem = ({ item }: { item: Note }) => {
    return <BaseListItem
        href={item.shortId!}
        title={item.title ?? ""}
        hasDraft={item.hasDraft}
        description={(item?.tags || []).length > 0 && (
            <TagsDisplay
                className={styles.tagsDisplay}
                tags={item?.tags ?? []}
                maxNumberOfTagsToDisplay={3}
            />
        )}
        sideInfo={
            <a className={styles.updatedAt} href={`/notes/${item.shortId}`}>
                <TimeDisplay isoDate={item.updatedAt ?? ""} />
            </a>
        } />
}
