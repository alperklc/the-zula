import * as React from 'react'
import { useNavigate } from 'react-router-dom'
import { Converter } from 'showdown'
import styles from './index.module.css'

interface MarkdownDisplayProps {
  className?: string
  content: string
}

const converter = new Converter({
  tables: true,
  simplifiedAutoLink: true,
  strikethrough: true,
  tasklists: true,
})

converter.setFlavor("github");

const getMarkdownPreview = (markdown: string) => Promise.resolve(converter.makeHtml(markdown))

const MarkdownDisplay = (props: MarkdownDisplayProps) => {
  const navigate = useNavigate()
  const divRef = React.useRef<HTMLDivElement>()

  const [__html, setHtml] = React.useState<string>('')

  React.useEffect(() => {
    (async function () {
      const html = await getMarkdownPreview(props.content)
      setHtml(html)
    })()
  }, [props.content])

  const handleClick = (event: React.MouseEvent) => {
    const clickedElement = event.target as HTMLAnchorElement
    if (clickedElement?.href?.includes(`/notes/`)) {
      event.preventDefault()
      navigate(clickedElement.pathname)
    }
  }

  return (
    <div
      className={`${props.className} ${styles.container}`}
      dangerouslySetInnerHTML={{ __html }}
      ref={divRef as React.RefObject<HTMLDivElement>}
      onClick={handleClick}
    />
  )
}

export default MarkdownDisplay
