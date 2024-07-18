import { FormattedMessage } from 'react-intl'
import Button from '../../form/button'
import modalStyles from '../../modal/index.module.css'

export interface DeleteNoteConfirmationModalProps {
  onConfirm: () => void
  onModalClosed?: () => void
}

const DeleteNoteConfirmation = (props: DeleteNoteConfirmationModalProps) => {
  return (
    <div>
      <div className={modalStyles.modalHeader}>&nbsp;</div>
      <div className={modalStyles.modalBody}>
        <FormattedMessage id='delete_confirmation_modal.title' />
      </div>
      <div className={modalStyles.modalButtons}>
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

export default DeleteNoteConfirmation
