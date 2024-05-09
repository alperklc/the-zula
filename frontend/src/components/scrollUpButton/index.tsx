import React from 'react'
import Animation from '../animations'
import icons from '../icons'

import styles from './index.module.css'

const ScrollUpButton = () => {
  const [visible, setVisibility] = React.useState(false)

  const rootElement = document.documentElement

  const onScrollUpClick = () => {
    rootElement.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
  }

  const onScroll = () => {
    if (rootElement.scrollTop > rootElement.clientHeight - 200) {
      setVisibility(true)
    } else {
      setVisibility(false)
    }
  }

  React.useEffect(() => {
    window.addEventListener('scroll', onScroll)

    return () => {
      window.removeEventListener('scroll', onScroll)
    }
  }, [])

  return (
    <Animation type='fadeIn' visible={visible}>
      <div className={styles.scrollUpButton} onClick={onScrollUpClick}>
        {icons.ArrowUp()}
      </div>
    </Animation>
  )
}

export default ScrollUpButton
