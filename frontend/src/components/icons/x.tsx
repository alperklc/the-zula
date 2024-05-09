import React from 'react'

function X(props: any) {
  return (
    <svg
      xmlns='http://www.w3.org/2000/svg'
      width='24'
      height='24'
      fill='none'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      strokeWidth='2'
      className='feather feather-x'
      viewBox='0 0 24 24'
      {...props}
    >
      <path d='M18 6L6 18'></path>
      <path d='M6 6L18 18'></path>
    </svg>
  )
}

export default X
