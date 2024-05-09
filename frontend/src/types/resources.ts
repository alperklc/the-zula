import { Note } from "./Api"

export const resourceTypes = ["note"] as const

export type ResourceType = typeof resourceTypes[number]

export type Resource = Note
