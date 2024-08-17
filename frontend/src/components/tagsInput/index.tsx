import React, { useState, ChangeEvent } from 'react'
import Input from '../form/input'
import Icons from '../icons'
import TagsDisplay from '../tagsDisplay'
import { searchTags } from './search-tags'
import { Tag } from '../../types/Api'
import { AutocompleteDropdown } from '../autocompleteDropdown'

import styles from './index.module.css'

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
          <AutocompleteDropdown
            listItems={searchResponse.foundTags}
            handleFoundItemClick={handleFoundTagClick}
            onItemHighlighted={setHighlightedTag}
          />
        )}
    </div>
  )
}

export { searchTags }

export default TagsInput
