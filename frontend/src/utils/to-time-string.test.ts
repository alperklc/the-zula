import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { toTimeString } from './to-time-string'

interface Test {
  dateTime: string
  expectedTimeSince: string
}

describe('toTimeString', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it.each`
    dateTime | expectedTimeSince
    ${'2024-05-10T13:11:18.139Z'} | ${'10.05.2024'}
    ${'2024-05-01T20:51:42.328Z'} | ${'01.05.2024'}
    ${'2024-05-13T13:20:00.328Z'} | ${'15:20'}
  `(
    'displays the given date correctly',
    ({ dateTime, expectedTimeSince }: Test) => {
      const date = new Date('2024-05-13T12:24:39.649Z')
      vi.setSystemTime(date)

      const timeSince = toTimeString(dateTime)

      expect(timeSince).toEqual(expectedTimeSince)
    })
})
