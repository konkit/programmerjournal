
export type WallDate = string;
export type WallMonth = string;
export type WallWeek = string;

export function addDay(date: WallDate, days: number): WallDate {
    let jsDate = new Date(date);
    jsDate.setDate(jsDate.getDate() + days);
    return toWallDate(jsDate)
}

export function addMonth(date: WallMonth, months: number): WallMonth {
  let year = parseInt(date.substring(0, 4))
  let month = parseInt(date.substring(5, 7))

  month += months;

  if (month > 12) {
    year++;
    month=1;
  } else if (month < 1) {
    year--;
    month=12;
  }

  return `${year}-${month < 10 ? "0" : "" }${month}`
}

export function addWeek(date: WallWeek, weeks: number): WallWeek {
  // Parse 2024-W01
  let year = parseInt(date.substring(0, 4))
  let week = parseInt(date.substring(6, 8))

  // Simple logic, might need more robust library like date-fns or luxon for correct ISO week math
  // But for now, let's try to approximate or use a library if available.
  // Since we don't have external libs easily, let's do a simple increment and handle overflow.
  // A year has 52 or 53 weeks.
  // This is tricky without a library.
  // Let's convert to a date (Monday of that week), add 7*weeks days, and convert back to week string.

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
    let jsDateA = new Date(a);
    let jsDateB = new Date(b);

    return jsDateA.getTime() <= jsDateB.getTime();
}

export function Today(): WallDate {
    return toWallDate(new Date());
}

export function StartOfThisWeek(): WallDate {
  return toWallDate(getMonday(new Date()));
}

export function ThisWeek(): WallWeek {
  return toWallWeek(new Date());
}

function getMonday(d: Date) {
  const day = d.getDay();
  const diff = d.getDate() - day + (day === 0 ? -6 : 1);
  return new Date(d.setDate(diff));
}

export function ThisMonth(): WallMonth {
  return toWallMonth(new Date());
}

export function getDayOfWeekFromDate(date: string) {
    const weekday = ["Sunday", "Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"];

    const d = new Date(date);
    return weekday[d.getDay()];
}

export function getMonthFromDate(date: string) {
  const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
  const year = parseInt(date.substring(0, 4))
  const month = parseInt(date.substring(5, 7))

  return months[month-1]
}

export function getYearFromDate(date: string) {
  return parseInt(date.substring(0, 4))
}

export function toWallDate(date: Date): WallDate {
    const offset = date.getTimezoneOffset()
    date = new Date(date.getTime() - (offset*60*1000))
    return date.toISOString().split('T')[0]
}

export function toWallMonth(date: Date): WallMonth {
  const offset = date.getTimezoneOffset()
  date = new Date(date.getTime() - (offset*60*1000))
  return date.toISOString().substring(0, 7)
}

export function toWallWeek(date: Date): WallWeek {
  // Copy date so don't modify original
  date = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  // Set to nearest Thursday: current date + 4 - current day number
  // Make Sunday's day number 7
  date.setUTCDate(date.getUTCDate() + 4 - (date.getUTCDay()||7));
  // Get first day of year
  var yearStart = new Date(Date.UTC(date.getUTCFullYear(),0,1));
  // Calculate full weeks to nearest Thursday
  var weekNo = Math.ceil(( ( (date.getTime() - yearStart.getTime()) / 86400000) + 1)/7);
  // Return array of year and week number
  return `${date.getUTCFullYear()}-W${weekNo < 10 ? '0' : ''}${weekNo}`;
}

export function getWeekString(date: Date): WallWeek {
  return toWallWeek(date);
}

export function toWeeklyDate(date: WallDate): WallWeek {
  return toWallWeek(new Date(date));
}

export function getWeekOfYear(date: Date): number {
  // Copy date so don't modify original
  date = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  // Set to nearest Thursday: current date + 4 - current day number
  // Make Sunday's day number 7
  date.setUTCDate(date.getUTCDate() + 4 - (date.getUTCDay()||7));
  // Get first day of year
  var yearStart = new Date(Date.UTC(date.getUTCFullYear(),0,1));
  // Calculate full weeks to nearest Thursday
  return Math.ceil((((date.getTime() - yearStart.getTime()) / 86400000) + 1) / 7);
}
