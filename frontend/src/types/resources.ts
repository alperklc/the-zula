import { Note } from "./Api.ts"

export const resourceTypes = ["note", "bookmark"] as const

export type ResourceType = typeof resourceTypes[number]

export type Resource = Note
