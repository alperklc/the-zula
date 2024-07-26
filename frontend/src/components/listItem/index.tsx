
import { Bookmark, Note } from "../../types/Api";
import { Resource, ResourceType } from "../../types/resources";
import { BookmarkListItem } from "./BookmarkListItem";
import { NoteListItem } from "./NoteListItem";

export const ListItem = ({ resourceType, item }: { resourceType: ResourceType, item: Resource }) => {
  switch (resourceType) {
    case "note":
      return <NoteListItem item={item as Note} />
    case "bookmark":
      return <BookmarkListItem item={item as Bookmark} />

  }
}
