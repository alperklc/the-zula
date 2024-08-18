import React from 'react'
import { useNavigate } from 'react-router-dom'
import ForceGraph2D, { ForceGraphMethods, NodeObject, GraphData, LinkObject } from 'react-force-graph-2d'

const shortenTitle = (title: string) => {
  if (title.length <= 8) {
    return title
  }

  return `${title.slice(0, 8)}...`
}

interface ExtendedNodeObject extends NodeObject {
  neighbors?: NodeObject[]
  links?: LinkObject[]
}

function addNeighbors(data: GraphData) {
  data.links?.forEach((link) => {
    const a: ExtendedNodeObject = data.nodes.find((el) => el.id === link?.source)!
    const b: ExtendedNodeObject = data.nodes.find((el) => el.id === link?.target)!

    a && !a?.neighbors && (a.neighbors = [])
    b && !b?.neighbors && (b.neighbors = [])
    a && a?.neighbors?.push(b)
    b && b?.neighbors?.push(a)

    a && !a?.links && (a.links = [])
    b && !b?.links && (b.links = [])
    a && a?.links?.push(link)
    b && b?.links?.push(link)
  })

  return data
}

export function ReferencesGraph({
  height,
  width,
  noteId,
  graphData,
}: {
  height?: number
  width?: number
  noteId: string
  graphData: GraphData
}) {
  const data = addNeighbors(graphData)
  const navigate = useNavigate()
  const fgRef = React.useRef<ForceGraphMethods>()

  const [highlightNodes, setHighlightNodes] = React.useState(new Set())
  const [highlightLinks, setHighlightLinks] = React.useState(new Set())
  const [selectedNodeUid, setSelectedNodeUid] = React.useState('')

  const updateHighlight = () => {
    setHighlightNodes(highlightNodes)
    setHighlightLinks(highlightLinks)
  }

  const handleNodeClick = (node: any) => {
    if (node.id === selectedNodeUid) {
      navigate(`/notes/${node.id}`)
      return
    }

    setSelectedNodeUid(node.id)

    highlightNodes.clear()
    highlightLinks.clear()
    if (node) {
      highlightNodes.add(node)

      node?.neighbors?.forEach((neighbor: any) => highlightNodes.add(neighbor))
      node?.links?.forEach((link: any) => highlightLinks.add(link))
    }

    updateHighlight()
  }

  const handleBackgroundClick = () => {
    setSelectedNodeUid('')

    highlightNodes.clear()
    highlightLinks.clear()
  }

  const nodeCanvasObj = (node: NodeObject, ctx: CanvasRenderingContext2D, globalScale: number) => {
    const label = shortenTitle((node as { title: string }).title)

    const fontSize = 12.5 / globalScale

    ctx.font = `${fontSize}px sans-serif`
    const textWidth = ctx.measureText(label as string).width

    const bckgDimensions = [textWidth + fontSize, fontSize * 1.5]

    ctx.fillStyle = node.id !== noteId ? 'rgba(44, 62, 80, 1)' : 'rgba(121, 30, 30, 1)'
    ctx.fillRect(
      node.x! - bckgDimensions[0] / 2,
      node.y! - bckgDimensions[1] / 2,
      bckgDimensions[0],
      bckgDimensions[1],
    )

    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillStyle = 'rgba(255,255,255,1)'

    ctx.fillText(label as string, node.x!, node.y!)

    const _node = node as NodeObject & { __bckgDimensions: number[] }

    _node.__bckgDimensions = bckgDimensions // to re-use in nodePointerAreaPaint
  }

  const nodePointerAreaPaint = (node: NodeObject, color: string, ctx: CanvasRenderingContext2D) => {
    ctx.fillStyle = color

    const _node = node as NodeObject & { __bckgDimensions: number[] }
    const bckgDimensions = _node.__bckgDimensions
    bckgDimensions &&
      ctx.fillRect(
        node.x! - bckgDimensions[0] / 2,
        node.y! - bckgDimensions[1] / 2,
        bckgDimensions[0],
        bckgDimensions[1],
      )
  }

  return (
    <ForceGraph2D
      ref={fgRef}
      width={width}
      height={height}
      minZoom={2.75}
      autoPauseRedraw={false}
      linkWidth={(link) => (highlightLinks.has(link) ? 1.5 : 1)}
      linkColor={(link) => (highlightLinks.has(link) ? '#2b2b2b' : '#ddd')}
      nodeLabel={(node: any) => (node.title.length > 8 ? node?.title : null)}
      nodeCanvasObject={nodeCanvasObj}
      nodePointerAreaPaint={nodePointerAreaPaint}
      onNodeClick={handleNodeClick}
      onBackgroundClick={handleBackgroundClick}
      graphData={data}
    />
  )
}

export default ReferencesGraph
