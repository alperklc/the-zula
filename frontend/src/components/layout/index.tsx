import React, { ReactNode } from 'react'
import { useIntl } from 'react-intl'
import StandardHeader from '../header'
import SideMenu from '../sideMenu'
import { useAuth } from '../../contexts/authContext'
import { useUI } from '../../contexts/uiContext'
import ScrollUpButton from '../scrollUpButton'

import styles from './index.module.css'

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
  const intl = useIntl()
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
          t={intl.formatMessage}
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
