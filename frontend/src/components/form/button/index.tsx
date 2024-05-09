import React from 'react'
import styles from './index.module.css'

type StandardButtonProps = React.DetailedHTMLProps<
  React.ButtonHTMLAttributes<HTMLButtonElement>,
  HTMLButtonElement
>

interface AdditionalButtonProps {
  muted?: boolean
  primary?: boolean
  danger?: boolean
}

type ButtonProps = AdditionalButtonProps & StandardButtonProps

// eslint-disable-next-line
const Button = React.forwardRef(
  ({ primary, danger, muted, ...props }: ButtonProps, ref?: React.LegacyRef<HTMLButtonElement>) => {
    return (
      <button
        className={`${styles.container}
          button
          ${muted ? 'muted-button' : ''}
          ${primary ? 'primary-button': ''}
          ${danger ? 'danger-button': ''}
          ${props.className}
        `}
        ref={ref}
        {...props}
      >
        {props.children}
      </button>
    )
  },
)

export { styles }

export default Button
