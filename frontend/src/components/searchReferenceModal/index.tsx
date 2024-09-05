import { useTranslation } from 'react-i18next'
import Button from '../form/button'
import styles from '../modal/index.module.css'

export interface SearchReferenceModalProps {
  onConfirm: () => void
  onModalClosed?: () => void
}

const SearchReferenceModal = (props: SearchReferenceModalProps) => {
  const { t } = useTranslation()

  return (
    <div>
      <div className={styles.modalHeader}>&nbsp;</div>
      <div className={styles.modalBody}>
        {t('delete_confirmation_modal.title')}
      </div>
      <div className={styles.modalButtons}>
        <Button danger onClick={props.onConfirm}>
          {t('common.buttons.delete')}
        </Button>
        <Button outline onClick={props.onModalClosed}>
          {t('common.buttons.cancel')}
        </Button>
      </div>
    </div>
  )
}

export default SearchReferenceModal
