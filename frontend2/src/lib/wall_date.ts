
export type WallDate = string;
export type WallMonth = string;

export function AddDay(date: WallDate, days: number): WallDate {
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

export function IsBeforeOrEqual(a: WallDate, b: WallDate): boolean {
    let jsDateA = new Date(a);
    let jsDateB = new Date(b);

    return jsDateA.getTime() <= jsDateB.getTime();
}

export function Today(): WallDate {
    return toWallDate(new Date());
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



