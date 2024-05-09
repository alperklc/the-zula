import React from 'react'
import { useIntl, FormattedMessage } from 'react-intl'
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
import { MDXEditor, headingsPlugin, BoldItalicUnderlineToggles, toolbarPlugin, BlockTypeSelect, InsertTable, ListsToggle, listsPlugin, tablePlugin, CodeToggle, InsertThematicBreak, thematicBreakPlugin } from '@mdxeditor/editor'
import { Api, NoteInput } from '../../../types/Api'
import { useAuth } from '../../../contexts/authContext'
import '@mdxeditor/editor/style.css'

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
  const intl = useIntl()
  const { isMobile } = useUI()
  const navigate = useNavigate()
  const { show: showToast } = useToast()

  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })
  
  const [note, setNote] = React.useState<Note>(initialNote)
  
  const [saving, setSaving] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const draft = NoteCache.read<Note>('new-note')
    if (draft) {
      setNote(draft)
    }
  }, [])

  React.useEffect(() => {
    NoteCache.save<NoteInput>('new-note', {
      title: note.title,
      content: note.content,
      tags: note.tags,
    })
  }, [note])

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

  const save = async () => {
    setSaving(true);
    setError(null);
    
    try {
      await api.api.v1NotesCreate({
        title: note?.title,
        tags: note?.tags,
        content: note?.content,
      })
      debugger
      setSaving(false);
      NoteCache.remove('new-note')
      navigate('/notes')
        
    } catch (e: any) {
      console.error(e);
      showToast(intl.formatMessage({ id: 'messages.notes.create_failure' }), 'error')
    }
    setSaving(false);
    debugger
  }

  return (
    <Layout
      fixedSubHeader={!isMobile}
      subHeaderContent={
        <>
          {!isMobile && <Breadcrumbs />}
          <div className={layoutStyles.subheader}>
            <div className={layoutStyles.subheaderTitleContainer}>
              <Button onClick={history.back} className={layoutStyles.backButton}>
                <icons.ArrowLeft />
              </Button>

              {!isMobile && (
                <Input
                  type='text'
                  placeholder={intl.formatMessage({ id: 'notes.form.label.title' })}
                  value={note.title}
                  onChange={handleTitleChange}
                  title={note.title}
                  className={layoutStyles.titleInSubHeader}
                />
              )}
            </div>
            <Button primary onClick={save}>
              <FormattedMessage id='common.buttons.save' />
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
                placeholder={intl.formatMessage({ id: 'notes.form.label.title' })}
                value={note.title}
                onChange={handleTitleChange}
                className={layoutStyles.title}
              />
            </div>
          )}
          <MDXEditor
            markdown={note.content ?? ''}
            onChange={handleContentChange}
            plugins={[
              headingsPlugin(),
              listsPlugin(),
              tablePlugin(),
              thematicBreakPlugin(),
              toolbarPlugin({
                toolbarContents: () => (
                  <>
                    <BoldItalicUnderlineToggles />
                    <CodeToggle />
                    <ListsToggle />
                    <BlockTypeSelect />
                    <InsertTable />
                    <InsertThematicBreak />
                  </>
                )
              })
      
            ]}
          />
        </>
        <>
          <TagsInput
            onSearch={searchTags('note')}
            tags={note.tags}
            onChange={handleTagsChange}
            label='Tags'
            placeholder={intl.formatMessage({ id: 'notes.form.label.tags' })}
          />
        </>
      </PageContent>
    </Layout>
  )
}

export default CreateNote
