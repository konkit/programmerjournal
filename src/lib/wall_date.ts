
export type WallDate = string;

export function AddDay(date: WallDate, days: number): WallDate {
    let jsDate = new Date(date);
    jsDate.setDate(jsDate.getDate() + days);
    return toWallDate(jsDate)
}

export function IsBeforeOrEqual(a: WallDate, b: WallDate): boolean {
    let jsDateA = new Date(a);
    let jsDateB = new Date(b);

    return jsDateA.getTime() <= jsDateB.getTime();
}

export function Today() {
    return toWallDate(new Date());
}

export function DayOfWeek(date: string) {
    const weekday = ["Sunday", "Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"];

    const d = new Date(date);
    return weekday[d.getDay()];
}

function toWallDate(date: Date) {
    const offset = date.getTimezoneOffset()
    date = new Date(date.getTime() - (offset*60*1000))
    return date.toISOString().split('T')[0]
}



