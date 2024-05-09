import React, { ReactElement, ReactNode } from 'react'
import { Animation } from '../animations'
import icons from '../icons'
import LoadingIndicator from '../loadingIndicator'

import styles from './index.module.css'

interface PageContentProps {
  loading?: boolean
  isMobile?: boolean
  tabs?: ReactElement | ReactElement[] | ReactNode | ReactNode[]
  children?: ReactElement | ReactElement[] | ReactNode | ReactNode[]
}

const PageContent = (props: PageContentProps) => {
  const [rightPanelVisible, setRightPanelVisibility] = React.useState<boolean>(false)

  const [content, sideContent] = React.Children.toArray(props.children) as ReactElement[]

  const handleRightPanelToggleClick = () => {
    setRightPanelVisibility(!rightPanelVisible)
  }

  return (
    <>
      {props?.tabs && <section className={styles.tabs}>{props.tabs}</section>}
      {props.loading ? (
        <LoadingIndicator className={styles.loadingSpinner} />
      ) : (
        <div
          className={`${styles.container} ${rightPanelVisible ? styles.preventScroll: ''}`}
        >
          {props.isMobile ? (
            <>
              {rightPanelVisible && (
                <div className={styles.rightPanelBackdrop} onClick={handleRightPanelToggleClick} />
              )}

              <section className={styles.mainContentMobile}>
                {sideContent && (
                  <div className={styles.rightPanelToggle}>
                    <span onClick={handleRightPanelToggleClick}>
                      <icons.ChevronLeft />
                    </span>
                  </div>
                )}
                {content}
              </section>

              <Animation type='slideFromRight' visible={rightPanelVisible}>
                <section className={`${styles.rightPanel} ${styles.rightPanelMobile}`}>
                  <div>
                    <span onClick={handleRightPanelToggleClick}>
                      <icons.ChevronRight />
                    </span>
                  </div>
                  {sideContent}
                </section>
              </Animation>
            </>
          ) : (
            <>
              <section className={`${styles.mainContent} ${!sideContent ? styles.mainContentSingle: ''}`}>
                {content}
              </section>
              {sideContent && <section className={styles.rightPanel}>{sideContent}</section>}
            </>
          )}
        </div>
      )}
    </>
  )
}

export default PageContent
