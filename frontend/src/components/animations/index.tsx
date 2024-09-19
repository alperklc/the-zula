import React from 'react'

import styles from './index.module.css'

interface AnimationProps {
  className?: string
  visible: boolean
  time?: number
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
  const { visible, className, time = 250 } = props

  const [render, setRender] = React.useState(visible);
  const [animate, setAnimate] = React.useState(visible);

  const [animationIn, animationOut] = getKeyframes(props.type)

  React.useEffect(() => {
    if (visible) {
      setTimeout(() => setRender(true), 1);
      setTimeout(() => setAnimate(true), 1);
    }  else {
      setAnimate(false);
      setTimeout(() => setRender(false), time);
    }
  }, [visible]);

  return render ? (
    <div
      className={`${styles.container} ${className}`}
      style={{
        animation: `${animate ? animationIn : animationOut} ${time / 1000}s`,
      }}
    >
      {props.children}
    </div>
  ) : null


}

export default Animation
