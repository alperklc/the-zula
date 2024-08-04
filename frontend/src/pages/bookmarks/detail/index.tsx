import React, { ChangeEvent } from 'react'
import classNames from 'classnames'
import { useIntl, FormattedMessage } from 'react-intl'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import Button from '../../../components/form/button'
import MessageBox from '../../../components/messageBox'
import TagsInput, { searchTags } from '../../../components/tagsInput'
import useModal from '../../../components/modal'
import { modalStyles } from '../../../components/modal'
import MarkdownDisplay from '../../../components/markdownDisplay'
import icons from '../../../components/icons'
import { useToast } from '../../../components/toast/toast-message-context'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import { useUI } from '../../../contexts/uiContext'
import { useNavigate, useParams } from 'react-router-dom'
import { useAuth } from '../../../contexts/authContext'
import { Api, Bookmark } from '../../../types/Api'
import TimeDisplay from '../../../components/timeDisplay'
import Input from '../../../components/form/input'

const initialBookmark: Bookmark = {
  url: '',
  tags: [],
  title: '',
  createdAt: '',
  updatedAt: '',
  pageContent: {
    mdContent: '',
    favicon: '',
  },
  shortId: ''
}

const DeleteBookmarkConfirmation = (props: {
  onConfirm: () => void
  onModalClosed?: () => void
}) => {
  return (
    <div>
      <div className={modalStyles.modalHeader}>&nbsp;</div>
      <div className={modalStyles.modalBody}>
        <FormattedMessage id='delete_confirmation_modal.title' />
      </div>
      <div className={modalStyles.modalButtons}>
        <Button danger onClick={props.onConfirm}>
          <FormattedMessage id='common.buttons.delete' />
        </Button>
        <Button outline onClick={props.onModalClosed}>
          <FormattedMessage id='common.buttons.cancel' />
        </Button>
      </div>
    </div>
  )
}

const BookmarkDetails = () => {
  const intl = useIntl()
  const { isMobile } = useUI()
  const navigate = useNavigate()
  const { shortId }  = useParams()

  const { show: showToast } = useToast()

  const [DeleteConfirmationModal, openDeleteModal, closeDeleteModal] = useModal<any>(DeleteBookmarkConfirmation)

  const [fetching, setFetching] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const [bookmark, setBookmark] = React.useState<Bookmark>(initialBookmark)
  const bookmarkRef = React.useRef(initialBookmark)

  const fetchBookmark = async () => {
    try {
      setFetching(true);
      setError(null);

      const { data, status } = await api.api.getBookmark(shortId ?? "")

      if (status === 200) {
        setBookmark(data);
      } else {
        console.error(data);
        setError("failed to load");
      }

    } catch (e: unknown) {
      console.error(e);
      setError(e as string);
    }
    setFetching(false);
  };

  React.useEffect(() => {
    fetchBookmark()
  }, [])

  React.useEffect(() => {
    bookmarkRef.current = bookmark
  }, [bookmark])

  const handleTagsChange = (tags: string[]) => {
    setBookmark({ ...bookmark, tags })
  }

  const handleTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
    setBookmark({ ...bookmark, title: event.target.value })
  }

  const save = async () => {
    setFetching(true);
    
    try {
      await api.api.updateBookmark(shortId!, {
        title: bookmark?.title,
        tags: bookmark?.tags,
        url: bookmark?.url,
      })
      
      navigate('/bookmarks')
    } catch (e) {
      console.error(e);
      showToast(intl.formatMessage({ id: 'messages.bookmarks.update_failure' }), 'error')
    }
    setFetching(false);
  }

  const deleteBookmark = async () => {
    setFetching(true);
    
    try {
      await api.api.deleteBookmark(shortId!)
      closeDeleteModal()
      
      navigate('/bookmarks')
    } catch (e) {
      console.error(e);
      showToast(intl.formatMessage({ id: 'messages.bookmarks.delete_failure' }), 'error')
    }
    setFetching(false);
  }

  const openDeleteConfirmationModal = () => openDeleteModal()

  const onCopyLinkClick = (event: React.MouseEvent) => {
    navigator?.clipboard?.writeText(bookmark?.url || '').then(() => {
      showToast(intl.formatMessage({ id: 'bookmarks.toast.copy_link' }), 'info')
    })

    event.stopPropagation()
  }

  const onShareClick = (event: React.MouseEvent) => {
    if (navigator?.share) {
      navigator
        ?.share({
          title: bookmark?.title || '',
          url: bookmark?.url || '',
        })
        .then(() => console.log('Share was successful.'))
        .catch((error) => console.log('Sharing failed', error))
    }

    event.stopPropagation()
  }

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
                <Input
                  type='text'
                  placeholder={intl.formatMessage({ id: 'notes.form.label.title' })}
                  value={bookmark.title}
                  onChange={handleTitleChange}
                  title={bookmark.title}
                  className={classNames(layoutStyles.title, layoutStyles.truncatedText)}
                />
              )}
            </div>
            {bookmark.url !== "" && (
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
      <PageContent loading={fetching} isMobile={isMobile}>
        <>
          {!fetching && error && <MessageBox type='error'>{error}</MessageBox>}
          {isMobile && bookmark.title && (
            <Input
              type='text'
              placeholder={intl.formatMessage({ id: 'notes.form.label.title' })}
              value={bookmark.title}
              onChange={handleTitleChange}
              title={bookmark.title}
              className={classNames(layoutStyles.title, layoutStyles.truncatedText)}
            />
          )}
          {bookmark?.pageContent?.mdContent ? (
            <MarkdownDisplay
              className={layoutStyles.htmlContentOfBookmark}
              content={bookmark?.pageContent?.mdContent ?? ""}
            />
          ) : (
            <MessageBox type='info'>Does not have content.</MessageBox>
          )}
        </>

        <>
          <label className={layoutStyles.rightPanelSectionTitle}>
            <FormattedMessage id='common.labels.updated_at' />
          </label>
          <div className={layoutStyles.secondaryText}>
            <span>{bookmark?.updatedBy}, </span>
            <TimeDisplay isoDate={bookmark?.updatedAt} />
          </div>
          <hr />
          <label className={layoutStyles.rightPanelSectionTitle}>
            <FormattedMessage id='bookmarks.form.label.url' />

            {!isMobile && (
              <a
                className={classNames(layoutStyles.copyLink, layoutStyles.labelLink)}
                onClick={onCopyLinkClick}
              >
                <FormattedMessage id='bookmarks.overflow_menu.copy_link' />
                <icons.Link height='.6rem' />
              </a>
            )}
          </label>
          <div style={{ display: 'flex', wordBreak: 'break-all' }}>
            <a
              href={bookmark.url}
              target='_blank'
              rel='noreferrer'
              className={classNames(layoutStyles.secondaryText)}
            >
              {bookmark.url}
            </a>
          </div>

          {isMobile && (
            <>
              <hr />
              <label className={layoutStyles.rightPanelSectionTitle} onClick={onShareClick}>
                <span className={layoutStyles.truncatedText}>
                  <FormattedMessage id='bookmarks.overflow_menu.share' />
                </span>
                <span>
                  <icons.Share width='1rem' height='1rem' />
                </span>
              </label>
            </>
          )}

          <hr />
          <TagsInput
            onSearch={searchTags('bookmark')}
            tags={bookmark.tags ?? []}
            onChange={handleTagsChange}
            label='Tags'
            placeholder={intl.formatMessage({ id: 'bookmarks.form.label.tags' })}
          />
        </>
      </PageContent>

      <DeleteConfirmationModal onConfirm={deleteBookmark} />
    </Layout>
  )
}

export default BookmarkDetails
