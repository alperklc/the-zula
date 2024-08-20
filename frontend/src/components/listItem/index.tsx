
import { Bookmark, Note, NoteChange, UserActivity } from "../../types/Api";
import { Resource, ResourceType } from "../../types/resources";
import { ActivityLogListItem } from "./ActivityLogListItem";
import { BookmarkListItem } from "./BookmarkListItem";
import { NoteChangeListItem } from "./NoteChangeListItem";
import { NoteListItem } from "./NoteListItem";

export const ListItem = ({ resourceType, item }: { resourceType: ResourceType, item: Resource }) => {
  switch (resourceType) {
    case "note":
      return <NoteListItem item={item as Note} />
    case "note-change":
        return <NoteChangeListItem item={item as NoteChange} />
    case "bookmark":
      return <BookmarkListItem item={item as Bookmark} />
    case "user-activity":
        return <ActivityLogListItem item={item as UserActivity} />
  }
}
