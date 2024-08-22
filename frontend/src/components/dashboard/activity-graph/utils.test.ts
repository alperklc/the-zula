import { describe, it, expect } from 'vitest';
import {
  getGraphDataKeyValuePairs,
  getDaysSequenceOfPastYear,
  getBoundariesOfQuartiles,
} from './utils'

describe('getDaysSequenceOfPastYear', () => {
  it.each`
    isoDate         | expectedFirstDate                          | expectedLastDate
    ${'2021-05-27'} | ${{ dayOfWeek: 3, isoDate: '2020-05-27' }} | ${{ dayOfWeek: 4, isoDate: '2021-05-27' }}
    ${'2021-05-31'} | ${{ dayOfWeek: 0, isoDate: '2020-05-31' }} | ${{ dayOfWeek: 1, isoDate: '2021-05-31' }}
  `(
    'WHEN a date object is given THEN days sequence of past year is generated',
    ({
      isoDate,
      expectedFirstDate,
      expectedLastDate,
    }: {
      isoDate: string
      expectedFirstDate: { isoDate: string; dayOfWeek: number }
      expectedLastDate: { isoDate: string; dayOfWeek: number }
    }) => {
      // arrange
      const date = new Date(isoDate)

      // act
      const result = getDaysSequenceOfPastYear(date)
      const [firstElement] = result
      const lastElement = result.pop()!

      // assert
      expect(firstElement.dayOfWeek).toEqual(expectedFirstDate.dayOfWeek)
      expect(lastElement.dayOfWeek).toEqual(expectedLastDate.dayOfWeek)
      expect(firstElement.isoDate).toEqual(expectedFirstDate.isoDate)
      expect(lastElement.isoDate).toEqual(expectedLastDate.isoDate)
    },
  )
})

describe('getBoundariesOfQuartiles', () => {
  it('doesnt calculate quartiles if the range is not wide enough', () => {
    // arrange
    const data = [
      { date: '2021-04-07', count: 18 },
      { date: '2021-04-08', count: 18 },
      { date: '2021-04-11', count: 21 },
      { date: '2021-04-12', count: 20 },
      { date: '2021-04-16', count: 21 },
      { date: '2021-04-17', count: 19 },
    ]

    // act
    const result = getBoundariesOfQuartiles(data)

    // assert
    expect(result).toStrictEqual([18])
  })

  it('calculates boundaries of quartiles correctly', () => {
    // arrange
    const data = [
      { date: '2021-04-07', count: 4 },
      { date: '2021-04-08', count: 1 },
      { date: '2021-04-11', count: 21 },
      { date: '2021-04-12', count: 7 },
      { date: '2021-04-16', count: 3 },
      { date: '2021-04-17', count: 5 },
    ]

    // act
    const result = getBoundariesOfQuartiles(data)

    // assert
    expect(result).toStrictEqual([6, 11, 16])
  })
})

describe('getGraphDataKeyValuePairs', () => {
  it('creates key value representation of date / count values', () => {
    // arrange
    const data = [
      { date: '2021-04-07', count: 4 },
      { date: '2021-04-08', count: 1 },
      { date: '2021-04-11', count: 21 },
      { date: '2021-04-12', count: 7 },
      { date: '2021-04-16', count: 3 },
      { date: '2021-04-17', count: 5 },
    ]

    // act
    const result = getGraphDataKeyValuePairs(data)

    // assert
    expect(result).toStrictEqual({
      '2021-04-07': { count: 4, quartile: 1 },
      '2021-04-08': { count: 1, quartile: 1 },
      '2021-04-11': { count: 21, quartile: 4 },
      '2021-04-12': { count: 7, quartile: 2 },
      '2021-04-16': { count: 3, quartile: 1 },
      '2021-04-17': { count: 5, quartile: 1 },
    })
  })
})
