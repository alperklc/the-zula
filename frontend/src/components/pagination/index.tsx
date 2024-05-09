
import styles from './index.module.css'

interface PaginationProps {
  onPageClicked: (_: number) => void
  numberOfPages: number
  currentPage: number
}

const range = (start: number, end: number) =>
  Array.from({ length: end - start }, (_, i) => i + start)

export const Pagination = ({ numberOfPages, onPageClicked, currentPage }: PaginationProps) => {
  const visibleWidth = 7
  const pageInTheMiddle = Math.ceil(visibleWidth / 2)

  let buttons: number[] = []

  if (numberOfPages <= visibleWidth) {
    buttons = range(1, numberOfPages + 1)
  }
  // If the current page is within the first half, renders an array like [1 2 *3 4 5 ... 167]
  else if (currentPage < pageInTheMiddle) {
    const fill = range(1, visibleWidth - 1)
    buttons = [...fill, -1, numberOfPages]
  }
  // If the current page is within the last half, renders an array like [1 ... 163 164 *165 166 167]
  else if (currentPage > numberOfPages - pageInTheMiddle) {
    const fill = range(numberOfPages - pageInTheMiddle, numberOfPages + 1)
    buttons = [1, -1, ...fill]
  }
  // If the current page is in the middle, renders dots on both sides like [1 ... 7 *8 9 ... 167]
  else {
    const quarterVisibleWidth = Math.floor(visibleWidth / 4)
    const fill = range(currentPage - quarterVisibleWidth, currentPage + quarterVisibleWidth + 1)
    buttons = [1, -1, ...fill, -1, numberOfPages]
  }

  // "-1" value will be rendered as three dots
  return numberOfPages > 0 ? (
    <div className={styles.pagination} data-testid='pagination'>
      {buttons.map((page: number, index: number) =>
        page === -1 ? (
          <span className={styles.dotsWrapper} key={index} data-testid='child-element'>
            . . .
          </span>
        ) : (
          <button
            key={index}
            className={`${styles.button} ${currentPage === page ? styles.active : ''}`}
            onClick={() => onPageClicked(page)}
            data-testid='child-element'
          >
            {page}
          </button>
        ),
      )}
    </div>
  ) : null
}

export default Pagination
