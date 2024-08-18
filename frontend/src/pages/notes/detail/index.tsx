import React from 'react'
import { FormattedMessage } from 'react-intl'
import { useNavigate, useParams } from 'react-router-dom'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import icons from '../../../components/icons'
import { useAuth } from '../../../contexts/authContext'
import { Api, Note } from '../../../types/Api.ts'
import TagsDisplay from '../../../components/tagsDisplay'
import TimeDisplay from '../../../components/timeDisplay'
import { useUI } from '../../../contexts/uiContext'
import MessageBox from '../../../components/messageBox/index.tsx'
import MarkdownDisplay from '../../../components/markdownDisplay/index.tsx'
import ReferencesGraph from '../../../components/referencesGraph/index.tsx'
import { ResizeWrapper } from '../../../components/referencesGraph/resizeWrapper.tsx'
import useModal from '../../../components/modal/index.tsx'
import { ReferencesModal, styles } from '../../../components/referencesModal/index.tsx'
import { GraphData } from 'react-force-graph-2d'

export const EditNote = () => {
  const navigate = useNavigate()
  const { shortId } = useParams()
  const { isMobile } = useUI()
  const { user, sessionId } = useAuth()

  const [note, setNote] = React.useState<Note>()

  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [ReferencesExpanded, openReferencesModal, closeReferencesModal] =
    useModal<{ noteUid: string; references: GraphData }>(ReferencesModal)

  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}`, sessionId } } })

  const fetchNote = async (shortId: string) => {
    try {
      setLoading(true);
      setError(null);

      const { data, status } = await api.api.getNote(shortId ?? "", { loadDraft: false, getHistory: true, getReferences: true } )

      if (status === 200) {
        setNote(data);
      } else {
        console.error(data);
        setError("could not fetch");
      }

    } catch (e: any) {
      console.error(e.error);
      setError(e.error.message as string);
    }
    setLoading(false);
  };

  React.useEffect(() => {
    closeReferencesModal()
    shortId && fetchNote(shortId)
  }, [shortId])


  const onEditClicked = () => {
    const shouldLoadDraft = note?.hasDraft
      ? confirm('You have an unsaved draft, would you like to load it?')
      : false

    navigate(`/notes/${shortId}/edit${shouldLoadDraft ? '?loadDraft=true' : ''}`)
  }

  const expandReferences = () => {
    openReferencesModal()
  }

  return (
    <Layout
      fixedSubHeader={!isMobile}
      subHeaderContent={
        <>
          {!isMobile && <Breadcrumbs />}
          <div className={layoutStyles.subheader}>
            <div className={layoutStyles.subheaderTitleContainer}>
              <Button onClick={() => navigate(-1)} className={layoutStyles.backButton}>
                <icons.ArrowLeft />
              </Button>
              {!isMobile && (
                <span
                  data-testid='title'
                  className={layoutStyles.titleInSubHeader}
                  title={note?.title}
                >
                  {note?.title}
                </span>
              )}
            </div>

            {!error && (
              <Button onClick={onEditClicked}>
                <FormattedMessage id='common.buttons.edit' />
              </Button>
            )}
          </div>
        </>
      }
    >
      {error ?
        <MessageBox type='error'>{error}</MessageBox> :
        <PageContent loading={loading} isMobile={isMobile}>
          <>
            {isMobile && (
              <div className={layoutStyles.flex}>
                <span data-testid='title' className={layoutStyles.title}>
                  {note?.title}
                </span>
              </div>
            )}
            <MarkdownDisplay
              className={layoutStyles.htmlContentOfBookmark}
              content={note?.content || ''}
            />
          </>
          <>
            <div>
              <label>
                <FormattedMessage id='common.labels.updated_at' />
              </label>
              <div className={layoutStyles.secondaryText}>
                <span>{note?.updatedBy}, </span>
                <TimeDisplay isoDate={note?.updatedAt ?? ""} />
              </div>
            </div>
            {note?.tags && note?.tags.length > 0 && (
              <>
                <hr />

                <label>
                  <FormattedMessage id='notes.form.label.tags' />
                </label>
                <TagsDisplay tags={note.tags} />
              </>
            )}
            {note?.references?.links && note?.references?.links.length > 0 && (
             <>
               <hr />
               <div className={layoutStyles.refernencesLabel}>
                 <label>
                   <FormattedMessage id='notes.form.label.referenced_by' />
                 </label>
                 <label className={layoutStyles.labelLink} onClick={expandReferences}>
                   expand
                 </label>
               </div>
               <ResizeWrapper>
                 {(props: any) => (
                   <ReferencesGraph
                     {...props}
                     noteId={shortId}
                     graphData={note?.references}
                   />
                 )}
               </ResizeWrapper>
             </>
             )}
          </>
        </PageContent>
      }
      <ReferencesExpanded
        className={styles.modal}
        noteUid={shortId as string}
        references={note?.references as GraphData}
      />
    </Layout>
  )
}

export default EditNote
