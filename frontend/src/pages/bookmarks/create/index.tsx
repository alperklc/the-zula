import React from 'react'
import { useIntl, FormattedMessage } from 'react-intl'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import Input from '../../../components/form/input'
import TagsInput, { searchTags } from '../../../components/tagsInput'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import icons from '../../../components/icons'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../../contexts/authContext'
import { Api } from '../../../types/Api'
import { useToast } from '../../../components/toast/toast-message-context'
import { useUI } from '../../../contexts/uiContext'

interface Bookmark {
  url: string
  title: string
  faviconUrl: string
  tags: string[]
}

const initialBookmark: Bookmark = {
  url: '',
  tags: [],
  title: '',
  faviconUrl: '',
}

const CreateBookmark = () => {
  const intl = useIntl()
  const navigate = useNavigate()
  const { show: showToast } = useToast()
  const { isMobile } = useUI()

  const [saving, setSaving] = React.useState(false);

  const { user, sessionId } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}`, sessionId } } })

  const [bookmark, setBookmark] = React.useState<Bookmark>(initialBookmark)

  const handleUrlChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target
    setBookmark({ ...bookmark, url: value })
  }

  const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target
    setBookmark({ ...bookmark, title: value })
  }

  const handleTagsChange = (tags: string[]) => {
    setBookmark({ ...bookmark, tags })
  }

  const save = async () => {
    setSaving(true);
    
    try {
      await api.api.createBookmark({
        title: bookmark?.title,
        tags: bookmark?.tags,
        url: bookmark?.url,
      })
      
      setSaving(false);
      navigate('/bookmarks')
        
    } catch (e) {
      console.error(e);
      showToast(intl.formatMessage({ id: 'messages.bookmarks.create_failure' }), 'error')
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
            <Button onClick={() => navigate(-1)} className={layoutStyles.backButton}>
              <icons.ArrowLeft />
            </Button>
            <Button primary onClick={save}>
              <FormattedMessage id='common.buttons.save' />
            </Button>
          </div>
        </>
      }
    >
      <PageContent loading={saving} isMobile={isMobile}>
        <>
          <Input
            type='text'
            value={bookmark.url}
            onChange={handleUrlChange}
            label={intl.formatMessage({ id: 'bookmarks.form.label.url' })}
          />
          <Input
            type='text'
            value={bookmark.title}
            onChange={handleTitleChange}
            label={intl.formatMessage({ id: 'bookmarks.form.label.title' })}
          />
          <TagsInput
            onSearch={searchTags('bookmark')}
            tags={bookmark.tags}
            onChange={handleTagsChange}
            label='Tags'
            placeholder={intl.formatMessage({ id: 'bookmarks.form.label.tags' })}
          />
        </>
      </PageContent>
    </Layout>
  )
}

export default CreateBookmark
