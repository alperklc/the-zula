import React from 'react'
import { FormattedMessage, useIntl } from 'react-intl'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import icons from '../../../components/icons'
import { useAuth } from '../../../contexts/authContext'
import { Api, Note } from '../../../types/Api.ts'
import { useUI } from '../../../contexts/uiContext'
import MessageBox from '../../../components/messageBox/index.tsx'
import useModal from '../../../components/modal/index.tsx'
import DeleteNoteConfirmation, { DeleteNoteConfirmationModalProps } from '../../../components/notes/deleteNoteModal/index.tsx'
import useDebounce from '../../../utils/useDebounce.tsx'
import Input from '../../../components/form/input/index.tsx'
import TagsInput, { searchTags } from '../../../components/tagsInput/index.tsx'
import Content from '../../../components/content/index.tsx'

const emptyNote: Note = {
  tags: [],
  title: '',
  content: '',
}

export const EditNote = () => {
  const intl = useIntl()

  const navigate = useNavigate()
  const { shortId } = useParams()
  const [searchParams] = useSearchParams()
  const loadDraft = searchParams.get("loadDraft") === "true"

  const { isMobile } = useUI()

  const [DeleteConfirmationModal, openDeleteModal, closeDeleteModal] =
    useModal<DeleteNoteConfirmationModalProps>(DeleteNoteConfirmation)

  const [note, setNote] = React.useState<Note>()
  const [loading, setLoading] = React.useState(true);
  const [errorLoading, setErrorLoading] = React.useState<string | null>(null);
  const [saving, setSaving] = React.useState(false);
  const [errorSaving, setErrorSaving] = React.useState<string | null>(null);

  const { user, sessionId } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}`, sessionId } } })

  const debouncedInput = useDebounce(note, 1000)

  React.useEffect(() => {
    if (note === emptyNote || note === undefined) {
      return
    }

    api.api.saveNoteDraft(shortId!,
      { title: note?.title, content: note?.content, tags: note?.tags },
    ).then(() => {
      console.log('saved draft')
    })
  }, [debouncedInput])


  const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target
    setNote({ ...note, title: value })
  }

  const handleTagsChange = (tags: string[]) => {
    setNote({ ...note, tags })
  }

  const handleContentChange = (content: string) => {
    setNote({ ...note, content })
  }

  const fetchNote = async (loadDraft: boolean) => {
    try {
      setLoading(true);
      setErrorLoading(null);

      const { data, status } = await api.api.getNote(shortId ?? "", loadDraft ? { loadDraft: true } : undefined)

      if (status === 200) {
        setNote(data);
      } else {
        console.error(data);
        setErrorLoading("could not fetch");
      }

    } catch (e: any) {
      console.error(e.error);
      setErrorLoading(e.error.message as string);
    }
    setLoading(false);
  };

  React.useEffect(() => {
    fetchNote(loadDraft ?? false)
  }, [])

  const save = async () => {
    setSaving(true);
    setErrorSaving(null);

    try {
      await api.api.updateNote(shortId!, {
        title: note?.title,
        tags: note?.tags,
        content: note?.content,
      })

      setSaving(false);
      navigate('/notes')

    } catch (e) {
      console.error(e);
    }
    setSaving(false);
  }

  const deleteNote = async () => {
    setSaving(true);
    setErrorSaving(null);

    try {
      await api.api.deleteNote(shortId!)
      
      setSaving(false);
      navigate('/notes')

    } catch (e) {
      console.error(e);
    }
    setSaving(false);
  }

  const openDeleteConfirmationModal = () => openDeleteModal()

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
                <Input
                  type='text'
                  placeholder={intl.formatMessage({ id: 'notes.form.label.title' })}
                  value={note?.title}
                  onChange={handleTitleChange}
                  title={note?.title}
                  className={layoutStyles.titleInSubHeader}
                />
              )}
            </div>

            {!errorLoading && (
              <div style={{ whiteSpace: 'nowrap' }}>
                <Button danger onClick={openDeleteConfirmationModal}>
                  <FormattedMessage id='common.buttons.delete' />
                </Button>
                <Button primary onClick={save}>
                  <FormattedMessage id='common.buttons.save' />
                </Button>
              </div>
            )}
          </div>
        </>
      }
    >
      {errorLoading ?
        <MessageBox type='error'>{errorLoading}</MessageBox> :
        <PageContent loading={loading} isMobile={isMobile}>
          <>
            {isMobile && (
              <Input
                type='text'
                placeholder={intl.formatMessage({ id: 'notes.form.label.title' })}
                value={note?.title}
                onChange={handleTitleChange}
                className={layoutStyles.title}
              />
            )}
            <Content
              data-testid='content'
              value={note?.content ?? ""}
              onChange={handleContentChange}
            />
          </>
          <>
            <TagsInput
              onSearch={searchTags('note')}
              tags={note?.tags || []}
              onChange={handleTagsChange}
              label='Tags'
              placeholder={intl.formatMessage({ id: 'notes.form.label.tags' })}
            />
          </>
        </PageContent>
      }
      <DeleteConfirmationModal onConfirm={deleteNote} />
    </Layout>
  )
}

export default EditNote
