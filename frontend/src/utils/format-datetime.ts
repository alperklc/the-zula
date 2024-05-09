export const formatDateTime = (dateTime: string, language: string, timeZone?: string) => {
  if (!dateTime) {
    return ''
  }

  return new Date(dateTime).toLocaleTimeString(language, {
    localeMatcher: 'best fit',
    hour12: false,
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    formatMatcher: 'best fit',
    timeZone,
  })
}
