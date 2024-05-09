import icons from '../icons'
import styles from './index.module.css'

export type ToastType = null | 'error' | 'success' | 'info'

export interface ToastProps {
  message?: string
  type: ToastType
}

const Toast = (props: ToastProps) => {
  const { message, type } = props

  return (
    <div
      className={`
        ${styles.container}
        ${type === 'error' ? styles.error : '' }
        ${type === 'success' ? styles.success : '' }
        ${type === 'info' ? styles.info : '' }
      )`}
    >
      <div className={styles.innerContainer}>
        <span className={styles.icon}>
          {type === 'success' && <icons.Checkmark height='1.2rem' width='1.2rem' />}
          {type === 'info' && <icons.Info height='1.2rem' width='1.2rem' />}
          {type === 'error' && <icons.X height='1.2rem' width='1.2rem' />}
        </span>
        <span className={styles.text}>{message}</span>
      </div>
    </div>
  )
}

export default Toast
