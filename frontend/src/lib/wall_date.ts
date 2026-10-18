export type WallDate = string;
export type WallMonth = string;
export type WallWeek = string;

// Parses "2024-12-31" as local midnight. `new Date("2024-12-31")` would parse it as UTC midnight,
// which is still the previous day in time zones west of UTC.
export function parseWallDate(date: WallDate): Date {
  const year = parseInt(date.substring(0, 4))
  const month = parseInt(date.substring(5, 7))
  const day = parseInt(date.substring(8, 10))
  return new Date(year, month - 1, day)
}

export function addDay(date: WallDate, days: number): WallDate {
  let jsDate = parseWallDate(date);
  jsDate.setDate(jsDate.getDate() + days);
  return toWallDate(jsDate)
}

export function addMonth(date: WallMonth, months: number): WallMonth {
  let year = parseInt(date.substring(0, 4))
  let month = parseInt(date.substring(5, 7))

  const monthIndex = year * 12 + (month - 1) + months
  year = Math.floor(monthIndex / 12)
  month = monthIndex % 12 + 1

  return `${year}-${month < 10 ? "0" : ""}${month}`
}

export function addWeek(date: WallWeek, weeks: number): WallWeek {
  // Parse 2024-W01
  let year = parseInt(date.substring(0, 4))
  let week = parseInt(date.substring(6, 8))

  // Go through the Monday of that week, so 52/53-week years are handled by toWallWeek
  let monday = getDateFromWeek(year, week);
  monday.setDate(monday.getDate() + (weeks * 7));
  return toWallWeek(monday);
}

export function getDateFromWeek(year: number, week: number): Date {
  const simple = new Date(year, 0, 1 + (week - 1) * 7);
  const dow = simple.getDay();
  const ISOweekStart = simple;
  if (dow <= 4)
    ISOweekStart.setDate(simple.getDate() - simple.getDay() + 1);
  else
    ISOweekStart.setDate(simple.getDate() + 8 - simple.getDay());
  return ISOweekStart;
}

export function IsBeforeOrEqual(a: WallDate, b: WallDate): boolean {
  let jsDateA = parseWallDate(a);
  let jsDateB = parseWallDate(b);

  return jsDateA.getTime() <= jsDateB.getTime();
}

export function Today(): WallDate {
  return toWallDate(new Date());
}

export function ThisWeek(): WallWeek {
  return toWallWeek(new Date());
}

export function ThisMonth(): WallMonth {
  return toWallMonth(new Date());
}

export function getDayOfWeekFromDate(date: string) {
  const weekday = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

  const d = parseWallDate(date);
  return weekday[d.getDay()];
}

export function getMonthFromDate(date: string) {
  const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
  const year = parseInt(date.substring(0, 4))
  const month = parseInt(date.substring(5, 7))

  return months[month - 1]
}

// "2026-10-18" -> "Sunday, 18 Oct 2026"
export function formatDayTitle(date: WallDate): string {
  const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
  const d = parseWallDate(date)
  return `${getDayOfWeekFromDate(date)}, ${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`
}

// "2026-W43" -> 43
export function getWeekNumber(date: WallWeek): number {
  return parseInt(date.substring(6, 8))
}

export function getYearFromDate(date: string) {
  return parseInt(date.substring(0, 4))
}

export function toWallDate(date: Date): WallDate {
  const offset = date.getTimezoneOffset()
  date = new Date(date.getTime() - (offset * 60 * 1000))
  return date.toISOString().split('T')[0]
}

export function toWallMonth(date: Date): WallMonth {
  const offset = date.getTimezoneOffset()
  date = new Date(date.getTime() - (offset * 60 * 1000))
  return date.toISOString().substring(0, 7)
}

export function toWallWeek(date: Date): WallWeek {
  // Copy date so don't modify original
  date = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  // Set to nearest Thursday: current date + 4 - current day number
  // Make Sunday's day number 7
  date.setUTCDate(date.getUTCDate() + 4 - (date.getUTCDay() || 7));
  // Get first day of year
  var yearStart = new Date(Date.UTC(date.getUTCFullYear(), 0, 1));
  // Calculate full weeks to nearest Thursday
  var weekNo = Math.ceil((((date.getTime() - yearStart.getTime()) / 86400000) + 1) / 7);
  // Return array of year and week number
  return `${date.getUTCFullYear()}-W${weekNo < 10 ? '0' : ''}${weekNo}`;
}

export function getWeekString(date: Date): WallWeek {
  return toWallWeek(date);
}

export function isMonthDate(date: string): boolean {
  return date.length === 7; // 2024-12
}

export function isDayDate(date: string): boolean {
  return date.length === 10; // 2024-12-12
}

export function isWeekDate(date: string): boolean {
  return date.length === 8 && date.includes('W'); // 2024-W01
}

export function toWeeklyDate(date: WallDate): WallWeek {
  if (isWeekDate(date)) {
    return date;
  }
  return toWallWeek(parseWallDate(date));
}

export function getWeekOfYear(date: Date): number {
  // Copy date so don't modify original
  date = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  // Set to nearest Thursday: current date + 4 - current day number
  // Make Sunday's day number 7
  date.setUTCDate(date.getUTCDate() + 4 - (date.getUTCDay() || 7));
  // Get first day of year
  var yearStart = new Date(Date.UTC(date.getUTCFullYear(), 0, 1));
  // Calculate full weeks to nearest Thursday
  return Math.ceil((((date.getTime() - yearStart.getTime()) / 86400000) + 1) / 7);
}
