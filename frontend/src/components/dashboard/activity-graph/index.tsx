import { useTranslation } from 'react-i18next'
import { getGraphDataKeyValuePairs, getDaysSequenceOfPastYear } from './utils'
import styles from './index.module.css'
import { ActivityOnDate } from '../../../types/Api'

const DaySquare = ({ data, isoDate }: any) => {
  const { t } = useTranslation()

  const title = data?.count
    ? t(
        'dashboard.activity_graph.day_square_with_count',
        { count: data.count, date: isoDate },
      )
    : t(
        'dashboard.activity_graph.day_square_without_count',
        { date: isoDate },
      )

  const evenMonth = new Date(isoDate).getUTCMonth() % 2

  return (
    <li
      data-testid={`daySquare-${isoDate}`}
      data-test-date={isoDate}
      data-level={data?.quartile}
      title={title}
      className={evenMonth ? styles.evenMonth : ''}
    ></li>
  )
}

export const ActivityGraph = ({
  data,
  currentDate = new Date(),
}: {
  data: ActivityOnDate[]
  currentDate?: Date
}) => {
  const { t } = useTranslation()

  const daysPastYear = getDaysSequenceOfPastYear(currentDate)
  const dateCountKeyValue = getGraphDataKeyValuePairs(data)

  // zero is sunday, therefore we have this ternary expression
  const lengthOfEmptySquares = daysPastYear[0].dayOfWeek === 0 ? 6 : daysPastYear[0].dayOfWeek - 1

  return (
    <div className={styles.container}>
      <div className={styles.graph}>
        <ul className={styles.days}>
          <li>{t('dashboard.activity_graph.days.monday')}</li>
          <li>{t('dashboard.activity_graph.days.tuesday')}</li>
          <li>{t('dashboard.activity_graph.days.wednesday')}</li>
          <li>{t('dashboard.activity_graph.days.thursday')}</li>
          <li>{t('dashboard.activity_graph.days.friday')}</li>
          <li>{t('dashboard.activity_graph.days.saturday')}</li>
          <li>{t('dashboard.activity_graph.days.sunday')}</li>
        </ul>
        <ul className={styles.squares} data-testid='squares'>
          {Array(lengthOfEmptySquares)
            .fill(0)
            .map((_, i) => (
              <li key={i} data-testid='emptySquare' className={styles.emptySquare} />
            ))}
          {daysPastYear.map((item, i) => (
            <DaySquare data={dateCountKeyValue?.[item.isoDate]} isoDate={item.isoDate} key={i} />
          ))}
        </ul>
      </div>
    </div>
  )
}
