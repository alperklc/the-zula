import React from 'react'

interface AdditionalProps {
  label?: string
  inputref?: React.RefObject<HTMLInputElement>
}

const Input = (
  props: React.DetailedHTMLProps<React.InputHTMLAttributes<HTMLInputElement>, HTMLInputElement> &
    AdditionalProps,
) => {
  return (
    <>
      {props.label && <label htmlFor={props.name}>{props.label}</label>}
      <input
        id={props.name}
        ref={props.inputref}
        className={props.className}
        placeholder={props.placeholder}
        {...props}
      />
    </>
  )
}

export default Input
