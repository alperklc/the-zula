import React, { useState } from 'react'
import { styles as layoutStyles } from '../layout'
import Logo from '../logo'
import SearchInput from '../search'
import Button from '../form/button'
import icons from '../icons'
import styles from './index.module.css'
import { Link, useNavigate } from 'react-router-dom'
import { QueryContext } from '../../contexts/queryContext'

const MobileHeader = (props: {
  linkTo?: string
  className?: string
  onMenuIconClicked?: () => void
  onFilterSelectionOpen?: () => void
}) => {
  const { query, updateQuery, tagsAsArray } = React.useContext(QueryContext)
  const [headerSearchBarVisible, setHeaderSearchBarVisibility] = useState<boolean>(false)
  const searchInputRef = React.createRef<HTMLInputElement>()
  const navigate = useNavigate()

  const handleOnLogoClick = () => {
    navigate('/')
  }

  const showSearchBar = () => {
    setHeaderSearchBarVisibility(true)
    searchInputRef.current?.focus()
  }
  return (
    <div className={`${styles.container} ${props.className}`}>
      <span className={styles.leftSide}>
        <Button onClick={props.onMenuIconClicked}>
          <icons.Menu />
        </Button>

        <div className={styles.logoContainer}>
          {!(query.q || headerSearchBarVisible) ? (
            <Logo width={154} onClick={handleOnLogoClick} />
          ) : (
            <SearchInput
              autoFocus
              inputref={searchInputRef}
              query={query?.q}
              onQueryUpdate={updateQuery}
              className={layoutStyles.mobileHeaderSearchBar}
            />
          )}
        </div>
      </span>
      <span className={styles.rightSide}>
        <Button outline onClick={showSearchBar} className={layoutStyles.filterButton}>
          <icons.Search />
        </Button>

        <Button outline onClick={props.onFilterSelectionOpen} className={layoutStyles.filterButton}>
          <icons.Filter />
          {tagsAsArray.length > 0 && (
            <span className={layoutStyles.badgeOnButton}>{tagsAsArray.length}</span>
          )}
        </Button>

        <Link to={props.linkTo ?? ''}>
          <Button className={layoutStyles.filterButton}>
            <icons.Plus />
          </Button>
        </Link>
      </span>
    </div>
  )
}

export default MobileHeader
