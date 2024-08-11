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
import { MDXEditor } from '@mdxeditor/editor'
import { useUI } from '../../../contexts/uiContext'
import MessageBox from '../../../components/messageBox/index.tsx'

export const EditNote = () => {
  const navigate = useNavigate()
  const { shortId } = useParams()
  const { isMobile } = useUI()

  const [note, setNote] = React.useState<Note>()

  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const fetchNote = async () => {
    try {
      setLoading(true);
      setError(null);

      const { data, status } = await api.api.getNote(shortId ?? "")

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
    fetchNote()
  }, [])


  const onEditClicked = () => {
    const shouldLoadDraft = note?.hasDraft
      ? confirm('You have an unsaved draft, would you like to load it?')
      : false

    navigate(`/notes/${shortId}/edit${shouldLoadDraft ? '?loadDraft=true' : ''}`)
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
            {note?.content && <MDXEditor markdown={note?.content} readOnly />}
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
          </>
        </PageContent>
      }
    </Layout>
  )
}

export default EditNote
