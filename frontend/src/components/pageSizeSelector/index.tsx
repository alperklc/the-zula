import styles from './index.module.css'

interface PageSizeSelectorProps {
  className?: string
  pageSize: number
  onPageSizeSelected: (_: number) => void
}

const PageSizeSelector = (props: PageSizeSelectorProps) => {
  const handlePageSizeChange = (nextPageSize: number) => () => {
    props.onPageSizeSelected(nextPageSize)
  }

  return (
    <span className={styles.pageSizeSelection}>
      {[10, 25, 50].map((i: number) => (
        <span
          key={i}
          className={`${styles.pageSizeSelector} ${props.pageSize === i ? styles.pageSizeSelectorActive : ''}`}
          onClick={handlePageSizeChange(i)}
        >
          {i}
        </span>
      ))}
    </span>
  )
}

export default PageSizeSelector
