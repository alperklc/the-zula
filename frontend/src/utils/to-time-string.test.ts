import { describe, it, expect } from 'vitest';
import { toTimeString } from './to-time-string'

interface Test {
  dateTime: string
  expectedTimeSince: string
}

// const mockReturnValue = new Date('2021-05-13T12:24:39.649Z').getTime()

describe('toTimeString', () => {
  [
    {
      dateTime: '2021-05-10T13:11:18.139Z',
      expectedTimeSince: '10.05',
    },
    {
      dateTime: '2021-05-01T20:51:42.328Z',
      expectedTimeSince: '01.05',
    },
    {
      dateTime: '2021-05-13T13:20:00.328Z',
      expectedTimeSince: '15:20',
    },
  ].map(({ dateTime, expectedTimeSince }: Test) => {
    it('displays the given date correctly', () => {
      const timeSince = toTimeString(dateTime)

      expect(timeSince).toEqual(expectedTimeSince)
    })
  })
})
