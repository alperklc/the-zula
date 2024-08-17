import React from 'react'
import { Toolbar, ToolbarButtonData } from './toolbar'
import { SvgIcon } from './icons'
import { CommandOrchestrator } from './commands/command-orchestrator'
import { TextArea } from './TextArea'
import MarkdownDisplay from '../markdownDisplay'
import {
  getDefaultToolbarCommands,
  getDefaultCommandMap,
} from './commands/default-commands/defaults'
import useDebounce from '../../utils/useDebounce'
import { useAuth } from '../../contexts/authContext'
import { Api } from '../../types/Api'
import styles from './index.module.css'
import { Suggestion } from './types'
import { filterEmptyValues } from '../../utils/filter'

interface ContentProps {
  value: string
  onChange: (_: string) => void
  label?: string
}

const Content = (props: ContentProps) => {
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const textareaRef = React.useRef<HTMLTextAreaElement>(null)
  const [topPosition, setTopPosition] = React.useState(100)
  const debouncedTopOfTextArea = useDebounce<number>(topPosition, 100)
  const [selectedTab, setSelectedTab] = React.useState<'write' | 'preview'>('write')

  const handleContentChange = (content: string) => {
    props.onChange(content)
  }

  const calculateTopOfTextArea = () => {
    const startPoint = textareaRef.current?.getBoundingClientRect()?.top || 200

    setTopPosition(Math.round(startPoint))
  }

  const setEditorHeight = () => {
    if (textareaRef?.current?.style) {
      textareaRef.current.style.height = `calc(100vh - ${topPosition}px - 1.2rem)`
      textareaRef.current.style.paddingBottom = `calc(calc(100vh - ${topPosition}px) - 3.2rem)`
    }
  }

  React.useEffect(setEditorHeight, [debouncedTopOfTextArea, topPosition])

  React.useEffect(() => {
    calculateTopOfTextArea()
    window.addEventListener('resize', calculateTopOfTextArea)

    return () => {
      window.removeEventListener('resize', calculateTopOfTextArea)
    }
  }, [])

  const commandOrchestrator = new CommandOrchestrator(getDefaultCommandMap(), textareaRef)

  const toolbarButtons = getDefaultToolbarCommands().map((group) => {
    return group.map((commandName) => {
      const command = commandOrchestrator.getCommand(commandName)
      return {
        commandName: commandName,
        buttonContent: <SvgIcon icon={commandName} />,
        buttonProps: command.buttonProps,
        buttonComponentClass: command.buttonComponentClass,
      }
    })
  })

  const handleCommand = async (commandName: string) => {
    await commandOrchestrator.executeCommand(commandName)
  }

  const fetchSuggestions = () => async (q: string, _:string): Promise<Suggestion[]> => {
    let response: Suggestion[] = [];

    const filteredQuery = filterEmptyValues({ q, page:1, pageSize: 5,  sortBy: 'title'});
    const { data } = await api.api.getNotes(filteredQuery);
    const { items } = data
        
    response = (items || []).map((item) => ({
      preview: item.title,
      value: `[${item.title}](/notes/${item.shortId})`,
    }));
      
    return Promise.resolve(response);
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <span>{props.label && <label>{props.label}</label>}</span>
      </div>
      <Toolbar
        className={styles.toolbar}
        buttons={toolbarButtons as ToolbarButtonData[][]}
        onCommand={handleCommand}
        onTabChange={setSelectedTab}
        tab={selectedTab}
        readOnly={false}
        disablePreview={false}
      />
      <div className={`${selectedTab !== 'write' ? styles.invisible: '' }`}>
        <TextArea
          className={styles.textArea}
          suggestionsAutoplace={true}
          refObject={textareaRef}
          onChange={handleContentChange}
          textAreaProps={{
            spellCheck: 'false',
          }}
          height={200}
          heightUnits={'px'}
          value={props.value}
          suggestionTriggerCharacters={['@']}
          loadSuggestions={fetchSuggestions()}
          onPossibleKeyCommand={undefined}
        />
      </div>

      {selectedTab !== 'write' && (
        <MarkdownDisplay className={styles.preview} content={props.value || ''} />
      )}
    </div>
  )
}

export default Content
