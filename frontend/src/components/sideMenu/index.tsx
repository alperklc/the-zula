import { MessageDescriptor } from 'react-intl'
import Button from '../form/button'
import styles from './index.module.css'
import { Link } from 'react-router-dom'

interface SideMenuProps {
  visible: boolean
  setVisibility: (_: boolean) => void
  onLogoutClicked: () => void
  t: (_: MessageDescriptor) => string
}

const SideMenu = ({ t, ...props }: SideMenuProps) => {
  return (
    <>
      {props.visible && (
        <div
          className={styles.backdrop}
          data-testid='backdrop'
          onClick={() => props.setVisibility(false)}
        />
      )}
      <div className={`${styles.sideNav} ${props.visible ? styles.visible : ''}`}>
        <div className={styles.linksToPages}>
          <label>
              <Link to='/'>{t({ id: 'dashboard.title' })}</Link>
          </label>
          <label>
              <Link to='/settings/profile'>{t({ id: 'settings.title' })}</Link>
          </label>
        </div>

        <div className={styles.bottom}>
          <Button onClick={props.onLogoutClicked} data-testid='logout-button'>
            {t({ id: 'common.buttons.logout' })}
          </Button>
        </div>
      </div>
    </>
  )
}

export default SideMenu
