import React from 'react'
import { FormattedMessage } from 'react-intl'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import MessageBox from '../../../components/messageBox'
import MarkdownDisplay from '../../../components/markdownDisplay'
import TimeDisplay from '../../../components/timeDisplay'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import icons from '../../../components/icons'
import { Api, NoteChange } from '../../../types/Api'
import { useAuth } from '../../../contexts/authContext'
import { useNavigate, useParams } from 'react-router-dom'
import { useUI } from '../../../contexts/uiContext'

export const NoteChangePage = () => {
  const navigate = useNavigate()
  const { isMobile } = useUI()
  const { shortId, timestamp } = useParams()
  const { user } = useAuth()

  const [noteChange, setNoteChange] = React.useState<NoteChange>()
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })
  const fetchNoteChange = async (shortId: string) => {
    try {
      setLoading(true);
      setError(null);

      const { data, status } = await api.api.getNotesChange(shortId ?? "", timestamp || "")

      if (status === 200) {
        setNoteChange(data);
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
    shortId && fetchNoteChange(shortId)
  }, [shortId])

  return (
    <Layout
      fixedSubHeader={!isMobile}
      showScrollUpButton
      subHeaderContent={
        <>
          {!isMobile && <Breadcrumbs />}

          <div className={layoutStyles.subheader}>
            <div className={layoutStyles.subheaderTitleContainer}>
              <Button onClick={() => navigate(-1)} className={layoutStyles.backButton}>
                <icons.ArrowLeft />
              </Button>
              {!isMobile && (
                <span data-testid='title' className={layoutStyles.titleInSubHeader}>
                  Change
                </span>
              )}
            </div>
            <span />
          </div>
        </>
      }
    >
      {error ? (
        <MessageBox type='error'>Error</MessageBox>
      ) : (
        <PageContent loading={loading} isMobile={isMobile}>
          <>
            {isMobile && (
              <div className={layoutStyles.flex}>
                <span data-testid='title' className={layoutStyles.title}>
                  Change
                </span>
              </div>
            )}
            <MarkdownDisplay
              className={layoutStyles.htmlContentOfBookmark}
              content={noteChange?.change || ''}
            />
          </>

          <>
            <div>
              <label>
                <FormattedMessage id='common.labels.updated_at' />
              </label>
              <div className={layoutStyles.secondaryText}>
                <span>{noteChange?.updatedBy}, </span>
                <TimeDisplay isoDate={noteChange?.updatedAt || ""} />
              </div>
            </div>
          </>
        </PageContent>
      )}
    </Layout>
  )
}

export default NoteChangePage
