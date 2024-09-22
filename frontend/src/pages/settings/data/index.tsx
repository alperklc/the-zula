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

  const [file, setFile] = React.useState(null);
  const [uploading, setUploading] = React.useState(false);

  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const handleSubmit = async () => {
    if (!file) {
      showToast('Please select a file before uploading.', "error");
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      setUploading(true);
      const response = await api.api.importData({ file });
      if (response.status === 200) {
        showToast('File uploaded successfully!', "success");
      } else {
        showToast('File upload failed.', "error");
      }
    } catch (error) {
      showToast('Error during upload: ' + error, "error");
    } finally {
      setUploading(false);
    }
  };

  const requestDataExport = () => {

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
        loading={false}
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
            disabled={uploading}
          />

          <div><Button 
            onClick={handleSubmit} 
            disabled={uploading} 
            style={{ marginTop: '20px' }}>
            {uploading ? 'Uploading...' : t('data.import.button')}
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
