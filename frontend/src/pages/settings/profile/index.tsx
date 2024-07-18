import React from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { FormattedMessage, useIntl } from 'react-intl'
import Layout, { styles as layoutStyles } from '../../../components/layout'
import MessageBox from '../../../components/messageBox'
import icons from '../../../components/icons'
import Input from '../../../components/form/input'
import Button from '../../../components/form/button'
import { Theme, availableThemes } from '../../../types/layout'
import PageContent from '../../../components/pageContent'
import Breadcrumbs from '../../../components/breadcrumbs'
import Tabs, { SettingsTabs } from '../../../components/tabs'
import { useUI } from '../../../contexts/uiContext'
import { Api, User } from '../../../types/Api'
import { useAuth } from '../../../contexts/authContext'

const Profile = () => {
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })
  
  const navigate = useNavigate()
  const { formatMessage } = useIntl()
  const { switchLanguage, switchTheme, isMobile } = useUI()

  const [loading, setLoading] = React.useState(true);
  const [userData, setUserData] = React.useState<User>();
  const [error, setError] = React.useState<string | null>(null);

  const [firstName, setFirstName] = React.useState<string>('')
  const [lastName, setLastName] = React.useState<string>('')
  const [displayName, setDisplayName] = React.useState<string>('')
  const [email, setEmail] = React.useState<string>('')
  const [language, setLanguage] = React.useState<string>('')
  const [theme, setTheme] = React.useState<string>('')

  const fetchUserProfile = async () => {
    if (!user) {
      return
    }

    try {
      setLoading(true);
      setError(null);

      const { data, status } = await api.api.getUser(user.profile.sub!)

      if (status === 200) {
        setUserData(data);
      } else {
        console.error(data);
        setError(data);
      }

    } catch (e: unknown) {
      console.error(e);
      setError(e as string);
    }
    setLoading(false);
  };

  React.useEffect(() => {
    fetchUserProfile()
  }, [])


  React.useEffect(() => {
    setFirstName(userData?.firstName || '')
    setLastName(userData?.lastName || '')
    setDisplayName(userData?.displayName || '')
    setEmail(userData?.email || '')
    setLanguage(userData?.language || '')
    setTheme(userData?.theme || '')
  }, [userData])

/*  const [{ fetching: saving, error }, saveSettings] = useMutation(
    `
      mutation(
        $Firstname: String
        $username: String
        $email: String
        $language: String
        $theme: String
        $profilePictureUid: String
      ) {
        updateSettings(
          Firstname: $Firstname
          username: $username
          email: $email
          language: $language
          theme: $theme
          profilePictureUid: $profilePictureUid
        )
      }
    `,
  )
*/
  const handleFirstnameInput = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFirstName(event.currentTarget.value)
  }

  const handleLastNameInput = (event: React.ChangeEvent<HTMLInputElement>) => {
    setLastName(event.currentTarget.value)
  }

  const handleEmailInput = (event: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(event.currentTarget.value)
  }

  const handleDisplayName = (event: React.ChangeEvent<HTMLInputElement>) => {
    setDisplayName(event.currentTarget.value)
  }
  
  const onSubmitClicked = () => {
   /* saveSettings({
      Firstname,
      username,
      email,
      language,
      theme,
    }) */
  }

  const handleLanguageInput = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setLanguage(event.target.value)
    switchLanguage(event.target.value)
  }

  const handleThemeInput = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setTheme(event.target.value)
    switchTheme(event.target.value)
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
            <Button primary onClick={onSubmitClicked}>
              <FormattedMessage id='common.buttons.save' />
            </Button>
          </div>
        </>
      }
    >
      {/* !saving && */ error && <MessageBox type='error'>Error</MessageBox>}
      <PageContent
        loading={loading}
        isMobile={isMobile}
        tabs={<Tabs selectedTab={SettingsTabs.PROFILE} />}
      >
        <>
          <Input
            name='firstname'
            type='text'
            label={formatMessage({ id: 'profile.form.firstname' })}
            placeholder={formatMessage({ id: 'profile.form.firstname' })}
            value={firstName}
            onChange={handleFirstnameInput}
          />
          <Input
            name='lastname'
            type='text'
            label={formatMessage({ id: 'profile.form.lastname' })}
            placeholder={formatMessage({ id: 'profile.form.lastname' })}
            value={lastName}
            onChange={handleLastNameInput}
          />
          <Input
            name='displayName'
            type='text'
            label={formatMessage({ id: 'profile.form.displayName' })}
            placeholder={formatMessage({ id: 'profile.form.displayName' })}
            value={displayName}
            onChange={handleDisplayName}
          />
          <Input
            name='email'
            type='text'
            label={formatMessage({ id: 'profile.form.email' })}
            placeholder={formatMessage({ id: 'profile.form.email' })}
            value={email}
            onChange={handleEmailInput}
          />
          {isMobile && (
            <>
              <label htmlFor='language'>Language</label>
              <select name='language' value={language} onChange={handleLanguageInput}>
                <option value='de'>de</option>
                <option value='en'>en</option>
              </select>
              <label htmlFor='theme'>Theme</label>
              <select name='theme' value={theme} onChange={handleThemeInput}>
                {availableThemes.map((theme: Theme, index: number) => (
                  <option value={theme.value} key={index}>
                    {theme.label}
                  </option>
                ))}
              </select>
              <label>&nbsp;</label>
              <Link to='/activity-log'>
                <Button className={layoutStyles.topMargin}>
                  <FormattedMessage id='activity-log.title' />
                  <icons.History height='.8rem' />
                </Button>
              </Link>
            </>
          )}
        </>

        {!isMobile && (
          <>
            <label htmlFor='language'>Language</label>
            <select name='language' value={language} onChange={handleLanguageInput}>
              <option value='de'>de</option>
              <option value='en'>en</option>
            </select>
            <label htmlFor='theme'>Theme</label>
            <select name='theme' value={theme} onChange={handleThemeInput}>
              {availableThemes.map((theme: Theme, index: number) => (
                <option value={theme.value} key={index}>
                  {theme.label}
                </option>
              ))}
            </select>

            <label>&nbsp;</label>

            <Link to='/activity-log'>
              <Button className={layoutStyles.truncatedText}>
                <FormattedMessage id='activity-log.title' />
                <icons.History height='.8rem' />
              </Button>
            </Link>
          </>
        )}
      </PageContent>
    </Layout>
  )
}

export default Profile
