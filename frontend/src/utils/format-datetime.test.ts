import { describe, it, expect } from 'vitest';
import { formatDateTime } from './format-datetime'

interface Test {
  language: string
  dateTime: string
  expectedResult: string
}

describe('formatDateTime', () => {
  [
    {
      language: 'de',
      dateTime: '2021-05-09T19:56:42.328Z',
      expectedResult: 'So., 09. Mai 2021, 21:56',
    },
    {
      language: 'de',
      dateTime: '2021-05-01T20:51:42.328Z',
      expectedResult: 'Sa., 01. Mai 2021, 22:51',
    },
    {
      language: 'en',
      dateTime: '2021-05-10T13:00:00.328Z',
      expectedResult: 'Mon, May 10, 2021, 15:00',
    },
  ].map(({ language, dateTime, expectedResult }: Test) => {
    it('displays DateTimeDisplay correctly', () => {
      const result = formatDateTime(dateTime, language, 'Europe/Berlin')

      expect(result).toEqual(expectedResult)
    })
  })
})
