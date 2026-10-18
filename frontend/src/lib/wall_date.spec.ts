import {
  addDay,
  addMonth,
  addWeek,
  formatDayTitle,
  getDateFromWeek,
  getDayOfWeekFromDate,
  getWeekNumber,
  IsBeforeOrEqual,
  parseWallDate,
  toWallDate,
  toWallWeek,
  toWeeklyDate
} from './wall_date';

describe('wall_date', () => {

  describe('parseWallDate', () => {
    it('parses as local midnight of the same calendar day', () => {
      const d = parseWallDate('2026-10-18');
      expect(d.getFullYear()).toBe(2026);
      expect(d.getMonth()).toBe(9);
      expect(d.getDate()).toBe(18);
      expect(d.getHours()).toBe(0);
    });

    it('round-trips through toWallDate', () => {
      expect(toWallDate(parseWallDate('2024-02-29'))).toBe('2024-02-29');
    });
  });

  describe('addDay', () => {
    it('moves within a month', () => {
      expect(addDay('2026-10-18', 1)).toBe('2026-10-19');
      expect(addDay('2026-10-18', -1)).toBe('2026-10-17');
    });

    it('crosses month and year boundaries', () => {
      expect(addDay('2026-10-31', 1)).toBe('2026-11-01');
      expect(addDay('2026-12-31', 1)).toBe('2027-01-01');
      expect(addDay('2027-01-01', -1)).toBe('2026-12-31');
    });

    it('handles leap years', () => {
      expect(addDay('2024-02-28', 1)).toBe('2024-02-29');
      expect(addDay('2026-02-28', 1)).toBe('2026-03-01');
      expect(addDay('2024-03-01', -1)).toBe('2024-02-29');
    });

    it('handles daylight saving time changes', () => {
      expect(addDay('2026-03-29', 1)).toBe('2026-03-30');
      expect(addDay('2026-10-25', 1)).toBe('2026-10-26');
    });
  });

  describe('addMonth', () => {
    it('moves by one month across year boundaries', () => {
      expect(addMonth('2026-10', 1)).toBe('2026-11');
      expect(addMonth('2026-12', 1)).toBe('2027-01');
      expect(addMonth('2026-01', -1)).toBe('2025-12');
    });

    it('moves by more than a year', () => {
      expect(addMonth('2026-01', -13)).toBe('2024-12');
      expect(addMonth('2026-11', 14)).toBe('2028-01');
    });
  });

  describe('ISO weeks', () => {
    it('assigns days around New Year to the right ISO week', () => {
      // 2026 starts on a Thursday, so it has 53 weeks
      expect(toWeeklyDate('2026-12-31')).toBe('2026-W53');
      expect(toWeeklyDate('2027-01-03')).toBe('2026-W53');
      expect(toWeeklyDate('2027-01-04')).toBe('2027-W01');
      // 29 Dec 2025 is the Monday of 2026-W01
      expect(toWeeklyDate('2025-12-28')).toBe('2025-W52');
      expect(toWeeklyDate('2025-12-29')).toBe('2026-W01');
    });

    it('puts Mondays and Sundays in the same week', () => {
      expect(toWeeklyDate('2026-10-19')).toBe('2026-W43');
      expect(toWeeklyDate('2026-10-25')).toBe('2026-W43');
      expect(toWeeklyDate('2026-10-26')).toBe('2026-W44');
    });

    it('leaves week strings unchanged', () => {
      expect(toWeeklyDate('2026-W43')).toBe('2026-W43');
    });

    it('finds the Monday of a week', () => {
      expect(toWallDate(getDateFromWeek(2026, 43))).toBe('2026-10-19');
      expect(toWallDate(getDateFromWeek(2026, 1))).toBe('2025-12-29');
      expect(toWallDate(getDateFromWeek(2026, 53))).toBe('2026-12-28');
      // 1 Jan 2023 is a Sunday, so week 1 starts the next day
      expect(toWallDate(getDateFromWeek(2023, 1))).toBe('2023-01-02');
    });

    it('steps through 52 and 53 week years', () => {
      expect(addWeek('2026-W52', 1)).toBe('2026-W53');
      expect(addWeek('2026-W53', 1)).toBe('2027-W01');
      expect(addWeek('2027-W01', -1)).toBe('2026-W53');
      expect(addWeek('2025-W52', 1)).toBe('2026-W01');
    });

    it('is consistent between toWallWeek and getDateFromWeek for every week of a year', () => {
      for (let week = 1; week <= 53; week++) {
        const weekString = `2026-W${week < 10 ? '0' : ''}${week}`;
        expect(toWallWeek(getDateFromWeek(2026, week))).toBe(weekString);
      }
    });

    it('reads the week number', () => {
      expect(getWeekNumber('2026-W07')).toBe(7);
      expect(getWeekNumber('2026-W53')).toBe(53);
    });
  });

  describe('formatting', () => {
    it('gives the day of the week', () => {
      expect(getDayOfWeekFromDate('2026-10-18')).toBe('Sunday');
      expect(getDayOfWeekFromDate('2026-10-19')).toBe('Monday');
    });

    it('formats day titles', () => {
      expect(formatDayTitle('2026-10-18')).toBe('Sunday, 18 Oct 2026');
      expect(formatDayTitle('2027-01-01')).toBe('Friday, 1 Jan 2027');
    });
  });

  it('compares wall dates', () => {
    expect(IsBeforeOrEqual('2026-10-18', '2026-10-18')).toBeTrue();
    expect(IsBeforeOrEqual('2026-10-17', '2026-10-18')).toBeTrue();
    expect(IsBeforeOrEqual('2026-10-19', '2026-10-18')).toBeFalse();
  });
});
