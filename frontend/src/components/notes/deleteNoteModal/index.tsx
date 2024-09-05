import { useTranslation } from 'react-i18next'
import Button from '../../form/button'
import modalStyles from '../../modal/index.module.css'

export interface DeleteNoteConfirmationModalProps {
  onConfirm: () => void
  onModalClosed?: () => void
}

const DeleteNoteConfirmation = (props: DeleteNoteConfirmationModalProps) => {
  const { t } = useTranslation()
  
  return (
    <div>
      <div className={modalStyles.modalHeader}>&nbsp;</div>
      <div className={modalStyles.modalBody}>
        {t('delete_confirmation_modal.title')}
      </div>
      <div className={modalStyles.modalButtons}>
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

export default DeleteNoteConfirmation
