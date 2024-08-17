import * as React from 'react'
import { Suggestion } from '../types/Suggestion'
import { CaretCoordinates } from '../util/TextAreaCaretPosition'
import styles from './index.module.css'

export interface SuggestionsDropdownProps {
  caret: CaretCoordinates
  suggestions: Suggestion[]
  suggestionsAutoplace: boolean
  onSuggestionSelected: (index: number) => void
  /**
   * Which item is focused by the keyboard
   */
  focusIndex: number
  textAreaRef: React.RefObject<HTMLTextAreaElement>
}

export const SuggestionsDropdown: React.FunctionComponent<SuggestionsDropdownProps> = ({
  suggestions,
  caret,
  onSuggestionSelected,
  suggestionsAutoplace,
  focusIndex,
  textAreaRef,
}) => {
  const handleSuggestionClick = (event: React.MouseEvent) => {
    event.preventDefault()
    const index = parseInt((event.currentTarget.attributes as any)['data-index'].value as string)
    onSuggestionSelected(index)
  }

  const handleMouseDown = (event: React.MouseEvent) => event.preventDefault()

  const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
  const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)

  const left = caret.left - (textAreaRef.current?.scrollLeft as number)
  const top = caret.top - (textAreaRef.current?.scrollTop as number)

  const style: React.CSSProperties = {}
  if (suggestionsAutoplace && top + textAreaRef.current!.getBoundingClientRect().top > vh / 2)
    style.bottom = textAreaRef.current!.offsetHeight - top
  else style.top = top

  if (suggestionsAutoplace && left + textAreaRef.current!.getBoundingClientRect().left > vw / 2)
    style.right = textAreaRef.current!.offsetWidth - left
  else style.left = left

  return (
    <ul className={styles.container} style={style}>
      {suggestions?.map((s, i) => (
        <li
          onClick={handleSuggestionClick}
          onMouseDown={handleMouseDown}
          key={i}
          aria-selected={focusIndex === i ? 'true' : 'false'}
          data-index={`${i}`}
        >
          {s.preview}
        </li>
      ))}
    </ul>
  )
}
