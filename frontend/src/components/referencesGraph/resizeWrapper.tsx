import React, { ReactNode } from 'react'
import { debounce } from '../../utils/useDebounce'

export function ResizeWrapper({
  children: graph,
  fullScreen,
}: {
  children: (_?: any) => ReactNode
  fullScreen?: boolean
}) {
  const graphWrapperRef = React.useRef<HTMLDivElement>(null)

  const [widthGraph, setWidthGraph] = React.useState(100)
  const [heightGraph, setHeightGraph] = React.useState(100)

  const calculateDimensionsOfGraph = debounce(() => {
    const dimensions = graphWrapperRef.current?.getBoundingClientRect()

    setWidthGraph(dimensions?.width!)
    setHeightGraph(fullScreen ? window.innerHeight * 0.8 : dimensions?.width! * 1.15)
  }, 100)

  React.useEffect(() => {
    calculateDimensionsOfGraph()
    window.addEventListener('resize', calculateDimensionsOfGraph)

    return () => {
      window.removeEventListener('resize', calculateDimensionsOfGraph)
    }
  }, [calculateDimensionsOfGraph])

  return (
    <div data-testid='graphwrapper' ref={graphWrapperRef}>
      {graph({ width: widthGraph, height: heightGraph })}
    </div>
  )
}
