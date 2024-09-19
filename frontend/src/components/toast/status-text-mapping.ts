export enum ResourceType {
  USER = 'USER',
  NOTE = 'NOTE',
  BOOKMARK = 'BOOKMARK',
}

export enum Action {
  CREATE = 'CREATE',
  UPDATE = 'UPDATE',
  DELETE = 'DELETE'
}

export const toStatusTextKey = (resource: ResourceType, action: Action): string => {
  if (resource === ResourceType.USER) {
    const statusLookup: Map<Action, string> = new Map([
      [Action.UPDATE, 'messages.user.update'],
    ])

    return statusLookup.get(action) || ''
  }

  if (resource === ResourceType.NOTE) {
    const statusLookup: Map<Action, string> = new Map([
      [Action.CREATE, 'messages.notes.create'],
      [Action.UPDATE, 'messages.notes.update'],
      [Action.DELETE, 'messages.notes.delete'],
    ])

    return statusLookup.get(action) || ''
  }

  if (resource === ResourceType.BOOKMARK) {
    const statusLookup: Map<Action, string> = new Map([
      [Action.CREATE, 'messages.bookmarks.create'],
      [Action.UPDATE, 'messages.bookmarks.update'],
      [Action.DELETE, 'messages.bookmarks.delete'],
    ])

    return statusLookup.get(action) || ''
  }

  return ''
}
