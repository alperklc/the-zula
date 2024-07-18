import React, { useState, ChangeEvent, useEffect } from 'react'
import Input from '../form/input'
import Icons from '../icons'
import TagsDisplay from '../tagsDisplay'
import { searchTags } from './search-tags'

import styles from './index.module.css'
import { Tag } from '../../types/Api'

const TagsDropdownList = ({ listItems, handleFoundTagClick, onTagHighlighted }: any) => {
  const [indexOfhighlightedItem, _setIndexOfhighlightedItem] = React.useState<number>(-1)
  const highlightedItemRef = React.useRef(indexOfhighlightedItem)
  const tagsList = React.useRef<any>([])

  const setIndexOfhighlightedItem = (indexOfNewHighlightedItem: number) => {
    highlightedItemRef.current = indexOfNewHighlightedItem
    _setIndexOfhighlightedItem(indexOfNewHighlightedItem)

    onTagHighlighted(tagsList.current?.[indexOfNewHighlightedItem]?.value || '')
  }

  useEffect(() => {
    tagsList.current = listItems
    setIndexOfhighlightedItem(-1)
  }, [listItems?.length])

  const highlightPreviousItem = () => {
    if (highlightedItemRef.current >= 0) {
      setIndexOfhighlightedItem(highlightedItemRef.current - 1)
    }
  }

  const highlightNextItem = () => {
    if (highlightedItemRef.current < tagsList.current.length - 1) {
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

  useEffect(() => {
    document.addEventListener('keydown', handleKeydown)

    return () => {
      document.removeEventListener('keydown', handleKeydown)
    }
  }, [])

  return listItems?.length > 0 ? (
    <div className={styles.dropdownList}>
      {listItems.map((tag: any, index: number) => (
        <span
          data-testid='dropdown-list-item'
          className={`${styles.dropdownListItem} ${
            indexOfhighlightedItem === index ? styles.highlightdDropdownListItem : ''
          }`}
          key={index}
          onMouseDown={handleFoundTagClick(tag)}
        >
          {tag.value}
        </span>
      ))}
    </div>
  ) : null
}

export interface TagsInputProps {
  tags: string[]
  onChange: (_: string[]) => void
  label: string
  placeholder: string
  onSearch: (_: string) => { fetching: boolean; foundTags: Tag[] }
}

export function TagsInput(props: TagsInputProps) {
  const [inputFieldActive, setInputFieldActive] = useState<boolean>(false)
  const [highlightedTag, setHighlightedTag] = useState<string>('')
  const [textInputValue, setTextInputValue] = useState('')

  const searchResponse = props.onSearch(textInputValue)

  const removeTag = (tag: string) => {
    const newTags = props.tags.filter((current: string) => tag !== current)
    props.onChange(newTags)
  }

  const addTag = (tag: string) => {
    if (props.tags.indexOf(tag) === -1) {
      props.onChange([...props.tags, tag])
      setTextInputValue('')
    }
  }

  const handleTextChange = (e: ChangeEvent<HTMLInputElement>) => {
    setTextInputValue(e?.target?.value)
  }

  const handleEnterDown = () => {
    addTag(highlightedTag || textInputValue)
  }

  const handleBackspaceDown = () => {
    if (textInputValue === '') {
      const [lastTag] = [...props.tags].reverse()
      removeTag(lastTag)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
      return
    }

    if (e.key === 'Enter') {
      e.preventDefault()
      handleEnterDown()
    } else if (e.key === 'Backspace') {
      handleBackspaceDown()
    }
    e.stopPropagation()
  }

  const onPlusClick = (e: React.MouseEvent) => {
    e.preventDefault()
    addTag(textInputValue)
  }

  const handleFoundTagClick = (tag: any) => () => {
    addTag(tag.value)
    setTextInputValue('')
  }

  return (
    <div className={styles.container}>
      {props?.label && <label>{props.label}</label>}
      {props?.tags?.length > 0 && <TagsDisplay tags={props.tags} removeTag={removeTag} />}
      <span className={styles.textInputContainer}>
        <Input
          data-testid='search-input'
          className={styles.textInput}
          onChange={handleTextChange}
          type='text'
          placeholder='Type here to add tags + hit enter'
          value={textInputValue}
          onKeyDown={handleKeyDown}
          onFocus={() => setInputFieldActive(true)}
          onBlur={() => setInputFieldActive(false)}
        />
        <span className={styles.addButton} onClick={onPlusClick}>
          <Icons.Plus width={18} height={18} />
        </span>
      </span>

      {textInputValue &&
        inputFieldActive &&
        !searchResponse?.fetching &&
        searchResponse?.foundTags.length > 0 && (
          <TagsDropdownList
            listItems={searchResponse.foundTags}
            handleFoundTagClick={handleFoundTagClick}
            onTagHighlighted={setHighlightedTag}
          />
        )}
    </div>
  )
}

export { searchTags }

export default TagsInput
