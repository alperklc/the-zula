import { Action, ResourceType } from '../../../types/enums'

export const toStatusTextKey = (resource: ResourceType, action: Action): string => {
  if (resource === ResourceType.USER) {
    const statusLookup: Map<Action, string> = new Map([
      [Action.CREATE, 'messages.user.create'],
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

  if (resource === ResourceType.FILE) {
    const statusLookup: Map<Action, string> = new Map([
      [Action.CREATE, 'messages.files.create'],
      [Action.UPDATE, 'messages.files.update'],
      [Action.MOVE, 'messages.files.move'],
      [Action.DELETE, 'messages.files.delete'],
    ])

    return statusLookup.get(action) || ''
  }

  if (resource === ResourceType.FOLDER) {
    const statusLookup: Map<Action, string> = new Map([
      [Action.CREATE, 'messages.folders.create'],
      [Action.UPDATE, 'messages.folders.update'],
      [Action.MOVE, 'messages.folders.move'],
      [Action.DELETE, 'messages.folders.delete'],
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
