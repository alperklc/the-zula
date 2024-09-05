import React from 'react'
import Button from '../form/button'
import TagsInput, { searchTags } from '../tagsInput'
import modalStyles from '../modal/index.module.css'
import styles from './index.module.css'
import { Filter, Sort } from '../../contexts/queryContext'
import { useTranslation } from 'react-i18next'

export interface FilterSelectorModalProps {
  typeOfParent: string
  tags?: string[]
  sortableFields?: string[]
  sortBy?: string
  sortDirection?: string
  onApply: (_: Partial<Filter & Sort>) => void
  onModalClosed?: () => void
}

const FilterSelector = (props: FilterSelectorModalProps) => {
  const { t } = useTranslation()

  const [tags, setTags] = React.useState<string[]>(props?.tags || [])
  const [sortBy, setSortBy] = React.useState<string>(props?.sortBy || '')
  const [sortDirection, setSortDirection] = React.useState<string>(props?.sortDirection || '')

  const onApplyClick = () => {
    props.onApply({ tags, sortBy, sortDirection })
    props.onModalClosed?.()
  }

  const handleTagsSelection = (tags: string[]) => {
    setTags(tags)
  }

  const onClearAllFiltersClick = () => {
    props.onApply({ tags: [] })
    props.onModalClosed?.()
  }

  const onSortByChanged = (nextSortBy: string) => () => {
    setSortBy(nextSortBy)
  }

  const onSortDirectionChangedTo = (nextSortDirection: string) => () => {
    setSortDirection(nextSortDirection)
  }

  return (
    <div>
      <div className={modalStyles.modalHeader}>
        {t('filter_selector_modal.title')}
      </div>
      <div className={modalStyles.modalBody}>
        {props?.tags && (
          <div>
            <TagsInput
              tags={tags}
              onSearch={searchTags(props.typeOfParent)}
              onChange={handleTagsSelection}
              label='Tags'
              placeholder={t('filter_selector_modal.tags.label')}
            />
          </div>
        )}
        {props?.sortableFields && (
          <div>
            <div>
              <label>
                {t('filter_selector_modal.sort_by_label')}
              </label>
              {props?.sortableFields.map((field: string, index: number) => (
                <span
                  key={index}
                  className={`${styles.choice} ${field === sortBy ? styles.selectedChoice: ''}`}
                  onClick={onSortByChanged(field)}
                >
                  {field}
                </span>
              ))}
            </div>
            <div>
              <label>
                {t('filter_selector_modal.sort_direction_label')}
              </label>
              <span
                className={`${styles.choice} ${'asc' === sortDirection ? styles.selectedChoice : ''}`}
                onClick={onSortDirectionChangedTo('asc')}
              >
                asc
              </span>
              <span
                className={`${styles.choice} ${'desc' === sortDirection ? styles.selectedChoice : ''}`}
                onClick={onSortDirectionChangedTo('desc')}
              >
                desc
              </span>
            </div>
          </div>
        )}
      </div>

      <div className={modalStyles.modalButtons}>
        <Button onClick={onApplyClick}>
          {t('common.buttons.apply')}
        </Button>
        {tags?.length > 0 ? (
          <Button outline onClick={onClearAllFiltersClick}>
            {t('filter_selector_modal.buttons.clear')}
          </Button>
        ) : (
          <Button outline onClick={props.onModalClosed}>
            {t('common.buttons.cancel')}
          </Button>
        )}
      </div>
    </div>
  )
}

export default FilterSelector
