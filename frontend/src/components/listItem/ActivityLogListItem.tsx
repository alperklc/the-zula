import { useIntl } from "react-intl";
import { BaseListItem } from "./BaseListItem";
import { UserActivity } from "../../types/Api";
import TimeDisplay from "../timeDisplay";
import styles from "./index.module.css";

const mapActivityLogText = (resourceType: string, action: string) => {
    if (resourceType === 'BOOKMARK') {
        const logMessageLookup: Map<string, string> = new Map([
            ['CREATE', 'activity-log.messages.bookmark_create'],
            ['READ', 'activity-log.messages.bookmark_read'],
            ['UPDATE', 'activity-log.messages.bookmark_update'],
            ['DELETE', 'activity-log.messages.bookmark_delete'],
        ])

        return logMessageLookup.get(action) || ''
    }

    if (resourceType === 'NOTE') {
        const logMessageLookup: Map<string, string> = new Map([
            ['CREATE', 'activity-log.messages.note_create'],
            ['READ', 'activity-log.messages.note_read'],
            ['UPDATE', 'activity-log.messages.note_update'],
            ['DELETE', 'activity-log.messages.note_delete'],
        ])
        return logMessageLookup.get(action) || ''
    }

    if (resourceType === 'USER') {
        const logMessageLookup: Map<string, string> = new Map([
            ['CREATE', 'activity-log.messages.user_create'],
            ['UPDATE', 'activity-log.messages.user_update'],
        ])
        return logMessageLookup.get(action) || ''
    }
}


export const ActivityLogListItem = ({ item }: { item: UserActivity }) => {
    const intl = useIntl()

    const getLink = (resourceType: string, action: string, uid: string) => {
        if (action === 'DELETE') {
            return null
        }

        if (resourceType === 'BOOKMARK') {
            return `/bookmarks/${uid}`
        }

        if (resourceType === 'NOTE') {
            return `/notes/${uid}`
        }

        return null
    }

    const translationKey = mapActivityLogText(item.resourceType!, item.action!)

    return (
        <BaseListItem
            href={getLink(item.resourceType!, item.action!, item.objectId!) || '#'}
            title={intl.formatMessage({ id: translationKey })}
            description={intl.formatMessage({ id: translationKey })}
            sideInfo={
                <a className={styles.updatedAt} href={getLink(item.resourceType!, item.action!, item.objectId!) || '#'}>
                    <TimeDisplay isoDate={item.timestamp ?? ""} />
                </a>
            } />
    )
}
