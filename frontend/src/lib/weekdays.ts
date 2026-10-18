// Weekday codes used by recurring tasks' freqByWeekDay, e.g. "MON,WED,FRI"
export const WEEKDAYS = [
  {value: 'MON', short: 'M', name: 'Monday'},
  {value: 'TUE', short: 'T', name: 'Tuesday'},
  {value: 'WED', short: 'W', name: 'Wednesday'},
  {value: 'THU', short: 'T', name: 'Thursday'},
  {value: 'FRI', short: 'F', name: 'Friday'},
  {value: 'SAT', short: 'S', name: 'Saturday'},
  {value: 'SUN', short: 'S', name: 'Sunday'},
]

export const WORKDAYS = ['MON', 'TUE', 'WED', 'THU', 'FRI']

// "MON,TUE,WED,THU,FRI" -> "Weekdays", "MON,WED" -> "Mon, Wed"
export function describeWeekdays(freqByWeekDay: string): string {
  const days = WEEKDAYS.filter(d => freqByWeekDay.includes(d.value)).map(d => d.value)
  if (days.length === 7) return 'Every day'
  if (days.length === 5 && WORKDAYS.every(d => days.includes(d))) return 'Weekdays'
  if (days.length === 2 && days.includes('SAT') && days.includes('SUN')) return 'Weekends'
  if (days.length === 0) return 'Never'
  return WEEKDAYS
    .filter(d => days.includes(d.value))
    .map(d => d.name.substring(0, 3))
    .join(', ')
}
