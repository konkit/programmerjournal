import {describeWeekdays} from './weekdays';

describe('describeWeekdays', () => {
  it('names common schedules', () => {
    expect(describeWeekdays('MON,TUE,WED,THU,FRI,SAT,SUN')).toBe('Every day');
    expect(describeWeekdays('MON,TUE,WED,THU,FRI')).toBe('Weekdays');
    expect(describeWeekdays('SAT,SUN')).toBe('Weekends');
  });

  it('lists other days in week order', () => {
    expect(describeWeekdays('FRI,MON')).toBe('Mon, Fri');
    expect(describeWeekdays('WED')).toBe('Wed');
  });

  it('handles an empty schedule', () => {
    expect(describeWeekdays('')).toBe('Never');
  });
});
