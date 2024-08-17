export interface Selection {
  start: number
  end: number
}

export interface TextSection {
  text: string
  selection: Selection
}
