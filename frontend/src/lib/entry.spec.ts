import {groupByDate} from './entry';

describe('groupByDate', () => {
  it('groups consecutive entries with the same date, keeping order', () => {
    const groups = groupByDate([
      {id: 1, createdDate: '2026-10-18'},
      {id: 2, createdDate: '2026-10-18'},
      {id: 3, createdDate: '2026-10-17'},
    ]);

    expect(groups.map(g => g.date)).toEqual(['2026-10-18', '2026-10-17']);
    expect(groups[0].entries.map(e => e.id)).toEqual([1, 2]);
    expect(groups[1].entries.map(e => e.id)).toEqual([3]);
  });

  it('returns no groups for no entries', () => {
    expect(groupByDate([])).toEqual([]);
  });
});
