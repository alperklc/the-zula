import icons from '../icons'

import styles from './index.module.css'

export type messageType = 'info' | 'error'

interface MessageBoxProps {
  type?: messageType
  className?: string
  children: JSX.Element | string
}

const MessageBox = ({ className, type = 'info', children }: MessageBoxProps) => (
  <div
    className={`${styles['message-box']} ${className} 
      ${type === 'error' ? styles['message-box-error']: ''}
    `}
  >
    <span className={styles['message-box-icon']}>
      {type === 'info' && <icons.Info />}
      {type === 'error' && <icons.Error />}
    </span>
    <span className={styles['message-box-content']}>{children}</span>
  </div>
)

export default MessageBox
