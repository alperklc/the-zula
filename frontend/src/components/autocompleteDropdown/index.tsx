import React from "react";

import styles from "./index.module.css";

export const AutocompleteDropdown = ({ listItems, handleFoundItemClick, onItemHighlighted }: any) => {
    const [indexOfhighlightedItem, _setIndexOfhighlightedItem] = React.useState<number>(-1)
    const highlightedItemRef = React.useRef(indexOfhighlightedItem)
    const itemsList = React.useRef<any>([])

    const setIndexOfhighlightedItem = (indexOfNewHighlightedItem: number) => {
        highlightedItemRef.current = indexOfNewHighlightedItem
        _setIndexOfhighlightedItem(indexOfNewHighlightedItem)

        onItemHighlighted(itemsList.current?.[indexOfNewHighlightedItem]?.value || '')
    }

    React.useEffect(() => {
        itemsList.current = listItems
        setIndexOfhighlightedItem(-1)
    }, [listItems?.length])

    const highlightPreviousItem = () => {
        if (highlightedItemRef.current >= 0) {
            setIndexOfhighlightedItem(highlightedItemRef.current - 1)
        }
    }

    const highlightNextItem = () => {
        if (highlightedItemRef.current < itemsList.current.length - 1) {
            setIndexOfhighlightedItem(highlightedItemRef.current + 1)
        }
    }

    const handleKeydown: any = (event: React.KeyboardEvent) => {
        if (event.key === 'ArrowDown') {
            highlightNextItem()
        } else if (event.key === 'ArrowUp') {
            highlightPreviousItem()
        }
    }

    React.useEffect(() => {
        document.addEventListener('keydown', handleKeydown)

        return () => {
            document.removeEventListener('keydown', handleKeydown)
        }
    }, [])

    return listItems?.length > 0 ? (
        <div className={styles.dropdownList}>
            {listItems.map((item: any, index: number) => (
                <span
                    data-testid='dropdown-list-item'
                    className={`${styles.dropdownListItem} ${indexOfhighlightedItem === index ? styles.highlightdDropdownListItem : ''
                        }`}
                    key={index}
                    onMouseDown={handleFoundItemClick(item)}
                >
                    {item.value}
                </span>
            ))}
        </div>
    ) : null
}
