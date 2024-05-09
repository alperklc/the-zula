
import { Note } from "../../types/Api";
import { Resource, ResourceType } from "../../types/resources";
import TagsDisplay from "../tagsDisplay";
import TimeDisplay from "../timeDisplay";
import styles from "./index.module.css";

const BaseListItem = ({ title, description, sideInfo }: { title: React.ReactNode, description: React.ReactNode, sideInfo: React.ReactNode }) => (
  <article className={styles.entry}>
    <div className={styles.content}>
      <span className={styles.title}>{title}</span>
      <span className={styles.rightSide}>{sideInfo}</span>
    </div>
    <div className={styles.description}>{description}</div>
  </article>
);

export const ListItem = ({ resourceType, item: initialItem }: { resourceType: ResourceType, item: Resource }) => {
  let item = initialItem

  switch (resourceType) {
    // add here other types later
    case "note":
      item = item as Note;
      return <BaseListItem
        title={item.title ?? ""}
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
}
