function Share(props: React.SVGProps<SVGSVGElement>) {
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
      className='feather feather-share'
      viewBox='0 0 24 24'
      {...props}
    >
      <path d='M4 12v8a2 2 0 002 2h12a2 2 0 002-2v-8'></path>
      <path d='M16 6L12 2 8 6'></path>
      <path d='M12 2L12 15'></path>
    </svg>
  )
}

export default Share
