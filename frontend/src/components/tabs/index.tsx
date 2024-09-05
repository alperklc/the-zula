import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import styles from './index.module.css'

export enum SettingsTabs {
  PROFILE = 'profile',
  DATA = 'data',
}

const Tabs = (props: { selectedTab: SettingsTabs }) => {
  const { t } = useTranslation()

  return (
    <nav className={styles.container}>
      <ol className={styles.tabs}>
        {Object.entries(SettingsTabs).map(([_, tab], index) => (
          <li
            className={`${styles.tab} ${tab === props.selectedTab ? styles.selectedTab : ""}`}
            key={index}
          >
            <Link to={`/settings/${tab}`}>
              {t(`settings_page.tabs.${tab}`)}
            </Link>
          </li>
        ))}
      </ol>
    </nav>
  )
}

export default Tabs
