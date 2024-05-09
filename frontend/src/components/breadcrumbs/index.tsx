import styles from './index.module.css'
import { useLocation, useNavigate } from 'react-router-dom'

const Breadcrumbs = (props: {
  id?: string
  className?: string
  children?: JSX.Element | JSX.Element[]
}) => {
  const navigate = useNavigate()
  const { pathname } = useLocation()

  const segmentsFromRouter = pathname
    .split('/')
    .map((href: string) => ({ href }))

  const handleSegmentClick = (index: number) => () => {
    const href = segmentsFromRouter
      .map((segment) => segment.href)
      .slice(0, index + 1)
      .join('/')

    navigate(`${href}`, { relative: "path" })
  }

  return (
    <ul id={props.id} className={styles.container}>
      {props?.children
        ? props.children
        : segmentsFromRouter.map((segment: { href: string }, index: number) => (
            <li key={index} onClick={handleSegmentClick(index)}>
              <span>{segment.href}</span>
            </li>
          ))}
    </ul>
  )
}

export default Breadcrumbs
