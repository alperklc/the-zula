import { addDays, getDaysInYear } from 'date-fns'
import { ActivityOnDate } from '../../../types/Api'

export const getDaysSequenceOfPastYear = (date: Date) => {
  const daysInYear = getDaysInYear(date)
  const aYearAgo = addDays(date, -daysInYear)

  return Array(daysInYear + 1)
    .fill(0)
    .map((_, j) => {
      const date = addDays(aYearAgo, j)

      return {
        isoDate: date.toISOString().substring(0, 10),
        dayOfWeek: date.getUTCDay(),
      }
    })
}

export const getBoundariesOfQuartiles = (data: ActivityGraphData[]) => {
  const values = data?.map((el) => el.count) || []

  const min = Math.min(...values)
  const max = Math.max(...values)

  const range = max - min

  return range < 4
    ? [min]
    : Array(3)
        .fill(min)
        .map((i, j) => i + (range * (j + 1)) / 4)
}

export const getQuartileOfValue = (value: number, boundaries: number[]) => {
  if (boundaries.length !== 3) {
    return 1
  }

  const [q2, q3, q4] = boundaries

  if (value < q2) {
    return 1
  } else if (q2 <= value && value < q3) {
    return 2
  } else if (q3 <= value && value < q4) {
    return 3
  }

  return 4
}

export type ActivityGraphData = { date: string; count: number }

export const getGraphDataKeyValuePairs = (data: ActivityOnDate[]) => {
  const quartileBoundaries = getBoundariesOfQuartiles(data as ActivityGraphData[])

  return (data as ActivityGraphData[])
    ?.map((element: ActivityGraphData) => ({
      [element.date]: {
        count: element.count,
        quartile: getQuartileOfValue(element.count, quartileBoundaries),
      },
    }))
    .reduce((prev, next) => ({ ...prev, ...next }), {})
}
