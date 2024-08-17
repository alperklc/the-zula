import { Command, CommandContext, CommandMap } from '../types'
import { extractKeyActivatedCommands } from './command-utils'
import * as React from 'react'
import { getDefaultCommandMap } from './default-commands/defaults'
import { TextApi, TextState } from '../types/CommandOptions'
import { insertText } from '../util/InsertTextAtPosition'
import { Selection } from '../types/Selection'

export class TextAreaTextApi implements TextApi {
  textAreaRef: React.RefObject<HTMLTextAreaElement>

  constructor(textAreaRef: React.RefObject<HTMLTextAreaElement>) {
    this.textAreaRef = textAreaRef
  }

  replaceSelection(text: string): TextState {
    const textArea = this.textAreaRef.current as HTMLTextAreaElement
    insertText(textArea, text)
    return getStateFromTextArea(textArea)
  }

  setSelectionRange(selection: Selection): TextState {
    const textArea = this.textAreaRef.current as HTMLTextAreaElement
    textArea.focus()
    textArea.selectionStart = selection.start
    textArea.selectionEnd = selection.end
    return getStateFromTextArea(textArea)
  }

  getState(): TextState {
    const textArea = this.textAreaRef.current as HTMLTextAreaElement
    return getStateFromTextArea(textArea)
  }
}

export function getStateFromTextArea(textArea: HTMLTextAreaElement): TextState {
  return {
    selection: {
      start: textArea.selectionStart,
      end: textArea.selectionEnd,
    },
    text: textArea.value,
    selectedText: textArea.value.slice(textArea.selectionStart, textArea.selectionEnd),
  }
}

export class CommandOrchestrator {
  private readonly textAreaRef: React.RefObject<HTMLTextAreaElement>
  private readonly textApi: TextApi
  private readonly commandMap: CommandMap
  /**
   * Names of commands that can be activated by the keyboard
   */
  keyActivatedCommands: string[]
  /**
   * Indicates whether there is a command currently executing
   */
  isExecuting = false

  constructor(customCommands: CommandMap, textArea: React.RefObject<HTMLTextAreaElement>) {
    this.commandMap = { ...getDefaultCommandMap(), ...(customCommands || {}) }
    this.keyActivatedCommands = extractKeyActivatedCommands(customCommands)
    this.textAreaRef = textArea
    this.textApi = new TextAreaTextApi(textArea)
  }

  getCommand = (name: string): Command => {
    const command = this.commandMap[name]
    if (!command) {
      throw new Error(`Cannot execute command. Command not found: ${name}`)
    }

    return command
  }

  /**
   * Tries to find a command the wants to handle the keyboard event.
   * If a command is found, it is executed and the function returns
   */
  handlePossibleKeyCommand = (e: React.KeyboardEvent<HTMLTextAreaElement>): boolean => {
    for (const commandName of this.keyActivatedCommands) {
      if (this.getCommand(commandName).handleKeyCommand?.(e)) {
        this.executeCommand(commandName).then((r: any) => ({}))
        return true
      }
    }
    return false
  }

  async executeCommand(commandName: string, context?: CommandContext): Promise<void> {
    if (this.isExecuting) {
      // The simplest thing to do is to ignore commands while
      // there is already a command executing. The alternative would be to queue commands
      // but there is no guarantee that the state after one command executes will still be compatible
      // with the next one. In fact, it is likely not to be.
      return
    }

    this.isExecuting = true
    const command = this.commandMap[commandName]
    const result = command.execute({
      initialState: getStateFromTextArea(this.textAreaRef.current as HTMLTextAreaElement),
      textApi: this.textApi,
      context,
    })
    await result
    this.isExecuting = false
  }

  /**
   * Returns a command by name
   * @param name
   */
  getCommandByName(name: string) {
    return this.commandMap[name]
  }
}
