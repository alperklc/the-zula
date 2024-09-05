import Button from '../form/button'
import styles from './index.module.css'
import { Link } from 'react-router-dom'

interface SideMenuProps {
  visible: boolean
  setVisibility: (_: boolean) => void
  onLogoutClicked: () => void
  t: (_: string) => string
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
              <Link to='/'>{t('dashboard.title')}</Link>
          </label>
          <label>
              <Link to='/settings/profile'>{t('settings.title')}</Link>
          </label>
        </div>

        <div className={styles.bottom}>
          <Button onClick={props.onLogoutClicked} data-testid='logout-button'>
            {t('common.buttons.logout')}
          </Button>
        </div>
      </div>
    </>
  )
}

export default SideMenu
