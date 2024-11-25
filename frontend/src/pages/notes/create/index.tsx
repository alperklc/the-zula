import React from 'react'
import { useTranslation } from 'react-i18next'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import Input from '../../../components/form/input'
import TagsInput, { searchTags } from '../../../components/tagsInput'
import { useToast } from '../../../components/toast/toast-message-context'
import NoteCache from '../../../components/noteCache'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import icons from '../../../components/icons'
import { useUI } from '../../../contexts/uiContext'
import { useNavigate } from 'react-router-dom'
import { Api } from '../../../types/Api.ts'
import { useAuth } from '../../../contexts/authContext'
import Content from '../../../components/content/index.tsx'

interface Note {
  tags: string[]
  title: string
  content: string
}

const initialNote: Note = {
  tags: [],
  title: '',
  content: '',
}

const CreateNote = () => {
  const { t } = useTranslation()
  const { isMobile } = useUI()
  const navigate = useNavigate()
  const { show: showToast } = useToast()

  const { user, sessionId } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}`, sessionId } } })
  
  const [note, setNote] = React.useState<Note>(initialNote)
  
  const [saving, setSaving] = React.useState(false);
  const [_, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const title = NoteCache.read<string>('new-note-title') ?? ""
    const content = NoteCache.read<string>('new-note-content') ?? ""
    const tags = NoteCache.read<string[]>('new-note-tags') ?? []

    if (title?.length || content?.length || tags?.length) {
      setNote({ title, content, tags})
    }
  }, [])

  const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target
    NoteCache.save<string>('new-note-title', value)
    setNote({ ...note, title: value })
  }

  const handleTagsChange = (tags: string[]) => {
    NoteCache.save<string[]>('new-note-tags', tags)
    setNote({ ...note, tags })
  }

  const handleContentChange = (content: string) => {
    NoteCache.save<string>('new-note-content', content)
    setNote({ ...note, content })
  }

  const save = async () => {
    setSaving(true);
    setError(null);
    
    try {
      await api.api.createNote({
        title: note?.title,
        tags: note?.tags,
        content: note?.content,
      })
      
      setSaving(false);
      NoteCache.remove('new-note-title')
      NoteCache.remove('new-note-content')
      NoteCache.remove('new-note-tags')
      navigate('/notes')
        
    } catch (e) {
      console.error(e);
      showToast(t('messages.notes.create_failure'), 'error')
    }
    setSaving(false);
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
                <Input
                  type='text'
                  placeholder={t('notes.form.label.title')}
                  value={note.title}
                  onChange={handleTitleChange}
                  title={note.title}
                  className={layoutStyles.titleInSubHeader}
                />
              )}
            </div>
            <Button primary onClick={save}>
              {t('common.buttons.save')}
            </Button>
          </div>
        </>
      }
    >
      <PageContent loading={saving} isMobile={isMobile}>
        <>
          {isMobile && (
            <div className={layoutStyles.flex}>
              <Input
                type='text'
                placeholder={t('notes.form.label.title')}
                value={note.title}
                onChange={handleTitleChange}
                className={layoutStyles.title}
              />
            </div>
          )}
          <Content value={note.content} onChange={handleContentChange} />
        </>
        <>
          <TagsInput
            onSearch={searchTags('note')}
            tags={note.tags}
            onChange={handleTagsChange}
            label='Tags'
            placeholder={t('notes.form.label.tags')}
          />
        </>
      </PageContent>
    </Layout>
  )
}

export default CreateNote
