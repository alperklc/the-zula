import React from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import icons from '../../../components/icons'
import Button from '../../../components/form/button'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import Tabs, { SettingsTabs } from '../../../components/tabs'
import { useUI } from '../../../contexts/uiContext'
import { Api } from '../../../types/Api'
import { useAuth } from '../../../contexts/authContext'
import Input from '../../../components/form/input'
import { useToast } from '../../../components/toast/toast-message-context'

const Data = () => {
  const { show: showToast } = useToast()
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })
  
  const navigate = useNavigate()
  const { t } = useTranslation()
  const { isMobile } = useUI()

  const [file, setFile] = React.useState<File | null>(null);
  const [busy, setBusy] = React.useState(false);

  const handleFileChange: React.ChangeEventHandler<HTMLInputElement> = (event) => {
    const [file] = event.target.files || []
    if (file) {
      setFile(file);
    }
  };

  const handleSubmit = async () => {
    if (!file) {
      showToast('Please select a file before uploading.', "error");
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      setBusy(true);
      const response = await api.api.importData({ file });
      if (response.status === 200) {
        showToast('File uploaded successfully!', "success");
      } else {
        showToast('File upload failed.', "error");
      }
    } catch (error) {
      showToast('Error during upload: ' + error, "error");
    } finally {
      setBusy(false);
    }
  };

  const requestDataExport = async () => {
    try {
      setBusy(true);
      const response = await api.api.exportData();
      if (response.status === 200) {
        showToast('Data dump successfully created!', "success");
      } else {
        showToast('Creation of data dump failed.', "error");
      }
    } catch (error) {
      showToast('Error during upload: ' + error, "error");
    } finally {
      setBusy(false);
    }
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
          </div>
        </>
      }
    >
      <PageContent
        loading={busy}
        isMobile={isMobile}
        tabs={<Tabs selectedTab={SettingsTabs.DATA} />}
      >
        <div>
          <label>
            {t('data.import.label')}
          </label>
          <br />
            {t('data.import.info')}
          <pre>
            <div>/</div>
            <div>/bookmarks</div>
            <div>/notes</div>
            <div>/notes_changes</div>
            <div>/notes_drafts</div>
            <div>/page-content</div>
            <div>/references</div>
            <div>/users-activity</div>
          </pre>
          
          {t('data.import.filenames')}
          <pre>
            [mongodb-object-id].json
          </pre>

          <Input
            type="file"
            id='file'
            accept=".zip"
            onChange={handleFileChange}
            disabled={busy}
          />

          <div><Button 
            onClick={handleSubmit} 
            disabled={busy} 
            style={{ marginTop: '20px' }}>
            {busy ? 'Uploading...' : t('data.import.button')}
          </Button></div>

          {isMobile && (<>
            <hr />
            <label>
              {t('data.export.label')}
            </label>
            <Button onClick={requestDataExport}>
              {t('data.export.button')}
            </Button>
          </>)}
        </div>

        {!isMobile && (<>
          <label>
            {t('data.export.label')}
          </label>
          <Button onClick={requestDataExport}>
            {t('data.export.button')}
          </Button>
        </>)}
      </PageContent>
    </Layout>
  )
}

export default Data
