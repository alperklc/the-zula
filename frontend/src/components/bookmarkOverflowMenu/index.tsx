import * as React from 'react'
import { FormattedMessage, useIntl } from 'react-intl'
import { useToast } from '../toast/toast-message-context'
import styles from './index.module.css'

const OverflowMenu = (props: any) => {
  const { show: showToast } = useToast()
  const intl = useIntl()

  const onCopyLinkClick = (event: React.MouseEvent) => {
    navigator?.clipboard?.writeText(props.link || '').then(() => {
      showToast(intl.formatMessage({ id: 'bookmarks.toast.copy_link' }), 'info')
    })

    event.stopPropagation()
    props.onOptionClick()
  }
  const onShareClick = (event: React.MouseEvent) => {
    if (navigator?.share) {
      navigator
        ?.share({
          title: props?.title || '',
          url: props?.link || '',
        })
        .then(() => console.log('Share was successful.'))
        .catch((error) => console.log('Sharing failed', error))
    }
    event.stopPropagation()
    props.onOptionClick()
  }

  return (
    <div className={`${props.className} ${styles.container}`}>
      <span className={styles.menuChoice} onClick={onCopyLinkClick}>
        <FormattedMessage id='bookmarks.overflow_menu.copy_link' />
      </span>

      <span className={styles.menuChoice} onClick={onShareClick}>
        <FormattedMessage id='bookmarks.overflow_menu.share' />
      </span>
    </div>
  )
}

export default OverflowMenu