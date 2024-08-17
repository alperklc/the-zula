import * as React from 'react'
import classNames from 'classnames'
import { ToolbarButton } from './ToolbarButton'

import styles from './index.module.css'

type Tab = 'write' | 'preview'

export interface ToolbarButtonData {
  commandName: string
  buttonContent: React.ReactNode
  buttonProps: any
  buttonComponentClass: React.ComponentClass | string
}

export interface ToolbarProps {
  className?: string
  buttons: ToolbarButtonData[][]
  onCommand: (commandName: string) => void
  onTabChange: (tab: Tab) => void
  readOnly: boolean
  disablePreview: boolean
  tab: Tab
  buttonProps?: any
  children?: React.ReactNode
}

export class Toolbar extends React.Component<ToolbarProps> {
  handleTabChange = (tab: Tab) => {
    const { onTabChange } = this.props
    onTabChange(tab)
  }

  render() {
    const { className, children, buttons, onCommand, readOnly, disablePreview, buttonProps } =
      this.props
    if ((!buttons || buttons.length === 0) && !children) {
      return null
    }

    const writePreviewTabs = (
      <div className={styles.tabs}>
        <button
          type='button'
          className={classNames({ [styles.selected]: this.props.tab === 'write' })}
          onClick={() => this.handleTabChange('write')}
        >
          {'write'}
        </button>
        <button
          type='button'
          className={classNames({ [styles.selected]: this.props.tab === 'preview' })}
          onClick={() => this.handleTabChange('preview')}
        >
          {'preview'}
        </button>
      </div>
    )

    return (
      <div className={classNames(styles.header, className)}>
        {!disablePreview && writePreviewTabs}
        {buttons.map((commandGroup: ToolbarButtonData[], i: number) => (
          <ul
            key={i}
            className={classNames(styles.headerGroup, {
              [styles.hidden]: this.props.tab === 'preview',
            })}
          >
            {commandGroup.map((c: ToolbarButtonData, j) => {
              return (
                <ToolbarButton
                  key={j}
                  name={c.commandName}
                  buttonContent={c.buttonContent}
                  buttonProps={{ ...(buttonProps || {}), ...c.buttonProps }}
                  onClick={() => onCommand(c.commandName)}
                  readOnly={readOnly}
                  buttonComponentClass={c.buttonComponentClass}
                />
              )
            })}
          </ul>
        ))}
      </div>
    )
  }
}
