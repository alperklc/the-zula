const Book = ({ className, width, height }: React.SVGProps<SVGSVGElement>) => {
  return (
    <svg
      xmlns='http://www.w3.org/2000/svg'
      width={width || 24}
      height={height || 24}
      fill='none'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      strokeWidth='2'
      className={className}
      viewBox='0 0 24 24'
    >
      <path d='M4 19.5A2.5 2.5 0 016.5 17H20'></path>
      <path d='M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z'></path>
    </svg>
  )
}

export default Book
