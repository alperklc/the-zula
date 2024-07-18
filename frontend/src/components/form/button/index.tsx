import React from 'react'
import styles from './index.module.css'

type StandardButtonProps = React.DetailedHTMLProps<
  React.ButtonHTMLAttributes<HTMLButtonElement>,
  HTMLButtonElement
>

interface AdditionalButtonProps {
  outline?: boolean
  primary?: boolean
  danger?: boolean
}

type ButtonProps = AdditionalButtonProps & StandardButtonProps

// eslint-disable-next-line
const Button = React.forwardRef(
  ({ primary, danger, outline, ...props }: ButtonProps, ref?: React.LegacyRef<HTMLButtonElement>) => {
    return (
      <button
        className={`${styles.container}
          button
          ${outline ? styles.outline : ''}
          ${primary ? styles.primary: ''}
          ${danger ? styles.danger: ''}
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
