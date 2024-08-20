import React from "react";
import { useNavigate } from "react-router-dom";
import { Bookmark } from "../../types/Api";
import TagsDisplay from "../tagsDisplay";
import TimeDisplay from "../timeDisplay";
import Icons from "../icons";
import OverflowMenu from "../bookmarkOverflowMenu";
import { useUI } from "../../contexts/uiContext";
import { styles as layoutStyles } from '../layout'
import styles from "./index.module.css";

const FaviconDisplay = (props: { faviconUrl: string }) => {
    const [faviconUrl, setFaviconUrl] = React.useState(props.faviconUrl)
    if (!faviconUrl) {
        return <Icons.Link className={styles.icon} />
    }

    return (
        <img
            className={styles.icon}
            src={faviconUrl}
            onError={(e: any) => {
                e.target.onerror = null
                setFaviconUrl('')
            }}
        />
    )
}

export const BookmarkListItem = ({ item }: { item: Bookmark }) => {
    const navigate = useNavigate()
    const { overflowMenuId, setOverflowMenu } = useUI()

    const openOverflowMenuFor = (rowUid: string) => (event: React.MouseEvent) => {
        setOverflowMenu(rowUid)
        event.stopPropagation()
    }

    const closeOverflowMenu = () => setOverflowMenu('')

    return (
        <article className={`${styles.entry} ${styles.bookmarkEntry}`} onClick={() => navigate(`/bookmarks/${item.shortId}`)}>
            <section className={styles.bookmarkContent}>
                <div className={styles.iconWrapper}>
                    <FaviconDisplay faviconUrl={item.faviconUrl!} />
                </div>
                <div className={styles.textContent}>
                    {item.title && <div>{item.title}</div>}
                    <div className={styles.secondaryText}>
                        <a className={styles.url} href={item.url} target='_blank' rel='noreferrer'>
                            {item.url}
                        </a>
                    </div>
                    {(item?.tags || []).length > 0 && (
                        <TagsDisplay
                            className={styles.tagsDisplay}
                            tags={item.tags!}
                            maxNumberOfTagsToDisplay={3}
                        />
                    )}
                </div>
            </section>
            <div className={styles.rightSide}>
                <a className={styles.permalink} href={`/bookmarks/${item.shortId}`}>
                    <TimeDisplay isoDate={item.updatedAt} />
                </a>
                <div className={styles.button}>
                    <Icons.ThreeDotsMenu width={16} height={16} onClick={openOverflowMenuFor(item.shortId)} />
                    {overflowMenuId === item.shortId && (
                        <div className={layoutStyles.overflowMenuContainer}>
                            <OverflowMenu
                                link={item?.url}
                                title={item?.title}
                                onOptionClick={closeOverflowMenu}
                                className={styles.overflowMenu}
                            />
                        </div>
                    )}
                </div>
            </div>
        </article>
    )
}

