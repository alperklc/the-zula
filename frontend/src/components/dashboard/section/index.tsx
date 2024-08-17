import { Link } from 'react-router-dom'
import { useIntl } from 'react-intl'
import icons from '../../icons'

import styles from './index.module.css'
import { Insights, VisitingStatistics } from '../../../types/Api'

export type Content = {
  typename: 'Note' | 'Bookmark'
  id: string
  title?: string
  name?: string
}

export type DashboardContentRow = { content: Content; count?: number }

export const DashboardSection = ({
  title,
  rows,
}: {
  title: string
  rows: VisitingStatistics[]
}) => {
  return (
    <div className={styles.dashboardSection}>
      <label className={styles.dashboardLabel} data-testid='title'>
        {title}
      </label>
      {rows?.map((item: VisitingStatistics, index: number) => (
        <ContentDisplay {...item} key={index} />
      ))}
    </div>
  )
}

export const ContentDisplay = (props: VisitingStatistics) => {  
  if (props.typename === 'NOTE') {
    return (
      <Link to={`/notes/${props.id}`}>
        <div data-testid='content' className={styles.truncatedText}>
          <span className={styles.dashboardListItemIcon}>
            <icons.Book height='1rem' width='1rem' />
          </span>
          <span className={styles.dashboardListItemTitle}>{props.title}</span>
        </div>
      </Link>
    )
  }
  if (props.typename === 'BOOKMARK') {
    return (
      <Link to={`/bookmarks/${props.id}`}>
        <div data-testid='content' className={styles.truncatedText}>
          <span className={styles.dashboardListItemIcon}>
            <icons.Link height='1rem' width='1rem' />
          </span>
          <span className={styles.dashboardListItemTitle}>{props.title}</span>
        </div>
      </Link>
    )
  }
  return <span />
}

export const AllContentSection = ({ data: {numberOfBookmarks, numberOfNotes} } : { data: Insights} ) => {
  const intl = useIntl()

  return (
    <div className={styles.dashboardSection}>
      <label className={styles.dashboardLabel}>{'All Content'}</label>
      <Link to='/notes'>
        <div className={styles.truncatedText}>
          <span className={styles.dashboardListItemIcon}>
            <icons.Book height='1rem' width='1rem' />
          </span>
          <span className={styles.dashboardListItemTitle}>
            {intl.formatMessage({ id: 'notes.title' })}
            <span data-testid='numberOfNotes'>
              {numberOfNotes ? ` (${numberOfNotes})` : ''}
            </span>
          </span>
        </div>
      </Link>
      
      <Link to='/bookmarks'>
        <div className={styles.truncatedText}>
          <span className={styles.dashboardListItemIcon}>
            <icons.Link height='1rem' width='1rem' />
          </span>
          <span className={styles.dashboardListItemTitle}>
            {intl.formatMessage({ id: 'bookmarks.title' })}
            <span data-testid='numberOfBookmarks'>
              {numberOfBookmarks ? ` (${numberOfBookmarks})` : ''}
            </span>
          </span>
        </div>
      </Link>
    </div>
  )
}

export default AllContentSection
