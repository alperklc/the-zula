import { FormattedMessage } from 'react-intl'
import Button from '../form/button'
import styles from '../modal/index.module.css'

export interface SearchReferenceModalProps {
  onConfirm: () => void
  onModalClosed?: () => void
}

const SearchReferenceModal = (props: SearchReferenceModalProps) => {
  return (
    <div>
      <div className={styles.modalHeader}>&nbsp;</div>
      <div className={styles.modalBody}>
        <FormattedMessage id='delete_confirmation_modal.title' />
      </div>
      <div className={styles.modalButtons}>
        <Button danger onClick={props.onConfirm}>
          <FormattedMessage id='common.buttons.delete' />
        </Button>
        <Button outline onClick={props.onModalClosed}>
          <FormattedMessage id='common.buttons.cancel' />
        </Button>
      </div>
    </div>
  )
}

export default SearchReferenceModal
