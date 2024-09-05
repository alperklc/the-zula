import React, { ReactNode } from 'react'
import StandardHeader from '../header'
import SideMenu from '../sideMenu'
import { useAuth } from '../../contexts/authContext'
import { useUI } from '../../contexts/uiContext'
import BaseToast from '../../components/toast/index'
import ScrollUpButton from '../scrollUpButton'
import { useToast } from '../toast/toast-message-context'
import Animation from '../animations'
import styles from './index.module.css'
import { useTranslation } from 'react-i18next'

interface LayoutProps {
  narrow?: boolean
  children?: ReactNode
  subHeaderContent?: ReactNode
  fixedSubHeader?: boolean
  customHeader?: (_: any) => JSX.Element
  showScrollUpButton?: boolean
}

const Layout = (props: LayoutProps) => {
  const { logout, user } = useAuth()
  const { t } = useTranslation()

  const { isVisible, message, type } = useToast()
  const { backdropActive, toggleBackdrop } = useUI()
  const [menuVisible, setVisibilityOfMenu] = React.useState(false)
  const toggleSideNav = () => setVisibilityOfMenu(!menuVisible)

  const Header = props?.customHeader ?? StandardHeader

  return (
    <>
      {user !== null && (
        <SideMenu
          visible={menuVisible}
          setVisibility={setVisibilityOfMenu}
          onLogoutClicked={logout}
          t={t}
        />
      )}
      {backdropActive && <div className={styles.backdropOverlay} onClick={toggleBackdrop} />}
      <div className={styles.container}>
        <Header className={styles.header} onMenuIconClicked={toggleSideNav} />
        <div
          className={`${styles.subHeader} ${
            !props?.fixedSubHeader || !props.subHeaderContent ? styles.emptySubHeader : ''
          }`}
        >
          {props?.fixedSubHeader ? props.subHeaderContent : <></>}
        </div>
        <Animation className={styles.toastContainer} type='fadeIn' visible={isVisible}>
          <BaseToast message={message} type={type} />
        </Animation>
        <div
          className={`${styles.pageContainer}
            ${props.narrow ? styles.narrowLayout : ''}
            ${!props?.fixedSubHeader ? styles.withSubHeader : ''}`
          }
        >
          {!props?.fixedSubHeader && props.subHeaderContent}
          {props.showScrollUpButton && <ScrollUpButton />}
          {props.children}
        </div>
      </div>
    </>
  )
}

export { styles }

export default Layout
