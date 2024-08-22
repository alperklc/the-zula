import { useIntl } from 'react-intl'
import { Link } from 'react-router-dom'

import styles from './index.module.css'

export enum SettingsTabs {
  PROFILE = 'profile',
}

const Tabs = (props: { selectedTab: SettingsTabs }) => {
  const intl = useIntl()

  return (
    <nav className={styles.container}>
      <ol className={styles.tabs}>
        {Object.entries(SettingsTabs).map(([_, tab], index) => (
          <li
            className={`${styles.tab} ${tab === props.selectedTab ? styles.selectedTab : ""}`}
            key={index}
          >
            <Link to={`/settings/${tab}`}>
              {intl.formatMessage({ id: `settings_page.tabs.${tab}` })}
            </Link>
          </li>
        ))}
      </ol>
    </nav>
  )
}

export default Tabs
