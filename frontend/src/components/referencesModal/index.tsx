import { GraphData } from 'react-force-graph-2d'
import { modalStyles } from '../modal'
import { ReferencesGraph } from '../referencesGraph'
import { ResizeWrapper } from '../referencesGraph/resizeWrapper'
import styles from './index.module.css'
import { useTranslation } from 'react-i18next'

export const ReferencesModal = (props: {
  onModalClosed?: () => void
  references: GraphData
  noteUid: string
}) => {
  const { t } = useTranslation()

  return (
    <div>
      <div className={modalStyles.modalHeader}>
        {t('notes.form.label.referenced_by')}
      </div>
      <div>
        <ResizeWrapper fullScreen>
          {(graphProps) => {
            return (
              <ReferencesGraph
                width={innerWidth > 600 ? graphProps.width : window.innerWidth * 0.95}
                height={innerWidth > 600 ? graphProps.height : window.innerHeight * 0.95}
                noteId={props.noteUid}
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
