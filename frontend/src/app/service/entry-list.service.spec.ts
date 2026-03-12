import { TestBed } from '@angular/core/testing';
import {EntryListService} from './entry-list.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('EntryListService', () => {
  let service: EntryListService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });
    service = TestBed.inject(EntryListService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
