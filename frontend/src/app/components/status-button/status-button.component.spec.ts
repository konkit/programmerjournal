import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { StatusButtonComponent } from './status-button.component';
import {Entry} from "../../../frontend-client";

describe('StatusButtonComponent', () => {
  let component: StatusButtonComponent;
  let fixture: ComponentFixture<StatusButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StatusButtonComponent, HttpClientTestingModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatusButtonComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('entry', {
      id: 1,
      title: 'Test Entry',
      createdDate: '2024-01-01',
      status: 'TODO',
      description: '',
      rank: 0,
      recurringTaskID: 0,
      taskID: '1',
      taskSnoozedUntil: '',
      taskUpdate: ''
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('date entry type checks', () => {
    function createEntry(createdDate: string): Entry {
      return {
        id: 1,
        title: 'Test Entry',
        createdDate: createdDate,
        status: 'TODO',
        description: '',
        rank: 0,
        recurringTaskID: 0,
        taskID: '1',
        taskSnoozedUntil: '',
        taskUpdate: ''
      };
    }

    it('isMonthEntry should return true for valid month entries', () => {
      expect(component.isMonthEntry(createEntry('2024-01'))).toBe(true);
      expect(component.isMonthEntry(createEntry('2024-12'))).toBe(true);
    });

    it('isMonthEntry should return false for invalid month entries', () => {
      expect(component.isMonthEntry(createEntry('2024-1'))).toBe(false);
      expect(component.isMonthEntry(createEntry('2024-01-01'))).toBe(false);
      expect(component.isMonthEntry(createEntry('2024-W01'))).toBe(false);
    });

    it('isDayEntry should return true for valid day entries', () => {
      expect(component.isDayEntry(createEntry('2024-01-01'))).toBe(true);
      expect(component.isDayEntry(createEntry('2024-12-31'))).toBe(true);
    });

    it('isDayEntry should return false for invalid day entries', () => {
      expect(component.isDayEntry(createEntry('2024-1-1'))).toBe(false);
      expect(component.isDayEntry(createEntry('2024-01'))).toBe(false);
      expect(component.isDayEntry(createEntry('2024-W01'))).toBe(false);
    });

    it('isWeekEntry should return true for valid week entries', () => {
      expect(component.isWeekEntry(createEntry('2024-W01'))).toBe(true);
      expect(component.isWeekEntry(createEntry('2024-W52'))).toBe(true);
    });

    it('isWeekEntry should return false for invalid week entries', () => {
      expect(component.isWeekEntry(createEntry('2024-W1'))).toBe(false);
      expect(component.isWeekEntry(createEntry('2024-01-01'))).toBe(false);
      expect(component.isWeekEntry(createEntry('2024-01'))).toBe(false);
    });
  });
});
