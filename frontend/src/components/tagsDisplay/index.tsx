import styles from './index.module.css'

interface TagsDisplayProps {
  className?: string
  tags: string[]
  maxNumberOfTagsToDisplay?: number
  removeTag?: (_: string) => void
}

const TagsDisplay = (props: TagsDisplayProps) => {
  const handleRemoveClick = (tag: string) => () => {
    props.removeTag?.(tag)
  }

  let tagsToDisplay = [...props.tags]

  if (props.maxNumberOfTagsToDisplay) {
    tagsToDisplay = tagsToDisplay.slice(0, props.maxNumberOfTagsToDisplay)
  }

  return (
    <div className={`${styles.tagsContainer} ${props.className}`}>
      {tagsToDisplay.map((tag: string, index: number) => (
        <div className={styles.tag} key={index}>
          <span className={styles.tagValue} data-testid='tag-value'>
            {tag}
          </span>
          {props?.removeTag && (
            <span
              className={styles.tagRemoveIcon}
              onClick={handleRemoveClick(tag)}
              data-testid='remove-tag-icon'
            >
              x
            </span>
          )}
        </div>
      ))}
    </div>
  )
}

export default TagsDisplay
