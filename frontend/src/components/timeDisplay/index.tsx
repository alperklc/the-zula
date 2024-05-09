import { useUI } from '../../contexts/uiContext';
import { formatDateTime } from '../../utils/format-datetime';
import { toTimeString } from '../../utils/to-time-string';

const TimeDisplay = ({ className, isoDate }: { className?: string; isoDate: string }) => {
  const { language } = useUI()
  const { timeZone } = Intl.DateTimeFormat().resolvedOptions()

  const displayValue = toTimeString(isoDate)
  const toolTipString = formatDateTime(isoDate, language, timeZone)

  return (
    <time className={className} title={toolTipString}>
      {displayValue}
    </time>
  )
}

export default TimeDisplay
