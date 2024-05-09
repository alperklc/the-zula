import React from 'react'

import styles from './index.module.css'

interface AnimationProps {
  className?: string
  visible: boolean
  type: 'fadeIn' | 'slideFromRight'
  children: React.ReactNode
}

function getKeyframes(type: string) {
  switch (type) {
    case 'fadeIn':
      return [styles.fadeIn, styles.fadeOut]
    case 'slideFromRight':
      return [styles.slideIn, styles.slideOut]
    default:
      return ['', '']
  }
}

export const Animation = (props: AnimationProps) => {
  const { visible, className } = props

  const [render, setRender] = React.useState(visible)

  React.useEffect(() => {
    if (visible) setRender(true)
  }, [visible])

  const onAnimationEnd = () => {
    if (!visible) setRender(false)
  }

  const [animationIn, animationOut] = getKeyframes(props.type)

  return render ? (
    <div
      className={`${styles.container} ${className}`}
      style={{
        animation: `${visible ? animationIn : animationOut} 0.3s`,
      }}
      onAnimationEnd={onAnimationEnd}
    >
      {props.children}
    </div>
  ) : null
}

export default Animation
