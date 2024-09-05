import React from 'react'
import { useTranslation } from 'react-i18next'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import MessageBox from '../../../components/messageBox'
import TimeDisplay from '../../../components/timeDisplay'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import icons from '../../../components/icons'
import { Api, NoteChange } from '../../../types/Api'
import { useAuth } from '../../../contexts/authContext'
import { useNavigate, useParams } from 'react-router-dom'
import { useUI } from '../../../contexts/uiContext'
import { Diff2HtmlUI, Diff2HtmlUIConfig } from 'diff2html/lib/ui/js/diff2html-ui-slim.js'
import 'diff2html/bundles/css/diff2html.min.css';

export const NoteChangePage = () => {
  const navigate = useNavigate()
  const { isMobile } = useUI()
  const { shortId, timestamp } = useParams()
  const { user } = useAuth()
  const { t } = useTranslation()

  const changeArea = React.useRef<HTMLDivElement>(null);
  const [noteChange, setNoteChange] = React.useState<NoteChange>()
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })
  const fetchNoteChange = async (shortId: string) => {
    try {
      setLoading(true);
      setError(null);

      const { data, error, status } = await api.api.getNotesChange(shortId ?? "", timestamp || "")

      if (status === 200) {
        setNoteChange(data);
      } else {
        console.error(error);
        setError(error);
      }

    } catch (e) {
      console.error(e);
      setError(e as string);
    }
    setLoading(false);
  };

  React.useEffect(() => {
    shortId && fetchNoteChange(shortId)
  }, [shortId])

  const configuration = {
    drawFileList: false,
    fileListToggle: false,
    fileListStartVisible: false,
    fileContentToggle: false,
    matching: 'lines',
    outputFormat: 'line-by-line',
    synchronisedScroll: true,
    highlight: true,
    renderNothingWhenEmpty: false,
  };

  React.useEffect(() => {    
    if (!loading) {
      const diff2htmlUi = new Diff2HtmlUI(changeArea.current!, noteChange?.change, configuration as Diff2HtmlUIConfig);
      diff2htmlUi.draw();
    }
  }, [noteChange?.change, loading])

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
            <div ref={changeArea}></div>
          </>

          <>
            <div>
              <label>
                {t('common.labels.updated_at')}
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
