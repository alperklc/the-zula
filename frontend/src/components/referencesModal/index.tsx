import { GraphData } from 'react-force-graph-2d'
import { FormattedMessage } from 'react-intl'
import { modalStyles } from '../modal'
import { ReferencesGraph } from '../referencesGraph'
import { ResizeWrapper } from '../referencesGraph/resizeWrapper'
import styles from './index.module.css'

export const ReferencesModal = (props: {
  onModalClosed?: () => void
  references: GraphData
  noteUid: string
}) => {
  return (
    <div>
      <div className={modalStyles.modalHeader}>
        <FormattedMessage id='notes.form.label.referenced_by' />
      </div>
      <div>
        <ResizeWrapper fullScreen>
          {(graphProps) => {
            return (
              <ReferencesGraph
                width={innerWidth > 600 ? graphProps.width : window.innerWidth * 0.95}
                height={innerWidth > 600 ? graphProps.height : window.innerHeight * 0.95}
                noteUid={props.noteUid}
                graphData={props.references}
              />
            )
          }}
        </ResizeWrapper>
      </div>
      <div className={modalStyles.modalButtons}>&nbsp;</div>
    </div>
  )
}

export { styles }
